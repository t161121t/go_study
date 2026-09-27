# 03: 関数・複数の戻り値・エラー

## 引数の型をまとめる

```go
func Add(a, b int) int {   // a int, b int と同じ
    return a + b
}
```

## 複数の戻り値

Go の関数は値を複数返せます。

```go
func Divmod(a, b int) (int, int) {
    return a / b, a % b
}

q, r := Divmod(7, 2)   // q = 3, r = 1
_, r2 := Divmod(9, 4)  // 使わない値は _ で捨てる
```

## エラー処理（Go で一番大事な書き方）

Go には例外（try/catch）がありません。**最後の戻り値で `error` を返す** のが決まりです。

```go
import (
    "errors"
    "fmt"
)

func Divide(a, b int) (int, error) {
    if b == 0 {
        return 0, errors.New("division by zero")   // 失敗: エラーを返す
    }
    return a / b, nil                              // 成功: エラーは nil
}

result, err := Divide(10, 0)
if err != nil {                // 呼び出し側は必ずチェックする
    fmt.Println("エラー:", err)
    return
}
fmt.Println(result)
```

- `nil` は「何もない」を表す値。`err == nil` なら成功
- 値を埋め込んだエラーを作るなら `fmt.Errorf("invalid age: %d", age)`

## 関数も値

関数を変数に入れたり、引数として渡したりできます。

```go
double := func(x int) int { return x * 2 }
fmt.Println(double(3))   // 6

func Apply(x int, f func(int) int) int {
    return f(x)
}
Apply(5, double)         // 10
```

## 可変長引数

```go
func Sum(nums ...int) int {   // nums はスライス（次のレッスンで詳しく）
    total := 0
    for _, n := range nums {
        total += n
    }
    return total
}
Sum(1, 2, 3)   // 6
```

```sh
go test -v ./lessons/03_functions/
```
