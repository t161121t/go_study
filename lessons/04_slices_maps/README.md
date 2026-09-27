# 04: スライスとマップ

## スライス（可変長の配列）

```go
nums := []int{1, 2, 3}       // 作成
nums = append(nums, 4)       // 末尾に追加（戻り値を受け取るのを忘れずに！）
fmt.Println(nums[0])         // 先頭の要素 → 1
fmt.Println(len(nums))       // 長さ → 4
part := nums[1:3]            // 1番目から3番目の手前まで → [2 3]

var empty []int              // ゼロ値は nil（長さ0のスライスとして使える）
s := make([]int, 0, 10)      // 長さ0・容量10で作成（大きさが分かるなら効率的）
```

### range でループ

```go
for i, v := range nums {     // i: インデックス, v: 値
    fmt.Println(i, v)
}
for _, v := range nums {     // インデックス不要なら _
    fmt.Println(v)
}
```

### 注意: スライスは中身を共有する

```go
a := []int{1, 2, 3}
b := a        // コピーではなく「同じ中身」を指す
b[0] = 99
fmt.Println(a) // [99 2 3]
```

## マップ（キーと値の組）

```go
ages := map[string]int{
    "alice": 20,
    "bob":   25,
}
ages["carol"] = 30           // 追加・更新
delete(ages, "bob")          // 削除
fmt.Println(ages["alice"])   // 20
fmt.Println(ages["nobody"])  // 存在しないキーはゼロ値 → 0

age, ok := ages["dave"]      // ok で存在チェック
if !ok {
    fmt.Println("いない")
}

counts := make(map[string]int)   // 空のマップを作る
counts["go"]++                   // ゼロ値に +1 できる（カウントに便利）

for key, value := range ages {   // ループ（順番はランダム！）
    fmt.Println(key, value)
}
```

⚠️ `var m map[string]int` のように宣言だけした nil マップに書き込むと panic します。`make` か `{}` で作りましょう。

## 文字列を単語に分ける

```go
import "strings"
words := strings.Fields("go is fun")   // → []string{"go", "is", "fun"}
```

```sh
go test -v ./lessons/04_slices_maps/
```
