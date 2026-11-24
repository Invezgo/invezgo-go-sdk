package invezgo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/coder/websocket"
)

// DefaultRealtimeURL is the production WebSocket endpoint (Elite plan).
const DefaultRealtimeURL = "wss://api.invezgo.com/ws"

// Realtime channels.
const (
	ChannelTrades    = "trades"
	ChannelOrderbook = "orderbook"
	ChannelStock     = "stock"
	ChannelIndex     = "index"
)

const (
	maxTradeCodes     = 200
	maxOrderbookCodes = 20
	maxStockCodes     = 50
	maxIndexCodes     = 100
)

// RealtimeTrade is one running-trade transaction.
type RealtimeTrade struct {
	// ID is the exchange trade ID (a string: it exceeds the JavaScript safe integer).
	ID string `json:"id"`
	// Time is HHMMSS (WIB).
	Time   string  `json:"time"`
	Code   string  `json:"code"`
	Price  float64 `json:"price"`
	Lot    int64   `json:"lot"`
	Value  int64   `json:"value"`
	Change float64 `json:"change"`
	// Side is "buy" (HAKA) or "sell" (HAKI).
	Side string `json:"side"`
}

// RealtimeTrades is a batch of trades.
type RealtimeTrades struct {
	// TS is the server send time in unix milliseconds.
	TS   int64           `json:"ts"`
	Data []RealtimeTrade `json:"data"`
}

// RealtimeOrderbook is the full set of levels of one side of a stock's order book.
type RealtimeOrderbook struct {
	Code string `json:"code"`
	// Side is "bid" or "offer".
	Side string `json:"side"`
	// Levels are kept as raw JSON (the element shape is not documented).
	Levels []json.RawMessage `json:"levels"`
	Seq    *int64            `json:"seq"`
}

// RealtimeStock is a full stock quote snapshot. Pointer fields are nil when the exchange has not
// provided them yet.
type RealtimeStock struct {
	Code      string   `json:"code"`
	Last      float64  `json:"last"`
	Prev      *float64 `json:"prev"`
	Open      *float64 `json:"open"`
	High      *float64 `json:"high"`
	Low       *float64 `json:"low"`
	Change    *float64 `json:"change"`
	ChangePct *float64 `json:"change_pct"`
	// Volume is in shares (not lots).
	Volume *int64 `json:"volume"`
	// Value is in Rupiah.
	Value *int64   `json:"value"`
	Freq  *int64   `json:"freq"`
	Avg   *float64 `json:"avg"`
	// IEP / IEV are the pre-opening Indicative Equilibrium Price and Volume (shares).
	IEP *float64 `json:"iep"`
	IEV *int64   `json:"iev"`
	// Time is HHMMSS (WIB).
	Time string `json:"time"`
	Seq  *int64 `json:"seq"`
}

// RealtimeIndex is an index snapshot. Code keeps the exchange's casing (e.g. "Investor33"), so
// compare it case-insensitively (strings.EqualFold).
type RealtimeIndex struct {
	Code      string   `json:"code"`
	Last      float64  `json:"last"`
	Prev      *float64 `json:"prev"`
	Open      *float64 `json:"open"`
	High      *float64 `json:"high"`
	Low       *float64 `json:"low"`
	Change    *float64 `json:"change"`
	ChangePct *float64 `json:"change_pct"`
	Volume    *int64   `json:"volume"`
	Value     *int64   `json:"value"`
	Freq      *int64   `json:"freq"`
	// Time is HHMMSS (WIB).
	Time string `json:"time"`
}

// RealtimeSubscribed is the server's current subscription state for one channel.
type RealtimeSubscribed struct {
	Channel string
	// Codes are the subscribed codes; empty when All is true.
	Codes []string
	// All is true when the server reports codes "*" (every stock or index).
	All bool
	// MinLot and MinValue are only present on the trades channel.
	MinLot, MinValue *int64
}

// RealtimeHandlers receives messages from Run. Every field is optional. Handlers run on the Run
// goroutine; they may call the Subscribe/Unsubscribe/Ping methods.
type RealtimeHandlers struct {
	// OnConnected is called after each successful connection, once subscriptions are replayed.
	OnConnected  func()
	OnSubscribed func(RealtimeSubscribed)
	OnTrades     func(RealtimeTrades)
	OnOrderbook  func(RealtimeOrderbook)
	OnStock      func(RealtimeStock)
	OnIndex      func(RealtimeIndex)
	OnPong       func()
	// OnLagged reports messages skipped because the client read too slowly.
	OnLagged func(missed int64)
	// OnError receives a server error message; the server then closes with code 4400.
	OnError func(message string)
	// OnReconnect is called before waiting delay and reconnecting after err. Backfill missed data
	// via REST (client.Analysis running trade / order book) once OnConnected fires again.
	OnReconnect func(attempt int, delay time.Duration, err error)
	// OnUnknown receives messages with an unknown type or an undecodable body.
	OnUnknown func(raw []byte)
}

// TradeFilter sets the trades size filters. A nil field leaves the current filter unchanged; 0 clears it.
type TradeFilter struct {
	MinLot   *int64
	MinValue *int64
}

// RealtimeOption configures a Realtime client.
type RealtimeOption func(*Realtime)

// WithReconnect enables or disables automatic reconnection (default enabled).
func WithReconnect(enabled bool) RealtimeOption {
	return func(r *Realtime) { r.reconnect = enabled }
}

// Realtime is a WebSocket client for running trades, order books, stock quotes and indices. It remembers the desired
// subscriptions and replays them after every (re)connection. Call Run once to start it.
type Realtime struct {
	apiKey     string
	url        string
	httpClient *http.Client
	reconnect  bool
	unit       time.Duration // backoff unit, 1s (shortened in tests)

	mu         sync.Mutex
	conn       *websocket.Conn // set once "connected" is received
	cancel     context.CancelFunc
	closed     bool
	tradesAll  bool
	tradeCodes []string
	minLot     *int64
	minValue   *int64
	bookCodes  []string
	stockCodes []string
	indexAll   bool
	indexCodes []string
}

// NewRealtime creates a standalone WebSocket client. url defaults to DefaultRealtimeURL.
func NewRealtime(apiKey, url string, opts ...RealtimeOption) (*Realtime, error) {
	if strings.TrimSpace(apiKey) == "" {
		return nil, errors.New("invezgo: apiKey is required")
	}
	if url == "" {
		url = DefaultRealtimeURL
	}
	r := &Realtime{apiKey: apiKey, url: url, httpClient: &http.Client{}, reconnect: true, unit: time.Second}
	for _, opt := range opts {
		opt(r)
	}
	return r, nil
}

// Realtime creates a WebSocket client using the client's primary API key, HTTP client and base URL
// (https → wss, http → ws, plus /ws).
func (c *Client) Realtime(opts ...RealtimeOption) *Realtime {
	url := c.baseURL + "/ws"
	if strings.HasPrefix(url, "https://") {
		url = "wss://" + strings.TrimPrefix(url, "https://")
	} else if strings.HasPrefix(url, "http://") {
		url = "ws://" + strings.TrimPrefix(url, "http://")
	}
	r, _ := NewRealtime(c.keys.keys[0], url, opts...)
	r.httpClient = c.httpClient
	return r
}

type wsCommand struct {
	Op       string   `json:"op"`
	Channel  string   `json:"channel,omitempty"`
	Codes    []string `json:"codes,omitempty"`
	MinLot   *int64   `json:"min_lot,omitempty"`
	MinValue *int64   `json:"min_value,omitempty"`
}

// SubscribeTrades subscribes to running trades of codes, or of every stock when codes is empty.
// filter may be nil. The subscription is sent now if connected, else after the next connection.
func (r *Realtime) SubscribeTrades(ctx context.Context, codes []string, filter *TradeFilter) error {
	codes = upper(codes)
	r.mu.Lock()
	defer r.mu.Unlock()
	merged := mergeCodes(r.tradeCodes, codes)
	if len(merged) > maxTradeCodes {
		return fmt.Errorf("invezgo: at most %d trade codes per connection", maxTradeCodes)
	}
	cmd := wsCommand{Op: "subscribe", Channel: ChannelTrades, Codes: codes}
	if len(codes) == 0 {
		r.tradesAll = true
	}
	r.tradeCodes = merged
	if filter != nil {
		cmd.MinLot, cmd.MinValue = filter.MinLot, filter.MinValue
		if filter.MinLot != nil {
			r.minLot = filter.MinLot
		}
		if filter.MinValue != nil {
			r.minValue = filter.MinValue
		}
	}
	return r.sendLocked(ctx, cmd)
}

// SubscribeOrderbook subscribes to the order book of codes (1 to 20 codes per connection).
func (r *Realtime) SubscribeOrderbook(ctx context.Context, codes []string) error {
	if len(codes) == 0 {
		return errors.New("invezgo: orderbook codes are required")
	}
	codes = upper(codes)
	r.mu.Lock()
	defer r.mu.Unlock()
	merged := mergeCodes(r.bookCodes, codes)
	if len(merged) > maxOrderbookCodes {
		return fmt.Errorf("invezgo: at most %d orderbook codes per connection", maxOrderbookCodes)
	}
	r.bookCodes = merged
	return r.sendLocked(ctx, wsCommand{Op: "subscribe", Channel: ChannelOrderbook, Codes: codes})
}

// SubscribeStock subscribes to the stock quote (OHLC, volume, value, freq, IEP/IEV) of codes
// (1 to 50 codes per connection).
func (r *Realtime) SubscribeStock(ctx context.Context, codes []string) error {
	if len(codes) == 0 {
		return errors.New("invezgo: stock codes are required")
	}
	codes = upper(codes)
	r.mu.Lock()
	defer r.mu.Unlock()
	merged := mergeCodes(r.stockCodes, codes)
	if len(merged) > maxStockCodes {
		return fmt.Errorf("invezgo: at most %d stock quote codes per connection", maxStockCodes)
	}
	r.stockCodes = merged
	return r.sendLocked(ctx, wsCommand{Op: "subscribe", Channel: ChannelStock, Codes: codes})
}

// SubscribeIndex subscribes to index quotes of codes (e.g. "COMPOSITE", "LQ45"), or of every index
// when codes is empty.
func (r *Realtime) SubscribeIndex(ctx context.Context, codes []string) error {
	codes = upper(codes)
	r.mu.Lock()
	defer r.mu.Unlock()
	merged := mergeCodes(r.indexCodes, codes)
	if len(merged) > maxIndexCodes {
		return fmt.Errorf("invezgo: at most %d index codes per connection", maxIndexCodes)
	}
	if len(codes) == 0 {
		r.indexAll = true
	}
	r.indexCodes = merged
	return r.sendLocked(ctx, wsCommand{Op: "subscribe", Channel: ChannelIndex, Codes: codes})
}

// Unsubscribe removes codes from channel (ChannelTrades, ChannelOrderbook, ChannelStock or
// ChannelIndex). Without codes the whole channel is unsubscribed (for trades this also clears the
// size filters).
func (r *Realtime) Unsubscribe(ctx context.Context, channel string, codes []string) error {
	codes = upper(codes)
	r.mu.Lock()
	defer r.mu.Unlock()
	all := len(codes) == 0
	switch channel {
	case ChannelTrades:
		if all {
			r.tradesAll, r.tradeCodes, r.minLot, r.minValue = false, nil, nil, nil
		} else {
			r.tradeCodes = removeCodes(r.tradeCodes, codes)
		}
	case ChannelIndex:
		if all {
			r.indexAll, r.indexCodes = false, nil
		} else {
			r.indexCodes = removeCodes(r.indexCodes, codes)
		}
	case ChannelOrderbook:
		r.bookCodes = removeCodes(r.bookCodes, codes)
		if all {
			r.bookCodes = nil
		}
	case ChannelStock:
		r.stockCodes = removeCodes(r.stockCodes, codes)
		if all {
			r.stockCodes = nil
		}
	default:
		return fmt.Errorf("invezgo: unknown channel %q", channel)
	}
	return r.sendLocked(ctx, wsCommand{Op: "unsubscribe", Channel: channel, Codes: codes})
}

// Ping sends {"op":"ping"}; the server answers with a pong (OnPong). It fails when not connected.
func (r *Realtime) Ping(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.conn == nil {
		return errors.New("invezgo: realtime not connected")
	}
	return r.sendLocked(ctx, wsCommand{Op: "ping"})
}

// Close closes the connection and stops Run (which then returns nil). It cannot be restarted.
func (r *Realtime) Close() error {
	r.mu.Lock()
	r.closed = true
	conn, cancel := r.conn, r.cancel
	r.mu.Unlock()
	var err error
	if conn != nil {
		err = conn.Close(websocket.StatusNormalClosure, "")
	}
	if cancel != nil {
		cancel()
	}
	return err
}

// sendLocked writes cmd when connected; otherwise the stored state is replayed on connect.
// r.mu must be held.
func (r *Realtime) sendLocked(ctx context.Context, cmd wsCommand) error {
	if r.conn == nil {
		return nil
	}
	b, _ := json.Marshal(cmd)
	return r.conn.Write(ctx, websocket.MessageText, b)
}

// replayLocked sends the desired subscription state. r.mu must be held.
func (r *Realtime) replayLocked(ctx context.Context) error {
	if r.tradesAll || len(r.tradeCodes) > 0 {
		cmd := wsCommand{Op: "subscribe", Channel: ChannelTrades, MinLot: r.minLot, MinValue: r.minValue}
		if !r.tradesAll {
			cmd.Codes = r.tradeCodes
		}
		if err := r.sendLocked(ctx, cmd); err != nil {
			return err
		}
	}
	if len(r.bookCodes) > 0 {
		if err := r.sendLocked(ctx, wsCommand{Op: "subscribe", Channel: ChannelOrderbook, Codes: r.bookCodes}); err != nil {
			return err
		}
	}
	if len(r.stockCodes) > 0 {
		if err := r.sendLocked(ctx, wsCommand{Op: "subscribe", Channel: ChannelStock, Codes: r.stockCodes}); err != nil {
			return err
		}
	}
	if r.indexAll || len(r.indexCodes) > 0 {
		cmd := wsCommand{Op: "subscribe", Channel: ChannelIndex}
		if !r.indexAll {
			cmd.Codes = r.indexCodes
		}
		return r.sendLocked(ctx, cmd)
	}
	return nil
}

// Run connects and delivers messages to h until ctx is done (returns ctx.Err()), Close is called
// (returns nil) or a fatal error occurs: handshake 401/403 (*APIError), server close code 4400, or
// any disconnect when reconnection is disabled. Network drops and handshake 429/5xx reconnect with
// backoff 1, 2, 4, 8, 16, 30s. Call Run once per Realtime.
func (r *Realtime) Run(ctx context.Context, h RealtimeHandlers) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	r.mu.Lock()
	if r.closed {
		r.mu.Unlock()
		return nil
	}
	r.cancel = cancel
	r.mu.Unlock()

	attempt := 0
	for {
		connected, err := r.session(ctx, h)
		r.mu.Lock()
		closed := r.closed
		r.mu.Unlock()
		switch {
		case closed:
			return nil
		case ctx.Err() != nil:
			return ctx.Err()
		case websocket.CloseStatus(err) == 4400:
			return fmt.Errorf("invezgo: realtime closed by server (invalid client message): %w", err)
		case IsAuthentication(err) || statusIs(err, func(s int) bool { return s == http.StatusForbidden }):
			return err
		case !r.reconnect:
			return err
		}
		if connected {
			attempt = 0
		}
		attempt++
		delay := r.unit << (attempt - 1)
		if attempt > 5 {
			delay = 30 * r.unit
		}
		if h.OnReconnect != nil {
			h.OnReconnect(attempt, delay, err)
		}
		select {
		case <-ctx.Done():
			r.mu.Lock()
			closed := r.closed
			r.mu.Unlock()
			if closed {
				return nil
			}
			return ctx.Err()
		case <-time.After(delay):
		}
	}
}

// session runs one connection; connected reports whether "connected" was received.
func (r *Realtime) session(ctx context.Context, h RealtimeHandlers) (connected bool, err error) {
	conn, resp, err := websocket.Dial(ctx, r.url, &websocket.DialOptions{
		HTTPClient: r.httpClient,
		HTTPHeader: http.Header{"Authorization": {"Bearer " + r.apiKey}, "User-Agent": {"invezgo-go-sdk/" + Version}},
	})
	if err != nil {
		if resp != nil && resp.StatusCode != http.StatusSwitchingProtocols {
			var body []byte
			if resp.Body != nil {
				body, _ = io.ReadAll(resp.Body)
			}
			if resp.StatusCode == http.StatusTooManyRequests {
				return false, rateLimitError(resp.StatusCode, resp.Header, body, time.Now())
			}
			return false, newAPIError(resp.StatusCode, body)
		}
		return false, fmt.Errorf("invezgo: realtime dial: %w", err)
	}
	conn.SetReadLimit(16 << 20) // trade batches can exceed the 32 KiB default
	defer func() {
		r.mu.Lock()
		r.conn = nil
		r.mu.Unlock()
		conn.CloseNow()
	}()
	for {
		_, data, err := conn.Read(ctx)
		if err != nil {
			return connected, err
		}
		if r.dispatch(ctx, conn, data, h) {
			connected = true
		}
	}
}

// dispatch decodes one server message; it returns true for "connected".
func (r *Realtime) dispatch(ctx context.Context, conn *websocket.Conn, data []byte, h RealtimeHandlers) bool {
	var env struct {
		Type     string          `json:"type"`
		Message  string          `json:"message"`
		Missed   int64           `json:"missed"`
		Channel  string          `json:"channel"`
		Codes    json.RawMessage `json:"codes"`
		MinLot   *int64          `json:"min_lot"`
		MinValue *int64          `json:"min_value"`
	}
	unknown := func() {
		if h.OnUnknown != nil {
			h.OnUnknown(data)
		}
	}
	if json.Unmarshal(data, &env) != nil {
		unknown()
		return false
	}
	switch env.Type {
	case "connected":
		r.mu.Lock()
		r.conn = conn
		err := r.replayLocked(ctx)
		r.mu.Unlock()
		if err != nil {
			conn.CloseNow() // the next Read fails and Run reconnects
			return true
		}
		if h.OnConnected != nil {
			h.OnConnected()
		}
		return true
	case "subscribed":
		s := RealtimeSubscribed{Channel: env.Channel, MinLot: env.MinLot, MinValue: env.MinValue}
		if string(env.Codes) == `"*"` {
			s.All = true
		} else if json.Unmarshal(env.Codes, &s.Codes) != nil && len(env.Codes) > 0 {
			unknown()
			return false
		}
		if h.OnSubscribed != nil {
			h.OnSubscribed(s)
		}
	case "pong":
		if h.OnPong != nil {
			h.OnPong()
		}
	case "error":
		if h.OnError != nil {
			h.OnError(env.Message)
		}
	case "lagged":
		if h.OnLagged != nil {
			h.OnLagged(env.Missed)
		}
	case "trades":
		var m RealtimeTrades
		if json.Unmarshal(data, &m) != nil {
			unknown()
		} else if h.OnTrades != nil {
			h.OnTrades(m)
		}
	case "orderbook":
		var m RealtimeOrderbook
		if json.Unmarshal(data, &m) != nil {
			unknown()
		} else if h.OnOrderbook != nil {
			h.OnOrderbook(m)
		}
	case "stock":
		var m RealtimeStock
		if json.Unmarshal(data, &m) != nil {
			unknown()
		} else if h.OnStock != nil {
			h.OnStock(m)
		}
	case "index":
		var m RealtimeIndex
		if json.Unmarshal(data, &m) != nil {
			unknown()
		} else if h.OnIndex != nil {
			h.OnIndex(m)
		}
	default:
		unknown()
	}
	return false
}

func upper(codes []string) []string {
	out := make([]string, 0, len(codes))
	for _, c := range codes {
		out = append(out, strings.ToUpper(strings.TrimSpace(c)))
	}
	return out
}

func mergeCodes(have, add []string) []string {
	out := append([]string(nil), have...)
	for _, c := range add {
		if !containsCode(out, c) {
			out = append(out, c)
		}
	}
	return out
}

func removeCodes(have, drop []string) []string {
	var out []string
	for _, c := range have {
		if !containsCode(drop, c) {
			out = append(out, c)
		}
	}
	return out
}

func containsCode(list []string, c string) bool {
	for _, x := range list {
		if x == c {
			return true
		}
	}
	return false
}
