package client

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

// testRecord is a minimal Record implementation used only for tests.
type testRecord struct {
	BaseResult
	Name string `json:"name,omitempty"`
}

// newTestClient creates a Client that points at an httptest.Server using the
// supplied handler. The server is closed automatically via t.Cleanup.
func newTestClient(t *testing.T, handler http.HandlerFunc) (*Client, *httptest.Server) {
	t.Helper()
	ts := httptest.NewServer(handler)
	t.Cleanup(ts.Close)
	return &Client{
		BaseURL:    ts.URL,
		Auth:       "Basic dGVzdDp0ZXN0", // base64("test:test")
		HTTPClient: ts.Client(),
		UserAgent:  "test",
	}, ts
}

func TestRequestJSON_BasicSuccess(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"records":[{"sys_id":"abc123","__status":"success"}]}`))
	})

	var record testRecord
	if err := client.GetObject(context.Background(), "test.do", "abc123", &record); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if record.GetID() != "abc123" {
		t.Fatalf("expected ID 'abc123', got %q", record.GetID())
	}
}

func TestRequestJSON_NotFound_EmptyList(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"records":[]}`))
	})

	var record testRecord
	err := client.GetObject(context.Background(), "test.do", "missing", &record)
	if err == nil {
		t.Fatal("expected NotFoundError, got nil")
	}
	if !IsNotFound(err) {
		t.Fatalf("expected IsNotFound(err) to be true, got err=%v", err)
	}
}

func TestRequestJSON_HTTP404(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	var record testRecord
	err := client.GetObject(context.Background(), "test.do", "missing", &record)
	if err == nil {
		t.Fatal("expected NotFoundError, got nil")
	}
	if !IsNotFound(err) {
		t.Fatalf("expected IsNotFound(err) to be true, got err=%v", err)
	}
}

func TestRequestJSON_HTTP500(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal error", http.StatusInternalServerError)
	})

	var record testRecord
	err := client.GetObject(context.Background(), "test.do", "abc", &record)
	if err == nil {
		t.Fatal("expected error for HTTP 500, got nil")
	}
	if IsNotFound(err) {
		t.Fatalf("expected error to NOT be NotFoundError, got %v", err)
	}
	msg := err.Error()
	if !strings.Contains(msg, "500") && !strings.Contains(msg, "Internal Server Error") {
		t.Fatalf("expected error to mention 500 or 'Internal Server Error', got %q", msg)
	}
}

func TestRequestJSON_AuthHeader(t *testing.T) {
	var receivedAuth string
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		receivedAuth = r.Header.Get("Authorization")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"records":[{"sys_id":"abc","__status":"success"}]}`))
	})

	var record testRecord
	if err := client.GetObject(context.Background(), "test.do", "abc", &record); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if receivedAuth != "Basic dGVzdDp0ZXN0" {
		t.Fatalf("expected Authorization header 'Basic dGVzdDp0ZXN0', got %q", receivedAuth)
	}
}

func TestRequestJSON_UserAgent(t *testing.T) {
	var receivedUA string
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		receivedUA = r.Header.Get("User-Agent")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"records":[{"sys_id":"abc","__status":"success"}]}`))
	})

	client.UserAgent = "test-ua/1.0"

	var record testRecord
	if err := client.GetObject(context.Background(), "test.do", "abc", &record); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if receivedUA != "test-ua/1.0" {
		t.Fatalf("expected User-Agent 'test-ua/1.0', got %q", receivedUA)
	}
}

func TestRequestJSON_StatusFieldNotSuccess(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"records":[{"sys_id":"x","__status":"failure","__error":{"reason":"bad","message":"go away"}}]}`))
	})

	var record testRecord
	err := client.GetObject(context.Background(), "test.do", "x", &record)
	if err == nil {
		t.Fatal("expected error from non-success status, got nil")
	}
	if !strings.Contains(err.Error(), "go away") {
		t.Fatalf("expected error message to contain 'go away', got %q", err.Error())
	}
}

func TestRequestJSON_ContextCancellation(t *testing.T) {
	// Server is harmless; the cancelled context should short-circuit the call.
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"records":[{"sys_id":"abc","__status":"success"}]}`))
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	var record testRecord
	err := client.GetObject(ctx, "test.do", "abc", &record)
	if err == nil {
		t.Fatal("expected context cancellation error, got nil")
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected errors.Is(err, context.Canceled) to be true, got %v", err)
	}
}

func TestRequestJSON_CreatePOSTBody(t *testing.T) {
	var receivedMethod string
	var receivedBody atomic.Value

	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("server failed to read request body: %v", err)
		}
		receivedBody.Store(string(bodyBytes))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"records":[{"sys_id":"new","__status":"success"}]}`))
	})

	record := &testRecord{Name: "foo"}
	if err := client.CreateObject(context.Background(), "test.do", record); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if receivedMethod != http.MethodPost {
		t.Fatalf("expected POST request, got %s", receivedMethod)
	}
	body, _ := receivedBody.Load().(string)
	if !strings.Contains(body, `"name":"foo"`) {
		t.Fatalf("expected request body to contain '\"name\":\"foo\"', got %q", body)
	}
}

// shortBackoffClient returns a client wired to the supplied test server, used by
// retry tests. It uses the same test client setup as newTestClient.
func TestRequestJSON_RetryOn503(t *testing.T) {
	var calls atomic.Int32
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		n := calls.Add(1)
		if n < 3 {
			http.Error(w, "unavailable", http.StatusServiceUnavailable)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"records":[{"sys_id":"abc","__status":"success"}]}`))
	})

	var record testRecord
	if err := client.GetObject(context.Background(), "test.do", "abc", &record); err != nil {
		t.Fatalf("expected success after retries, got %v", err)
	}
	if got := calls.Load(); got != 3 {
		t.Fatalf("expected 3 server calls, got %d", got)
	}
	if record.GetID() != "abc" {
		t.Fatalf("expected ID 'abc', got %q", record.GetID())
	}
}

func TestRequestJSON_RetryOn429_RetryAfter(t *testing.T) {
	var calls atomic.Int32
	var firstCallAt, secondCallAt atomic.Int64
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		n := calls.Add(1)
		now := time.Now().UnixNano()
		switch n {
		case 1:
			firstCallAt.Store(now)
			w.Header().Set("Retry-After", "1")
			http.Error(w, "rate limited", http.StatusTooManyRequests)
		default:
			secondCallAt.Store(now)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"records":[{"sys_id":"abc","__status":"success"}]}`))
		}
	})

	var record testRecord
	if err := client.GetObject(context.Background(), "test.do", "abc", &record); err != nil {
		t.Fatalf("expected success after retry, got %v", err)
	}
	if got := calls.Load(); got != 2 {
		t.Fatalf("expected 2 server calls, got %d", got)
	}
	elapsed := time.Duration(secondCallAt.Load() - firstCallAt.Load())
	if elapsed < time.Second {
		t.Fatalf("expected at least 1s between attempts (Retry-After: 1), got %v", elapsed)
	}
}

func TestRequestJSON_NoRetryOn400(t *testing.T) {
	var calls atomic.Int32
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		http.Error(w, "bad request", http.StatusBadRequest)
	})

	var record testRecord
	err := client.GetObject(context.Background(), "test.do", "abc", &record)
	if err == nil {
		t.Fatal("expected error for HTTP 400, got nil")
	}
	if IsNotFound(err) {
		t.Fatalf("expected non-NotFound error, got %v", err)
	}
	if got := calls.Load(); got != 1 {
		t.Fatalf("expected exactly 1 server call (no retries on 4xx), got %d", got)
	}
}

func TestRequestJSON_MaxRetriesExhausted(t *testing.T) {
	var calls atomic.Int32
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		http.Error(w, "unavailable", http.StatusServiceUnavailable)
	})

	var record testRecord
	err := client.GetObject(context.Background(), "test.do", "abc", &record)
	if err == nil {
		t.Fatal("expected error after retries exhausted, got nil")
	}
	if got := calls.Load(); got != 4 {
		t.Fatalf("expected 4 total attempts, got %d", got)
	}
	if !strings.Contains(err.Error(), "503") && !strings.Contains(err.Error(), "Service Unavailable") {
		t.Fatalf("expected last error to mention 503/Service Unavailable, got %q", err.Error())
	}
}

func TestRequestJSON_ContextCancelledDuringBackoff(t *testing.T) {
	var calls atomic.Int32
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		http.Error(w, "unavailable", http.StatusServiceUnavailable)
	})

	ctx, cancel := context.WithCancel(context.Background())
	// Cancel after 100ms — this falls during the first backoff (1s) and aborts
	// the retry loop before the next HTTP call is made.
	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	var record testRecord
	start := time.Now()
	err := client.GetObject(ctx, "test.do", "abc", &record)
	elapsed := time.Since(start)
	if err == nil {
		t.Fatal("expected context cancellation error, got nil")
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected errors.Is(err, context.Canceled), got %v", err)
	}
	if got := calls.Load(); got >= 4 {
		t.Fatalf("expected fewer than 4 attempts due to cancellation, got %d", got)
	}
	// Sanity check: the backoff loop should exit promptly after cancellation,
	// not run to the full 1s+2s+4s schedule.
	if elapsed > 2*time.Second {
		t.Fatalf("expected early exit on cancellation, took %v", elapsed)
	}
}
