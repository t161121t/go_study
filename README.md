# Go 学習リポジトリ

`lessons/` 以下を番号順に進めていきます。各レッスンは次の3ファイルでできています。

| ファイル | 役割 |
| --- | --- |
| `README.md` | そのレッスンの解説。まずここを読む |
| `exercise.go` | 課題。`// TODO` の部分を自分で実装する |
| `exercise_test.go` | 採点用テスト。**編集しない** |

## 進め方

```sh
# 1. 解説を読む
cat lessons/01_basics/README.md

# 2. exercise.go の TODO を実装する

# 3. テストを実行して採点（最初は全部 FAIL するのが正常）
go test ./lessons/01_basics/

# 詳しく見たいとき
go test -v ./lessons/01_basics/

# 全レッスンまとめて
go test ./lessons/...
```

`ok` と表示されたらクリア。次のレッスンへ進みましょう。

## レッスン一覧

### 基礎文法

- [00_hello](lessons/00_hello/) — プログラムを動かしてみる（`go run`）
- [01_basics](lessons/01_basics/) — 変数・型・定数
- [02_control](lessons/02_control/) — if / for / switch
- [03_functions](lessons/03_functions/) — 関数・複数の戻り値・エラー
- [04_slices_maps](lessons/04_slices_maps/) — スライスとマップ
- [05_structs](lessons/05_structs/) — 構造体とメソッド

### 実戦編: Web API サーバー

- [06_http_basics](lessons/06_http_basics/) — net/http・ルーティング・httptest
- [07_json_api](lessons/07_json_api/) — JSONのエンコード/デコード・バリデーション
- [08_rest_crud](lessons/08_rest_crud/) — リポジトリパターンでのCRUD API・sync.RWMutex
- [09_middleware](lessons/09_middleware/) — ミドルウェア・認証・context によるタイムアウト
- [10_concurrency](lessons/10_concurrency/) — goroutine・channel・ワーカープール
- [11_capstone_bookmarks_api](lessons/11_capstone_bookmarks_api/) — 総仕上げ課題（採点テストなし、自分で設計・実装）

## 困ったら

- `go fmt ./...` でコードを自動整形できます
- 公式チュートリアル: https://go.dev/tour/ （日本語版: https://go-tour-jp.appspot.com/）
