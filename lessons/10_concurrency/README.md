# 10: goroutine と channel（並行処理）

APIサーバーでは「複数の外部APIを同時に呼ぶ」「複数リクエストを裏で処理する」といった
場面で並行処理が必須になります。

## goroutine（軽量スレッド）

`go` を関数呼び出しの前に付けるだけで、別の goroutine として実行されます。

```go
go doSomething()   // 呼び出し元はここで待たずに次の行へ進む
```

## WaitGroup で完了を待つ

```go
var wg sync.WaitGroup
for _, url := range urls {
    wg.Add(1)                 // これから1つ増える
    go func(u string) {
        defer wg.Done()       // 終わったら1つ減らす
        fetch(u)
    }(url)                    // ⚠️ url をそのまま使うと全goroutineが同じ変数を参照してしまう罠がある
}
wg.Wait()                     // 全部の Done() が呼ばれるまでブロック
```

Go 1.22 以降は for ループの変数がイテレーションごとに新しくなったので、
`for _, url := range urls { go func(){ fetch(url) }() }` と書いても実は安全です。
ただし引数として渡す書き方も広く使われるので両方見て慣れておきましょう。

## Mutex で共有データを守る

複数の goroutine が同じ変数を読み書きするなら `sync.Mutex` で保護します
（08 で使った `sync.RWMutex` はこの発展形）。

```go
type Counter struct {
    mu    sync.Mutex
    value int
}

func (c *Counter) Inc() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.value++
}
```

保護しないとどうなるか: `go test -race` を付けて実行すると、データ競合を検出してくれます。

```sh
go test -race ./lessons/10_concurrency/
```

## channel でやり取りする

channel は goroutine 同士がデータを送り合うための「パイプ」です。

```go
ch := make(chan int)     // バッファなしの channel
go func() {
    ch <- 42              // 送信（受け取られるまでブロック）
}()
v := <-ch                // 受信
```

## ワーカープール（channel + goroutine の定番パターン）

決まった数の goroutine（ワーカー）で、大量のジョブを分担して処理します。

```go
jobs := make(chan int, len(items))
results := make(chan int, len(items))

for w := 0; w < numWorkers; w++ {
    go func() {
        for j := range jobs {      // jobs が close されるまで受け取り続ける
            results <- process(j)
        }
    }()
}
for _, item := range items {
    jobs <- item
}
close(jobs)   // これ以上送らないことを伝える → ワーカーの for range が終わる
```

```sh
go test -v ./lessons/10_concurrency/
```
