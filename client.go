// Package invezgo is the official Go SDK for the Invezgo API (Indonesia stock market data).
//
// Create a client once and reuse it; it is safe for concurrent use:
//
//	client, err := invezgo.NewClient(os.Getenv("INVEZGO_API_KEY"))
//	stocks, err := client.Analysis.GetStockList(ctx)
//
// Every GET response is cached in an in-memory LRU (see CacheConfig) and HTTP 429 responses
// automatically rotate to the keys given with WithAdditionalAPIKeys.
package invezgo

import (
	"errors"
	"net/http"
	"strings"
	"time"
)

// Version is the SDK version, sent in the User-Agent header.
const Version = "2.0.0"

// DefaultBaseURL is the production API endpoint.
const DefaultBaseURL = "https://api.invezgo.com"

// Client is the Invezgo API client. Endpoints are grouped by resource:
// client.Analysis, client.Batch, client.Search, client.Screener, client.Usage,
// client.Watchlists, client.Journals, client.Portfolios, client.Alerts, client.Posts, client.Trades.
type Client struct {
	services

	baseURL    string
	httpClient *http.Client
	timeout    time.Duration
	keys       *keyRing
	cacheCfg   CacheConfig
	cache      *lruCache // nil when caching is disabled
	now        func() time.Time
}

// CacheConfig tunes the in-memory LRU response cache. A zero field means "use the default";
// a negative TTL disables caching for that tier.
type CacheConfig struct {
	// TTL for regular GET endpoints. Default 60s.
	TTL time.Duration
	// RealtimeTTL for [REALTIME] endpoints (order book, running trade, intraday, ...). Default 5s.
	RealtimeTTL time.Duration
	// MaxEntries caps the number of cached responses. Default 1000.
	MaxEntries int
	// MaxBytes caps the total size of cached response bodies. Default 32 MiB.
	MaxBytes int64
}

// DefaultCacheConfig returns the cache settings used when no cache option is given.
func DefaultCacheConfig() CacheConfig {
	return CacheConfig{TTL: 60 * time.Second, RealtimeTTL: 5 * time.Second, MaxEntries: 1000, MaxBytes: 32 << 20}
}

// Option configures a Client.
type Option func(*Client)

// WithAdditionalAPIKeys adds fallback API keys. When a key receives HTTP 429 (per-minute rate limit
// or monthly quota), the request is retried with the next available key.
func WithAdditionalAPIKeys(keys ...string) Option {
	return func(c *Client) { c.keys.add(keys...) }
}

// WithBaseURL overrides the API base URL (default https://api.invezgo.com).
func WithBaseURL(baseURL string) Option {
	return func(c *Client) { c.baseURL = strings.TrimRight(baseURL, "/") }
}

// WithHTTPClient uses a custom *http.Client (proxies, transports, tracing...).
func WithHTTPClient(h *http.Client) Option {
	return func(c *Client) { c.httpClient = h }
}

// WithTimeout sets the timeout of each HTTP attempt (default 30s). Zero or negative disables it;
// the request context still applies.
func WithTimeout(d time.Duration) Option {
	return func(c *Client) { c.timeout = d }
}

// WithCache customises the in-memory response cache (enabled by default).
func WithCache(cfg CacheConfig) Option {
	return func(c *Client) {
		def := DefaultCacheConfig()
		if cfg.TTL == 0 {
			cfg.TTL = def.TTL
		}
		if cfg.RealtimeTTL == 0 {
			cfg.RealtimeTTL = def.RealtimeTTL
		}
		if cfg.MaxEntries <= 0 {
			cfg.MaxEntries = def.MaxEntries
		}
		if cfg.MaxBytes <= 0 {
			cfg.MaxBytes = def.MaxBytes
		}
		c.cacheCfg = cfg
	}
}

// WithoutCache disables the response cache entirely.
func WithoutCache() Option {
	return func(c *Client) { c.cacheCfg = CacheConfig{TTL: -1, RealtimeTTL: -1} }
}

// NewClient creates a client. apiKey is required; get one at https://invezgo.com/id/setting/api.
func NewClient(apiKey string, opts ...Option) (*Client, error) {
	if strings.TrimSpace(apiKey) == "" {
		return nil, errors.New("invezgo: apiKey is required")
	}
	c := &Client{
		baseURL:    DefaultBaseURL,
		httpClient: &http.Client{},
		timeout:    30 * time.Second,
		keys:       newKeyRing(apiKey),
		cacheCfg:   DefaultCacheConfig(),
		now:        time.Now,
	}
	for _, opt := range opts {
		opt(c)
	}
	if c.httpClient == nil {
		c.httpClient = &http.Client{}
	}
	if c.cacheCfg.TTL > 0 || c.cacheCfg.RealtimeTTL > 0 {
		c.cache = newLRUCache(c.cacheCfg.MaxEntries, c.cacheCfg.MaxBytes)
	}
	c.initServices()
	return c, nil
}

// ClearCache drops every cached response.
func (c *Client) ClearCache() {
	if c.cache != nil {
		c.cache.clear()
	}
}
