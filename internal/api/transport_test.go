package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync/atomic"
	"testing"
	"time"
)

func TestTransportAddsAuthAndUserAgent(t *testing.T) {
	var (
		gotAuth string
		gotUA   string
	)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		gotUA = r.Header.Get("User-Agent")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("{}"))
	}))
	defer srv.Close()

	c, err := New(Options{BaseURL: srv.URL, Token: "secret", UserAgent: "test-ua"})
	if err != nil {
		t.Fatal(err)
	}
	// Any endpoint with no required params works for this transport check.
	_, err = c.BillingGetPaymentMethodListV1WithResponse(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if gotAuth != "Bearer secret" {
		t.Errorf("Authorization header = %q, want Bearer secret", gotAuth)
	}
	if gotUA != "test-ua" {
		t.Errorf("User-Agent header = %q, want test-ua", gotUA)
	}
}

func TestTransportRetriesOn429(t *testing.T) {
	var hits atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := hits.Add(1)
		if n < 3 {
			w.Header().Set("Retry-After", "0")
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("{}"))
	}))
	defer srv.Close()

	c, err := New(Options{BaseURL: srv.URL, Token: "t", MaxRetries: 3})
	if err != nil {
		t.Fatal(err)
	}
	resp, err := c.BillingGetPaymentMethodListV1WithResponse(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode() != http.StatusOK {
		t.Errorf("status = %d, want 200", resp.StatusCode())
	}
	if got := hits.Load(); got != 3 {
		t.Errorf("hits = %d, want 3", got)
	}
}

func TestParseRetryAfter(t *testing.T) {
	if d := parseRetryAfter(""); d != 0 {
		t.Errorf("empty: got %v", d)
	}
	if d := parseRetryAfter("2"); d != 2*time.Second {
		t.Errorf("seconds: got %v", d)
	}
	if d := parseRetryAfter(strconv.Itoa(0)); d != 0 {
		t.Errorf("zero: got %v", d)
	}
}
