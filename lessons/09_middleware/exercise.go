package middleware

import (
	"context"
	"net/http"
	"time"
)

// Middleware はハンドラをラップする関数
type Middleware func(http.Handler) http.Handler

// 課題1: mws を後ろから順に適用する1つの Middleware にまとめて返す
// 例: Chain(A, B)(h) は A(B(h)) と同じ動きになる
// （つまりリクエストは A → B → h の順に通る）
func Chain(mws ...Middleware) Middleware {
	return func(final http.Handler) http.Handler {
		for i := len(mws) - 1; i >= 0; i-- {
			final = mws[i](final)
		}
		return final
	}
}

// 課題2: リクエストごとに calls をインクリメントするロギング用ミドルウェア
// 本来のログ出力はせず、代わりに *int（呼び出し回数のカウンタ）を増やすだけでよい
// ヒント: 前処理として *calls++ してから next.ServeHTTP を呼ぶ
func CountingMiddleware(calls *int) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			*calls++
			next.ServeHTTP(w, r)
		})
	}
}

// 課題3: Authorization ヘッダが "Bearer "+token と一致しなければ
// 401 Unauthorized（本文は自由）を返し、次のハンドラを呼ばずに止める
// 一致すれば next をそのまま呼ぶ
func Auth(token string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Authorization") != "Bearer "+token {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// 課題4: リクエストの context に d の制限時間をつけて次のハンドラへ渡す
// ヒント: context.WithTimeout, defer cancel(), r.WithContext(ctx)
func WithTimeout(d time.Duration) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), d)
			defer cancel()
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// 課題5: 200ms かかる処理を模したハンドラ
// - 200ms 経過する前に ctx が Done になったら 504 Gateway Timeout を返す
// - 何事もなく200ms 経過したら 200 と本文 "done" を返す
// ヒント: time.After(200*time.Millisecond) と r.Context().Done() を select で待つ
func SlowHandler(w http.ResponseWriter, r *http.Request) {
	select {
	case <-time.After(200 * time.Millisecond):
		w.Write([]byte("done"))
	case <-r.Context().Done():
		w.WriteHeader(http.StatusGatewayTimeout)
	}
}
