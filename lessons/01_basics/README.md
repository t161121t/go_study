# 01: 変数・型・定数

## 変数の宣言

```go
var age int = 20      // 基本形: var 名前 型 = 値
var name = "Gopher"   // 型は値から推論できる
count := 3            // 関数の中ではこの短い書き方が一般的（:= は「宣言して代入」）
count = 4             // 2回目以降の代入は = （:= ではない）
```

宣言だけして値を入れないと **ゼロ値** になります。

| 型 | ゼロ値 |
| --- | --- |
| `int`, `float64` | `0` |
| `string` | `""`（空文字） |
| `bool` | `false` |

## よく使う型

- `int` … 整数
- `float64` … 小数
- `string` … 文字列（`"..."` で書く）
- `bool` … `true` / `false`

## 型変換

Go は **暗黙の型変換をしません**。`int` と `float64` はそのまま足せません。

```go
a := 3
b := 1.5
c := float64(a) + b   // 型名(値) で変換する
d := int(b)           // 小数点以下は切り捨て → 1
```

## 定数

```go
const Pi = 3.14
```

## 文字列

```go
s := "Go" + "lang"    // + で連結
n := len(s)           // len でバイト数（英数字なら文字数と同じ）
```

## 関数の読み方（課題で使う）

```go
func Add(a int, b int) int {  // 引数 a, b（int）を受け取って int を返す
    return a + b
}
```

`exercise.go` には関数の「枠」だけあります。中身を書いて `return` を正しい値にしてください。

```sh
go test -v ./lessons/01_basics/
```
