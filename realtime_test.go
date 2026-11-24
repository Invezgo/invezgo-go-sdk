package invezgo

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/coder/websocket"
)

// wsServer starts a mock /ws endpoint; serve runs for each accepted connection n (1-based).
func wsServer(t *testing.T, serve func(t *testing.T, c *websocket.Conn, n int)) (*Client, *int32) {
	t.Helper()
	var conns int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/ws" || r.Header.Get("Authorization") != "Bearer k1" ||
			r.Header.Get("User-Agent") != "invezgo-go-sdk/"+Version {
			http.Error(w, `{"message":"bad upgrade request"}`, http.StatusBadRequest)
			return
		}
		c, err := websocket.Accept(w, r, nil)
		if err != nil {
			return
		}
		defer c.CloseNow()
		serve(t, c, int(atomic.AddInt32(&conns, 1)))
	}))
	t.Cleanup(srv.Close)
	client, err := NewClient("k1", WithBaseURL(srv.URL))
	if err != nil {
		t.Fatal(err)
	}
	return client, &conns
}

func wsSend(c *websocket.Conn, msg string) {
	_ = c.Write(context.Background(), websocket.MessageText, []byte(msg))
}

func wsRead(t *testing.T, c *websocket.Conn) string {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, b, err := c.Read(ctx)
	if err != nil {
		t.Errorf("server read: %v", err)
	}
	return string(b)
}

func runRealtime(t *testing.T, rt *Realtime, h RealtimeHandlers) chan error {
	t.Helper()
	rt.unit = 10 * time.Millisecond
	done := make(chan error, 1)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	t.Cleanup(cancel)
	go func() { done <- rt.Run(ctx, h) }()
	return done
}

func wait(t *testing.T, done chan error) error {
	t.Helper()
	select {
	case err := <-done:
		return err
	case <-time.After(10 * time.Second):
		t.Fatal("Run did not return")
		return nil
	}
}

func TestRealtimeReplaysQueuedSubscriptionsAndDecodes(t *testing.T) {
	ctx := context.Background()
	client, _ := wsServer(t, func(t *testing.T, c *websocket.Conn, n int) {
		wsSend(c, `{"type":"connected"}`)
		if got, want := wsRead(t, c), `{"op":"subscribe","channel":"trades","codes":["BBCA","TLKM"],"min_lot":500}`; got != want {
			t.Errorf("trades replay %s, want %s", got, want)
		}
		if got, want := wsRead(t, c), `{"op":"subscribe","channel":"orderbook","codes":["BBCA"]}`; got != want {
			t.Errorf("orderbook replay %s, want %s", got, want)
		}
		if got, want := wsRead(t, c), `{"op":"subscribe","channel":"stock","codes":["BBCA"]}`; got != want {
			t.Errorf("stock replay %s, want %s", got, want)
		}
		if got, want := wsRead(t, c), `{"op":"subscribe","channel":"index","codes":["COMPOSITE"]}`; got != want {
			t.Errorf("index replay %s, want %s", got, want)
		}
		wsSend(c, `{"type":"subscribed","channel":"trades","codes":"*","min_lot":0,"min_value":0}`)
		wsSend(c, `{"type":"trades","ts":1793240653184,"data":[{"id":"202610290000255958","time":"092413","code":"ISAT","price":1465.0,"lot":9,"value":1318500,"change":-105.0,"side":"buy"}]}`)
		wsSend(c, `{"type":"orderbook","code":"BBCA","side":"bid","levels":[[9875,120],{"p":9850}],"seq":null}`)
		wsSend(c, `{"type":"stock","code":"BBCA","last":5650.0,"prev":5550.0,"open":null,"change_pct":1.8,"volume":114428700,"iev":7927800,"time":"110436","seq":4408664}`)
		wsSend(c, `{"type":"index","code":"Investor33","last":5221.553,"change_pct":1.462,"volume":20865566335,"time":"110512"}`)
		wsSend(c, `{"type":"lagged","missed":120}`)
		wsSend(c, `{"type":"brand_new","x":1}`)
		if got := wsRead(t, c); got != `{"op":"ping"}` {
			t.Errorf("ping %s", got)
		}
		wsSend(c, `{"type":"pong"}`)
		_, _, _ = c.Read(context.Background()) // wait for the client close frame
	})
	rt := client.Realtime()
	if rt.url[:5] != "ws://" {
		t.Fatalf("url %s", rt.url)
	}
	if err := rt.SubscribeTrades(ctx, []string{"bbca", "TLKM"}, &TradeFilter{MinLot: Ptr(int64(500))}); err != nil {
		t.Fatal(err)
	}
	if err := rt.SubscribeOrderbook(ctx, []string{"bbca"}); err != nil {
		t.Fatal(err)
	}
	if err := rt.SubscribeStock(ctx, []string{"bbca"}); err != nil {
		t.Fatal(err)
	}
	if err := rt.SubscribeIndex(ctx, []string{"composite", "LQ45"}); err != nil {
		t.Fatal(err)
	}
	if err := rt.Unsubscribe(ctx, ChannelIndex, []string{"lq45"}); err != nil {
		t.Fatal(err)
	}

	var mu sync.Mutex
	var got []string
	var trades RealtimeTrades
	var book RealtimeOrderbook
	var sub RealtimeSubscribed
	var stock RealtimeStock
	var index RealtimeIndex
	note := func(s string) { mu.Lock(); got = append(got, s); mu.Unlock() }
	done := runRealtime(t, rt, RealtimeHandlers{
		OnConnected:  func() { note("connected") },
		OnSubscribed: func(s RealtimeSubscribed) { sub = s; note("subscribed") },
		OnTrades:     func(m RealtimeTrades) { trades = m; note("trades") },
		OnOrderbook:  func(m RealtimeOrderbook) { book = m; note("orderbook") },
		OnStock:      func(m RealtimeStock) { stock = m; note("stock") },
		OnIndex:      func(m RealtimeIndex) { index = m; note("index") },
		OnLagged:     func(n int64) { note("lagged") },
		OnUnknown: func(raw []byte) {
			note("unknown")
			if err := rt.Ping(ctx); err != nil {
				t.Error(err)
			}
		},
		OnPong: func() { note("pong"); go rt.Close() },
	})
	if err := wait(t, done); err != nil {
		t.Fatalf("Run: %v", err)
	}
	if want := "[connected subscribed trades orderbook stock index lagged unknown pong]"; fmt.Sprint(got) != want {
		t.Fatalf("events %v, want %s", got, want)
	}
	if !sub.All || sub.MinLot == nil || *sub.MinLot != 0 {
		t.Errorf("subscribed %+v", sub)
	}
	tr := trades.Data[0]
	if trades.TS != 1793240653184 || tr.ID != "202610290000255958" || tr.Price != 1465 || tr.Lot != 9 || tr.Value != 1318500 || tr.Side != "buy" {
		t.Errorf("trade %+v", trades)
	}
	if book.Code != "BBCA" || book.Side != "bid" || len(book.Levels) != 2 || book.Seq != nil {
		t.Errorf("orderbook %+v", book)
	}
	if stock.Last != 5650 || *stock.Prev != 5550 || stock.Open != nil || *stock.Volume != 114428700 || *stock.IEV != 7927800 || *stock.Seq != 4408664 || stock.Time != "110436" {
		t.Errorf("stock %+v", stock)
	}
	if index.Code != "Investor33" || index.Last != 5221.553 || *index.ChangePct != 1.462 || index.Prev != nil {
		t.Errorf("index %+v", index)
	}
}

func TestRealtimeReconnectReplaysSubscriptions(t *testing.T) {
	ctx := context.Background()
	client, conns := wsServer(t, func(t *testing.T, c *websocket.Conn, n int) {
		wsSend(c, `{"type":"connected"}`)
		if got, want := wsRead(t, c), `{"op":"subscribe","channel":"trades"}`; got != want {
			t.Errorf("conn %d trades %s", n, got)
		}
		if n == 1 {
			if got := wsRead(t, c); got != `{"op":"subscribe","channel":"orderbook","codes":["BBCA"]}` {
				t.Errorf("conn 1 orderbook %s", got)
			}
			if got := wsRead(t, c); got != `{"op":"subscribe","channel":"orderbook","codes":["TLKM"]}` {
				t.Errorf("conn 1 live orderbook %s", got)
			}
			return // drop the connection: the client must reconnect
		}
		if got := wsRead(t, c); got != `{"op":"subscribe","channel":"orderbook","codes":["BBCA","TLKM"]}` {
			t.Errorf("conn 2 orderbook replay %s", got)
		}
		_, _, _ = c.Read(context.Background()) // wait for the client close frame
	})
	rt := client.Realtime()
	_ = rt.SubscribeTrades(ctx, nil, nil)
	_ = rt.SubscribeOrderbook(ctx, []string{"BBCA"})
	var connected, reconnects int32
	done := runRealtime(t, rt, RealtimeHandlers{
		OnConnected: func() {
			if atomic.AddInt32(&connected, 1) == 1 {
				// live subscription while connected: sent now and remembered for replay
				_ = rt.SubscribeOrderbook(ctx, []string{"tlkm"})
			} else {
				go rt.Close()
			}
		},
		OnReconnect: func(attempt int, delay time.Duration, err error) {
			atomic.AddInt32(&reconnects, 1)
			if attempt != 1 || delay != 10*time.Millisecond {
				t.Errorf("reconnect attempt %d delay %s", attempt, delay)
			}
		},
	})
	if err := wait(t, done); err != nil {
		t.Fatal(err)
	}
	if atomic.LoadInt32(conns) != 2 || atomic.LoadInt32(&reconnects) != 1 {
		t.Fatalf("conns %d reconnects %d", *conns, reconnects)
	}
}

func TestRealtimeStopsOn4400(t *testing.T) {
	client, conns := wsServer(t, func(t *testing.T, c *websocket.Conn, n int) {
		wsSend(c, `{"type":"connected"}`)
		wsSend(c, `{"type":"error","message":"unknown op"}`)
		c.Close(4400, "invalid message")
	})
	var msg string
	err := wait(t, runRealtime(t, client.Realtime(), RealtimeHandlers{OnError: func(m string) { msg = m }}))
	if websocket.CloseStatus(err) != 4400 || msg != "unknown op" {
		t.Fatalf("err %v msg %q", err, msg)
	}
	if atomic.LoadInt32(conns) != 1 {
		t.Fatalf("reconnected after 4400: %d connections", *conns)
	}
}

func TestRealtimeHandshakeErrors(t *testing.T) {
	var calls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"statusCode":401,"message":"Invalid API key","error":"Unauthorized"}`))
	}))
	defer srv.Close()
	rt, err := NewRealtime("bad", "ws"+srv.URL[4:]+"/ws")
	if err != nil {
		t.Fatal(err)
	}
	err = wait(t, runRealtime(t, rt, RealtimeHandlers{}))
	if !IsAuthentication(err) || atomic.LoadInt32(&calls) != 1 {
		t.Fatalf("err %v calls %d", err, calls)
	}
}

func TestRealtimeGuards(t *testing.T) {
	ctx := context.Background()
	rt, _ := NewRealtime("k", "")
	if rt.url != DefaultRealtimeURL {
		t.Fatal(rt.url)
	}
	if rt.SubscribeOrderbook(ctx, nil) == nil {
		t.Error("empty orderbook codes accepted")
	}
	many := make([]string, 21)
	for i := range many {
		many[i] = fmt.Sprintf("C%d", i)
	}
	if rt.SubscribeOrderbook(ctx, many) == nil {
		t.Error("21 orderbook codes accepted")
	}
	trades := make([]string, 201)
	for i := range trades {
		trades[i] = fmt.Sprintf("C%d", i)
	}
	if rt.SubscribeTrades(ctx, trades, nil) == nil {
		t.Error("201 trade codes accepted")
	}
	stocks := make([]string, 51)
	for i := range stocks {
		stocks[i] = fmt.Sprintf("C%d", i)
	}
	if rt.SubscribeStock(ctx, nil) == nil || rt.SubscribeStock(ctx, stocks) == nil {
		t.Error("empty or 51 stock codes accepted")
	}
	if rt.SubscribeIndex(ctx, trades[:101]) == nil {
		t.Error("101 index codes accepted")
	}
	if rt.Unsubscribe(ctx, "quotes", nil) == nil {
		t.Error("unknown channel accepted")
	}
}
