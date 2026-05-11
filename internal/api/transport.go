package api

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/rizaleow/hostinger-cli/internal/version"
)

// DefaultBaseURL is the production Hostinger API host.
const DefaultBaseURL = "https://developers.hostinger.com"

// Options controls the construction of a Hostinger API client.
type Options struct {
	BaseURL    string
	Token      string
	UserAgent  string
	HTTPClient *http.Client
	MaxRetries int
}

// New builds a typed ClientWithResponses with auth, UA, and retry middleware.
func New(opts Options) (*ClientWithResponses, error) {
	if opts.BaseURL == "" {
		opts.BaseURL = DefaultBaseURL
	}
	if opts.UserAgent == "" {
		opts.UserAgent = version.UserAgent()
	}
	if opts.MaxRetries == 0 {
		opts.MaxRetries = 3
	}
	base := opts.HTTPClient
	if base == nil {
		base = &http.Client{Timeout: 60 * time.Second}
	}
	base.Transport = &transport{
		base:       roundTripper(base.Transport),
		token:      opts.Token,
		userAgent:  opts.UserAgent,
		maxRetries: opts.MaxRetries,
	}

	return NewClientWithResponses(opts.BaseURL, WithHTTPClient(base))
}

func roundTripper(rt http.RoundTripper) http.RoundTripper {
	if rt != nil {
		return rt
	}
	return http.DefaultTransport
}

type transport struct {
	base       http.RoundTripper
	token      string
	userAgent  string
	maxRetries int
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.token != "" {
		req.Header.Set("Authorization", "Bearer "+t.token)
	}
	req.Header.Set("User-Agent", t.userAgent)
	if req.Header.Get("Accept") == "" {
		req.Header.Set("Accept", "application/json")
	}

	var (
		resp *http.Response
		err  error
	)
	for attempt := 0; attempt <= t.maxRetries; attempt++ {
		// Reset request body for retries.
		if attempt > 0 && req.GetBody != nil {
			body, gerr := req.GetBody()
			if gerr != nil {
				return nil, gerr
			}
			req.Body = body
		}

		resp, err = t.base.RoundTrip(req)
		if err != nil {
			if !isRetryable(req.Context(), err) || attempt == t.maxRetries {
				return nil, err
			}
			sleep(req.Context(), backoff(attempt, 0))
			continue
		}

		if !shouldRetryStatus(resp.StatusCode) || attempt == t.maxRetries {
			return resp, nil
		}

		retryAfter := parseRetryAfter(resp.Header.Get("Retry-After"))
		drainAndClose(resp.Body)
		sleep(req.Context(), backoff(attempt, retryAfter))
	}
	return resp, err
}

func shouldRetryStatus(code int) bool {
	switch code {
	case http.StatusTooManyRequests,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout:
		return true
	}
	return false
}

func isRetryable(ctx context.Context, err error) bool {
	if ctx.Err() != nil {
		return false
	}
	return !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded)
}

func parseRetryAfter(h string) time.Duration {
	if h == "" {
		return 0
	}
	if secs, err := strconv.Atoi(h); err == nil && secs >= 0 {
		return time.Duration(secs) * time.Second
	}
	if t, err := http.ParseTime(h); err == nil {
		if d := time.Until(t); d > 0 {
			return d
		}
	}
	return 0
}

func backoff(attempt int, hint time.Duration) time.Duration {
	if hint > 0 {
		return hint
	}
	base := time.Duration(1<<attempt) * 200 * time.Millisecond
	if base > 5*time.Second {
		base = 5 * time.Second
	}
	//nolint:gosec // jitter doesn't need crypto randomness
	jitter := time.Duration(rand.Int63n(int64(base / 3)))
	return base + jitter
}

func sleep(ctx context.Context, d time.Duration) {
	if d <= 0 {
		return
	}
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
	case <-t.C:
	}
}

func drainAndClose(r io.ReadCloser) {
	if r == nil {
		return
	}
	_, _ = io.Copy(io.Discard, r)
	_ = r.Close()
}

// APIError is a structured error extracted from a non-2xx response body.
type APIError struct {
	StatusCode    int    `json:"status_code"`
	Code          string `json:"code,omitempty"`
	Message       string `json:"message"`
	CorrelationID string `json:"correlation_id,omitempty"`
}

func (e *APIError) Error() string {
	if e.CorrelationID != "" {
		return fmt.Sprintf("hostinger api: %s (status %d, correlation_id %s)", e.Message, e.StatusCode, e.CorrelationID)
	}
	return fmt.Sprintf("hostinger api: %s (status %d)", e.Message, e.StatusCode)
}
