# Invezgo Go SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/invezgo/invezgo-go-sdk/v2.svg)](https://pkg.go.dev/github.com/invezgo/invezgo-go-sdk/v2)
[![Go](https://img.shields.io/badge/go-%3E%3D1.21-00ADD8.svg)](https://go.dev)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

SDK resmi untuk [Invezgo API](https://invezgo.com), penyedia data pasar modal Indonesia (BEI/IDX). SDK ini menyediakan akses bertipe ke seluruh endpoint Invezgo API, dilengkapi cache in-memory dan rotasi API key otomatis.

[Dokumentasi API](https://api.invezgo.com/documentation) · [Referensi SDK](docs/API.md) · [Changelog](CHANGELOG.md)

## Daftar Isi

- [Persyaratan](#persyaratan)
- [Instalasi](#instalasi)
- [Penggunaan](#penggunaan)
- [Konfigurasi](#konfigurasi)
- [Memanggil Endpoint](#memanggil-endpoint)
- [Tipe Data](#tipe-data)
- [Caching](#caching)
- [Rotasi API Key](#rotasi-api-key)
- [Penanganan Error](#penanganan-error)
- [Pagination](#pagination)
- [Praktik yang Disarankan](#praktik-yang-disarankan)
- [Pengembangan](#pengembangan)
- [Dukungan](#dukungan)
- [Lisensi](#lisensi)

## Persyaratan

- Go 1.21 atau lebih baru
- API key Invezgo, dapat dibuat di [halaman pengaturan API](https://invezgo.com/id/setting/api)

SDK ini hanya menggunakan standard library dan aman digunakan secara konkuren oleh banyak goroutine.

## Instalasi

```bash
go get github.com/invezgo/invezgo-go-sdk/v2
```

## Penggunaan

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/invezgo/invezgo-go-sdk/v2"
)

func main() {
	client, err := invezgo.NewClient(os.Getenv("INVEZGO_API_KEY"))
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	stocks, err := client.Analysis.GetStockList(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(len(stocks), stocks[0].Code)

	top, err := client.Analysis.GetTopChange(ctx, invezgo.TopChangeParams{Date: "2024-12-30"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(top.Gain[0].Code, top.Gain[0].Change)
}
```

`NewClient` mengembalikan error jika API key kosong. Contoh lengkap tersedia di [example/main.go](example/main.go).

## Konfigurasi

```go
client, err := invezgo.NewClient(
	os.Getenv("INVEZGO_API_KEY"),
	invezgo.WithAdditionalAPIKeys(os.Getenv("INVEZGO_API_KEY_2")),
	invezgo.WithTimeout(30*time.Second),
	invezgo.WithCache(invezgo.CacheConfig{TTL: 60 * time.Second, RealtimeTTL: 5 * time.Second}),
)
```

| Option | Default | Keterangan |
|---|---|---|
| `apiKey` (argumen) | wajib | API key utama. |
| `WithAdditionalAPIKeys(keys...)` | tidak ada | API key cadangan, digunakan saat terjadi `429`. Lihat [Rotasi API Key](#rotasi-api-key). |
| `WithBaseURL(url)` | `https://api.invezgo.com` | Base URL API. |
| `WithTimeout(d)` | `30s` | Timeout per request. Nilai `0` menonaktifkan timeout; `context` tetap berlaku. |
| `WithHTTPClient(*http.Client)` | `&http.Client{}` | HTTP client alternatif, misalnya untuk proxy atau tracing. |
| `WithCache(CacheConfig)` | aktif | Konfigurasi cache. |
| `WithoutCache()` | — | Menonaktifkan cache. |

| Field `CacheConfig` | Default | Keterangan |
|---|---|---|
| `TTL` | `60s` | TTL endpoint reguler. |
| `RealtimeTTL` | `5s` | TTL endpoint realtime. |
| `MaxEntries` | `1000` | Jumlah maksimum entri cache. |
| `MaxBytes` | `32 << 20` | Ukuran maksimum cache (32 MiB). |

Field bernilai `0` menggunakan nilai default. TTL negatif menonaktifkan cache untuk tier tersebut.

Buat satu instance client dan gunakan kembali di seluruh aplikasi. Cache dan status rotasi API key disimpan per instance.

## Memanggil Endpoint

Endpoint dikelompokkan berdasarkan resource:

| Resource | Cakupan |
|---|---|
| `Analysis` | Daftar saham, broker, dan indeks; harga, chart, intraday, order book, running trade; broker summary; pemegang saham; laporan keuangan; kalender aksi korporasi; berita dan keterbukaan informasi |
| `Batch` | Order book dan data intraday untuk beberapa saham dalam satu request |
| `Search` | Pencarian saham dan pengguna |
| `Screener` | Eksekusi screener dan pengelolaan preset |
| `Usage` | Pemakaian dan sisa kuota API |
| `Watchlists` | Watchlist dan grup watchlist |
| `Journals` | Jurnal transaksi |
| `Portfolios` | Portofolio |
| `Alerts` | Alert berbasis formula |
| `Posts` | Post, komentar, like, dan vote komunitas |
| `Trades` | Realized trade |

Setiap method mengikuti pola berikut:

```
client.<Resource>.<Method>(ctx, pathParams..., body, params)
```

| Argumen | Keterangan |
|---|---|
| `ctx` | `context.Context` untuk timeout dan pembatalan request. |
| `pathParams` | Parameter path sebagai argumen posisional, sesuai urutan di URL. |
| `body` | Request body untuk `POST`, `PUT`, `PATCH`, dan `DELETE`, bertipe `*<Nama>Request`. |
| `params` | Query parameter dalam struct `<Nama>Params`. Dikirim sebagai nilai jika endpoint memiliki query parameter wajib, dan sebagai pointer (boleh `nil`) jika seluruhnya opsional. |

Field opsional bertipe pointer dan dapat diisi menggunakan `invezgo.Ptr(v)`. Parameter enum memiliki tipe dan konstanta tersendiri.

```go
// Query parameter wajib dan opsional
top, err := client.Analysis.GetTopChange(ctx, invezgo.TopChangeParams{
	Date:           "2024-12-30",
	FilterColumn:   invezgo.Ptr(invezgo.TopChangeFilterColumnChange),
	FilterOperator: invezgo.Ptr(invezgo.TopChangeFilterOperatorGt),
	FilterValue:    invezgo.Ptr(5.0),
})

// Path parameter
stalker, err := client.Analysis.GetBrokerStalker(ctx, "YP", "BBCA", invezgo.BrokerStalkerParams{
	From:     "2024-12-01",
	To:       "2024-12-31",
	Investor: invezgo.BrokerStalkerInvestorAll,
	Market:   invezgo.BrokerStalkerMarketRG,
})

// Request body
rows, err := client.Screener.Screen(ctx, &invezgo.ScreenRequest{
	Formula:  "close > prev AND volume > 1000000",
	Category: []string{"SYARIAH"},
})

// Beberapa saham sekaligus, dihitung sebagai satu request
books, err := client.Batch.GetOrderBook(ctx, "BBCA|BBRI|BMRI", nil)

// Timeout per request
reqCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
defer cancel()
book, err := client.Analysis.GetOrderBook(reqCtx, "BBCA", invezgo.OrderBookParams{Market: invezgo.OrderBookMarketRG})
```

Daftar lengkap method, parameter, nilai default, dan tipe kembalian tersedia di [docs/API.md](docs/API.md).

## Tipe Data

Seluruh model request, query parameter, dan response dihasilkan dari spesifikasi OpenAPI dan dilengkapi dokumentasi per field.

| Kasus | Representasi |
|---|---|
| Parameter enum | Tipe khusus dengan konstanta, misalnya `invezgo.StockRunningTradeTypeBUY` |
| Enum pada response | `string` |
| Field wajib | Nilai biasa (`string`, `int64`, `float64`, `bool`, struct) |
| Field nullable atau opsional | Pointer, misalnya `*float64` |
| Response yang dapat bernilai `null` | Pointer `nil` dengan `err == nil`, misalnya `GetShareholderClassifyTable` |
| Angka yang dapat dikirim sebagai string | `invezgo.FlexFloat` atau `invezgo.FlexInt` |
| Object dengan key dinamis | `map[string]T` |

Kolom dinamis hasil screener tersedia pada field `Values`:

```go
for _, row := range rows {
	if v := row.Values["close"]; v != nil {
		fmt.Println(row.Code, *v)
	}
}
```

Response berbentuk union menggunakan satu struct gabungan; gunakan field `Type` untuk membedakan jenisnya:

```go
graph, err := client.Analysis.GetShareholderRelation(ctx, &invezgo.ShareholderRelationParams{Code: invezgo.Ptr("BBCA")})
if err != nil {
	log.Fatal(err)
}
for _, node := range graph.Nodes {
	if node.Type == "stock" && node.Code != nil {
		fmt.Println(*node.Code)
	} else {
		fmt.Println(node.Label)
	}
}
```

## Caching

Cache in-memory aktif secara default untuk mengurangi pemakaian kuota dan latensi.

- Hanya request `GET` yang disimpan. Cache key terdiri dari path dan query string yang telah dinormalisasi.
- Cache menggunakan algoritma LRU dengan batas jumlah entri (`MaxEntries`) dan total ukuran (`MaxBytes`). Entri yang paling lama tidak digunakan akan dihapus saat salah satu batas terlampaui.
- Endpoint realtime (order book, running trade, intraday, chart, batch, dan sejenisnya) menggunakan `RealtimeTTL`. Endpoint lainnya menggunakan `TTL`.
- Endpoint `Usage` dan request selain `GET` tidak disimpan di cache.
- Setiap request `POST`, `PUT`, `PATCH`, atau `DELETE` yang berhasil akan menghapus cache data personal (`watchlists`, `journals`, `portfolios`, `trades`, `alerts`, `screener`).

```go
invezgo.NewClient(key, invezgo.WithoutCache())                                          // nonaktifkan cache
invezgo.NewClient(key, invezgo.WithCache(invezgo.CacheConfig{RealtimeTTL: -1}))         // nonaktifkan cache untuk endpoint realtime
invezgo.NewClient(key, invezgo.WithCache(invezgo.CacheConfig{TTL: 5 * time.Minute}))    // TTL endpoint reguler 5 menit

client.ClearCache() // hapus seluruh cache
```

## Rotasi API Key

Setiap API key memiliki batas request per menit untuk setiap endpoint dan kuota bulanan. Jika Anda memiliki lebih dari satu API key, tambahkan sebagai cadangan:

```go
client, err := invezgo.NewClient(
	os.Getenv("INVEZGO_API_KEY"),
	invezgo.WithAdditionalAPIKeys(os.Getenv("INVEZGO_API_KEY_2"), os.Getenv("INVEZGO_API_KEY_3")),
)
```

Ketika API mengembalikan `429 Too Many Requests`:

1. Lama penundaan diambil dari header `Retry-After`. Jika tidak tersedia, digunakan `X-RateLimit-Reset`, kemudian 60 detik.
2. Untuk batas per menit, API key ditangguhkan hanya pada endpoint tersebut. Untuk kuota bulanan, API key ditangguhkan pada seluruh endpoint.
3. Request dikirim ulang menggunakan API key berikutnya yang tersedia.
4. API key utama digunakan kembali setelah masa penangguhan berakhir.
5. Jika seluruh API key sedang ditangguhkan, SDK mengembalikan `*invezgo.RateLimitError` tanpa mengirim request.

Rotasi hanya dilakukan untuk status `429`. Untuk endpoint data personal seperti watchlist dan portofolio, gunakan API key dari akun yang sama.

Sisa kuota dapat diperiksa tanpa mengurangi kuota:

```go
usage, err := client.Usage.GetAPIUsage(ctx)
if err != nil {
	log.Fatal(err)
}
fmt.Println(usage.Usage, usage.Limit, usage.Remaining, usage.Expire)
```

## Penanganan Error

Response dengan status selain 2xx dikembalikan sebagai `*invezgo.APIError`. Status `429` dikembalikan sebagai `*invezgo.RateLimitError`, yang juga dapat dibaca sebagai `*invezgo.APIError` melalui `errors.As`.

| Helper | Status HTTP | Penyebab |
|---|---|---|
| `IsValidation` | 400, 422 | Parameter atau body tidak valid |
| `IsAuthentication` | 401 | API key tidak valid atau kedaluwarsa |
| `IsPaymentRequired` | 402 | Endpoint tidak termasuk dalam paket langganan |
| `IsNotFound` | 404 | Data tidak ditemukan |
| `IsRateLimited` | 429 | Batas request atau kuota tercapai pada seluruh API key |
| `IsServerError` | 5xx | Kesalahan pada server |

Field `APIError`: `StatusCode`, `Message`, `ErrorText`, `RequestID`, `Path`, dan `Body`. `RateLimitError` juga menyediakan `RetryAfter`, `Limit`, `Remaining`, `Reset`, dan `QuotaExceeded`. Timeout dan kegagalan jaringan dikembalikan sebagai error standar Go, misalnya `context.DeadlineExceeded`.

```go
_, err := client.Analysis.GetSummaryStock(ctx, "BBCA", invezgo.SummaryStockParams{
	From:     "2024-12-01",
	To:       "2024-12-31",
	Investor: invezgo.SummaryStockInvestorAll,
	Market:   invezgo.SummaryStockMarketRG,
})
if err != nil {
	var rateErr *invezgo.RateLimitError
	var apiErr *invezgo.APIError
	switch {
	case errors.As(err, &rateErr):
		log.Printf("rate limit tercapai, coba lagi dalam %s", rateErr.RetryAfter)
	case invezgo.IsValidation(err):
		log.Printf("parameter tidak valid: %v", err)
	case errors.As(err, &apiErr):
		log.Printf("%d %s (request id: %s)", apiErr.StatusCode, apiErr.Message, apiErr.RequestID)
	default:
		log.Printf("kegagalan jaringan: %v", err)
	}
}
```

## Pagination

Endpoint dengan pagination mengembalikan `Page`, `TotalPage`, dan `NextPage`. Nilai `NextPage` adalah `nil` pada halaman terakhir.

```go
page := int64(1)
for {
	res, err := client.Analysis.GetNews(ctx, &invezgo.NewsParams{Page: invezgo.Ptr(page), Limit: invezgo.Ptr(int64(50))})
	if err != nil {
		log.Fatal(err)
	}
	for _, item := range res.Data {
		fmt.Println(item.Title)
	}
	if res.NextPage == nil {
		break
	}
	page = *res.NextPage
}
```

## Praktik yang Disarankan

- Gunakan satu instance client untuk seluruh aplikasi.
- Teruskan `context` dari HTTP handler atau job agar request ikut dibatalkan.
- Simpan API key di environment variable atau secret manager.
- Gunakan format tanggal `YYYY-MM-DD`. Data end-of-day tersedia setelah pukul 18.00 WIB.
- Field waktu bersufiks `Z` menggunakan waktu WIB (UTC+7); jangan dikonversi ke zona waktu lain.
- Gunakan resource `Batch` untuk mengambil data beberapa saham sekaligus.
- Saat menerima `RateLimitError`, tunggu sesuai `RetryAfter` sebelum mencoba kembali.
- Untuk polling data realtime, sesuaikan `RealtimeTTL` dengan interval polling.

## Pengembangan

File `*_gen.go` dan `docs/API.md` dihasilkan dari spesifikasi OpenAPI dan tidak boleh diubah secara manual.

```bash
python codegen/generate.py go   # generate ulang kode dari OpenAPI (dari root repositori)
go vet ./...
go test ./...
```

## Dukungan

Sertakan nilai `RequestID` dari error saat menghubungi tim dukungan Invezgo.

- Website: [invezgo.com](https://invezgo.com)
- Dokumentasi API: [api.invezgo.com/documentation](https://api.invezgo.com/documentation)

## Lisensi

Dirilis di bawah [MIT License](LICENSE).
