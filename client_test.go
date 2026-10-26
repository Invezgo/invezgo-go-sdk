package invezgo

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

// mock is a scripted API server that records every request it receives.
type mock struct {
	mu      sync.Mutex
	calls   []string // "KEY METHOD /path?query"
	handler func(w http.ResponseWriter, r *http.Request, key string)
}

func newMock(t *testing.T, h func(w http.ResponseWriter, r *http.Request, key string), opts ...Option) (*Client, *mock, *time.Time) {
	t.Helper()
	m := &mock{handler: h}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		m.mu.Lock()
		m.calls = append(m.calls, key+" "+r.Method+" "+r.URL.RequestURI())
		m.mu.Unlock()
		m.handler(w, r, key)
	}))
	t.Cleanup(srv.Close)
	c, err := NewClient("k1", append([]Option{WithBaseURL(srv.URL + "/")}, opts...)...)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Unix(1_700_000_000, 0)
	c.now = func() time.Time { return now }
	return c, m, &now
}

func (m *mock) count() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.calls)
}

func (m *mock) last() string {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.calls[len(m.calls)-1]
}

func jsonOK(body string) func(http.ResponseWriter, *http.Request, string) {
	return func(w http.ResponseWriter, _ *http.Request, _ string) { io.WriteString(w, body) }
}

var ctx = context.Background()

func TestNewClientRequiresKey(t *testing.T) {
	if _, err := NewClient("  "); err == nil {
		t.Fatal("expected error for empty api key")
	}
}

func TestLRUEvictsByEntriesBytesAndTTL(t *testing.T) {
	now := time.Unix(0, 0)
	c := newLRUCache(2, 10)
	c.set("a", []byte("1234"), time.Minute, now)
	c.set("b", []byte("1234"), time.Minute, now)
	c.get("a", now) // a becomes MRU
	c.set("c", []byte("12"), time.Minute, now)
	if _, ok := c.get("b", now); ok {
		t.Fatal("b should be evicted (entry cap, LRU)")
	}
	c.set("d", []byte("123456789"), time.Minute, now) // 4+2+9 > 10 bytes
	if n, b := c.len(); n != 1 || b != 9 {
		t.Fatalf("byte cap: got %d entries %d bytes", n, b)
	}
	c.set("huge", []byte("12345678901"), time.Minute, now)
	if _, ok := c.get("huge", now); ok {
		t.Fatal("bodies larger than maxBytes must not be cached")
	}
	if _, ok := c.get("d", now.Add(time.Minute)); ok {
		t.Fatal("entry should expire at TTL")
	}
	if n, b := c.len(); n != 0 || b != 0 {
		t.Fatalf("expired entry not removed: %d %d", n, b)
	}
}

func TestCacheTiersAndClear(t *testing.T) {
	c, m, now := newMock(t, func(w http.ResponseWriter, r *http.Request, _ string) {
		if r.URL.Path == "/analysis/list/stock" {
			io.WriteString(w, `[]`)
			return
		}
		io.WriteString(w, `{"code":"COMPOSITE","open":"7000.5"}`)
	})
	for i := 0; i < 3; i++ {
		if _, err := c.Analysis.GetStockList(ctx); err != nil {
			t.Fatal(err)
		}
		if _, err := c.Analysis.GetIntradayIndex(ctx, "COMPOSITE"); err != nil {
			t.Fatal(err)
		}
	}
	if m.count() != 2 {
		t.Fatalf("expected 2 network calls, got %d", m.count())
	}
	*now = now.Add(6 * time.Second) // realtime TTL (5s) expired, default (60s) not
	c.Analysis.GetStockList(ctx)
	c.Analysis.GetIntradayIndex(ctx, "COMPOSITE")
	if m.count() != 3 {
		t.Fatalf("only the realtime entry should refetch, calls=%d", m.count())
	}
	c.ClearCache()
	c.Analysis.GetStockList(ctx)
	if m.count() != 4 {
		t.Fatalf("ClearCache should force a refetch, calls=%d", m.count())
	}
}

func TestWithoutCacheAndUncachedEndpoints(t *testing.T) {
	c, m, _ := newMock(t, jsonOK(`{"usage":1,"remaining":2,"limit":3,"isBlocked":false,"expire":"x","role":"ADVANCE"}`), WithoutCache())
	c.Analysis.GetStockList(ctx)
	c.Analysis.GetStockList(ctx)
	if m.count() != 2 {
		t.Fatalf("WithoutCache must not cache, calls=%d", m.count())
	}
	c2, m2, _ := newMock(t, jsonOK(`{"usage":1}`))
	c2.Usage.GetAPIUsage(ctx)
	u, _ := c2.Usage.GetAPIUsage(ctx)
	if m2.count() != 2 || u.Usage != 1 {
		t.Fatalf("usage endpoints are never cached, calls=%d", m2.count())
	}
}

func TestMutationPurgesPersonalCacheOnly(t *testing.T) {
	c, m, _ := newMock(t, func(w http.ResponseWriter, r *http.Request, _ string) {
		if r.Method == http.MethodPost {
			if ct := r.Header.Get("Content-Type"); ct != "application/json" {
				t.Errorf("content-type = %q", ct)
			}
			b, _ := io.ReadAll(r.Body)
			if !strings.Contains(string(b), `"code":"BBCA"`) || strings.Contains(string(b), "note") {
				t.Errorf("unexpected body %s", b)
			}
			io.WriteString(w, `{"message":"ok"}`)
			return
		}
		io.WriteString(w, `[]`)
	})
	c.Watchlists.List(ctx, nil)
	c.Portfolios.List(ctx)
	c.Analysis.GetStockList(ctx)
	if _, err := c.Watchlists.Add(ctx, &AddWatchlistRequest{Group: "g", Code: "BBCA", Price: 9000, Scope: []string{}}); err != nil {
		t.Fatal(err)
	}
	before := m.count()
	c.Watchlists.List(ctx, nil)
	c.Portfolios.List(ctx)
	c.Analysis.GetStockList(ctx)
	if got := m.count() - before; got != 2 {
		t.Fatalf("expected watchlists+portfolios refetch only, got %d calls", got)
	}
}

func TestQueryCanonicalAndPathEscaped(t *testing.T) {
	c, m, _ := newMock(t, jsonOK(`{"data":[]}`))
	_, err := c.Analysis.GetStockRunningTrade(ctx, "BB CA", StockRunningTradeParams{
		Date: "2024-12-30", Type: Ptr(StockRunningTradeTypeBUY), Minimum: Ptr(1.0), Maximum: Ptr(2.5), Page: Ptr(int64(3)),
	})
	if err != nil {
		t.Fatal(err)
	}
	want := "k1 GET /analysis/running-trade/BB%20CA?date=2024-12-30&maximum=2.5&minimum=1&page=3&type=BUY"
	if m.last() != want {
		t.Fatalf("got  %s\nwant %s", m.last(), want)
	}
}

func TestNullBodyFlexAndAdditionalProperties(t *testing.T) {
	c, _, _ := newMock(t, func(w http.ResponseWriter, r *http.Request, _ string) {
		switch {
		case strings.HasPrefix(r.URL.Path, "/analysis/shareholder/classify-table/"):
			io.WriteString(w, "null")
		case r.URL.Path == "/analysis/intraday-index/COMPOSITE":
			io.WriteString(w, `{"code":"COMPOSITE","open":"7000.5","high":7100,"low":null,"close":""}`)
		default:
			io.WriteString(w, `[{"code":"BBCA","matched":true,"close":9875,"rsi(14)":null,"bad":"x"}]`)
		}
	})
	ct, err := c.Analysis.GetShareholderClassifyTable(ctx, "XXXX")
	if err != nil || ct != nil {
		t.Fatalf("null body should give nil, nil; got %v %v", ct, err)
	}
	ix, err := c.Analysis.GetIntradayIndex(ctx, "COMPOSITE")
	if err != nil {
		t.Fatal(err)
	}
	if *ix.Open != 7000.5 || *ix.High != 7100 || ix.Low != nil || *ix.Close != 0 {
		t.Fatalf("flex decode: %+v", ix)
	}
	rows, err := c.Screener.Screen(ctx, &ScreenRequest{Formula: "close > prev"})
	if err != nil {
		t.Fatal(err)
	}
	r := rows[0]
	if r.Code != "BBCA" || !r.Matched || *r.Values["close"] != 9875 || r.Values["rsi(14)"] != nil || len(r.Values) != 2 {
		t.Fatalf("additional properties: %+v", r)
	}
	if _, err := c.Screener.Screen(ctx, nil); err == nil {
		t.Fatal("nil required body must error")
	}
}

func TestErrorMapping(t *testing.T) {
	cases := []struct {
		status int
		body   string
		is     func(error) bool
		msg    string
	}{
		{401, `{"statusCode":401,"message":"Authentication required","error":"Unauthorized","requestId":"r1"}`, IsAuthentication, "Authentication required"},
		{402, `{"statusCode":402,"message":"Starter role user only","error":"Payment Required"}`, IsPaymentRequired, "Starter role user only"},
		{404, `{"statusCode":404,"message":"Not found","error":"Not Found"}`, IsNotFound, "Not found"},
		{422, `{"statusCode":422,"message":["a is bad","b is bad"],"error":"Unprocessable Entity"}`, IsValidation, "a is bad; b is bad"},
		{400, `{"message":"Unknown formula variable","error":"Unknown formula variable"}`, IsValidation, "Unknown formula variable"},
		{503, `gateway down`, IsServerError, "gateway down"},
	}
	for _, tc := range cases {
		c, _, _ := newMock(t, func(w http.ResponseWriter, _ *http.Request, _ string) {
			w.WriteHeader(tc.status)
			io.WriteString(w, tc.body)
		})
		_, err := c.Analysis.GetStockList(ctx)
		var apiErr *APIError
		if !errors.As(err, &apiErr) || !tc.is(err) || apiErr.StatusCode != tc.status || apiErr.Message != tc.msg {
			t.Fatalf("%d: got %#v", tc.status, err)
		}
	}
}

func TestRotationPerRouteMinuteLimit(t *testing.T) {
	c, m, _ := newMock(t, func(w http.ResponseWriter, r *http.Request, key string) {
		if key == "k1" && r.URL.Path == "/analysis/list/stock" {
			w.Header().Set("Retry-After", "30")
			w.Header().Set("X-RateLimit-Limit", "250")
			w.WriteHeader(429)
			io.WriteString(w, `{"statusCode":429,"message":"Too many requests: 250 requests per minute on this endpoint.","error":"Too Many Requests"}`)
			return
		}
		io.WriteString(w, `[]`)
	}, WithAdditionalAPIKeys("k2", "", "k1"), WithoutCache())

	if _, err := c.Analysis.GetStockList(ctx); err != nil {
		t.Fatal(err)
	}
	if m.count() != 2 || !strings.HasPrefix(m.last(), "k2 ") {
		t.Fatalf("expected retry with k2, calls=%v", m.calls)
	}
	c.Analysis.GetBrokerList(ctx) // other endpoint: k1 still usable
	if !strings.HasPrefix(m.last(), "k1 ") {
		t.Fatalf("minute limit must be per route, got %s", m.last())
	}
	c.Analysis.GetStockList(ctx) // blocked route goes straight to k2
	if !strings.HasPrefix(m.last(), "k2 ") || m.count() != 4 {
		t.Fatalf("blocked route should skip k1, calls=%v", m.calls)
	}
}

func TestRotationQuotaBlocksWholeKeyAndAllBlocked(t *testing.T) {
	brokerHits := 0
	c, m, now := newMock(t, func(w http.ResponseWriter, r *http.Request, key string) {
		if key == "k1" {
			w.Header().Set("Retry-After", "3600")
			w.WriteHeader(429)
			io.WriteString(w, `{"statusCode":429,"message":"Monthly quota exceeded, resets 2026-11-01","error":"Too Many Requests"}`)
			return
		}
		if r.URL.Path == "/analysis/list/broker" && brokerHits == 0 { // per-minute limit on k2, no Retry-After → X-RateLimit-Reset
			brokerHits++
			w.Header().Set("X-RateLimit-Reset", fmt.Sprint(time.Unix(1_700_000_000, 0).Add(20*time.Second).Unix()))
			w.WriteHeader(429)
			io.WriteString(w, `{"statusCode":429,"message":"Too many requests per minute","error":"Too Many Requests"}`)
			return
		}
		io.WriteString(w, `[]`)
	}, WithAdditionalAPIKeys("k2"), WithoutCache())

	c.Analysis.GetStockList(ctx) // k1 quota → k2
	c.Analysis.GetIndexList(ctx) // k1 skipped entirely now
	if m.count() != 3 || !strings.HasPrefix(m.last(), "k2 ") {
		t.Fatalf("quota should block k1 for every route, calls=%v", m.calls)
	}
	_, err := c.Analysis.GetBrokerList(ctx) // k2 minute-limited, k1 quota-blocked
	var rl *RateLimitError
	if !errors.As(err, &rl) || rl.QuotaExceeded || rl.RetryAfter != 20*time.Second || !IsRateLimited(err) {
		t.Fatalf("want minute RateLimitError with 20s, got %#v", err)
	}
	calls := m.count()
	_, err = c.Analysis.GetBrokerList(ctx)
	if !errors.As(err, &rl) || m.count() != calls || rl.RetryAfter != 20*time.Second {
		t.Fatalf("all keys blocked must fail without a request: %v calls %d→%d", err, calls, m.count())
	}
	*now = now.Add(21 * time.Second)
	if _, err := c.Analysis.GetBrokerList(ctx); err != nil || !strings.HasPrefix(m.last(), "k2 ") {
		t.Fatalf("k2 should be usable after its block expires: %v", err)
	}
}

func TestSingleKeyRateLimitDefaults60s(t *testing.T) {
	c, m, _ := newMock(t, func(w http.ResponseWriter, _ *http.Request, _ string) {
		w.WriteHeader(429)
		io.WriteString(w, `{"statusCode":429,"message":"slow down","error":"Too Many Requests"}`)
	})
	_, err := c.Analysis.GetStockList(ctx)
	var rl *RateLimitError
	if !errors.As(err, &rl) || rl.RetryAfter != 60*time.Second || rl.Message != "slow down" || m.count() != 1 {
		t.Fatalf("got %#v calls=%d", err, m.count())
	}
}

func TestConcurrentUse(t *testing.T) {
	c, _, _ := newMock(t, jsonOK(`[{"code":"BBCA","name":"Bank Central Asia"}]`), WithAdditionalAPIKeys("k2"))
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if i%10 == 0 {
				c.ClearCache()
			}
			s, err := c.Analysis.GetStockList(ctx)
			if err != nil || len(s) != 1 || s[0].Code != "BBCA" {
				t.Errorf("got %v %v", s, err)
			}
		}(i)
	}
	wg.Wait()
}
