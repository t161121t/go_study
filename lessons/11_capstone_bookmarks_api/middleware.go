package main

import (
	"log"
	"net/http"
	"time"
)

// Middleware は09と同じ形（ハンドラを受け取ってハンドラを返す）
type Middleware func(http.Handler) http.Handler

// Chain は09と同じ。複数のMiddlewareを1つにまとめる
func Chain(mws ...Middleware) Middleware {
	return func(final http.Handler) http.Handler {
		for i := len(mws) - 1; i >= 0; i-- {
			final = mws[i](final)
		}
		return final
	}
}

// statusRecorder は http.ResponseWriter を埋め込み、
// WriteHeader で渡されたステータスコードだけをこっそりメモしておく
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

// Logging は method, path, ステータスコード, 処理時間 をログに出す
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(rec, r)

		log.Printf("%s %s %d %s", r.Method, r.URL.Path, rec.status, time.Since(start))
	})
}

// AuthExceptGet は GET 以外のメソッドにだけ Bearer トークン認証を要求する
// （GET は誰でも見られる。09のAuthをメソッドで場合分けするよう拡張したもの）
func AuthExceptGet(token string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet && r.Header.Get("Authorization") != "Bearer "+token {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
