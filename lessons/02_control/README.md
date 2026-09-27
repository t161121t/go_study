# 02: if / for / switch

## if

条件を `()` で囲みません。`{}` は必須です。

```go
if x > 10 {
    // ...
} else if x > 5 {
    // ...
} else {
    // ...
}
```

比較: `==` `!=` `<` `<=` `>` `>=`　論理: `&&`（かつ） `||`（または） `!`（否定）

## for（Go のループはこれだけ。while はない）

```go
// 基本形
for i := 0; i < 5; i++ {
    fmt.Println(i)   // 0, 1, 2, 3, 4
}

// while のような書き方
n := 1
for n < 100 {
    n *= 2
}

// 無限ループ（break で抜ける）
for {
    if done {
        break
    }
}

// 回数だけ回す（Go 1.22 以降）
for i := range 5 {
    fmt.Println(i)   // 0〜4
}
```

`continue` でそのループの残りをスキップして次へ進みます。

## switch

`break` を書かなくても、一致した case だけ実行されます。

```go
switch day {
case "Sat", "Sun":
    return "休日"
default:
    return "平日"
}

// 条件式を書くこともできる
switch {
case score >= 80:
    return "A"
case score >= 60:
    return "B"
}
```

## 割り算の余り

`%` で余りが求まります。`n%2 == 0` なら偶数です。

```sh
go test -v ./lessons/02_control/
```
