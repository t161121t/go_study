# 09: ミドルウェア・context — まとめ

## ミドルウェアとは

「ハンドラを箱で包む関数」。全リクエスト共通の処理（ログ・認証・タイムアウト）を
差し込むための仕組み。

```go
type Middleware func(http.Handler) http.Handler
```

```
お客さん → [ミドルウェア] → [本来のハンドラ]
```

## Chain — 複数を重ねる（マトリョーシカ）

`Chain(A, B)(h)` を `A(B(h))` にしたい。**後ろから**包んでいく。

```go
func Chain(mws ...Middleware) Middleware {
	return func(final http.Handler) http.Handler {
		for i := len(mws) - 1; i >= 0; i-- {
			final = mws[i](final)
		}
		return final
	}
}
```

`mws = [A, B]` なら: `final = B(h)` → `final = A(B(h))`。
一番外側がA、内側がBになるので、実行順は A→B→本来のハンドラ。

## ポインタで外の変数を書き換える

```go
func CountingMiddleware(calls *int) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			*calls++             // 05で習ったポインタの中身書き換え
			next.ServeHTTP(w, r)
		})
	}
}
```

## 認証ミドルウェア — next を呼ばずに止める

```go
func Auth(token string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Authorization") != "Bearer "+token {
				w.WriteHeader(http.StatusUnauthorized)
				return // next を呼ばない = 先に進ませない
			}
			next.ServeHTTP(w, r)
		})
	}
}
```

## context.Context — 「有効期限つきの整理券」

リクエストには必ず`r.Context()`（今の整理券）が付いている。

```go
ctx, cancel := context.WithTimeout(r.Context(), d) // d時間後に期限切れになる新しい整理券
defer cancel()                                       // 使い終わったら後片付け（Unlockと同じ理由）
next.ServeHTTP(w, r.WithContext(ctx))                // 新しい整理券を持たせて渡す
```

## select — 複数のchannelを同時に待って先着勝負

```go
func SlowHandler(w http.ResponseWriter, r *http.Request) {
	select {
	case <-time.After(200 * time.Millisecond): // 200ms後に届く通知
		w.Write([]byte("done"))
	case <-r.Context().Done(): // 整理券が期限切れになったら届く通知
		w.WriteHeader(http.StatusGatewayTimeout) // 504
	}
}
```

- 制限時間 > 200ms → `time.After`が先着 → `"done"`
- 制限時間 < 200ms → `ctx.Done()`が先着 → 504タイムアウト

## 実装した関数

- `Chain` — 複数のMiddlewareを1つにまとめる
- `CountingMiddleware` — リクエスト回数を数える最小のミドルウェア例
- `Auth` — トークン認証（不一致なら401で止める）
- `WithTimeout` — リクエストのcontextに制限時間を付ける
- `SlowHandler` — `select`でタイムアウトと処理完了のどちらが先か競わせる

→ `go test -race -v ./lessons/09_middleware/` でPASS確認済み。

## 次: 10_concurrency

goroutine・channel・sync.Mutex・ワーカープール。09で少し触れたchannelを本格的に使う。
