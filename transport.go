package invezgo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type cacheTier int

const (
	tierNone cacheTier = iota
	tierDefault
	tierRealtime
)

// request describes one API call; built by the generated resource methods.
type request struct {
	method string
	route  string // path template, scope of per-minute rate limits
	path   string
	query  url.Values
	body   any
	tier   cacheTier
}

var quotaRe = regexp.MustCompile(`(?i)quota|kuota|monthly|bulanan`)

// keyRing holds the API keys and their rate-limit blocks.
type keyRing struct {
	mu     sync.Mutex
	keys   []string
	blocks map[string]time.Time // "i" = whole key (monthly quota), "i route" = per-minute limit of one endpoint
}

func newKeyRing(primary string) *keyRing {
	return &keyRing{keys: []string{primary}, blocks: map[string]time.Time{}}
}

func (k *keyRing) add(keys ...string) {
	for _, key := range keys {
		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}
		dup := false
		for _, have := range k.keys {
			dup = dup || have == key
		}
		if !dup {
			k.keys = append(k.keys, key)
		}
	}
}

// blockedUntil must be called with mu held.
func (k *keyRing) blockedUntil(i int, route string) time.Time {
	a, b := k.blocks[strconv.Itoa(i)], k.blocks[strconv.Itoa(i)+" "+route]
	if b.After(a) {
		return b
	}
	return a
}

// pick returns the first untried, unblocked key index, or -1 and the shortest remaining wait.
func (k *keyRing) pick(route string, tried map[int]bool, now time.Time) (int, time.Duration) {
	k.mu.Lock()
	defer k.mu.Unlock()
	var wait time.Duration = -1
	for i := range k.keys {
		until := k.blockedUntil(i, route)
		if !now.Before(until) {
			if !tried[i] {
				return i, 0
			}
			continue
		}
		if d := until.Sub(now); wait < 0 || d < wait {
			wait = d
		}
	}
	if wait < 0 {
		wait = 0
	}
	return -1, wait
}

func (k *keyRing) block(i int, route string, quota bool, until time.Time) {
	k.mu.Lock()
	defer k.mu.Unlock()
	id := strconv.Itoa(i)
	if !quota {
		id += " " + route
	}
	k.blocks[id] = until
}

func (c *Client) ttl(t cacheTier) time.Duration {
	switch t {
	case tierDefault:
		return c.cacheCfg.TTL
	case tierRealtime:
		return c.cacheCfg.RealtimeTTL
	}
	return 0
}

// do executes r and decodes the JSON response into out (may be nil).
func (c *Client) do(ctx context.Context, r request, out any) error {
	key := r.path
	if q := r.query.Encode(); q != "" { // Encode sorts keys: canonical cache key
		key += "?" + q
	}
	ttl := c.ttl(r.tier)
	cacheable := c.cache != nil && r.method == http.MethodGet && ttl > 0
	if cacheable {
		if body, ok := c.cache.get(key, c.now()); ok {
			return decode(body, out)
		}
	}
	var payload []byte
	if r.body != nil {
		var err error
		if payload, err = json.Marshal(r.body); err != nil {
			return fmt.Errorf("invezgo: encode request body: %w", err)
		}
	}

	tried := map[int]bool{}
	var last *RateLimitError
	for {
		idx, wait := c.keys.pick(r.route, tried, c.now())
		if idx < 0 {
			if last == nil {
				last = &RateLimitError{APIError: APIError{StatusCode: 429, Message: "all API keys are rate limited",
					ErrorText: "Too Many Requests"}, Limit: -1, Remaining: -1}
			}
			last.RetryAfter = wait.Round(time.Second)
			if last.RetryAfter < wait {
				last.RetryAfter += time.Second
			}
			return last
		}
		tried[idx] = true
		status, header, body, err := c.send(ctx, r, key, payload, c.keys.keys[idx])
		if err != nil {
			return err
		}
		if status == http.StatusTooManyRequests {
			last = rateLimitError(status, header, body, c.now())
			c.keys.block(idx, r.route, last.QuotaExceeded, c.now().Add(last.RetryAfter))
			continue
		}
		if status < 200 || status > 299 {
			return newAPIError(status, body)
		}
		if r.method != http.MethodGet && c.cache != nil {
			c.cache.purge(func(k string) bool {
				seg := strings.SplitN(strings.TrimPrefix(k, "/"), "/", 2)[0]
				return personalSegments[strings.SplitN(seg, "?", 2)[0]]
			})
		}
		if cacheable {
			c.cache.set(key, body, ttl, c.now())
		}
		return decode(body, out)
	}
}

func (c *Client) send(ctx context.Context, r request, pathAndQuery string, payload []byte, apiKey string) (int, http.Header, []byte, error) {
	if c.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.timeout)
		defer cancel()
	}
	var reader io.Reader
	if payload != nil {
		reader = bytes.NewReader(payload)
	}
	req, err := http.NewRequestWithContext(ctx, r.method, c.baseURL+pathAndQuery, reader)
	if err != nil {
		return 0, nil, nil, fmt.Errorf("invezgo: build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "invezgo-go-sdk/"+Version)
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, nil, nil, fmt.Errorf("invezgo: %s %s: %w", r.method, r.path, err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, nil, fmt.Errorf("invezgo: read response: %w", err)
	}
	return resp.StatusCode, resp.Header, body, nil
}

func rateLimitError(status int, h http.Header, body []byte, now time.Time) *RateLimitError {
	e := &RateLimitError{APIError: *newAPIError(status, body), Limit: headerInt(h, "X-RateLimit-Limit"),
		Remaining: headerInt(h, "X-RateLimit-Remaining")}
	if reset := headerInt(h, "X-RateLimit-Reset"); reset > 0 {
		e.Reset = time.Unix(int64(reset), 0)
	}
	e.QuotaExceeded = quotaRe.MatchString(e.Message)
	switch {
	case headerInt(h, "Retry-After") > 0:
		e.RetryAfter = time.Duration(headerInt(h, "Retry-After")) * time.Second
	case !e.Reset.IsZero() && e.Reset.After(now):
		e.RetryAfter = e.Reset.Sub(now)
	default:
		e.RetryAfter = 60 * time.Second
	}
	return e
}

func headerInt(h http.Header, name string) int {
	v, err := strconv.Atoi(strings.TrimSpace(h.Get(name)))
	if err != nil {
		return -1
	}
	return v
}

func decode(body []byte, out any) error {
	b := bytes.TrimSpace(body)
	if out == nil || len(b) == 0 || string(b) == "null" {
		return nil
	}
	if err := json.Unmarshal(b, out); err != nil {
		return fmt.Errorf("invezgo: decode response: %w", err)
	}
	return nil
}
