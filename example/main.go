// Command example shows typical Invezgo Go SDK usage.
//
//	INVEZGO_API_KEY=xxx INVEZGO_API_KEYS_EXTRA=yyy,zzz go run ./example
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/invezgo/invezgo-go-sdk/v2"
)

func main() {
	var extra []string
	if v := os.Getenv("INVEZGO_API_KEYS_EXTRA"); v != "" {
		extra = strings.Split(v, ",")
	}

	client, err := invezgo.NewClient(os.Getenv("INVEZGO_API_KEY"),
		invezgo.WithAdditionalAPIKeys(extra...), // rotated automatically on HTTP 429
		invezgo.WithCache(invezgo.CacheConfig{TTL: 2 * time.Minute}),
	)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 1. Simple list (cached for 2 minutes; the second call never hits the network).
	stocks, err := client.Analysis.GetStockList(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d stocks listed, first: %s - %s\n", len(stocks), stocks[0].Code, stocks[0].Name)

	// 2. Required + optional query params with typed enums.
	top, err := client.Analysis.GetTopChange(ctx, invezgo.TopChangeParams{
		Date:           "2024-12-30",
		FilterColumn:   invezgo.Ptr(invezgo.TopChangeFilterColumnValue),
		FilterOperator: invezgo.Ptr(invezgo.TopChangeFilterOperatorGt),
		FilterValue:    invezgo.Ptr(1e9),
	})
	if err != nil {
		log.Fatal(err)
	}
	for _, g := range top.Gain {
		fmt.Printf("gainer %s %+.2f%%\n", g.Code, g.Change)
	}

	// 3. Path param + typed error handling.
	info, err := client.Analysis.GetInformation(ctx, "BBCA")
	var rl *invezgo.RateLimitError
	switch {
	case errors.As(err, &rl):
		fmt.Printf("all keys rate limited, retry in %s (quota exceeded: %v)\n", rl.RetryAfter, rl.QuotaExceeded)
	case invezgo.IsPaymentRequired(err):
		fmt.Println("this endpoint is not included in your plan")
	case err != nil:
		log.Fatal(err)
	default:
		fmt.Printf("%s: %s\n", info.Code, info.Name)
	}

	// 4. Quota monitoring (never cached, does not consume quota).
	usage, err := client.Usage.GetAPIUsage(ctx)
	if err == nil {
		fmt.Printf("plan %s: %d/%d requests used\n", usage.Role, usage.Usage, usage.Limit)
	}
}
