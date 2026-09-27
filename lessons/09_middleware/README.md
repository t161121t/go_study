# 09: ミドルウェア・context

実務のAPIサーバーには「全リクエスト共通の処理」（ログ・認証・タイムアウト）が必ずあります。
Go ではこれを **ミドルウェア** というパターンで書きます。

## ミドルウェアの正体

「ハンドラを受け取って、別のハンドラを返す関数」です。

```go
type Middleware func(http.Handler) http.Handler

func Logging(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        fmt.Println(r.Method, r.URL.Path)  // 前処理
        next.ServeHTTP(w, r)               // 本来のハンドラを呼ぶ
        // ここに後処理も書ける
    })
}

handler = Logging(handler)   // 元のハンドラをラップする
```

`http.Handler` は `ServeHTTP(http.ResponseWriter, *http.Request)` を持つインターフェース。
`http.HandlerFunc` はただの関数をこのインターフェースに適合させるアダプタです。

## 複数のミドルウェアをまとめる

```go
func Chain(mws ...Middleware) Middleware {
    return func(final http.Handler) http.Handler {
        for i := len(mws) - 1; i >= 0; i-- {   // 後ろから包んでいく
            final = mws[i](final)
        }
        return final
    }
}

handler = Chain(Logging, Auth)(handler)
```

## 認証ミドルウェア

```go
func Auth(token string) Middleware {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            got := r.Header.Get("Authorization")
            if got != "Bearer "+token {
                http.Error(w, "unauthorized", http.StatusUnauthorized)
                return   // next を呼ばない = 処理を止める
            }
            next.ServeHTTP(w, r)
        })
    }
}
```

## context によるタイムアウト

`context.Context` は「このリクエストはいつまで有効か」「キャンセルされたか」を
下流に伝える仕組みです。リクエストは必ず `r.Context()` を持っています。

```go
func WithTimeout(d time.Duration) Middleware {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            ctx, cancel := context.WithTimeout(r.Context(), d)
            defer cancel()
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
```

時間のかかる処理側は `ctx.Done()` を監視して、タイムアウトしたら諦めます。

```go
select {
case <-time.After(2 * time.Second):
    w.Write([]byte("done"))
case <-r.Context().Done():
    http.Error(w, "timeout", http.StatusGatewayTimeout)
}
```

`select` は複数の channel を同時に待ち、先に来た方を実行する Go の構文です。

```sh
go test -v ./lessons/09_middleware/
```
