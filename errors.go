package invezgo

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// APIError is returned for every non-2xx response. Use errors.As to inspect it:
//
//	var apiErr *invezgo.APIError
//	if errors.As(err, &apiErr) { log.Println(apiErr.StatusCode, apiErr.RequestID) }
type APIError struct {
	// StatusCode is the HTTP status (400, 401, 402, 404, 422, 429, 5xx...).
	StatusCode int
	// Message is the human-readable message (array messages are joined with "; ").
	Message string
	// ErrorText is the HTTP status text from the body, e.g. "Unauthorized".
	ErrorText string
	// RequestID identifies the request; quote it when contacting support.
	RequestID string
	// Path is the request path reported by the server.
	Path string
	// Body is the raw response body (JSON error envelope or text).
	Body []byte
}

func (e *APIError) Error() string {
	if e.RequestID != "" {
		return fmt.Sprintf("invezgo: %d %s (request %s)", e.StatusCode, e.Message, e.RequestID)
	}
	return fmt.Sprintf("invezgo: %d %s", e.StatusCode, e.Message)
}

// RateLimitError is returned on HTTP 429 once every API key is rate limited. errors.As works with
// both *RateLimitError and *APIError.
type RateLimitError struct {
	APIError
	// RetryAfter is how long until the earliest key can be used again.
	RetryAfter time.Duration
	// Limit and Remaining mirror the X-RateLimit-Limit / X-RateLimit-Remaining headers (-1 if absent).
	Limit, Remaining int
	// Reset mirrors X-RateLimit-Reset (zero if absent).
	Reset time.Time
	// QuotaExceeded is true when the monthly quota (not the per-minute limit) is used up.
	QuotaExceeded bool
}

// Unwrap exposes the embedded APIError to errors.As.
func (e *RateLimitError) Unwrap() error { return &e.APIError }

func statusIs(err error, match func(int) bool) bool {
	var e *APIError
	return errors.As(err, &e) && match(e.StatusCode)
}

// IsValidation reports a 400/422 error (invalid parameters or formula).
func IsValidation(err error) bool {
	return statusIs(err, func(s int) bool { return s == 400 || s == 422 })
}

// IsAuthentication reports a 401 error (missing, invalid or expired API key).
func IsAuthentication(err error) bool { return statusIs(err, func(s int) bool { return s == 401 }) }

// IsPaymentRequired reports a 402 error (endpoint not included in the subscription plan).
func IsPaymentRequired(err error) bool { return statusIs(err, func(s int) bool { return s == 402 }) }

// IsNotFound reports a 404 error.
func IsNotFound(err error) bool { return statusIs(err, func(s int) bool { return s == 404 }) }

// IsRateLimited reports a 429 error (all API keys rate limited or out of quota).
func IsRateLimited(err error) bool { return statusIs(err, func(s int) bool { return s == 429 }) }

// IsServerError reports a 5xx error.
func IsServerError(err error) bool { return statusIs(err, func(s int) bool { return s >= 500 }) }

func newAPIError(status int, body []byte) *APIError {
	e := &APIError{StatusCode: status, Body: body}
	var env struct {
		Message   json.RawMessage `json:"message"`
		Error     string          `json:"error"`
		RequestID string          `json:"requestId"`
		Path      string          `json:"path"`
	}
	if json.Unmarshal(body, &env) == nil {
		e.ErrorText, e.RequestID, e.Path = env.Error, env.RequestID, env.Path
		var s string
		var list []string
		if json.Unmarshal(env.Message, &s) == nil {
			e.Message = s
		} else if json.Unmarshal(env.Message, &list) == nil {
			e.Message = strings.Join(list, "; ")
		}
	} else if t := strings.TrimSpace(string(body)); t != "" && len(t) < 500 {
		e.Message = t
	}
	if e.Message == "" {
		e.Message = http.StatusText(status)
	}
	return e
}
