# 08: REST CRUD（リポジトリパターン）

06（ルーティング）と07（JSON）を組み合わせて、本物っぽい CRUD API を作ります。
実務でも定番の **「データの保存場所」と「HTTPハンドラ」を分ける** 設計にします。

## 全体像

```
リクエスト → ハンドラ（HTTPの世話） → Store（データの保存・検索・更新・削除）
```

Store をハンドラから分離しておくと、後で「メモリ上の map」から「本物のDB」に
差し替えてもハンドラ側は変更不要になります。これを **リポジトリパターン** と呼びます。

## 複数 goroutine から同時アクセスされる問題

HTTP サーバーは、リクエストごとに **別々の goroutine** でハンドラが実行されます。
つまり複数のリクエストが同時に同じ map を読み書きする可能性があります。

```go
type Store struct {
    mu    sync.RWMutex     // 読み書きを保護するロック
    items map[int]Item
    next  int
}

func (s *Store) Get(id int) (Item, bool) {
    s.mu.RLock()           // 読み取り: 複数の goroutine が同時にOK
    defer s.mu.RUnlock()
    item, ok := s.items[id]
    return item, ok
}

func (s *Store) Create(item Item) Item {
    s.mu.Lock()            // 書き込み: 1つの goroutine だけ
    defer s.mu.Unlock()
    s.next++
    item.ID = s.next
    s.items[item.ID] = item
    return item
}
```

`defer s.mu.Unlock()` は「関数を抜ける直前に必ず実行される」予約です。
早期 return があっても解除し忘れません。

## ステータスコードの使い分け（REST の慣習）

| 操作 | 成功時 | 見つからない時 |
| --- | --- | --- |
| POST（作成） | 201 Created | - |
| GET（取得・一覧） | 200 OK | 404 Not Found |
| PUT（更新） | 200 OK | 404 Not Found |
| DELETE（削除） | 204 No Content（ボディなし） | 404 Not Found |

## strconv でパスパラメータを数値に変換

```go
id, err := strconv.Atoi(r.PathValue("id"))
if err != nil {
    // "abc" のような不正な id → 400
}
```

```sh
go test -v ./lessons/08_rest_crud/
```
