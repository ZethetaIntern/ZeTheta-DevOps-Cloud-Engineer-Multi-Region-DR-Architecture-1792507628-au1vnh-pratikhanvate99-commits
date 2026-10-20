CREATE TABLE IF NOT EXISTS payments (
  payment_id UUID PRIMARY KEY,
  merchant_id TEXT NOT NULL,
  idempotency_key TEXT NOT NULL,
  amount_minor BIGINT NOT NULL CHECK (amount_minor > 0),
  currency CHAR(3) NOT NULL,
  state TEXT NOT NULL CHECK (state IN ('PENDING','SIMULATED_APPROVED','SIMULATED_DECLINED')),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (merchant_id, idempotency_key)
);
CREATE TABLE IF NOT EXISTS outbox (
  event_id UUID PRIMARY KEY,
  aggregate_id UUID NOT NULL REFERENCES payments(payment_id),
  event_type TEXT NOT NULL,
  payload JSONB NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  published_at TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS outbox_unpublished_idx ON outbox(created_at) WHERE published_at IS NULL;
CREATE TABLE IF NOT EXISTS consumed_events (
  consumer_name TEXT NOT NULL,
  event_id UUID NOT NULL,
  PRIMARY KEY (consumer_name, event_id),
  consumed_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE IF NOT EXISTS settlement_simulations (
  payment_id UUID PRIMARY KEY REFERENCES payments(payment_id),
  status TEXT NOT NULL CHECK (status IN ('QUEUED','RECONCILED')),
  queued_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

