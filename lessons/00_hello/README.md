# 00: Hello, World

まずは Go のプログラムを動かしてみましょう。

```sh
go run ./lessons/00_hello/
```

`Hello, World!` と表示されれば OK です。

## main.go の読み方

```go
package main        // 実行できるプログラムは必ず package main

import "fmt"        // 標準ライブラリの fmt（表示用）を使う宣言

func main() {       // プログラムはここから始まる
    fmt.Println("Hello, World!")
}
```

- `package main` と `func main()` の組み合わせが「実行可能なプログラム」の目印です
- `fmt.Println` は引数を表示して改行します
- Go ではセミコロン `;` は書きません
- 使っていない `import` があるとコンパイルエラーになります（Go は厳しめ）

## やってみよう（テストなし）

1. 表示する文字列を自分の名前に変えて `go run` してみる
2. `fmt.Println("1 + 2 =", 1+2)` を追加してみる（カンマ区切りで複数表示できる）
3. `import "fmt"` を消して `go run` し、エラーメッセージを読んでみる
