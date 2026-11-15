package invezgo

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
	"testing"
	"time"
)

const whSecret = "whsec_test"

func whSign(ts string, body []byte) string {
	mac := hmac.New(sha256.New, []byte(whSecret))
	mac.Write([]byte(ts + "." + string(body)))
	return "sha256=" + hex.EncodeToString(mac.Sum(nil))
}

func whReason(t *testing.T, err error) string {
	t.Helper()
	var we *WebhookError
	if !errors.As(err, &we) {
		t.Fatalf("want *WebhookError, got %v", err)
	}
	return we.Reason
}

func TestWebhookVerify(t *testing.T) {
	body := []byte(`{"id":"evt_1","domain":"alert","type":"triggered","created_at":1762412400,
		"data":{"alert_id":"a1","alert_name":"Breakout","stocks_in":[{"code":"BBCA","price":9875}],"stocks_out":[]}}`)
	now := time.Unix(1762412400, 0)
	clock := WithWebhookNow(func() time.Time { return now })
	ts := "1762412400"

	// lower-case, non-canonical header keys must still be found
	h := http.Header{"x-invezgo-timestamp": {ts}, "x-invezgo-signature": {whSign(ts, body)}}
	ev, err := VerifyWebhook(body, h, whSecret, clock)
	if err != nil {
		t.Fatal(err)
	}
	if ev.Name() != "alert.triggered" || ev.ID != "evt_1" || ev.CreatedAt != 1762412400 {
		t.Fatalf("event %+v", ev)
	}
	d, err := ev.AlertTriggeredData()
	if err != nil || d.AlertID != "a1" || len(d.StocksIn) != 1 || d.StocksIn[0].Price != 9875 {
		t.Fatalf("data %+v %v", d, err)
	}

	cases := []struct {
		name, ts, sig string
		body          []byte
		reason        string
	}{
		{"bad signature", ts, whSign(ts, body) + "0", body, WebhookReasonSignature},
		{"other secret body", ts, whSign(ts, []byte("x")), body, WebhookReasonSignature},
		{"stale", "1762412000", whSign("1762412000", body), body, WebhookReasonTimestamp},
		{"future", "1762412800", whSign("1762412800", body), body, WebhookReasonTimestamp},
		{"non-digit", "+1762412400", whSign("+1762412400", body), body, WebhookReasonTimestamp},
		{"empty timestamp", "", whSign("", body), body, WebhookReasonTimestamp},
		{"bad json", ts, whSign(ts, []byte("{bad")), []byte("{bad"), WebhookReasonPayload},
	}
	for _, c := range cases {
		_, err := VerifyWebhookSignature(c.body, c.sig, c.ts, whSecret, clock)
		if got := whReason(t, err); got != c.reason {
			t.Errorf("%s: reason %q, want %q", c.name, got, c.reason)
		}
	}

	// tolerance option widens the window
	old := "1762412000"
	if _, err := VerifyWebhookSignature(body, whSign(old, body), old, whSecret, clock, WithWebhookTolerance(time.Hour)); err != nil {
		t.Fatal(err)
	}
}

func TestParseWebhookTypedData(t *testing.T) {
	ev, err := ParseWebhook([]byte(`{"id":"q","domain":"quota","type":"low","created_at":1,
		"data":{"plan":"PRIME","limit":90000,"used":85500,"remaining":4500,"reset_at":null}}`))
	if err != nil {
		t.Fatal(err)
	}
	q, err := ev.QuotaData()
	if err != nil || q.Plan != "PRIME" || q.Remaining != 4500 || q.ResetAt != nil {
		t.Fatalf("quota %+v %v", q, err)
	}
	ev, _ = ParseWebhook([]byte(`{"domain":"update","type":"news","data":{"title":"T","date":"2026-11-06 16:32:19","path":["/analysis/news"]}}`))
	u, err := ev.UpdateData()
	if err != nil || u.Title != "T" || len(u.Path) != 1 {
		t.Fatalf("update %+v %v", u, err)
	}
	if _, err := ParseWebhook([]byte("nope")); whReason(t, err) != WebhookReasonPayload {
		t.Fatal("want payload error")
	}
}
