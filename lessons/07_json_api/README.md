# 07: JSON API

Web API では JSON でやり取りするのが基本です。`encoding/json` を使います。

## 構造体 ⇔ JSON

```go
type Item struct {
    ID    int     `json:"id"`
    Name  string  `json:"name"`
    Price float64 `json:"price"`
}
```

`` `json:"id"` `` は **構造体タグ**。JSON にする時のキー名を指定します
（省略するとフィールド名そのまま `"ID"` になり、他言語の慣習と合わずやりづらい）。

## レスポンスとして JSON を返す

```go
func WriteItem(w http.ResponseWriter, item Item) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(item)   // 構造体 → JSON にして書き込む
}
```

`w.Header().Set` は **必ず `w.WriteHeader` より前** に呼びます（ヘッダを送った後は変更できません）。

## リクエストボディから JSON を読む

```go
func ReadItem(r *http.Request) (Item, error) {
    var item Item
    err := json.NewDecoder(r.Body).Decode(&item)   // JSON → 構造体（ポインタを渡す）
    return item, err
}
```

不正な JSON が来たら `err != nil` になるので、400 Bad Request を返すのが定石です。

## エラーもJSONで返す

```go
type ErrorResponse struct {
    Error string `json:"error"`
}

func WriteError(w http.ResponseWriter, status int, message string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}
```

## バリデーション

JSON としては正しくても、値がおかしいことがあります（名前が空、価格がマイナス、など）。
デコードした後に自分でチェックするのが基本です。

```go
if item.Name == "" {
    WriteError(w, http.StatusBadRequest, "name is required")
    return
}
```

```sh
go test -v ./lessons/07_json_api/
```
