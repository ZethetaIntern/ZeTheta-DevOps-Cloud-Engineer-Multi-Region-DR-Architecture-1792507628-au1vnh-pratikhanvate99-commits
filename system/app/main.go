package main

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/segmentio/kafka-go"
)

type Config struct {
	DatabaseURL  string
	Brokers      []string
	Topic        string
	Group        string
	HTTPAddr     string
	Region       string
	MerchantKeys map[string]string
}

type Payment struct {
	ID          string    `json:"payment_id"`
	MerchantID  string    `json:"merchant_id"`
	AmountMinor int64     `json:"amount_minor"`
	Currency    string    `json:"currency"`
	State       string    `json:"state"`
	CreatedAt   time.Time `json:"created_at"`
}

type Event struct {
	EventID      string    `json:"event_id"`
	PaymentID    string    `json:"payment_id"`
	MerchantID   string    `json:"merchant_id"`
	AmountMinor  int64     `json:"amount_minor"`
	Currency     string    `json:"currency"`
	EventType    string    `json:"event_type"`
	Outcome      string    `json:"outcome,omitempty"`
	SourceRegion string    `json:"source_region"`
	CreatedAt    time.Time `json:"created_at"`
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func loadConfig() Config {
	keys := map[string]string{}
	for _, pair := range strings.Split(os.Getenv("MERCHANT_API_KEYS"), ",") {
		parts := strings.SplitN(strings.TrimSpace(pair), "=", 2)
		if len(parts) == 2 && parts[0] != "" && parts[1] != "" {
			keys[parts[0]] = parts[1]
		}
	}
	brokers := strings.Split(env("KAFKA_BROKERS", "kafka:9092"), ",")
	return Config{DatabaseURL: env("DATABASE_URL", ""), Brokers: brokers, Topic: env("KAFKA_TOPIC", "paysecure.payment-events"), Group: env("KAFKA_GROUP", "paysecure-demo-processor"), HTTPAddr: env("HTTP_ADDR", ":8080"), Region: env("REGION", "local"), MerchantKeys: keys}
}

func openDB(ctx context.Context, cfg Config) (*sql.DB, error) {
	if cfg.DatabaseURL == "" {
		return nil, errors.New("DATABASE_URL is required")
	}
	db, err := sql.Open("pgx", cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(40)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(30 * time.Minute)
	if err = db.PingContext(ctx); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

func uuid() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	h := hex.EncodeToString(b)
	return h[:8] + "-" + h[8:12] + "-" + h[12:16] + "-" + h[16:20] + "-" + h[20:], nil
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	cfg := loadConfig()
	if len(os.Args) < 2 {
		slog.Error("usage: paysecure api|publisher|processor")
		os.Exit(2)
	}
	ctx := context.Background()
	db, err := openDB(ctx, cfg)
	if err != nil {
		slog.Error("database connection failed", "error", err)
		os.Exit(1)
	}
	defer db.Close()
	switch os.Args[1] {
	case "api":
		err = runAPI(ctx, db, cfg)
	case "publisher":
		err = runPublisher(ctx, db, cfg)
	case "processor":
		err = runProcessor(ctx, db, cfg)
	case "settlement":
		err = runSettlement(ctx, db, cfg)
	default:
		err = fmt.Errorf("unknown command %q", os.Args[1])
	}
	if err != nil {
		slog.Error("service stopped", "error", err)
		os.Exit(1)
	}
}

func authorized(cfg Config, r *http.Request) string {
	merchant := r.Header.Get("X-Merchant-ID")
	key := r.Header.Get("X-API-Key")
	if merchant == "" || key == "" {
		return ""
	}
	expected, ok := cfg.MerchantKeys[merchant]
	if !ok || expected != key {
		return ""
	}
	return merchant
}

func runAPI(ctx context.Context, db *sql.DB, cfg Config) error {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health/live", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok\n"))
	})
	mux.HandleFunc("GET /health/ready", func(w http.ResponseWriter, r *http.Request) {
		if err := db.PingContext(r.Context()); err != nil {
			http.Error(w, "not ready", http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ready\n"))
	})
	mux.HandleFunc("POST /v1/payments", func(w http.ResponseWriter, r *http.Request) {
		merchant := authorized(cfg, r)
		if merchant == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		key := strings.TrimSpace(r.Header.Get("Idempotency-Key"))
		if len(key) < 8 || len(key) > 128 {
			http.Error(w, "Idempotency-Key must contain 8-128 characters", http.StatusBadRequest)
			return
		}
		var input struct {
			AmountMinor int64  `json:"amount_minor"`
			Currency    string `json:"currency"`
		}
		r.Body = http.MaxBytesReader(w, r.Body, 16*1024)
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()
		if err := dec.Decode(&input); err != nil {
			http.Error(w, "invalid JSON body", http.StatusBadRequest)
			return
		}
		if input.AmountMinor <= 0 || input.AmountMinor > 100_000_000_000 || input.Currency != "INR" {
			http.Error(w, "amount_minor must be positive and currency must be INR", http.StatusBadRequest)
			return
		}
		payment, created, err := createPayment(r.Context(), db, cfg, merchant, key, input.AmountMinor, input.Currency)
		if errors.Is(err, errIdempotencyConflict) {
			http.Error(w, "idempotency key reused with different request", http.StatusConflict)
			return
		}
		if err != nil {
			slog.Error("create payment failed", "error", err)
			http.Error(w, "payment unavailable", http.StatusServiceUnavailable)
			return
		}
		status := http.StatusAccepted
		if !created {
			status = http.StatusOK
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(payment)
	})
	mux.HandleFunc("GET /v1/payments/{id}", func(w http.ResponseWriter, r *http.Request) {
		merchant := authorized(cfg, r)
		if merchant == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		var p Payment
		err := db.QueryRowContext(r.Context(), `SELECT payment_id::text, merchant_id, amount_minor, currency, state, created_at FROM payments WHERE payment_id=$1 AND merchant_id=$2`, r.PathValue("id"), merchant).Scan(&p.ID, &p.MerchantID, &p.AmountMinor, &p.Currency, &p.State, &p.CreatedAt)
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		if err != nil {
			http.Error(w, "payment unavailable", http.StatusServiceUnavailable)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(p)
	})
	srv := &http.Server{Addr: cfg.HTTPAddr, Handler: securityHeaders(mux), ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 10 * time.Second, WriteTimeout: 10 * time.Second, IdleTimeout: 60 * time.Second}
	go func() {
		<-ctx.Done()
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = srv.Shutdown(c)
	}()
	slog.Info("api listening", "addr", cfg.HTTPAddr, "region", cfg.Region)
	err := srv.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}

func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}

var errIdempotencyConflict = errors.New("idempotency key reused with different request")

func createPayment(ctx context.Context, db *sql.DB, cfg Config, merchant, key string, amount int64, currency string) (Payment, bool, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return Payment{}, false, err
	}
	defer tx.Rollback()
	id, err := uuid()
	if err != nil {
		return Payment{}, false, err
	}
	eventID, err := uuid()
	if err != nil {
		return Payment{}, false, err
	}
	var p Payment
	err = tx.QueryRowContext(ctx, `INSERT INTO payments(payment_id,merchant_id,idempotency_key,amount_minor,currency,state) VALUES($1,$2,$3,$4,$5,'PENDING') ON CONFLICT(merchant_id,idempotency_key) DO NOTHING RETURNING payment_id::text,merchant_id,amount_minor,currency,state,created_at`, id, merchant, key, amount, currency).Scan(&p.ID, &p.MerchantID, &p.AmountMinor, &p.Currency, &p.State, &p.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		err = tx.QueryRowContext(ctx, `SELECT payment_id::text,merchant_id,amount_minor,currency,state,created_at FROM payments WHERE merchant_id=$1 AND idempotency_key=$2`, merchant, key).Scan(&p.ID, &p.MerchantID, &p.AmountMinor, &p.Currency, &p.State, &p.CreatedAt)
		if err != nil {
			return Payment{}, false, err
		}
		if p.AmountMinor != amount || p.Currency != currency {
			return Payment{}, false, errIdempotencyConflict
		}
		if err = tx.Commit(); err != nil {
			return Payment{}, false, err
		}
		return p, false, nil
	}
	if err != nil {
		return Payment{}, false, err
	}
	event := Event{EventID: eventID, PaymentID: p.ID, MerchantID: merchant, AmountMinor: amount, Currency: currency, EventType: "payment.requested", SourceRegion: cfg.Region, CreatedAt: time.Now().UTC()}
	b, err := json.Marshal(event)
	if err != nil {
		return Payment{}, false, err
	}
	if _, err = tx.ExecContext(ctx, `INSERT INTO outbox(event_id,aggregate_id,event_type,payload) VALUES($1,$2,$3,$4::jsonb)`, eventID, p.ID, event.EventType, string(b)); err != nil {
		return Payment{}, false, err
	}
	if err = tx.Commit(); err != nil {
		return Payment{}, false, err
	}
	return p, true, nil
}

func runPublisher(ctx context.Context, db *sql.DB, cfg Config) error {
	w := &kafka.Writer{Addr: kafka.TCP(cfg.Brokers...), Topic: cfg.Topic, Balancer: &kafka.Hash{}, RequiredAcks: kafka.RequireAll, Async: false, WriteTimeout: 10 * time.Second, ReadTimeout: 10 * time.Second}
	defer w.Close()
	ticker := time.NewTicker(250 * time.Millisecond)
	defer ticker.Stop()
	for {
		if err := publishBatch(ctx, db, w); err != nil {
			slog.Error("outbox publish failed", "error", err)
		}
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
		}
	}
}

func publishBatch(ctx context.Context, db *sql.DB, w *kafka.Writer) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	rows, err := tx.QueryContext(ctx, `SELECT event_id::text,aggregate_id::text,payload::text FROM outbox WHERE published_at IS NULL ORDER BY created_at LIMIT 100 FOR UPDATE SKIP LOCKED`)
	if err != nil {
		tx.Rollback()
		return err
	}
	type row struct{ id, payment, payload string }
	batch := []row{}
	for rows.Next() {
		var x row
		if err = rows.Scan(&x.id, &x.payment, &x.payload); err != nil {
			rows.Close()
			tx.Rollback()
			return err
		}
		batch = append(batch, x)
	}
	if err = rows.Err(); err != nil {
		rows.Close()
		tx.Rollback()
		return err
	}
	rows.Close()
	for _, x := range batch {
		if err = w.WriteMessages(ctx, kafka.Message{Key: []byte(x.payment), Value: []byte(x.payload), Time: time.Now().UTC()}); err != nil {
			tx.Rollback()
			return err
		}
		if _, err = tx.ExecContext(ctx, `UPDATE outbox SET published_at=now() WHERE event_id=$1`, x.id); err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func runProcessor(ctx context.Context, db *sql.DB, cfg Config) error {
	r := kafka.NewReader(kafka.ReaderConfig{Brokers: cfg.Brokers, Topic: cfg.Topic, GroupID: cfg.Group, MinBytes: 1, MaxBytes: 10e6, CommitInterval: 0, StartOffset: kafka.FirstOffset})
	defer r.Close()
	for {
		msg, err := r.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return nil
			}
			return err
		}
		var e Event
		if err = json.Unmarshal(msg.Value, &e); err != nil {
			slog.Error("invalid event; not committing offset", "error", err)
			return err
		}
		if err = consumeEvent(ctx, db, e, "transaction-processor"); err != nil {
			slog.Error("event processing failed", "event_id", e.EventID, "error", err)
			return err
		}
		if err = r.CommitMessages(ctx, msg); err != nil {
			return err
		}
	}
}

func consumeEvent(ctx context.Context, db *sql.DB, e Event, consumer string) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	res, err := tx.ExecContext(ctx, `INSERT INTO consumed_events(consumer_name,event_id) VALUES($1,$2) ON CONFLICT DO NOTHING`, consumer, e.EventID)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return tx.Commit()
	}
	if e.EventType == "payment.requested" {
		// Demo-only outcome: this is not an authorization from a bank, network or acquirer.
		_, err = tx.ExecContext(ctx, `UPDATE payments SET state='SIMULATED_APPROVED',updated_at=now() WHERE payment_id=$1 AND state='PENDING'`, e.PaymentID)
		if err != nil {
			return err
		}
		resultID, err := uuid()
		if err != nil {
			return err
		}
		result := e
		result.EventID = resultID
		result.EventType = "payment.simulated_result"
		result.Outcome = "SIMULATED_APPROVED"
		result.CreatedAt = time.Now().UTC()
		b, err := json.Marshal(result)
		if err != nil {
			return err
		}
		if _, err = tx.ExecContext(ctx, `INSERT INTO outbox(event_id,aggregate_id,event_type,payload) VALUES($1,$2,$3,$4::jsonb)`, resultID, e.PaymentID, result.EventType, string(b)); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func runSettlement(ctx context.Context, db *sql.DB, cfg Config) error {
	r := kafka.NewReader(kafka.ReaderConfig{Brokers: cfg.Brokers, Topic: cfg.Topic, GroupID: "paysecure-demo-settlement", MinBytes: 1, MaxBytes: 10e6, CommitInterval: 0, StartOffset: kafka.FirstOffset})
	defer r.Close()
	for {
		msg, err := r.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return nil
			}
			return err
		}
		var e Event
		if err = json.Unmarshal(msg.Value, &e); err != nil {
			return err
		}
		if e.EventType == "payment.simulated_result" {
			if err = queueSettlement(ctx, db, e); err != nil {
				slog.Error("settlement simulation failed", "event_id", e.EventID, "error", err)
				return err
			}
		}
		if err = r.CommitMessages(ctx, msg); err != nil {
			return err
		}
	}
}

func queueSettlement(ctx context.Context, db *sql.DB, e Event) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	res, err := tx.ExecContext(ctx, `INSERT INTO consumed_events(consumer_name,event_id) VALUES('settlement-engine',$1) ON CONFLICT DO NOTHING`, e.EventID)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return tx.Commit()
	}
	if e.Outcome != "SIMULATED_APPROVED" {
		return tx.Commit()
	}
	if _, err = tx.ExecContext(ctx, `INSERT INTO settlement_simulations(payment_id,status) VALUES($1,'QUEUED') ON CONFLICT(payment_id) DO NOTHING`, e.PaymentID); err != nil {
		return err
	}
	return tx.Commit()
}
