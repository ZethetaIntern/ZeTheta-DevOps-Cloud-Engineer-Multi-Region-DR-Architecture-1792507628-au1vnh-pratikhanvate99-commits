package main

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

func TestUUIDFormatAndUniqueness(t *testing.T) {
	pattern := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
	first, err := uuid()
	if err != nil {
		t.Fatal(err)
	}
	second, err := uuid()
	if err != nil {
		t.Fatal(err)
	}
	if !pattern.MatchString(first) || !pattern.MatchString(second) {
		t.Fatalf("UUID does not match RFC 4122 v4 format: %q, %q", first, second)
	}
	if first == second {
		t.Fatal("uuid() returned a duplicate value")
	}
}

func TestAuthorizedRequiresMatchingMerchantAndKey(t *testing.T) {
	cfg := Config{MerchantKeys: map[string]string{"m-1": "key-1"}}
	cases := []struct {
		name, merchant, key, want string
	}{
		{"valid", "m-1", "key-1", "m-1"},
		{"unknown merchant", "m-2", "key-1", ""},
		{"wrong key", "m-1", "wrong", ""},
		{"missing merchant", "", "key-1", ""},
		{"missing key", "m-1", "", ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest("POST", "/v1/payments", nil)
			r.Header.Set("X-Merchant-ID", tc.merchant)
			r.Header.Set("X-API-Key", tc.key)
			if got := authorized(cfg, r); got != tc.want {
				t.Fatalf("authorized() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestSecurityHeaders(t *testing.T) {
	h := securityHeaders(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(204) }))
	r := httptest.NewRequest("GET", "/health/live", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != 204 {
		t.Fatalf("status = %d, want 204", w.Code)
	}
	for name, want := range map[string]string{"X-Content-Type-Options": "nosniff", "Cache-Control": "no-store"} {
		if got := w.Header().Get(name); got != want {
			t.Errorf("%s = %q, want %q", name, got, want)
		}
	}
}
