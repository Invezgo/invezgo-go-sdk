package invezgo

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Reasons reported by WebhookError.
const (
	WebhookReasonTimestamp = "timestamp"
	WebhookReasonSignature = "signature"
	WebhookReasonPayload   = "payload"
)

// WebhookError is returned when a webhook delivery fails verification or decoding.
// It is not an *APIError.
type WebhookError struct {
	// Reason is WebhookReasonTimestamp, WebhookReasonSignature or WebhookReasonPayload.
	Reason string
	// Err is the underlying decode error, if any.
	Err error
}

func (e *WebhookError) Error() string {
	if e.Err != nil {
		return "invezgo: webhook " + e.Reason + ": " + e.Err.Error()
	}
	return "invezgo: invalid webhook " + e.Reason
}

// Unwrap returns the underlying error.
func (e *WebhookError) Unwrap() error { return e.Err }

// WebhookEvent is one webhook delivery. Domain and Type are plain strings so new events never
// break decoding; decode Data with the typed helpers (UpdateData, AlertTriggeredData, ...).
type WebhookEvent struct {
	// ID is the delivery ID (same as X-Invezgo-Delivery); stable across retries, use it to deduplicate.
	ID string `json:"id"`
	// Domain is the event group: test, quota, update or alert.
	Domain string `json:"domain"`
	// Type is the event type within the domain, e.g. summary or news.
	Type string `json:"type"`
	// CreatedAt is the unix time (seconds) the event was created.
	CreatedAt int64 `json:"created_at"`
	// Data is the event-specific payload.
	Data json.RawMessage `json:"data"`
}

// Name returns "{domain}.{type}", e.g. "update.summary" (same as the X-Invezgo-Event header).
func (e *WebhookEvent) Name() string { return e.Domain + "." + e.Type }

// WebhookUpdateData is the data of every update.* event.
type WebhookUpdateData struct {
	// Date is YYYY-MM-DD (YYYY-MM-DD HH:mm:ss WIB for update.news).
	Date string `json:"date"`
	// Path lists the API endpoints whose data changed.
	Path  []string `json:"path"`
	Code  string   `json:"code,omitempty"`
	Title string   `json:"title,omitempty"`
}

// WebhookAlertStock is a stock that entered or left a Live Alert formula.
type WebhookAlertStock struct {
	Code  string  `json:"code"`
	Price float64 `json:"price"`
}

// WebhookAlertTriggeredData is the data of alert.triggered.
type WebhookAlertTriggeredData struct {
	AlertID   string              `json:"alert_id"`
	AlertName string              `json:"alert_name"`
	StocksIn  []WebhookAlertStock `json:"stocks_in"`
	StocksOut []WebhookAlertStock `json:"stocks_out"`
}

// WebhookAlertDisabledData is the data of alert.disabled.
type WebhookAlertDisabledData struct {
	AlertID   string `json:"alert_id"`
	AlertName string `json:"alert_name"`
	Reason    string `json:"reason"`
}

// WebhookTestData is the data of test.ping.
type WebhookTestData struct {
	Message string `json:"message"`
}

// WebhookQuotaData is the data of quota.low and quota.exhausted.
type WebhookQuotaData struct {
	Plan      string `json:"plan"`
	Limit     int64  `json:"limit"`
	Used      int64  `json:"used"`
	Remaining int64  `json:"remaining"`
	// ResetAt is the unix time (seconds) the quota resets; may be nil.
	ResetAt *int64 `json:"reset_at"`
}

// UpdateData decodes the data of an update.* event.
func (e *WebhookEvent) UpdateData() (*WebhookUpdateData, error) {
	d := &WebhookUpdateData{}
	return d, e.decodeData(d)
}

// AlertTriggeredData decodes the data of alert.triggered.
func (e *WebhookEvent) AlertTriggeredData() (*WebhookAlertTriggeredData, error) {
	d := &WebhookAlertTriggeredData{}
	return d, e.decodeData(d)
}

// AlertDisabledData decodes the data of alert.disabled.
func (e *WebhookEvent) AlertDisabledData() (*WebhookAlertDisabledData, error) {
	d := &WebhookAlertDisabledData{}
	return d, e.decodeData(d)
}

// TestData decodes the data of test.ping.
func (e *WebhookEvent) TestData() (*WebhookTestData, error) {
	d := &WebhookTestData{}
	return d, e.decodeData(d)
}

// QuotaData decodes the data of quota.low and quota.exhausted.
func (e *WebhookEvent) QuotaData() (*WebhookQuotaData, error) {
	d := &WebhookQuotaData{}
	return d, e.decodeData(d)
}

func (e *WebhookEvent) decodeData(v any) error {
	if len(e.Data) == 0 || string(e.Data) == "null" {
		return nil
	}
	if err := json.Unmarshal(e.Data, v); err != nil {
		return &WebhookError{Reason: WebhookReasonPayload, Err: err}
	}
	return nil
}

// WebhookOption configures VerifyWebhook.
type WebhookOption func(*webhookConfig)

type webhookConfig struct {
	tolerance time.Duration
	now       func() time.Time
}

// WithWebhookTolerance sets the maximum age (and clock skew) of X-Invezgo-Timestamp. Default 5 minutes.
func WithWebhookTolerance(d time.Duration) WebhookOption {
	return func(c *webhookConfig) { c.tolerance = d }
}

// WithWebhookNow overrides the clock used to check the timestamp (tests).
func WithWebhookNow(now func() time.Time) WebhookOption {
	return func(c *webhookConfig) { c.now = now }
}

// VerifyWebhook checks the X-Invezgo-Timestamp and X-Invezgo-Signature headers (looked up
// case-insensitively) against rawBody, the exact bytes received, and returns the decoded event.
// It returns a *WebhookError on failure.
//
//	body, _ := io.ReadAll(r.Body)
//	event, err := invezgo.VerifyWebhook(body, r.Header, os.Getenv("INVEZGO_WEBHOOK_SECRET"))
func VerifyWebhook(rawBody []byte, header http.Header, secret string, opts ...WebhookOption) (*WebhookEvent, error) {
	return VerifyWebhookSignature(rawBody, headerFold(header, "X-Invezgo-Signature"),
		headerFold(header, "X-Invezgo-Timestamp"), secret, opts...)
}

// VerifyWebhookSignature is VerifyWebhook with the signature and timestamp header values given directly.
func VerifyWebhookSignature(rawBody []byte, signature, timestamp, secret string, opts ...WebhookOption) (*WebhookEvent, error) {
	cfg := webhookConfig{tolerance: 5 * time.Minute, now: time.Now}
	for _, opt := range opts {
		opt(&cfg)
	}
	ts, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil || strings.Trim(timestamp, "0123456789") != "" {
		return nil, &WebhookError{Reason: WebhookReasonTimestamp}
	}
	if age := cfg.now().Sub(time.Unix(ts, 0)); age > cfg.tolerance || -age > cfg.tolerance {
		return nil, &WebhookError{Reason: WebhookReasonTimestamp}
	}
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(timestamp + "."))
	mac.Write(rawBody)
	expected := "sha256=" + hex.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(expected), []byte(signature)) {
		return nil, &WebhookError{Reason: WebhookReasonSignature}
	}
	return ParseWebhook(rawBody)
}

// ParseWebhook decodes rawBody without verifying it (tests or already verified bodies).
func ParseWebhook(rawBody []byte) (*WebhookEvent, error) {
	e := &WebhookEvent{}
	if err := json.Unmarshal(rawBody, e); err != nil {
		return nil, &WebhookError{Reason: WebhookReasonPayload, Err: err}
	}
	return e, nil
}

// headerFold looks a header up case-insensitively, also for non-canonical keys.
func headerFold(h http.Header, name string) string {
	if v := h.Get(name); v != "" {
		return v
	}
	for k, v := range h {
		if strings.EqualFold(k, name) && len(v) > 0 {
			return v[0]
		}
	}
	return ""
}
