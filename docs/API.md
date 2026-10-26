# Invezgo Go SDK — API Reference

Generated from the official OpenAPI spec (Invezgo API 1.2) by `codegen/gen_go.py`. Do not edit by hand.

All methods take `ctx context.Context` first and return `(result, error)`. Cache tiers: **default** (`CacheConfig.TTL`, 60s), **realtime** (`CacheConfig.RealtimeTTL`, 5s), **none** (never cached).

## Resources

- [`client.Analysis`](#clientanalysis) — Indonesia stock market data: lists, prices, charts, broker flow, shareholders, fundamentals.
- [`client.Batch`](#clientbatch) — Batch endpoints: fetch many stocks in one request (counted as 1 request).
- [`client.Search`](#clientsearch) — Search stocks and users.
- [`client.Screener`](#clientscreener) — Run stock screeners and manage saved screener presets.
- [`client.Usage`](#clientusage) — API usage and quota monitoring (does not consume quota).
- [`client.Watchlists`](#clientwatchlists) — Personal watchlists and watchlist groups.
- [`client.Journals`](#clientjournals) — Personal trading journal transactions.
- [`client.Portfolios`](#clientportfolios) — Personal portfolio built from journal transactions.
- [`client.Alerts`](#clientalerts) — Live formula alerts.
- [`client.Posts`](#clientposts) — Community posts, comments, likes and votes.
- [`client.Trades`](#clienttrades) — Personal realized trades.

## client.Analysis

Indonesia stock market data: lists, prices, charts, broker flow, shareholders, fundamentals.

### GetStockList

**Stock List** — `GET /analysis/list/stock` — cache: **default**

```go
func (s *AnalysisService) GetStockList(ctx context.Context) ([]StockListItem, error)
```

Get a complete list of all companies listed on the Indonesia Stock Exchange (IDX).

Returns: `[]StockListItem`

### GetBrokerList

**Broker List** — `GET /analysis/list/broker` — cache: **default**

```go
func (s *AnalysisService) GetBrokerList(ctx context.Context) ([]BrokerListItem, error)
```

Get a complete list of all brokers/securities registered as exchange members at IDX.

Returns: `[]BrokerListItem`

### GetIndexList

**Index List** — `GET /analysis/list/index` — cache: **default**

```go
func (s *AnalysisService) GetIndexList(ctx context.Context) ([]IndexListItem, error)
```

Get a complete list of all indexes available at IDX.

Returns: `[]IndexListItem`

### GetUma

**Stock UMA List** — `GET /analysis/uma` — cache: **default**

```go
func (s *AnalysisService) GetUma(ctx context.Context, params *UmaParams) (*UmaResponse, error)
```

Get the list of UMA (Unusual Market Activity) stock announcements from the Indonesia Stock Exchange (IDX).

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `params.Page` | `*int64` | no | `1` | Nomor halaman, dimulai dari 1. Default: 1. Range: 1..-. |
| `params.Limit` | `*int64` | no | `20` | Jumlah data per halaman, 1-100. Default: 20. Range: 1..100. |
| `params.Code` | `*string` | no |  | Filter kode saham, mis. `BBCA` (opsional). |

Returns: `*UmaResponse`

### GetSuspend

**Stock Suspend List** — `GET /analysis/suspend` — cache: **default**

```go
func (s *AnalysisService) GetSuspend(ctx context.Context, params *SuspendParams) (*SuspendResponse, error)
```

Get the list of stock trading suspensions from the Indonesia Stock Exchange (IDX).

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `params.Page` | `*int64` | no | `1` | Nomor halaman, dimulai dari 1. Default: 1. Range: 1..-. |
| `params.Limit` | `*int64` | no | `20` | Jumlah data per halaman, 1-100. Default: 20. Range: 1..100. |
| `params.Code` | `*string` | no |  | Filter kode saham, mis. `BBCA` (opsional). |
| `params.Active` | `*bool` | no |  | `true`: hanya suspensi yang masih berlaku (end_date null); `false`: hanya yang sudah dicabut; kosongkan untuk semua (opsional). |

Returns: `*SuspendResponse`

### GetNotation

**Stock Notation List** — `GET /analysis/notation` — cache: **default**

```go
func (s *AnalysisService) GetNotation(ctx context.Context) ([]NotationItem, error)
```

Get a complete list of special notations for all stocks listed on the Indonesia Stock Exchange (IDX), grouped by stock code and effective date.

Returns: `[]NotationItem`

### GetInformation

**Company Information** — `GET /analysis/information/{code}` — cache: **default**

```go
func (s *AnalysisService) GetInformation(ctx context.Context, code string) (*InformationResponse, error)
```

Get complete and detailed information about a company based on stock code.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `code (path)` | `string` | yes |  | Kode saham emiten yang tercatat di BEI, mis. `BBCA`. Tidak case-sensitive. Daftar kode tersedia di `/analysis/list/stock`. |

Returns: `*InformationResponse`

### GetTopChange

**Top Gainer & Loser** — `GET /analysis/top/change` — cache: **default**

```go
func (s *AnalysisService) GetTopChange(ctx context.Context, params TopChangeParams) (*TopChangeResponse, error)
```

Get a list of stocks with the largest price increases and decreases on a specific day, complete with 5-day historical graph data.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `params.Date` | `string` | yes |  | Tanggal data (`YYYY-MM-DD`). Gunakan hari bursa; tanggal libur bursa dapat menghasilkan data kosong. |
| `params.FilterColumn` | `*TopChangeFilterColumn` | no |  | Kolom yang difilter: `change` = perubahan harga (%), `value` = nilai transaksi, `volume` = volume, `ratio` = rasio. Filter aktif hanya jika `filter_column`, `filter_operator`, dan `filter_value` dikirim bersamaan. One of: change, value, volume, ratio. |
| `params.FilterOperator` | `*TopChangeFilterOperator` | no |  | Operator pembanding untuk filter. Contoh: `filter_column=change&filter_operator=>&filter_value=5` menampilkan saham yang naik lebih dari 5%. One of: <, >, =, >=, <=, !=. |
| `params.FilterValue` | `*float64` | no |  | Nilai pembanding untuk filter. Mendukung bilangan desimal dan negatif, mis. `5` atau `-2.5`. |

Returns: `*TopChangeResponse`

### GetTopForeign

**Top Foreign Flow** — `GET /analysis/top/foreign` — cache: **default**

```go
func (s *AnalysisService) GetTopForeign(ctx context.Context, params TopForeignParams) (*TopForeignResponse, error)
```

Get a list of stocks with the largest foreign accumulation and distribution on a specific day, complete with 5-day historical graph data.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `params.Date` | `string` | yes |  | Tanggal data (`YYYY-MM-DD`). Gunakan hari bursa; tanggal libur bursa dapat menghasilkan data kosong. |
| `params.FilterColumn` | `*TopForeignFilterColumn` | no |  | Kolom yang difilter: `change` = perubahan harga (%), `value` = nilai transaksi, `volume` = volume, `ratio` = rasio. Filter aktif hanya jika `filter_column`, `filter_operator`, dan `filter_value` dikirim bersamaan. One of: change, value, volume, ratio. |
| `params.FilterOperator` | `*TopForeignFilterOperator` | no |  | Operator pembanding untuk filter. Contoh: `filter_column=change&filter_operator=>&filter_value=5` menampilkan saham yang naik lebih dari 5%. One of: <, >, =, >=, <=, !=. |
| `params.FilterValue` | `*float64` | no |  | Nilai pembanding untuk filter. Mendukung bilangan desimal dan negatif, mis. `5` atau `-2.5`. |

Returns: `*TopForeignResponse`

### GetTopAccumulation

**Top BDM Flow** — `GET /analysis/top/accumulation` — cache: **default**

```go
func (s *AnalysisService) GetTopAccumulation(ctx context.Context, params TopAccumulationParams) (*TopAccumulationResponse, error)
```

Get a list of stocks with the largest bandarmology accumulation and distribution on a specific day.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `params.Date` | `string` | yes |  | Tanggal data (`YYYY-MM-DD`). Gunakan hari bursa; tanggal libur bursa dapat menghasilkan data kosong. |
| `params.FilterColumn` | `*TopAccumulationFilterColumn` | no |  | Kolom yang difilter: `change` = perubahan harga (%), `value` = nilai transaksi, `volume` = volume, `ratio` = rasio. Filter aktif hanya jika `filter_column`, `filter_operator`, dan `filter_value` dikirim bersamaan. One of: change, value, volume, ratio. |
| `params.FilterOperator` | `*TopAccumulationFilterOperator` | no |  | Operator pembanding untuk filter. Contoh: `filter_column=change&filter_operator=>&filter_value=5` menampilkan saham yang naik lebih dari 5%. One of: <, >, =, >=, <=, !=. |
| `params.FilterValue` | `*float64` | no |  | Nilai pembanding untuk filter. Mendukung bilangan desimal dan negatif, mis. `5` atau `-2.5`. |

Returns: `*TopAccumulationResponse`

### GetTopRitel

**Top Ritel Flow** — `GET /analysis/top/ritel` — cache: **default**

```go
func (s *AnalysisService) GetTopRitel(ctx context.Context, params TopRitelParams) (*TopRitelResponse, error)
```

Get a list of stocks with the largest retail accumulation and distribution on a specific day.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `params.Date` | `string` | yes |  | Tanggal data (`YYYY-MM-DD`). Gunakan hari bursa; tanggal libur bursa dapat menghasilkan data kosong. |
| `params.FilterColumn` | `*TopRitelFilterColumn` | no |  | Kolom yang difilter: `change` = perubahan harga (%), `value` = nilai transaksi, `volume` = volume, `ratio` = rasio. Filter aktif hanya jika `filter_column`, `filter_operator`, dan `filter_value` dikirim bersamaan. One of: change, value, volume, ratio. |
| `params.FilterOperator` | `*TopRitelFilterOperator` | no |  | Operator pembanding untuk filter. Contoh: `filter_column=change&filter_operator=>&filter_value=5` menampilkan saham yang naik lebih dari 5%. One of: <, >, =, >=, <=, !=. |
| `params.FilterValue` | `*float64` | no |  | Nilai pembanding untuk filter. Mendukung bilangan desimal dan negatif, mis. `5` atau `-2.5`. |

Returns: `*TopRitelResponse`

### GetStockRunningTrade

**Stock Running Trade** — `GET /analysis/running-trade/{code}` — cache: **realtime**

```go
func (s *AnalysisService) GetStockRunningTrade(ctx context.Context, code string, params StockRunningTradeParams) (*StockRunningTradeResponse, error)
```

Get realtime/historical running trade transaction data for a specific stock code.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `code (path)` | `string` | yes |  | Kode saham emiten yang tercatat di BEI, mis. `BBCA`. Tidak case-sensitive. Daftar kode tersedia di `/analysis/list/stock`. |
| `params.Date` | `string` | yes |  | Tanggal data (`YYYY-MM-DD`). Gunakan hari bursa; tanggal libur bursa dapat menghasilkan data kosong. |
| `params.Page` | `*int64` | no | `1` | Nomor halaman, dimulai dari 1. Default: 1. Range: 1..-. |
| `params.Limit` | `*int64` | no | `50` | Jumlah transaksi per halaman. Default: 50. Range: 1..1000. |
| `params.Type` | `*StockRunningTradeType` | no | `ALL` | Sisi transaksi: `BUY` = match di harga offer (inisiatif beli), `SELL` = match di harga bid (inisiatif jual), `ALL` = seluruh transaksi. Diabaikan saat jam bursa (tanggal hari ini, 08:45–16:15 WIB pada hari bursa); hanya berlaku pada pengambilan data sehari penuh. One of: ALL, BUY, SELL. Default: ALL. |
| `params.Orderby` | `*StockRunningTradeOrderby` | no | `TIME` | Kolom pengurutan: `TIME` = waktu transaksi, `PRICE` = harga, `VOLUME` = volume. Diabaikan saat jam bursa (tanggal hari ini, 08:45–16:15 WIB pada hari bursa); hanya berlaku pada pengambilan data sehari penuh. One of: TIME, PRICE, VOLUME. Default: TIME. |
| `params.Minimum` | `*float64` | no |  | Volume minimal transaksi yang ditampilkan, dalam lembar (1 lot = 100 lembar). Berguna untuk menyaring transaksi besar. Diabaikan saat jam bursa (tanggal hari ini, 08:45–16:15 WIB pada hari bursa); hanya berlaku pada pengambilan data sehari penuh. Range: 0..-. |
| `params.Maximum` | `*float64` | no |  | Volume maksimal transaksi yang ditampilkan, dalam lembar (1 lot = 100 lembar). Diabaikan saat jam bursa (tanggal hari ini, 08:45–16:15 WIB pada hari bursa); hanya berlaku pada pengambilan data sehari penuh. Range: 0..-. |
| `params.Market` | `*StockRunningTradeMarket` | no | `ALL` | Jenis pasar: `ALL` = seluruh pasar, `RG` = Pasar Reguler, `NG` = Pasar Negosiasi, `TN` = Pasar Tunai. Diabaikan saat jam bursa (tanggal hari ini, 08:45–16:15 WIB pada hari bursa); hanya berlaku pada pengambilan data sehari penuh. One of: ALL, RG, NG, TN. Default: ALL. |
| `params.Pricefrom` | `*float64` | no |  | Harga terendah transaksi yang ditampilkan (Rupiah). Menampilkan transaksi dengan harga >= `pricefrom`. Diabaikan saat jam bursa (tanggal hari ini, 08:45–16:15 WIB pada hari bursa); hanya berlaku pada pengambilan data sehari penuh. Range: 1..-. |
| `params.Priceto` | `*float64` | no |  | Harga tertinggi transaksi yang ditampilkan (Rupiah). Menampilkan transaksi dengan harga <= `priceto`. Diabaikan saat jam bursa (tanggal hari ini, 08:45–16:15 WIB pada hari bursa); hanya berlaku pada pengambilan data sehari penuh. Range: 1..-. |
| `params.Timefrom` | `*string` | no |  | Jam awal transaksi dalam WIB, format `HH:mm` atau `HH:mm:ss`. Diabaikan saat jam bursa (tanggal hari ini, 08:45–16:15 WIB pada hari bursa); hanya berlaku pada pengambilan data sehari penuh. |
| `params.Timeto` | `*string` | no |  | Jam akhir transaksi dalam WIB, format `HH:mm` atau `HH:mm:ss`. Tanpa detik, dihitung hingga detik ke-59. Diabaikan saat jam bursa (tanggal hari ini, 08:45–16:15 WIB pada hari bursa); hanya berlaku pada pengambilan data sehari penuh. |
| `params.Sort` | `*StockRunningTradeSort` | no | `ASC` | Arah pengurutan: `ASC` = menaik (terlama lebih dulu), `DESC` = menurun (terbaru lebih dulu). Diabaikan saat jam bursa (tanggal hari ini, 08:45–16:15 WIB pada hari bursa); hanya berlaku pada pengambilan data sehari penuh. One of: ASC, DESC. Default: ASC. |

Returns: `*StockRunningTradeResponse`

### GetIntradayChart

**Stock/Index Intraday Chart** — `GET /analysis/intraday/{code}` — cache: **realtime**

```go
func (s *AnalysisService) GetIntradayChart(ctx context.Context, code string, params IntradayChartParams) ([]IntradayChartItem, error)
```

Get intraday chart data for a stock with the latest date.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `code (path)` | `string` | yes |  | Kode saham emiten yang tercatat di BEI, mis. `BBCA`. Tidak case-sensitive. Daftar kode tersedia di `/analysis/list/stock`. |
| `params.Market` | `IntradayChartMarket` | yes |  | Jenis pasar: `RG` = Pasar Reguler, `NG` = Pasar Negosiasi, `TN` = Pasar Tunai. Gunakan `RG` untuk analisis pergerakan harga pada umumnya. One of: RG, NG, TN. |

Returns: `[]IntradayChartItem`

### GetOrderBook

**Stock Order Book** — `GET /analysis/order-book/{code}` — cache: **realtime**

```go
func (s *AnalysisService) GetOrderBook(ctx context.Context, code string, params OrderBookParams) (*OrderBookResponse, error)
```

Get order book data for a stock with the latest date or historical data at specific date/time.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `code (path)` | `string` | yes |  | Kode saham emiten yang tercatat di BEI, mis. `BBCA`. Tidak case-sensitive. Daftar kode tersedia di `/analysis/list/stock`. |
| `params.Market` | `OrderBookMarket` | yes |  | Jenis pasar: `RG` = Pasar Reguler, `NG` = Pasar Negosiasi, `TN` = Pasar Tunai. Gunakan `RG` untuk analisis pergerakan harga pada umumnya. One of: RG, NG, TN. |
| `params.Date` | `*string` | no |  | Tanggal order book historis (`YYYY-MM-DD`), tersedia mulai `2026-03-25`. Wajib dikirim bersama `time`; kosongkan keduanya untuk data realtime. |
| `params.Time` | `*string` | no |  | Jam snapshot order book dalam WIB, format `HH:mm` pada jam perdagangan (09:00–15:59). Wajib dikirim bersama `date`. |

Returns: `*OrderBookResponse`

### GetOrderQueue

**Stock Order Queue (Tracking)** — `GET /analysis/queue/{code}` — cache: **realtime**

```go
func (s *AnalysisService) GetOrderQueue(ctx context.Context, code string, params OrderQueueParams) ([]OrderQueueItem, error)
```

Get order queue (order tracking) for a stock at a specific price, side, page, and limit.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `code (path)` | `string` | yes |  | Kode saham emiten yang tercatat di BEI, mis. `BBCA`. Tidak case-sensitive. Daftar kode tersedia di `/analysis/list/stock`. |
| `params.Price` | `int64` | yes |  | Level harga order book (Rupiah) yang ingin dilihat detail antreannya. Harus sesuai fraksi harga saham. Range: 1..-. |
| `params.Side` | `OrderQueueSide` | yes |  | Sisi antrean: `BUY` = antrean beli (bid), `SELL` = antrean jual (offer). One of: BUY, SELL, B, S, buy, sell. |
| `params.Page` | `*int64` | no | `0` | Nomor halaman, dimulai dari 0. Default: 0. Range: 0..-. |
| `params.Limit` | `*int64` | no | `50` | Jumlah order per halaman. Default: 50. Range: 1..50. |

Returns: `[]OrderQueueItem`

### GetIntradayData

**Stock Intraday Data** — `GET /analysis/intraday-data/{code}` — cache: **realtime**

```go
func (s *AnalysisService) GetIntradayData(ctx context.Context, code string, params IntradayDataParams) (*IntradayDataResponse, error)
```

Get intraday stock data with the latest date or historical data on a specific date.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `code (path)` | `string` | yes |  | Kode saham emiten yang tercatat di BEI, mis. `BBCA`. Tidak case-sensitive. Daftar kode tersedia di `/analysis/list/stock`. |
| `params.Market` | `IntradayDataMarket` | yes |  | Jenis pasar: `RG` = Pasar Reguler, `NG` = Pasar Negosiasi, `TN` = Pasar Tunai. Gunakan `RG` untuk analisis pergerakan harga pada umumnya. One of: RG, NG, TN. |
| `params.Date` | `*string` | no |  | Tanggal data historis (`YYYY-MM-DD`). Kosongkan untuk data realtime hari ini. |

Returns: `*IntradayDataResponse`

### GetIntradayIndex

**Index Intraday Data** — `GET /analysis/intraday-index/{code}` — cache: **realtime**

```go
func (s *AnalysisService) GetIntradayIndex(ctx context.Context, code string) (*IntradayIndexResponse, error)
```

Get intraday data for an index with the latest date.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `code (path)` | `string` | yes |  | Kode indeks BEI, mis. `COMPOSITE` (IHSG), `LQ45`, atau `IDX30`. Daftar kode tersedia di `/analysis/list/index`. |

Returns: `*IntradayIndexResponse`

### GetAdvanceChart

**Stock Chart** — `GET /analysis/chart/stock/{code}` — cache: **realtime**

```go
func (s *AnalysisService) GetAdvanceChart(ctx context.Context, code string, params *AdvanceChartParams) ([]AdvanceChartItem, error)
```

Get complete stock price chart data with OHLCV data for a specific period.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `code (path)` | `string` | yes |  | Kode saham emiten yang tercatat di BEI, mis. `BBCA`. Tidak case-sensitive. Daftar kode tersedia di `/analysis/list/stock`. |
| `params.From` | `*string` | no |  | Tanggal awal periode (`YYYY-MM-DD`). Harus lebih kecil atau sama dengan `to`, dengan rentang maksimal 730 hari. Opsional; default 1 tahun sebelum hari ini (WIB). |
| `params.To` | `*string` | no |  | Tanggal akhir periode (`YYYY-MM-DD`). Harus lebih besar atau sama dengan `from`. Opsional; default hari ini (WIB). |

Returns: `[]AdvanceChartItem`

### GetIndexChart

**Index Chart** — `GET /analysis/chart/index/{code}` — cache: **realtime**

```go
func (s *AnalysisService) GetIndexChart(ctx context.Context, code string, params IndexChartParams) ([]IndexChartItem, error)
```

Get complete index price chart data with OHLCV data for a specific period.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `code (path)` | `string` | yes |  | Kode indeks BEI, mis. `COMPOSITE` (IHSG), `LQ45`, atau `IDX30`. Daftar kode tersedia di `/analysis/list/index`. |
| `params.From` | `string` | yes |  | Tanggal awal periode (`YYYY-MM-DD`). Harus lebih kecil atau sama dengan `to`, dengan rentang maksimal 730 hari. |
| `params.To` | `string` | yes |  | Tanggal akhir periode (`YYYY-MM-DD`). Harus lebih besar atau sama dengan `from`. |

Returns: `[]IndexChartItem`

### GetMultiTimeChart

**Stock/Index Multi Timeframe Chart** — `GET /analysis/chart/multi-time/{code}` — cache: **realtime**

```go
func (s *AnalysisService) GetMultiTimeChart(ctx context.Context, code string, params MultiTimeChartParams) ([]MultiTimeChartItem, error)
```

Get multi-timeframe chart data for stocks or indexes.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `code (path)` | `string` | yes |  | Kode saham (mis. `BBCA`) atau kode indeks (mis. `COMPOSITE`). Daftar kode tersedia di `/analysis/list/stock` dan `/analysis/list/index`. |
| `params.From` | `string` | yes |  | Tanggal awal periode (`YYYY-MM-DD`). Harus lebih kecil atau sama dengan `to`. Batas data historis mengikuti `timeframe`. |
| `params.To` | `string` | yes |  | Tanggal akhir periode (`YYYY-MM-DD`). Harus lebih besar atau sama dengan `from`. |
| `params.Timeframe` | `*MultiTimeChartTimeframe` | no | `D` | Timeframe candle: `1`, `5`, `15`, `30`, `60` = menit; `D` = harian, `W` = mingguan, `M` = bulanan. Batas data historis: `1` = 3 bulan, `5` = 6 bulan, `15` = 12 bulan, lainnya 24 bulan. Opsional, default `D`; harus huruf kapital. One of: 1, 5, 15, 30, 60, D, W, M. Default: D. |

Returns: `[]MultiTimeChartItem`

### GetIndicatorChart

**Stock Chart Indicator** — `GET /analysis/chart/stock/{indicator}/{code}` — cache: **default**

```go
func (s *AnalysisService) GetIndicatorChart(ctx context.Context, indicator IndicatorChartIndicator, code string, params IndicatorChartParams) ([]IndicatorChartItem, error)
```

Get chart data with technical indicators for stock price analysis in a specific period.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `indicator (path)` | `IndicatorChartIndicator` | yes |  | Indikator yang ditampilkan pada chart: `bdm` = aliran dana bandar, `foreign` = net beli/jual asing, `ritel` = aliran dana ritel, `freq` = frekuensi transaksi, `ratio` = rasio, `lpm` = Liquidity Pressure Model, `corr_f` = Correlation Foreign, `par_f` = Participation Foreign, `intensity` = Intensity, `vol_ratio` = Volume Rotation. Wajib huruf kecil. One of: bdm, foreign, ratio, ritel, freq, lpm, corr_f, par_f, intensity, vol_ratio. |
| `code (path)` | `string` | yes |  | Kode saham emiten yang tercatat di BEI, mis. `BBCA`. Tidak case-sensitive. Daftar kode tersedia di `/analysis/list/stock`. |
| `params.From` | `string` | yes |  | Tanggal awal periode (`YYYY-MM-DD`). Harus lebih kecil atau sama dengan `to`, dengan rentang maksimal 730 hari. |
| `params.To` | `string` | yes |  | Tanggal akhir periode (`YYYY-MM-DD`). Harus lebih besar atau sama dengan `from`. |

Returns: `[]IndicatorChartItem`

### GetShareholderDetail

**Stock Shareholder 5%** — `GET /analysis/shareholder-detail` — cache: **default**

```go
func (s *AnalysisService) GetShareholderDetail(ctx context.Context, params *ShareholderDetailParams) ([]ShareholderDetailItem, error)
```

Get company stock ownership data for a specific period.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `params.Code` | `*string` | no |  | Kode saham yang ingin dilihat daftar pemegang sahamnya, mis. `BBCA`. Wajib diisi jika `name` kosong. |
| `params.Name` | `*string` | no |  | Nama pemegang saham untuk melihat seluruh saham yang dimilikinya. Wajib diisi jika `code` kosong. |

Returns: `[]ShareholderDetailItem`

### GetShareholderOneDetail

**Stock Shareholder 1%** — `GET /analysis/shareholder-detail-one` — cache: **default**

```go
func (s *AnalysisService) GetShareholderOneDetail(ctx context.Context, params *ShareholderOneDetailParams) ([]ShareholderOneDetailItem, error)
```

Get company stock ownership data based on stock code or shareholder name.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `params.Code` | `*string` | no |  | Kode saham yang ingin dilihat daftar pemegang sahamnya, mis. `BBCA`. Wajib diisi jika `name` kosong. |
| `params.Name` | `*string` | no |  | Nama pemegang saham untuk melihat seluruh saham yang dimilikinya. Wajib diisi jika `code` kosong. |

Returns: `[]ShareholderOneDetailItem`

### GetShareholderHigh

**Stock High Concentration** — `GET /analysis/shareholder/high` — cache: **default**

```go
func (s *AnalysisService) GetShareholderHigh(ctx context.Context) ([]ShareholderHighItem, error)
```

Get high concentration stock ownership data for all issuers.

Returns: `[]ShareholderHighItem`

### GetShareholderNumber

**Stock Shareholder Number** — `GET /analysis/shareholder/number/{code}` — cache: **default**

```go
func (s *AnalysisService) GetShareholderNumber(ctx context.Context, code string) ([]ShareholderNumberItem, error)
```

Get the number of investors/shareholders of a company for a specific period.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `code (path)` | `string` | yes |  | Kode saham emiten yang tercatat di BEI, mis. `BBCA`. Tidak case-sensitive. Daftar kode tersedia di `/analysis/list/stock`. |

Returns: `[]ShareholderNumberItem`

### GetShareholderRelation

**Stock Shareholder Relation** — `GET /analysis/shareholder/relation` — cache: **default**

```go
func (s *AnalysisService) GetShareholderRelation(ctx context.Context, params *ShareholderRelationParams) (*ShareholderRelationResponse, error)
```

Get bidirectional shareholder relation graph data based on the latest stock ownership data.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `params.Code` | `*string` | no |  | Kode saham sebagai titik awal graf relasi kepemilikan. Wajib diisi jika `name` kosong. |
| `params.Name` | `*string` | no |  | Nama pemegang saham sebagai titik awal graf relasi kepemilikan. Wajib diisi jika `code` kosong. |
| `params.Depth` | `*int64` | no | `3` | Kedalaman penelusuran relasi (jumlah lompatan). `1` = relasi langsung, `2` = termasuk relasi tingkat kedua, dan seterusnya. Default: 3. Range: 1..4. |
| `params.MaxNodes` | `*int64` | no | `120` | Jumlah node maksimal (saham dan pemegang saham) dalam graf. Nilai lebih besar menghasilkan graf yang lebih lengkap dengan waktu proses lebih lama. Default: 120. Range: 10..300. |
| `params.Neighbors` | `*int64` | no | `25` | Jumlah relasi maksimal yang diambil dari setiap node. Default: 25. Range: 1..50. |
| `params.MinPercentage` | `*float64` | no | `1` | Persentase kepemilikan minimal yang ditampilkan. Mendukung desimal, mis. `5` untuk pemegang saham minimal 5%. Default: 1. Range: 0..100. |

Returns: `*ShareholderRelationResponse`

### GetShareholder

**Stock Shareholder Composition** — `GET /analysis/shareholder/{code}` — cache: **default**

```go
func (s *AnalysisService) GetShareholder(ctx context.Context, code string) ([]ShareholderItem, error)
```

Get company stock ownership composition data based on the largest shareholders on the latest date.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `code (path)` | `string` | yes |  | Kode saham emiten yang tercatat di BEI, mis. `BBCA`. Tidak case-sensitive. Daftar kode tersedia di `/analysis/list/stock`. |

Returns: `[]ShareholderItem`

### GetShareholderKSEI

**Stock Shareholder KSEI** — `GET /analysis/shareholder/ksei/{code}` — cache: **default**

```go
func (s *AnalysisService) GetShareholderKSEI(ctx context.Context, code string, params ShareholderKSEIParams) ([]ShareholderKSEIItem, error)
```

Get shareholder data based on KSEI classification (foreign and domestic investors in various groups) for a specific period.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `code (path)` | `string` | yes |  | Kode saham emiten yang tercatat di BEI, mis. `BBCA`. Tidak case-sensitive. Daftar kode tersedia di `/analysis/list/stock`. |
| `params.Range` | `int64` | yes |  | Jumlah bulan terakhir yang ditampilkan, mis. `6` untuk 6 bulan terakhir. Akun TRIAL dibatasi maksimal 12 bulan. Range: 1..36. |

Returns: `[]ShareholderKSEIItem`

### GetShareholderClassification

**Stock Shareholder Full Classification** — `GET /analysis/shareholder/classification/{code}` — cache: **default**

```go
func (s *AnalysisService) GetShareholderClassification(ctx context.Context, code string, params ShareholderClassificationParams) ([]ShareholderClassificationItem, error)
```

Get complete shareholder data with all investor classifications for a specific period.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `code (path)` | `string` | yes |  | Kode saham emiten yang tercatat di BEI, mis. `BBCA`. Tidak case-sensitive. Daftar kode tersedia di `/analysis/list/stock`. |
| `params.Range` | `int64` | yes |  | Jumlah bulan terakhir yang ditampilkan, mis. `6` untuk 6 bulan terakhir. Akun TRIAL dibatasi maksimal 12 bulan. Range: 1..36. |

Returns: `[]ShareholderClassificationItem`

### GetShareholderClassifyTable

**Stock Shareholder Classification Table** — `GET /analysis/shareholder/classify-table/{code}` — cache: **default**

```go
func (s *AnalysisService) GetShareholderClassifyTable(ctx context.Context, code string) (*ShareholderClassifyTableResponse, error)
```

Get latest shareholder data with investor classifications in table format.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `code (path)` | `string` | yes |  | Kode saham emiten yang tercatat di BEI, mis. `BBCA`. Tidak case-sensitive. Daftar kode tersedia di `/analysis/list/stock`. |

Returns: `*ShareholderClassifyTableResponse`

### GetSummaryStock

**Stock Broker Summary** — `GET /analysis/summary/stock/{code}` — cache: **default**

```go
func (s *AnalysisService) GetSummaryStock(ctx context.Context, code string, params SummaryStockParams) ([]SummaryStockItem, error)
```

Get detailed trading broker activity analysis for a specific stock within a specific time period.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `code (path)` | `string` | yes |  | Kode saham emiten yang tercatat di BEI, mis. `BBCA`. Tidak case-sensitive. Daftar kode tersedia di `/analysis/list/stock`. |
| `params.From` | `string` | yes |  | Tanggal awal periode (`YYYY-MM-DD`). Harus lebih kecil atau sama dengan `to`, dengan rentang maksimal 730 hari. |
| `params.To` | `string` | yes |  | Tanggal akhir periode (`YYYY-MM-DD`). Harus lebih besar atau sama dengan `from`. |
| `params.Investor` | `SummaryStockInvestor` | yes |  | Jenis investor yang dihitung: `all` = seluruh investor, `f` = investor asing, `d` = investor domestik. One of: all, f, d. |
| `params.Market` | `SummaryStockMarket` | yes |  | Jenis pasar: `RG` = Pasar Reguler, `NG` = Pasar Negosiasi, `TN` = Pasar Tunai. Gunakan `RG` untuk analisis pergerakan harga pada umumnya. One of: RG, NG, TN. |

Returns: `[]SummaryStockItem`

### GetSummaryBroker

**Broker Broker Summary** — `GET /analysis/summary/broker/{code}` — cache: **default**

```go
func (s *AnalysisService) GetSummaryBroker(ctx context.Context, code string, params SummaryBrokerParams) ([]SummaryBrokerItem, error)
```

Get detailed trading activity analysis for a specific broker for all stocks within a specific time period.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `code (path)` | `string` | yes |  | Kode broker anggota bursa (2 huruf), mis. `AK`. Dapat berisi beberapa kode yang dipisahkan koma, mis. `AK,BK,YP`. Daftar broker tersedia di `/analysis/list/broker`. Kode yang tidak terdiri dari 2 karakter akan diabaikan. |
| `params.From` | `string` | yes |  | Tanggal awal periode (`YYYY-MM-DD`). Harus lebih kecil atau sama dengan `to`, dengan rentang maksimal 730 hari. |
| `params.To` | `string` | yes |  | Tanggal akhir periode (`YYYY-MM-DD`). Harus lebih besar atau sama dengan `from`. |
| `params.Investor` | `SummaryBrokerInvestor` | yes |  | Jenis investor yang dihitung: `all` = seluruh investor, `f` = investor asing, `d` = investor domestik. One of: all, f, d. |
| `params.Market` | `SummaryBrokerMarket` | yes |  | Jenis pasar: `RG` = Pasar Reguler, `NG` = Pasar Negosiasi, `TN` = Pasar Tunai. Gunakan `RG` untuk analisis pergerakan harga pada umumnya. One of: RG, NG, TN. |

Returns: `[]SummaryBrokerItem`

### GetSectorStalker

**Sector Stalker** — `GET /analysis/stalker/sector` — cache: **realtime**

```go
func (s *AnalysisService) GetSectorStalker(ctx context.Context, params SectorStalkerParams) (*SectorStalkerResponse, error)
```

Get sectoral index or stock price movement data within a specific time range.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `params.From` | `string` | yes |  | Tanggal awal periode (`YYYY-MM-DD`). Harus lebih kecil atau sama dengan `to`, dengan rentang maksimal 730 hari. |
| `params.To` | `string` | yes |  | Tanggal akhir periode (`YYYY-MM-DD`). Harus lebih besar atau sama dengan `from`. |
| `params.Base` | `*SectorStalkerBase` | no | `COMPOSITE` | Indeks acuan. `COMPOSITE` (IHSG) menampilkan rotasi antar sektor, sedangkan `SYARIAH` (saham syariah, dibandingkan dengan indeks ISSI) atau indeks sektoral IDX-IC (`IDXENERGY`, `IDXFINANCE`, dst.) menampilkan rotasi saham anggota kategori tersebut. `ISSI` diterima sebagai alias `SYARIAH`. Wajib huruf kapital. One of: COMPOSITE, SYARIAH, IDXENERGY, IDXBASIC, IDXINDUST, IDXNONCYC, IDXCYCLIC, IDXHEALTH, IDXFINANCE, IDXPROPERT, IDXTECHNO, IDXINFRA, IDXTRANS. Default: COMPOSITE. |
| `params.Limit` | `*int64` | no |  | Jumlah saham maksimal yang ditampilkan. Hanya berlaku jika `base` bukan `COMPOSITE`. Range: 1..100. |
| `params.Filter` | `*string` | no |  | Pencarian saham berdasarkan sebagian kode atau nama, mis. `BB` atau `BANK`. Hanya berlaku jika `base` bukan `COMPOSITE`. |
| `params.FilterColumn` | `*SectorStalkerFilterColumn` | no |  | Kolom yang difilter: `change` = perubahan harga (%), `value` = nilai transaksi, `volume` = volume, `foreign` = net asing, `freq` = frekuensi, `bdm` = aliran dana bandar, `ritel` = aliran dana ritel, `ratio` = rasio, `open`/`high`/`low`/`close` = harga OHLC. Filter aktif hanya jika `filter_column`, `filter_operator`, dan `filter_value` dikirim bersamaan. One of: change, value, volume, foreign, freq, bdm, ritel, ratio, open, high, low, close. |
| `params.FilterOperator` | `*SectorStalkerFilterOperator` | no |  | Operator pembanding untuk filter. Contoh: `filter_column=change&filter_operator=>&filter_value=5` menampilkan saham yang naik lebih dari 5%. One of: <, >, =, >=, <=, !=. |
| `params.FilterValue` | `*float64` | no |  | Nilai pembanding untuk filter. Mendukung bilangan desimal dan negatif, mis. `5` atau `-2.5`. |

Returns: `*SectorStalkerResponse`

### GetSectorRotation

**Sector Rotation Chart** — `GET /analysis/sector/rotation` — cache: **realtime**

```go
func (s *AnalysisService) GetSectorRotation(ctx context.Context, params SectorRotationParams) (*SectorRotationResponse, error)
```

Get sector rotation data for visualizing relative movement of sectoral indices or stocks.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `params.From` | `string` | yes |  | Tanggal awal periode (`YYYY-MM-DD`). Harus lebih kecil atau sama dengan `to`, dengan rentang maksimal 730 hari. |
| `params.To` | `string` | yes |  | Tanggal akhir periode (`YYYY-MM-DD`). Harus lebih besar atau sama dengan `from`. |
| `params.Base` | `*SectorRotationBase` | no | `COMPOSITE` | Indeks acuan. `COMPOSITE` (IHSG) menampilkan rotasi antar sektor, sedangkan `SYARIAH` (saham syariah, dibandingkan dengan indeks ISSI) atau indeks sektoral IDX-IC (`IDXENERGY`, `IDXFINANCE`, dst.) menampilkan rotasi saham anggota kategori tersebut. `ISSI` diterima sebagai alias `SYARIAH`. Wajib huruf kapital. One of: COMPOSITE, SYARIAH, IDXENERGY, IDXBASIC, IDXINDUST, IDXNONCYC, IDXCYCLIC, IDXHEALTH, IDXFINANCE, IDXPROPERT, IDXTECHNO, IDXINFRA, IDXTRANS. Default: COMPOSITE. |
| `params.Length` | `*int64` | no | `10` | Periode smoothing untuk menghitung RS-Ratio dan RS-Momentum. Nilai kecil lebih responsif, nilai besar lebih halus. Default: 10. Range: 5..50. |
| `params.Interval` | `*SectorRotationInterval` | no | `weekly` | Interval data pada Relative Rotation Graph (RRG): `daily` = harian, `weekly` = mingguan. One of: daily, weekly. Default: weekly. |
| `params.Tail` | `*int64` | no | `5` | Jumlah titik jejak (tail) per sektor atau saham untuk menunjukkan arah pergerakan. Default: 5. Range: 1..52. |
| `params.Limit` | `*int64` | no |  | Jumlah saham maksimal yang ditampilkan. Hanya berlaku jika `base` bukan `COMPOSITE`. Range: 1..50. |
| `params.Filter` | `*string` | no |  | Pencarian saham berdasarkan sebagian kode atau nama, mis. `BB` atau `BANK`. Hanya berlaku jika `base` bukan `COMPOSITE`. |
| `params.FilterColumn` | `*SectorRotationFilterColumn` | no |  | Kolom yang difilter: `change` = perubahan harga (%), `value` = nilai transaksi, `volume` = volume, `foreign` = net asing, `freq` = frekuensi, `bdm` = aliran dana bandar, `ritel` = aliran dana ritel, `ratio` = rasio, `open`/`high`/`low`/`close` = harga OHLC. Filter aktif hanya jika `filter_column`, `filter_operator`, dan `filter_value` dikirim bersamaan. One of: change, value, volume, foreign, freq, bdm, ritel, ratio, open, high, low, close. |
| `params.FilterOperator` | `*SectorRotationFilterOperator` | no |  | Operator pembanding untuk filter. Contoh: `filter_column=change&filter_operator=>&filter_value=5` menampilkan saham yang naik lebih dari 5%. One of: <, >, =, >=, <=, !=. |
| `params.FilterValue` | `*float64` | no |  | Nilai pembanding untuk filter. Mendukung bilangan desimal dan negatif, mis. `5` atau `-2.5`. |

Returns: `*SectorRotationResponse`

### GetBrokerStalker

**Stock Broker Stalker** — `GET /analysis/stalker/broker/{broker}/{stock}` — cache: **default**

```go
func (s *AnalysisService) GetBrokerStalker(ctx context.Context, broker string, stock string, params BrokerStalkerParams) (*BrokerStalkerResponse, error)
```

Track a specific broker's trading activity on a specific stock in daily calendar format.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `broker (path)` | `string` | yes |  | Kode broker anggota bursa (2 huruf), mis. `AK`. Dapat berisi beberapa kode yang dipisahkan koma, mis. `AK,BK,YP`. Daftar broker tersedia di `/analysis/list/broker`. Kode yang tidak terdiri dari 2 karakter akan diabaikan. |
| `stock (path)` | `string` | yes |  | Kode saham emiten yang tercatat di BEI, mis. `BBCA`. Tidak case-sensitive. Daftar kode tersedia di `/analysis/list/stock`. |
| `params.From` | `string` | yes |  | Tanggal awal periode (`YYYY-MM-DD`). Harus lebih kecil atau sama dengan `to`, dengan rentang maksimal 730 hari. |
| `params.To` | `string` | yes |  | Tanggal akhir periode (`YYYY-MM-DD`). Harus lebih besar atau sama dengan `from`. |
| `params.Investor` | `BrokerStalkerInvestor` | yes |  | Jenis investor yang dihitung: `all` = seluruh investor, `f` = investor asing, `d` = investor domestik. One of: all, f, d. |
| `params.Market` | `BrokerStalkerMarket` | yes |  | Jenis pasar: `RG` = Pasar Reguler, `NG` = Pasar Negosiasi, `TN` = Pasar Tunai. Gunakan `RG` untuk analisis pergerakan harga pada umumnya. One of: RG, NG, TN. |

Returns: `*BrokerStalkerResponse`

### GetBrokerStalkerList

**Broker Stalker List** — `GET /analysis/stalker/list/{code}` — cache: **default**

```go
func (s *AnalysisService) GetBrokerStalkerList(ctx context.Context, code string, params BrokerStalkerListParams) (*BrokerStalkerListResponse, error)
```

Get list of all stocks traded by a specific broker within a specific time period.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `code (path)` | `string` | yes |  | Kode broker anggota bursa (2 huruf), mis. `AK`. Dapat berisi beberapa kode yang dipisahkan koma, mis. `AK,BK,YP`. Daftar broker tersedia di `/analysis/list/broker`. Kode yang tidak terdiri dari 2 karakter akan diabaikan. |
| `params.From` | `string` | yes |  | Tanggal awal periode (`YYYY-MM-DD`). Harus lebih kecil atau sama dengan `to`, dengan rentang maksimal 730 hari. |
| `params.To` | `string` | yes |  | Tanggal akhir periode (`YYYY-MM-DD`). Harus lebih besar atau sama dengan `from`. |
| `params.Investor` | `BrokerStalkerListInvestor` | yes |  | Jenis investor yang dihitung: `all` = seluruh investor, `f` = investor asing, `d` = investor domestik. One of: all, f, d. |
| `params.Market` | `BrokerStalkerListMarket` | yes |  | Jenis pasar: `RG` = Pasar Reguler, `NG` = Pasar Negosiasi, `TN` = Pasar Tunai. Gunakan `RG` untuk analisis pergerakan harga pada umumnya. One of: RG, NG, TN. |

Returns: `*BrokerStalkerListResponse`

### GetSummaryStockChart

**Stock Broker Summary Chart** — `GET /analysis/summary-chart/stock/{code}` — cache: **default**

```go
func (s *AnalysisService) GetSummaryStockChart(ctx context.Context, code string, params SummaryStockChartParams) ([]SummaryStockChartItem, error)
```

Get broker summary infographic data for a specific stock with attractive visualization for a specific period.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `code (path)` | `string` | yes |  | Kode saham emiten yang tercatat di BEI, mis. `BBCA`. Tidak case-sensitive. Daftar kode tersedia di `/analysis/list/stock`. |
| `params.From` | `string` | yes |  | Tanggal awal periode (`YYYY-MM-DD`). Harus lebih kecil atau sama dengan `to`, dengan rentang maksimal 730 hari. |
| `params.To` | `string` | yes |  | Tanggal akhir periode (`YYYY-MM-DD`). Harus lebih besar atau sama dengan `from`. |
| `params.Scope` | `SummaryStockChartScope` | yes |  | Metrik yang dihitung: `volume` = jumlah lembar saham, `value` = nilai transaksi (Rupiah), `freq` = frekuensi transaksi. One of: volume, value, freq. |
| `params.Market` | `SummaryStockChartMarket` | yes |  | Jenis pasar: `RG` = Pasar Reguler, `NG` = Pasar Negosiasi, `TN` = Pasar Tunai. Gunakan `RG` untuk analisis pergerakan harga pada umumnya. One of: RG, NG, TN. |

Returns: `[]SummaryStockChartItem`

### GetSummaryBrokerChart

**Broker Broker Summary Chart** — `GET /analysis/summary-chart/broker/{code}` — cache: **default**

```go
func (s *AnalysisService) GetSummaryBrokerChart(ctx context.Context, code string, params SummaryBrokerChartParams) ([]SummaryBrokerChartItem, error)
```

Get broker summary infographic data for a specific broker with attractive visualization for a specific period.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `code (path)` | `string` | yes |  | Kode broker anggota bursa (2 huruf), mis. `AK`. Dapat berisi beberapa kode yang dipisahkan koma, mis. `AK,BK,YP`. Daftar broker tersedia di `/analysis/list/broker`. Kode yang tidak terdiri dari 2 karakter akan diabaikan. |
| `params.From` | `string` | yes |  | Tanggal awal periode (`YYYY-MM-DD`). Harus lebih kecil atau sama dengan `to`, dengan rentang maksimal 730 hari. |
| `params.To` | `string` | yes |  | Tanggal akhir periode (`YYYY-MM-DD`). Harus lebih besar atau sama dengan `from`. |
| `params.Scope` | `SummaryBrokerChartScope` | yes |  | Metrik yang dihitung: `volume` = jumlah lembar saham, `value` = nilai transaksi (Rupiah), `freq` = frekuensi transaksi. One of: volume, value, freq. |
| `params.Market` | `SummaryBrokerChartMarket` | yes |  | Jenis pasar: `RG` = Pasar Reguler, `NG` = Pasar Negosiasi, `TN` = Pasar Tunai. Gunakan `RG` untuk analisis pergerakan harga pada umumnya. One of: RG, NG, TN. |

Returns: `[]SummaryBrokerChartItem`

### GetInventoryStockChart

**Stock Inventory Chart** — `GET /analysis/inventory-chart/stock/{code}` — cache: **default**

```go
func (s *AnalysisService) GetInventoryStockChart(ctx context.Context, code string, params InventoryStockChartParams) (*InventoryStockChartResponse, error)
```

Get inventory chart data for a specific stock with broker activity visualization for a specific period.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `code (path)` | `string` | yes |  | Kode saham emiten yang tercatat di BEI, mis. `BBCA`. Tidak case-sensitive. Daftar kode tersedia di `/analysis/list/stock`. |
| `params.From` | `string` | yes |  | Tanggal awal periode (`YYYY-MM-DD`). Harus lebih kecil atau sama dengan `to`, dengan rentang maksimal 730 hari. |
| `params.To` | `string` | yes |  | Tanggal akhir periode (`YYYY-MM-DD`). Harus lebih besar atau sama dengan `from`. |
| `params.Scope` | `InventoryStockChartScope` | yes |  | Metrik inventori: `vol` = volume (lembar saham), `val` = nilai transaksi (Rupiah), `freq` = frekuensi transaksi. One of: vol, val, freq. |
| `params.Investor` | `InventoryStockChartInvestor` | yes |  | Jenis investor yang dihitung: `all` = seluruh investor, `f` = investor asing, `d` = investor domestik. One of: all, f, d. |
| `params.Limit` | `*int64` | no | `20` | Jumlah broker teratas yang ditampilkan. Default: 20. Range: 1..50. |
| `params.Market` | `InventoryStockChartMarket` | yes |  | Jenis pasar: `ALL` = seluruh pasar, `RG` = Pasar Reguler, `NG` = Pasar Negosiasi, `TN` = Pasar Tunai. One of: ALL, RG, NG, TN. |
| `params.Filter` | `*string` | no |  | Kode broker yang ingin ditampilkan, dipisahkan koma, mis. `AK,BK`. Kosongkan untuk menampilkan broker teratas. |
| `params.FilterOperator` | `*InventoryStockChartFilterOperator` | no |  | Saat ini tidak digunakan pada endpoint ini. One of: <, >, =, >=, <=, !=. |
| `params.FilterValue` | `*float64` | no |  | Saat ini tidak digunakan pada endpoint ini. |

Returns: `*InventoryStockChartResponse`

### GetInventoryBrokerChart

**Broker Inventory Chart** — `GET /analysis/inventory-chart/broker/{code}` — cache: **default**

```go
func (s *AnalysisService) GetInventoryBrokerChart(ctx context.Context, code string, params InventoryBrokerChartParams) ([]InventoryBrokerChartItem, error)
```

Get inventory chart data for a specific broker with stock activity visualization for a specific period.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `code (path)` | `string` | yes |  | Kode broker anggota bursa (2 huruf), mis. `AK`. Dapat berisi beberapa kode yang dipisahkan koma, mis. `AK,BK,YP`. Daftar broker tersedia di `/analysis/list/broker`. Setiap kode wajib 2 karakter. |
| `params.From` | `string` | yes |  | Tanggal awal periode (`YYYY-MM-DD`). Harus lebih kecil atau sama dengan `to`, dengan rentang maksimal 730 hari. |
| `params.To` | `string` | yes |  | Tanggal akhir periode (`YYYY-MM-DD`). Harus lebih besar atau sama dengan `from`. |
| `params.Scope` | `InventoryBrokerChartScope` | yes |  | Metrik inventori: `vol` = volume (lembar saham), `val` = nilai transaksi (Rupiah), `freq` = frekuensi transaksi. One of: vol, val, freq. |
| `params.Investor` | `InventoryBrokerChartInvestor` | yes |  | Jenis investor yang dihitung: `all` = seluruh investor, `f` = investor asing, `d` = investor domestik. One of: all, f, d. |
| `params.Limit` | `*int64` | no | `20` | Jumlah saham teratas yang ditampilkan. Default: 20. Range: 1..50. |
| `params.Market` | `InventoryBrokerChartMarket` | yes |  | Jenis pasar: `ALL` = seluruh pasar, `RG` = Pasar Reguler, `NG` = Pasar Negosiasi, `TN` = Pasar Tunai. One of: ALL, RG, NG, TN. |
| `params.Filter` | `*string` | no |  | Kode saham yang ingin ditampilkan, dipisahkan koma, mis. `BBCA,BBRI`. Kosongkan untuk menampilkan saham teratas. |
| `params.FilterOperator` | `*InventoryBrokerChartFilterOperator` | no |  | Operator pembanding untuk memfilter nilai `scope`. Digunakan bersama `filter_value`. One of: <, >, =, >=, <=, !=. |
| `params.FilterValue` | `*float64` | no |  | Nilai pembanding untuk `filter_operator`. Contoh: `>` dengan nilai `0` menampilkan saham yang terakumulasi (net beli). |

Returns: `[]InventoryBrokerChartItem`

### GetMomentumChart

**Stock Momentum Chart** — `GET /analysis/momentum-chart/{code}` — cache: **realtime**

```go
func (s *AnalysisService) GetMomentumChart(ctx context.Context, code string, params MomentumChartParams) ([]MomentumChartItem, error)
```

Get momentum chart data for real-time transaction analysis with HAKA HAKI details for a specific period.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `code (path)` | `string` | yes |  | Kode saham emiten yang tercatat di BEI, mis. `BBCA`. Tidak case-sensitive. Daftar kode tersedia di `/analysis/list/stock`. |
| `params.Date` | `string` | yes |  | Tanggal data (`YYYY-MM-DD`). Gunakan hari bursa; tanggal libur bursa dapat menghasilkan data kosong. |
| `params.Range` | `int64` | yes |  | Interval waktu dalam menit, mis. `5` untuk per 5 menit atau `60` untuk per jam. Range: 1..365. |
| `params.Scope` | `MomentumChartScope` | yes |  | Basis perhitungan: `value` = nilai transaksi (Rupiah), `volume` = jumlah lembar saham. One of: value, volume. |

Returns: `[]MomentumChartItem`

### GetIntradayInventoryChart

**Stock Intraday Inventory Chart** — `GET /analysis/intraday-inventory-chart/{code}` — cache: **default**

```go
func (s *AnalysisService) GetIntradayInventoryChart(ctx context.Context, code string, params IntradayInventoryChartParams) (*IntradayInventoryChartResponse, error)
```

Get intraday inventory chart data for transaction analysis with broker activity details for a specific period.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `code (path)` | `string` | yes |  | Kode saham emiten yang tercatat di BEI, mis. `BBCA`. Tidak case-sensitive. Daftar kode tersedia di `/analysis/list/stock`. |
| `params.Date` | `string` | yes |  | Tanggal data (`YYYY-MM-DD`). Gunakan hari bursa; tanggal libur bursa dapat menghasilkan data kosong. |
| `params.Range` | `int64` | yes |  | Interval waktu dalam menit, mis. `5` untuk per 5 menit atau `60` untuk per jam. Range: 1..365. |
| `params.Type` | `IntradayInventoryChartType` | yes |  | Basis perhitungan: `value` = nilai transaksi (Rupiah), `volume` = jumlah lembar saham. One of: value, volume. |
| `params.Total` | `int64` | yes |  | Jumlah broker teratas yang ditampilkan. Range: 1..100. |
| `params.Buyer` | `IntradayInventoryChartBuyer` | yes |  | Jenis investor di sisi pembeli: `ALL` = seluruh investor, `F` = asing, `D` = domestik. One of: ALL, F, D. |
| `params.Seller` | `IntradayInventoryChartSeller` | yes |  | Jenis investor di sisi penjual: `ALL` = seluruh investor, `F` = asing, `D` = domestik. One of: ALL, F, D. |
| `params.Filter` | `*string` | no |  | Kode broker yang ingin ditampilkan, dipisahkan koma, mis. `AK,BK`. Kosongkan untuk menampilkan broker teratas. |
| `params.Market` | `IntradayInventoryChartMarket` | yes |  | Jenis pasar: `RG` = Pasar Reguler, `NG` = Pasar Negosiasi, `TN` = Pasar Tunai. Gunakan `RG` untuk analisis pergerakan harga pada umumnya. One of: RG, NG, TN. |

Returns: `*IntradayInventoryChartResponse`

### GetSankeyChart

**Stock Distribution Chart** — `GET /analysis/sankey-chart/{code}` — cache: **default**

```go
func (s *AnalysisService) GetSankeyChart(ctx context.Context, code string, params SankeyChartParams) (*SankeyChartResponse, error)
```

Get sankey chart data for crossing transaction visualization with broker flow details for a specific period.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `code (path)` | `string` | yes |  | Kode saham emiten yang tercatat di BEI, mis. `BBCA`. Tidak case-sensitive. Daftar kode tersedia di `/analysis/list/stock`. |
| `params.Date` | `string` | yes |  | Tanggal data (`YYYY-MM-DD`). Gunakan hari bursa; tanggal libur bursa dapat menghasilkan data kosong. |
| `params.Type` | `SankeyChartType` | yes |  | Basis perhitungan: `value` = nilai transaksi (Rupiah), `volume` = jumlah lembar saham. One of: value, volume. |
| `params.Buyer` | `SankeyChartBuyer` | yes |  | Jenis investor di sisi pembeli: `ALL` = seluruh investor, `F` = asing, `D` = domestik. One of: ALL, F, D. |
| `params.Seller` | `SankeyChartSeller` | yes |  | Jenis investor di sisi penjual: `ALL` = seluruh investor, `F` = asing, `D` = domestik. One of: ALL, F, D. |
| `params.Market` | `SankeyChartMarket` | yes |  | Jenis pasar: `RG` = Pasar Reguler, `NG` = Pasar Negosiasi, `TN` = Pasar Tunai. Gunakan `RG` untuk analisis pergerakan harga pada umumnya. One of: RG, NG, TN. |

Returns: `*SankeyChartResponse`

### GetPriceTable

**Stock Price Table** — `GET /analysis/price-table/{code}` — cache: **realtime**

```go
func (s *AnalysisService) GetPriceTable(ctx context.Context, code string, params PriceTableParams) ([]PriceTableItem, error)
```

Get daily stock trading price transaction table with volume and frequency details for a specific period.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `code (path)` | `string` | yes |  | Kode saham emiten yang tercatat di BEI, mis. `BBCA`. Tidak case-sensitive. Daftar kode tersedia di `/analysis/list/stock`. |
| `params.Date` | `string` | yes |  | Tanggal data (`YYYY-MM-DD`). Gunakan hari bursa; tanggal libur bursa dapat menghasilkan data kosong. |

Returns: `[]PriceTableItem`

### GetTimeTable

**Stock Time Table** — `GET /analysis/time-table/{code}` — cache: **realtime**

```go
func (s *AnalysisService) GetTimeTable(ctx context.Context, code string, params TimeTableParams) ([]TimeTableItem, error)
```

Get daily stock trading time transaction table with HAKA/HAKI volume (lots) and percentage details per time interval for a specific period.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `code (path)` | `string` | yes |  | Kode saham emiten yang tercatat di BEI, mis. `BBCA`. Tidak case-sensitive. Daftar kode tersedia di `/analysis/list/stock`. |
| `params.Date` | `string` | yes |  | Tanggal data (`YYYY-MM-DD`). Gunakan hari bursa; tanggal libur bursa dapat menghasilkan data kosong. |
| `params.Range` | `int64` | yes |  | Interval waktu dalam menit, mis. `5` untuk per 5 menit atau `60` untuk per jam. Range: 1..365. |

Returns: `[]TimeTableItem`

### GetPriceDiary

**Stock Price Diary** — `GET /analysis/price-diary/{code}` — cache: **default**

```go
func (s *AnalysisService) GetPriceDiary(ctx context.Context, code string) ([]PriceDiaryItem, error)
```

Get daily stock price change table for the last 30 trading days (no period parameter).

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `code (path)` | `string` | yes |  | Kode saham emiten yang tercatat di BEI, mis. `BBCA`. Tidak case-sensitive. Daftar kode tersedia di `/analysis/list/stock`. |

Returns: `[]PriceDiaryItem`

### GetPriceSeasonality

**Stock Price Seasonality** — `GET /analysis/price-seasonality/{code}` — cache: **default**

```go
func (s *AnalysisService) GetPriceSeasonality(ctx context.Context, code string, params PriceSeasonalityParams) ([]PriceSeasonalityItem, error)
```

Get monthly stock price change table with seasonality analysis for a specific period.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `code (path)` | `string` | yes |  | Kode saham emiten yang tercatat di BEI, mis. `BBCA`. Tidak case-sensitive. Daftar kode tersedia di `/analysis/list/stock`. |
| `params.Range` | `int64` | yes |  | Jumlah tahun data historis untuk menghitung pola musiman harga bulanan, mis. `3` untuk 3 tahun terakhir. Nilai efektif maksimal 20. Range: 1..120. |

Returns: `[]PriceSeasonalityItem`

### GetShareholderAbove

**Stock Insider KSEI** — `GET /analysis/shareholder-above` — cache: **default**

```go
func (s *AnalysisService) GetShareholderAbove(ctx context.Context, params *ShareholderAboveParams) (*ShareholderAboveResponse, error)
```

Get a list of ownership changes above 5% based on exchange reports for a specific period.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `params.Code` | `*string` | no |  | Filter berdasarkan kode saham. Dapat berisi beberapa kode yang dipisahkan koma, mis. `BBCA,BBRI,BMRI`. |
| `params.Broker` | `*string` | no |  | Filter berdasarkan broker/kustodian tempat saham disimpan (kode 2 huruf). Dapat berisi beberapa kode yang dipisahkan koma. |
| `params.Name` | `*string` | no |  | Nama atau sebagian nama pemegang saham yang dicari. Kosongkan untuk menampilkan seluruh pemegang saham. |
| `params.From` | `*string` | no |  | Tanggal awal periode laporan (`YYYY-MM-DD`). Opsional; filter tanggal hanya berlaku bila `from` dan `to` sama-sama diisi. |
| `params.To` | `*string` | no |  | Tanggal akhir periode laporan (`YYYY-MM-DD`). Opsional; filter tanggal hanya berlaku bila `from` dan `to` sama-sama diisi. |
| `params.Page` | `*int64` | no | `1` | Nomor halaman, dimulai dari 1. Default: 1. Range: 1..-. |
| `params.Limit` | `*int64` | no | `20` | Jumlah transaksi per halaman. Default: 20. Range: 1..50. |

Returns: `*ShareholderAboveResponse`

### GetShareholderAboveChart

**Stock Insider KSEI Chart** — `GET /analysis/shareholder-above-chart/{code}` — cache: **default**

```go
func (s *AnalysisService) GetShareholderAboveChart(ctx context.Context, code string, params ShareholderAboveChartParams) ([]ShareholderAboveChartItem, error)
```

Get chart data for ownership visualization above 5% for a specific period.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `code (path)` | `string` | yes |  | Kode saham emiten yang tercatat di BEI, mis. `BBCA`. Tidak case-sensitive. Daftar kode tersedia di `/analysis/list/stock`. |
| `params.Broker` | `*string` | no |  | Filter berdasarkan broker/kustodian (kode 2 huruf). Gunakan `INSTITUSI` untuk menampilkan kepemilikan institusi. |
| `params.Name` | `*string` | no |  | Nama pemegang saham dengan kepemilikan di atas 5% yang ingin ditampilkan. |
| `params.Date` | `string` | yes |  | Tanggal data kepemilikan (`YYYY-MM-DD`). Wajib diisi. |

Returns: `[]ShareholderAboveChartItem`

### GetShareholderOne

**Stock Insider 1%** — `GET /analysis/shareholder-one` — cache: **default**

```go
func (s *AnalysisService) GetShareholderOne(ctx context.Context, params *ShareholderOneParams) (*ShareholderOneResponse, error)
```

Get a list of ownership changes above 1% based on exchange reports for a specific period.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `params.Code` | `*string` | no |  | Filter berdasarkan kode saham. Dapat berisi beberapa kode yang dipisahkan koma, mis. `BBCA,BBRI,BMRI`. |
| `params.Name` | `*string` | no |  | Nama atau sebagian nama pemegang saham yang dicari. Kosongkan untuk menampilkan seluruh pemegang saham. |
| `params.From` | `*string` | no |  | Tanggal awal periode laporan (`YYYY-MM-DD`). Opsional; filter tanggal hanya berlaku bila `from` dan `to` sama-sama diisi. |
| `params.To` | `*string` | no |  | Tanggal akhir periode laporan (`YYYY-MM-DD`). Opsional; filter tanggal hanya berlaku bila `from` dan `to` sama-sama diisi. |
| `params.Page` | `*int64` | no | `1` | Nomor halaman, dimulai dari 1. Default: 1. Range: 1..-. |
| `params.Limit` | `*int64` | no | `20` | Jumlah transaksi per halaman. Default: 20. Range: 1..50. |

Returns: `*ShareholderOneResponse`

### GetShareholderOneChart

**Stock Insider 1% Chart** — `GET /analysis/shareholder-one-chart/{code}` — cache: **default**

```go
func (s *AnalysisService) GetShareholderOneChart(ctx context.Context, code string, params ShareholderOneChartParams) ([]ShareholderOneChartItem, error)
```

Get chart data for ownership visualization above 1% for a specific period.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `code (path)` | `string` | yes |  | Kode saham emiten yang tercatat di BEI, mis. `BBCA`. Tidak case-sensitive. Daftar kode tersedia di `/analysis/list/stock`. |
| `params.Name` | `*string` | no |  | Nama pemegang saham dengan kepemilikan 1%–5% yang ingin ditampilkan. |
| `params.Date` | `string` | yes |  | Tanggal data kepemilikan (`YYYY-MM-DD`). Wajib diisi. |

Returns: `[]ShareholderOneChartItem`

### GetInsider

**Stock Insider IDX** — `GET /analysis/shareholder-insider` — cache: **default**

```go
func (s *AnalysisService) GetInsider(ctx context.Context, params *InsiderParams) (*InsiderResponse, error)
```

Get a list of insider ownership changes (directors, commissioners, majority shareholders) for a specific period.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `params.Code` | `*string` | no |  | Filter berdasarkan kode saham. Dapat berisi beberapa kode yang dipisahkan koma, mis. `BBCA,BBRI,BMRI`. |
| `params.Name` | `*string` | no |  | Nama atau sebagian nama pemegang saham yang dicari. Kosongkan untuk menampilkan seluruh pemegang saham. |
| `params.From` | `*string` | no |  | Tanggal awal periode laporan (`YYYY-MM-DD`). Opsional; filter tanggal hanya berlaku bila `from` dan `to` sama-sama diisi. |
| `params.To` | `*string` | no |  | Tanggal akhir periode laporan (`YYYY-MM-DD`). Opsional; filter tanggal hanya berlaku bila `from` dan `to` sama-sama diisi. |
| `params.Page` | `*int64` | no | `1` | Nomor halaman, dimulai dari 1. Default: 1. Range: 1..-. |
| `params.Limit` | `*int64` | no | `20` | Jumlah transaksi per halaman. Default: 20. Range: 1..50. |

Returns: `*InsiderResponse`

### GetInsiderChart

**Stock Insider IDX Chart** — `GET /analysis/insider-chart/{code}` — cache: **default**

```go
func (s *AnalysisService) GetInsiderChart(ctx context.Context, code string, params InsiderChartParams) ([]InsiderChartItem, error)
```

Get chart data for insider trading visualization for a specific period.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `code (path)` | `string` | yes |  | Kode saham emiten yang tercatat di BEI, mis. `BBCA`. Tidak case-sensitive. Daftar kode tersedia di `/analysis/list/stock`. |
| `params.Name` | `*string` | no |  | Nama insider (direksi, komisaris, atau pemegang saham pengendali) yang ingin ditampilkan. |
| `params.Date` | `string` | yes |  | Tanggal laporan transaksi insider di BEI (`YYYY-MM-DD`). Wajib diisi. |

Returns: `[]InsiderChartItem`

### GetFinancialStatement

**Stock Financial Statement** — `GET /analysis/financial-statement/{code}` — cache: **default**

```go
func (s *AnalysisService) GetFinancialStatement(ctx context.Context, code string, params FinancialStatementParams) (*FinancialStatementResponse, error)
```

Get company financial statement data based on report type and specific period.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `code (path)` | `string` | yes |  | Kode saham emiten yang tercatat di BEI, mis. `BBCA`. Tidak case-sensitive. Daftar kode tersedia di `/analysis/list/stock`. |
| `params.Statement` | `FinancialStatementStatement` | yes |  | Jenis laporan keuangan: `BS` = neraca (balance sheet), `IS` = laba rugi (income statement), `CF` = arus kas (cash flow). One of: BS, IS, CF. |
| `params.Type` | `FinancialStatementType` | yes |  | Periode laporan: `FY` = tahunan, `Q` = seluruh kuartal, `Q1`–`Q4` = kuartal tertentu. One of: FY, Q, Q1, Q2, Q3, Q4. |
| `params.Limit` | `*int64` | no | `8` | Jumlah periode laporan terakhir yang ditampilkan. Default: 8. Range: 1..8. |

Returns: `*FinancialStatementResponse`

### GetFinancialStatementChart

**Stock Financial Statement Chart** — `GET /analysis/financial-statement-chart/{code}` — cache: **default**

```go
func (s *AnalysisService) GetFinancialStatementChart(ctx context.Context, code string, params FinancialStatementChartParams) (*FinancialStatementChartResponse, error)
```

Get chart data for financial statement visualization for a specific period.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `code (path)` | `string` | yes |  | Kode saham emiten yang tercatat di BEI, mis. `BBCA`. Tidak case-sensitive. Daftar kode tersedia di `/analysis/list/stock`. |
| `params.Statement` | `FinancialStatementChartStatement` | yes |  | Jenis laporan keuangan: `BS` = neraca (balance sheet), `IS` = laba rugi (income statement), `CF` = arus kas (cash flow). One of: BS, IS, CF. |
| `params.Type` | `FinancialStatementChartType` | yes |  | Periode laporan: `FY` = tahunan, `Q` = seluruh kuartal, `Q1`–`Q4` = kuartal tertentu. One of: FY, Q, Q1, Q2, Q3, Q4. |
| `params.Limit` | `*int64` | no | `8` | Jumlah periode laporan terakhir yang ditampilkan. Default: 8. Range: 1..8. |
| `params.Account` | `string` | yes |  | ID akun laporan keuangan (string 32 karakter) yang akan ditampilkan sebagai grafik. Gunakan nilai `rows[].id` dari response `/analysis/financial-statement/{code}` dengan `statement` yang sama; ID yang tidak ditemukan menghasilkan `account_id`/`name` null dengan `amount` 0. |

Returns: `*FinancialStatementChartResponse`

### GetKeystat

**Stock Key Statistics** — `GET /analysis/keystat/{code}` — cache: **default**

```go
func (s *AnalysisService) GetKeystat(ctx context.Context, code string, params KeystatParams) (*KeystatResponse, error)
```

Get key statistics and important financial metrics of a company for a specific period.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `code (path)` | `string` | yes |  | Kode saham emiten yang tercatat di BEI, mis. `BBCA`. Tidak case-sensitive. Daftar kode tersedia di `/analysis/list/stock`. |
| `params.Type` | `KeystatType` | yes |  | Periode laporan: `FY` = tahunan, `Q` = seluruh kuartal, `Q1`–`Q4` = kuartal tertentu. One of: FY, Q, Q1, Q2, Q3, Q4. |
| `params.Limit` | `*int64` | no | `8` | Jumlah periode terakhir yang ditampilkan. Default: 8. Range: 1..8. |

Returns: `*KeystatResponse`

### GetKeystatChart

**Stock Key Statistics Chart** — `GET /analysis/keystat-chart/{code}` — cache: **default**

```go
func (s *AnalysisService) GetKeystatChart(ctx context.Context, code string, params KeystatChartParams) (*KeystatChartResponse, error)
```

Get chart data for key statistics visualization for a specific period.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `code (path)` | `string` | yes |  | Kode saham emiten yang tercatat di BEI, mis. `BBCA`. Tidak case-sensitive. Daftar kode tersedia di `/analysis/list/stock`. |
| `params.Type` | `KeystatChartType` | yes |  | Periode laporan: `FY` = tahunan, `Q` = seluruh kuartal, `Q1`–`Q4` = kuartal tertentu. One of: FY, Q, Q1, Q2, Q3, Q4. |
| `params.Limit` | `*int64` | no | `8` | Jumlah periode terakhir yang ditampilkan (default 8). Maksimal 8 untuk periode kuartalan (`Q`, `Q1`-`Q4`) dan 3 untuk `type=FY`; karena default 8 melebihi batas FY, `limit` (1-3) wajib diisi saat `type=FY`. Default: 8. Range: 1..8. |
| `params.Name` | `string` | yes |  | Nama rasio atau statistik kunci yang akan ditampilkan sebagai grafik, mis. `ROE`, `PER`, atau `PBV`. Gunakan nama yang sama persis dengan response `/analysis/keystat/{code}`. |

Returns: `*KeystatChartResponse`

### GetCalendar

**Stock Corporate Action Calendar** — `GET /analysis/calendar` — cache: **default**

```go
func (s *AnalysisService) GetCalendar(ctx context.Context, params *CalendarParams) (*CalendarResponse, error)
```

Get the corporate action schedule of IDX issuers, newest first, with optional filters.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `params.Code` | `*string` | no |  | Filter hasil untuk satu kode saham tertentu. |
| `params.Type` | `*CalendarType` | no |  | Jenis aksi korporasi: `IPO`, `PUBLIC_EXPOSE` = paparan publik, `RUPS_SCHEDULE` / `RUPS_RESULT` = jadwal / hasil RUPS, `DIVIDEND`, `RIGHT` = rights issue, `WARRANT`, `BONUS`, `SPLIT` = stock split, `REVERSE` = reverse split, `CONVERTION` = konversi. One of: IPO, PUBLIC_EXPOSE, REVERSE, RIGHT, RUPS_RESULT, RUPS_SCHEDULE, SPLIT, WARRANT, BONUS, CONVERTION, DIVIDEND. |
| `params.Page` | `*int64` | no | `1` | Nomor halaman, dimulai dari 1. Default: 1. Range: 1..-. |
| `params.Limit` | `*int64` | no | `20` | Jumlah data per halaman. Default: 20. Range: 1..50. |

Returns: `*CalendarResponse`

### GetNews

**Stock News** — `GET /analysis/news` — cache: **default**

```go
func (s *AnalysisService) GetNews(ctx context.Context, params *NewsParams) (*NewsResponse, error)
```

Get the stock news listing (NEWS-category posts).

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `params.Page` | `*int64` | no | `1` | Nomor halaman, dimulai dari 1. Default: 1. Range: 1..-. |
| `params.Limit` | `*int64` | no | `20` | Jumlah data per halaman, 1-100. Default: 20. Range: 1..100. |
| `params.Code` | `*string` | no |  | Filter kode saham, mis. `BBCA` (opsional). |

Returns: `*NewsResponse`

### GetDisclosure

**Stock Disclosure Information** — `GET /analysis/disclosure` — cache: **default**

```go
func (s *AnalysisService) GetDisclosure(ctx context.Context, params *DisclosureParams) (*DisclosureResponse, error)
```

Get the IDX disclosure listing (REPORT-category posts).

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `params.Page` | `*int64` | no | `1` | Nomor halaman, dimulai dari 1. Default: 1. Range: 1..-. |
| `params.Limit` | `*int64` | no | `20` | Jumlah data per halaman, 1-100. Default: 20. Range: 1..100. |
| `params.Code` | `*string` | no |  | Filter kode saham, mis. `BBCA` (opsional). |

Returns: `*DisclosureResponse`


## client.Batch

Batch endpoints: fetch many stocks in one request (counted as 1 request).

### GetOrderBook

**Batch Stock Order Book** — `GET /batch/order-book/{code}` — cache: **realtime**

```go
func (s *BatchService) GetOrderBook(ctx context.Context, code string, params *BatchOrderBookParams) ([]BatchOrderBookItem, error)
```

Get order book data for multiple stocks at once in a single request (up to 10 codes on the MAX plan, 25 on the ELITE plan).

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `code (path)` | `string` | yes |  | Daftar kode saham yang dipisahkan `\|` atau koma, mis. `BBCA\|GOTO\|TLKM`. Maksimal 10 kode (paket MAX) atau 25 kode (paket ELITE) per request. |
| `params.Market` | `*BatchOrderBookMarket` | no | `RG` | Jenis pasar: `RG` = Pasar Reguler, `NG` = Pasar Negosiasi, `TN` = Pasar Tunai. Gunakan `RG` untuk analisis pergerakan harga pada umumnya. Default `RG` bila tidak diisi. One of: RG, NG, TN. Default: RG. |
| `params.Date` | `*string` | no |  | Tanggal order book historis (`YYYY-MM-DD`), data tersedia mulai `2026-03-25`. Kirim bersama `time`; kosongkan keduanya untuk data realtime. |
| `params.Time` | `*string` | no |  | Jam order book dalam WIB, format `HH:mm` dengan presisi menit. Kirim bersama `date`. |

Returns: `[]BatchOrderBookItem`

### GetIntradayData

**Batch Stock Intraday Data** — `GET /batch/intraday-data/{code}` — cache: **realtime**

```go
func (s *BatchService) GetIntradayData(ctx context.Context, code string, params *BatchIntradayDataParams) ([]BatchIntradayDataItem, error)
```

Get compact intraday price data for multiple stocks at once in a single request (up to 10 codes on the MAX plan, 25 on the ELITE plan).

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `code (path)` | `string` | yes |  | Daftar kode saham yang dipisahkan `\|` atau koma, mis. `BBCA\|GOTO\|TLKM`. Maksimal 10 kode (paket MAX) atau 25 kode (paket ELITE) per request. |
| `params.Market` | `*BatchIntradayDataMarket` | no | `RG` | Jenis pasar: `RG` = Pasar Reguler, `NG` = Pasar Negosiasi, `TN` = Pasar Tunai. Gunakan `RG` untuk analisis pergerakan harga pada umumnya. Default `RG` bila tidak diisi. One of: RG, NG, TN. Default: RG. |
| `params.Date` | `*string` | no |  | Tanggal data historis (`YYYY-MM-DD`). Kosongkan untuk data realtime hari ini. |

Returns: `[]BatchIntradayDataItem`

### GetIntradayIndex

**Batch Index Intraday Data** — `GET /batch/intraday-index/{code}` — cache: **realtime**

```go
func (s *BatchService) GetIntradayIndex(ctx context.Context, code string) ([]BatchIntradayIndexItem, error)
```

Get compact intraday data for multiple indexes at once in a single request (up to 10 codes on the MAX plan, 25 on the ELITE plan).

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `code (path)` | `string` | yes |  | Daftar kode indeks yang dipisahkan `\|` atau koma, mis. `COMPOSITE\|LQ45\|IDX30`. Maksimal 10 kode (paket MAX) atau 25 kode (paket ELITE) per request. Daftar indeks tersedia di `/analysis/list/index`. |

Returns: `[]BatchIntradayIndexItem`


## client.Search

Search stocks and users.

### Search

**Search Stock or User** — `GET /search` — cache: **default**

```go
func (s *SearchService) Search(ctx context.Context, params SearchParams) (*SearchResponse, error)
```

Search for stocks or users simultaneously based on query keywords.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `params.Query` | `string` | yes |  | Kata kunci pencarian berupa kode saham, nama perusahaan, atau username. |

Returns: `*SearchResponse`

### SearchStock

**Search Stock only** — `GET /search/stock` — cache: **default**

```go
func (s *SearchService) SearchStock(ctx context.Context, params SearchStockParams) (*SearchStockResponse, error)
```

Search for stock data listed on the Indonesia Stock Exchange (IDX) based on query keywords with pagination support.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `params.Query` | `string` | yes |  | Kata kunci pencarian berupa kode saham atau nama perusahaan. |
| `params.Cursor` | `*int64` | no | `1` | Nomor halaman hasil pencarian, dimulai dari 1. Setiap halaman berisi 20 hasil. Default: 1. Range: 1..-. |

Returns: `*SearchStockResponse`

### SearchUser

**Search User only** — `GET /search/user` — cache: **default**

```go
func (s *SearchService) SearchUser(ctx context.Context, params SearchUserParams) (*SearchUserResponse, error)
```

Search for Invezgo app users based on query keywords with pagination support.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `params.Query` | `string` | yes |  | Kata kunci pencarian berupa username atau nama pengguna. |
| `params.Cursor` | `*int64` | no | `1` | Nomor halaman hasil pencarian, dimulai dari 1. Setiap halaman berisi 20 hasil. Default: 1. Range: 1..-. |

Returns: `*SearchUserResponse`


## client.Screener

Run stock screeners and manage saved screener presets.

### Screen

**Screener** — `POST /screener/screen` — cache: **none**

```go
func (s *ScreenerService) Screen(ctx context.Context, body *ScreenRequest) ([]ScreenerScreenItem, error)
```

Get a list of stock data matching the specified formula.

Body: `*ScreenRequest` (required)

Returns: `[]ScreenerScreenItem`

### List

**Screener Preset List** — `GET /screener` — cache: **default**

```go
func (s *ScreenerService) List(ctx context.Context) ([]ScreenerListItem, error)
```

Get the list of screener presets, grouped into your own presets, public presets, and presets shared with your plan.

Returns: `[]ScreenerListItem`

### Save

**Save Screener Preset** — `POST /screener` — cache: **none**

```go
func (s *ScreenerService) Save(ctx context.Context, body *SaveScreenerRequest) (*ScreenerSaveResponse, error)
```

Save a new screener preset. The formula is tested first; if it is invalid, the preset is not saved and the 400 response contains the formula error details.

Body: `*SaveScreenerRequest` (required)

Returns: `*ScreenerSaveResponse`

### Update

**Update Screener Preset** — `PUT /screener/{id}` — cache: **none**

```go
func (s *ScreenerService) Update(ctx context.Context, id string, body *UpdateScreenerRequest) (*ScreenerUpdateResponse, error)
```

Update one of your screener presets by ID. All fields (name, description, scope, formula, category) are overwritten with the new values, and the formula is tested first.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `id (path)` | `string` | yes |  | ID preset screener, diambil dari field `id` pada response `GET /screener`. |

Body: `*UpdateScreenerRequest` (required)

Returns: `*ScreenerUpdateResponse`

### Delete

**Delete Screener Preset** — `DELETE /screener/{id}` — cache: **none**

```go
func (s *ScreenerService) Delete(ctx context.Context, id string) (*ScreenerDeleteResponse, error)
```

Delete one of your screener presets by ID.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `id (path)` | `string` | yes |  | ID preset screener, diambil dari field `id` pada response `GET /screener`. |

Returns: `*ScreenerDeleteResponse`


## client.Usage

API usage and quota monitoring (does not consume quota).

### GetAPIUsage

**API Usage** — `GET /usage/api` — cache: **none**

```go
func (s *UsageService) GetAPIUsage(ctx context.Context) (*ApiUsageResponse, error)
```

Information on API usage.

Returns: `*ApiUsageResponse`

### GetLogs

**API Usage Logs** — `GET /usage/logs` — cache: **none**

```go
func (s *UsageService) GetLogs(ctx context.Context, params *UsageLogsParams) ([]UsageLogsItem, error)
```

Get the API usage log history (requests made with an API key) of the authenticated user, newest first.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `params.From` | `*string` | no |  | Waktu awal dalam WIB, format `YYYY-MM-DD HH:mm:ss` atau `YYYY-MM-DD`. Maksimal 30 hari ke belakang. Default 24 jam terakhir. |
| `params.To` | `*string` | no |  | Waktu akhir dalam WIB, format `YYYY-MM-DD HH:mm:ss` atau `YYYY-MM-DD`. Harus lebih besar atau sama dengan `from`. Default waktu saat ini. |
| `params.Status` | `*int64` | no |  | Filter log berdasarkan HTTP status code, mis. `200`, `422`, atau `429`. Range: 100..599. |
| `params.Limit` | `*int64` | no |  | Jumlah log maksimal yang dikembalikan. Kosongkan untuk mengambil seluruh log. Range: 1..320000. |

Returns: `[]UsageLogsItem`


## client.Watchlists

Personal watchlists and watchlist groups.

### List

**Watchlist** — `GET /watchlists` — cache: **default**

```go
func (s *WatchlistsService) List(ctx context.Context, params *WatchlistsListParams) ([]WatchlistsListItem, error)
```

Get the stocks in the user's watchlist with their last price and percentage change against the saved reference price.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `params.Group` | `*string` | no |  | ID grup watchlist yang ingin ditampilkan, diambil dari field `id` pada response `GET /watchlists/group`. Kosongkan atau isi `null` untuk menampilkan seluruh watchlist dari semua grup. |

Returns: `[]WatchlistsListItem`

### Add

**Add Stock to Watchlist** — `POST /watchlists` — cache: **none**

```go
func (s *WatchlistsService) Add(ctx context.Context, body *AddWatchlistRequest) (*WatchlistsAddResponse, error)
```

Add a new stock to the user's watchlist.

Body: `*AddWatchlistRequest` (required)

Returns: `*WatchlistsAddResponse`

### Delete

**Delete Watchlist Items** — `DELETE /watchlists` — cache: **none**

```go
func (s *WatchlistsService) Delete(ctx context.Context, body *DeleteIdsRequest) (*WatchlistsDeleteResponse, error)
```

Delete one or more stocks from the user's watchlist by item ID.

Body: `*DeleteIdsRequest` (required)

Returns: `*WatchlistsDeleteResponse`

### ListGroups

**Watchlist Group List** — `GET /watchlists/group` — cache: **default**

```go
func (s *WatchlistsService) ListGroups(ctx context.Context) ([]WatchlistsListGroupsItem, error)
```

Get the user's watchlist groups.

Returns: `[]WatchlistsListGroupsItem`

### AddGroup

**Create Watchlist Group** — `POST /watchlists/group` — cache: **none**

```go
func (s *WatchlistsService) AddGroup(ctx context.Context, body *WatchlistGroupRequest) (*WatchlistsAddGroupResponse, error)
```

Create a new watchlist group for the user.

Body: `*WatchlistGroupRequest` (required)

Returns: `*WatchlistsAddGroupResponse`

### UpdateGroup

**Update Watchlist Group** — `PUT /watchlists/group/{id}` — cache: **none**

```go
func (s *WatchlistsService) UpdateGroup(ctx context.Context, id string, body *WatchlistGroupRequest) (*WatchlistsUpdateGroupResponse, error)
```

Rename one of the user's watchlist groups.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `id (path)` | `string` | yes |  | ID grup watchlist, diambil dari field `id` pada response `GET /watchlists/group`. |

Body: `*WatchlistGroupRequest` (required)

Returns: `*WatchlistsUpdateGroupResponse`

### DeleteGroup

**Delete Watchlist Group** — `DELETE /watchlists/group/{id}` — cache: **none**

```go
func (s *WatchlistsService) DeleteGroup(ctx context.Context, id string) (*WatchlistsDeleteGroupResponse, error)
```

Delete one of the user's watchlist groups.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `id (path)` | `string` | yes |  | ID grup watchlist, diambil dari field `id` pada response `GET /watchlists/group`. |

Returns: `*WatchlistsDeleteGroupResponse`

### Update

**Update Watchlist Item** — `PUT /watchlists/{id}` — cache: **none**

```go
func (s *WatchlistsService) Update(ctx context.Context, id string, body *UpdateWatchlistRequest) (*WatchlistsUpdateResponse, error)
```

Move a watchlist item to another group and/or change its reference price.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `id (path)` | `string` | yes |  | ID item watchlist, diambil dari field `id` pada response `GET /watchlists`. |

Body: `*UpdateWatchlistRequest` (required)

Returns: `*WatchlistsUpdateResponse`

### UpdateNote

**Update Watchlist Note** — `PATCH /watchlists/{id}` — cache: **none**

```go
func (s *WatchlistsService) UpdateNote(ctx context.Context, id string, body *NoteRequest) (*WatchlistsUpdateNoteResponse, error)
```

Update the note on a watchlist item.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `id (path)` | `string` | yes |  | ID item watchlist, diambil dari field `id` pada response `GET /watchlists`. |

Body: `*NoteRequest` (required)

Returns: `*WatchlistsUpdateNoteResponse`


## client.Journals

Personal trading journal transactions.

### Add

**Add Journal Transaction** — `POST /journals` — cache: **none**

```go
func (s *JournalsService) Add(ctx context.Context, body *AddJournalRequest) (*JournalsAddResponse, error)
```

Add one BUY or SELL transaction to the user's journal and update the portfolio. A BUY adds lots and recalculates the average price (incl. buy fee) for the stock code + broker pair; a SELL reduces lots and records a realized trade (profit/loss) shown in `GET /trades`.

Body: `*AddJournalRequest` (required)

Returns: `*JournalsAddResponse`

### List

**Journal Transaction List** — `GET /journals` — cache: **default**

```go
func (s *JournalsService) List(ctx context.Context, params JournalsListParams) ([]JournalsListItem, error)
```

Get the user's journal transactions (BUY/SELL) dated between `from` and `to` (inclusive), newest first.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `params.From` | `string` | yes |  | Tanggal awal periode (`YYYY-MM-DD`). |
| `params.To` | `string` | yes |  | Tanggal akhir periode (`YYYY-MM-DD`). Harus lebih besar atau sama dengan `from`. |

Returns: `[]JournalsListItem`

### Delete

**Delete Journal Transactions** — `DELETE /journals` — cache: **none**

```go
func (s *JournalsService) Delete(ctx context.Context, body *DeleteIdsRequest) (*JournalsDeleteResponse, error)
```

Delete one or more journal transactions by ID (body `{"id": ["<id>", ...]}`) and reverse their effect on the portfolio: deleting a BUY reduces the position's lots, deleting a SELL returns the lots to the position and removes the related realized trade.

Body: `*DeleteIdsRequest` (required)

Returns: `*JournalsDeleteResponse`

### GetSummary

**Journal Summary** — `GET /journals/summary` — cache: **default**

```go
func (s *JournalsService) GetSummary(ctx context.Context, params JournalsSummaryParams) (*JournalsSummaryResponse, error)
```

Get a summary of the user's journal transactions dated between `from` and `to` (inclusive): counts and percentages per transaction type, total and average value, and a per-stock recap.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `params.From` | `string` | yes |  | Tanggal awal periode (`YYYY-MM-DD`). |
| `params.To` | `string` | yes |  | Tanggal akhir periode (`YYYY-MM-DD`). Harus lebih besar atau sama dengan `from`. |

Returns: `*JournalsSummaryResponse`


## client.Portfolios

Personal portfolio built from journal transactions.

### List

**Portfolio** — `GET /portfolios` — cache: **default**

```go
func (s *PortfoliosService) List(ctx context.Context) ([]PortfoliosListItem, error)
```

Get the stock positions the user currently holds (built from journal entries), combined per stock code across all brokers, most recently updated first.

Returns: `[]PortfoliosListItem`

### GetSummary

**Portfolio Summary** — `GET /portfolios/summary` — cache: **default**

```go
func (s *PortfoliosService) GetSummary(ctx context.Context) (*PortfoliosSummaryResponse, error)
```

Get the user's portfolio summary: total cost basis, market value, unrealized profit/loss, and allocation per stock, sector and subsector.

Returns: `*PortfoliosSummaryResponse`


## client.Alerts

Live formula alerts.

### Create

**Create Alert** — `POST /alerts` — cache: **none**

```go
func (s *AlertsService) Create(ctx context.Context, body *AlertRequest) (*AlertsCreateResponse, error)
```

Create an alert to get notified when stocks meet a given formula.

Body: `*AlertRequest` (required)

Returns: `*AlertsCreateResponse`

### List

**Alert List** — `GET /alerts` — cache: **default**

```go
func (s *AlertsService) List(ctx context.Context) ([]AlertsListItem, error)
```

List of created alerts.

Returns: `[]AlertsListItem`

### Test

**Test Alert Formula** — `POST /alerts/test` — cache: **none**

```go
func (s *AlertsService) Test(ctx context.Context, body *TestAlertRequest) (*AlertsTestResponse, error)
```

Test an alert formula before creating the alert.

Body: `*TestAlertRequest` (required)

Returns: `*AlertsTestResponse`

### Update

**Update Alert** — `PUT /alerts/{id}` — cache: **none**

```go
func (s *AlertsService) Update(ctx context.Context, id string, body *UpdateAlertRequest) (*AlertsUpdateResponse, error)
```

Update a created alert.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `id (path)` | `string` | yes |  | ID alert, diambil dari field `id` pada response `GET /alerts`. |

Body: `*UpdateAlertRequest` (required)

Returns: `*AlertsUpdateResponse`

### Delete

**Delete Alert** — `DELETE /alerts/{id}` — cache: **none**

```go
func (s *AlertsService) Delete(ctx context.Context, id string) (*AlertsDeleteResponse, error)
```

Delete a created alert.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `id (path)` | `string` | yes |  | ID alert, diambil dari field `id` pada response `GET /alerts`. |

Returns: `*AlertsDeleteResponse`


## client.Posts

Community posts, comments, likes and votes.

### List

**Post List** — `GET /posts` — cache: **default**

```go
func (s *PostsService) List(ctx context.Context, params PostsListParams) (*PostsListResponse, error)
```

Get the list of all Invezgo community posts with pagination, sorted newest first.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `params.Page` | `int64` | yes |  | Nomor halaman, dimulai dari 1. Range: 1..-. |
| `params.Limit` | `int64` | yes |  | Jumlah postingan per halaman. Range: 1..20. |

Returns: `*PostsListResponse`

### ListByCategory

**Posts by Category** — `GET /posts/category/{category}` — cache: **default**

```go
func (s *PostsService) ListByCategory(ctx context.Context, category PostsListByCategoryCategory, params PostsListByCategoryParams) (*PostsListByCategoryResponse, error)
```

Get the community posts of one category with pagination, sorted newest first.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `category (path)` | `PostsListByCategoryCategory` | yes |  | Kategori postingan: `POST` = postingan umum, `ARTICLE` = artikel, `VIDEO` = video, `STOCKPICK` = rekomendasi saham, `POLLING` = polling. One of: POST, ARTICLE, VIDEO, STOCKPICK, POLLING. |
| `params.Page` | `int64` | yes |  | Nomor halaman, dimulai dari 1. Range: 1..-. |
| `params.Limit` | `int64` | yes |  | Jumlah postingan per halaman. Range: 1..20. |

Returns: `*PostsListByCategoryResponse`

### ListByStock

**Stock Posts** — `GET /posts/space/{code}` — cache: **default**

```go
func (s *PostsService) ListByStock(ctx context.Context, code string, params PostsListByStockParams) (*PostsListByStockResponse, error)
```

Get the posts in a stock's discussion space (posts that mention the stock code) with pagination, sorted newest first.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `code (path)` | `string` | yes |  | Kode saham untuk menampilkan postingan pada ruang diskusi saham tersebut. |
| `params.Page` | `int64` | yes |  | Nomor halaman, dimulai dari 1. Range: 1..-. |
| `params.Limit` | `int64` | yes |  | Jumlah postingan per halaman. Range: 1..20. |

Returns: `*PostsListByStockResponse`

### ListByStockCategory

**Stock Posts by Category** — `GET /posts/space/category/{code}/{category}` — cache: **default**

```go
func (s *PostsService) ListByStockCategory(ctx context.Context, code string, category PostsListByStockCategoryCategory, params PostsListByStockCategoryParams) (*PostsListByStockCategoryResponse, error)
```

Get the posts in a stock's discussion space for one category with pagination, sorted newest first.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `code (path)` | `string` | yes |  | Kode saham untuk menampilkan postingan pada ruang diskusi saham tersebut. |
| `category (path)` | `PostsListByStockCategoryCategory` | yes |  | Kategori postingan: `POST` = postingan umum, `ARTICLE` = artikel, `VIDEO` = video, `STOCKPICK` = rekomendasi saham, `POLLING` = polling. One of: POST, ARTICLE, VIDEO, STOCKPICK, POLLING. |
| `params.Page` | `int64` | yes |  | Nomor halaman, dimulai dari 1. Range: 1..-. |
| `params.Limit` | `int64` | yes |  | Jumlah postingan per halaman. Range: 1..20. |

Returns: `*PostsListByStockCategoryResponse`

### Get

**Post Detail** — `GET /posts/detail/{id}` — cache: **default**

```go
func (s *PostsService) Get(ctx context.Context, id string) (*PostDetailResponse, error)
```

Get the full content of one post by ID.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `id (path)` | `string` | yes |  | ID postingan, diambil dari field `id` pada response daftar postingan (mis. `GET /posts`). |

Returns: `*PostDetailResponse`

### GetComments

**Post Comments** — `GET /posts/comment/{id}` — cache: **default**

```go
func (s *PostsService) GetComments(ctx context.Context, id string, params PostsCommentsParams) (*PostsCommentsResponse, error)
```

Get the comments on a post with pagination, sorted by oldest update first.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `id (path)` | `string` | yes |  | ID postingan, diambil dari field `id` pada response daftar postingan (mis. `GET /posts`). |
| `params.Page` | `int64` | yes |  | Nomor halaman, dimulai dari 1. Range: 1..-. |
| `params.Limit` | `int64` | yes |  | Jumlah postingan per halaman. Range: 1..20. |

Returns: `*PostsCommentsResponse`

### GetLikes

**Liked Posts** — `GET /posts/like` — cache: **default**

```go
func (s *PostsService) GetLikes(ctx context.Context, params PostsLikesParams) (*PostsLikesResponse, error)
```

Get the posts and comments you have liked with pagination, sorted by most recently liked.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `params.Page` | `int64` | yes |  | Nomor halaman, dimulai dari 1. Range: 1..-. |
| `params.Limit` | `int64` | yes |  | Jumlah postingan per halaman. Range: 1..20. |

Returns: `*PostsLikesResponse`

### GetFavorites

**Favorite Posts** — `GET /posts/favorite` — cache: **default**

```go
func (s *PostsService) GetFavorites(ctx context.Context, params PostsFavoritesParams) (*PostsFavoritesResponse, error)
```

Get the posts and comments you have favorited with pagination, sorted by most recently favorited.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `params.Page` | `int64` | yes |  | Nomor halaman, dimulai dari 1. Range: 1..-. |
| `params.Limit` | `int64` | yes |  | Jumlah postingan per halaman. Range: 1..20. |

Returns: `*PostsFavoritesResponse`

### GetVoters

**Post Votes** — `GET /posts/vote/{id}` — cache: **default**

```go
func (s *PostsService) GetVoters(ctx context.Context, id string) (*PostsVotersResponse, error)
```

Get the vote summary (agree/disagree) of a stock pick (stockpick) post.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `id (path)` | `string` | yes |  | ID postingan, diambil dari field `id` pada response daftar postingan (mis. `GET /posts`). |

Returns: `*PostsVotersResponse`


## client.Trades

Personal realized trades.

### List

**Realized Trade List** — `GET /trades` — cache: **default**

```go
func (s *TradesService) List(ctx context.Context, params TradesListParams) ([]TradesListItem, error)
```

Get the user's realized trades (from SELL journal entries) with a sell date between `from` and `to` (inclusive), newest first.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `params.From` | `string` | yes |  | Tanggal awal periode (`YYYY-MM-DD`). |
| `params.To` | `string` | yes |  | Tanggal akhir periode (`YYYY-MM-DD`). Harus lebih besar atau sama dengan `from`. |

Returns: `[]TradesListItem`

### Delete

**Delete Realized Trades** — `DELETE /trades` — cache: **none**

```go
func (s *TradesService) Delete(ctx context.Context, body *DeleteIdsRequest) (*TradesDeleteResponse, error)
```

Delete one or more realized trades by ID (body `{"id": ["<id>", ...]}`). Only the realized-trade records are removed; journal transactions and the portfolio are unchanged.

Body: `*DeleteIdsRequest` (required)

Returns: `*TradesDeleteResponse`

### GetSummary

**Realized Trade Summary** — `GET /trades/summary` — cache: **default**

```go
func (s *TradesService) GetSummary(ctx context.Context, params TradesSummaryParams) (*TradesSummaryResponse, error)
```

Get a performance summary of the user's realized trades with a sell date between `from` and `to` (inclusive): profit/loss/BEP counts, totals, win rate, and the best and worst trades.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `params.From` | `string` | yes |  | Tanggal awal periode (`YYYY-MM-DD`). |
| `params.To` | `string` | yes |  | Tanggal akhir periode (`YYYY-MM-DD`). Harus lebih besar atau sama dengan `from`. |

Returns: `*TradesSummaryResponse`

### GetSummaryChart

**Realized Trade Summary Chart** — `GET /trades/summary-chart` — cache: **default**

```go
func (s *TradesService) GetSummaryChart(ctx context.Context, params TradesSummaryChartParams) ([]TradesSummaryChartItem, error)
```

Get daily chart data of the user's realized profit/loss with a sell date between `from` and `to` (inclusive), with its running total and the comparison index close on the same day. Only days with realized trades are returned, in ascending date order.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `params.From` | `string` | yes |  | Tanggal awal periode (`YYYY-MM-DD`). |
| `params.To` | `string` | yes |  | Tanggal akhir periode (`YYYY-MM-DD`). Harus lebih besar atau sama dengan `from`. |
| `params.Code` | `string` | yes |  | Kode indeks pembanding pada grafik, mis. `COMPOSITE`. Bila dikirim kosong (`code=`), setiap tanggal berisi satu baris per indeks. |

Returns: `[]TradesSummaryChartItem`

### UpdateNote

**Update Realized Trade Note** — `PATCH /trades/{id}` — cache: **none**

```go
func (s *TradesService) UpdateNote(ctx context.Context, id string, body *NoteTradeRequest) (*TradesUpdateNoteResponse, error)
```

Update the note of one of the user's realized trades. Body: `{"note": "..."}` (optional, max 150 characters); when `note` is omitted the note is cleared.

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `id (path)` | `string` | yes |  | ID transaksi terealisasi, diambil dari field `id` pada response `GET /trades`. |

Body: `*NoteTradeRequest` (required)

Returns: `*TradesUpdateNoteResponse`

