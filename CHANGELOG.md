# Changelog

All notable changes to this project are documented in this file.

## v2.0.0

Major release: the SDK is rebuilt and generated from the official OpenAPI specification (Invezgo API 1.2).

### Added
- Module path `github.com/invezgo/invezgo-go-sdk/v2` (Go 1.21+).
- 107 endpoints in 11 resources (`client.Analysis`, `client.Batch`, `client.Search`, `client.Screener`,
  `client.Usage`, `client.Watchlists`, `client.Journals`, `client.Portfolios`, `client.Alerts`,
  `client.Posts`, `client.Trades`) with typed responses, params, request bodies and enums.
- In-memory LRU cache enabled by default (TTL 60s, realtime 5s, max 1000 entries / 32 MiB),
  automatic purge of personal data after mutations, `client.ClearCache()`.
- Automatic API key rotation on HTTP 429 (`WithAdditionalAPIKeys`): per-minute limits block the key for
  that endpoint only, monthly quota blocks the whole key; honours `Retry-After` / `X-RateLimit-Reset`.
- Typed errors `*APIError` and `*RateLimitError` with helpers `IsAuthentication`, `IsPaymentRequired`,
  `IsNotFound`, `IsValidation`, `IsRateLimited`, `IsServerError`.
- `FlexFloat`/`FlexInt` for numbers that may be sent as strings, `Ptr()` for optional fields.
- `docs/API.md` reference.

### Changed (breaking)
- `NewClient` returns `(*Client, error)`.
- Methods are grouped by resource: `client.GetStockList(ctx)` → `client.Analysis.GetStockList(ctx)`.
- Query parameters use `<Name>Params` structs; dates are `YYYY-MM-DD` strings.
- `APIError.ErrorType` → `APIError.ErrorText`.

### Removed
- Endpoints not in the official spec (AI, profile, membership, recommendation, journal file extraction).

### Migrating from v1

1. Update the import path and dependencies:

   ```go
   // v1
   import "github.com/invezgo/invezgo-go-sdk"
   // v2
   import "github.com/invezgo/invezgo-go-sdk/v2"
   ```

   ```bash
   go get github.com/invezgo/invezgo-go-sdk/v2 && go mod tidy
   ```
2. `NewClient` now returns an error: `client, err := invezgo.NewClient(key)`.
3. Methods are grouped by resource: `client.GetStockList(ctx)` → `client.Analysis.GetStockList(ctx)`,
   `client.ListWatchlist(ctx)` → `client.Watchlists.List(ctx, nil)`, and so on.
4. Query parameters are `<Name>Params` structs. Dates are `YYYY-MM-DD` strings instead of `time.Time`;
   optional fields are pointers set with `invezgo.Ptr(...)`.
5. Response type names follow the OpenAPI spec (`[]StockListItem`, `*TopChangeResponse`, ...).
6. Errors: `*invezgo.APIError` (`ErrorType` → `ErrorText`), `*invezgo.RateLimitError`, and the `IsXxx` helpers.
7. Caching and key rotation are enabled by default; use `invezgo.WithoutCache()` to keep v1 behaviour.
8. Endpoints that are not in the official spec were removed (AI, profile, membership, recommendation,
   journal file extraction).

## v1.0.0
- Initial release.
