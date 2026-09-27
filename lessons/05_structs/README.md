# 05: 構造体とメソッド

## 構造体（struct）

関連するデータをひとまとめにする型です。Go にはクラスがなく、構造体 + メソッドで表現します。

```go
type User struct {
    Name string
    Age  int
}

u := User{Name: "Alice", Age: 20}   // 作成
fmt.Println(u.Name)                  // フィールドにアクセス
u.Age = 21                           // 書き換え
var empty User                       // ゼロ値: Name="", Age=0
```

名前が **大文字で始まる** もの（`User`, `Name`）はパッケージの外から使えます（public）。小文字なら外から見えません（private）。

## メソッド

型に紐づいた関数です。`func` と関数名の間に **レシーバ** を書きます。

```go
func (u User) Greet() string {       // 値レシーバ: u はコピー
    return "Hi, " + u.Name
}

func (u *User) Birthday() {          // ポインタレシーバ: 元の値を書き換えられる
    u.Age++
}

u := User{Name: "Bob", Age: 30}
u.Greet()      // "Hi, Bob"
u.Birthday()   // u.Age が 31 になる
```

**使い分け**: フィールドを変更したいメソッドは `*User`（ポインタレシーバ）にします。
値レシーバで `u.Age++` しても、コピーが変わるだけで元は変わりません。

## ポインタ

```go
x := 10
p := &x     // & で「x の場所（アドレス）」を取る
*p = 20     // * で「その場所の中身」を書き換える
fmt.Println(x)  // 20
```

## コンストラクタ関数

Go には `new User()` のようなコンストラクタ構文はないので、`NewXxx` という関数を作るのが慣習です。

```go
func NewUser(name string) *User {
    return &User{Name: name, Age: 0}
}
```

## 構造体のスライス

```go
users := []User{{Name: "A"}, {Name: "B"}}
for i := range users {
    users[i].Age++           // for _, u := range だと u はコピーなので注意
}
```

```sh
go test -v ./lessons/05_structs/
```
