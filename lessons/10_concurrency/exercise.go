package concurrency

import "sync"

// Counter はスレッドセーフなカウンタ
type Counter struct {
	mu    sync.Mutex
	value int
}

// 課題1: value をロックしてから +1 する
func (c *Counter) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
}

// 課題2: 現在の value をロックしてから返す
func (c *Counter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

// 課題3: 複数の URL を並行に fetch し、結果を「入力と同じ順番」で返す
// ヒント: 結果を格納する []string をあらかじめ len(urls) で作り、
// 各 goroutine は自分の index に書き込む（index ごとに場所が違うので Mutex は不要）
// sync.WaitGroup で全部の完了を待つこと
func FetchAll(urls []string, fetch func(string) string) []string {
	results := make([]string, len(urls))
	var wg sync.WaitGroup
	for i, u := range urls {
		wg.Add(1)
		go func(i int, u string) {
			defer wg.Done()
			results[i] = fetch(u)
		}(i, u)
	}
	wg.Wait()
	return results
}

// 課題4: ワーカープールで jobs を処理し、結果の合計を返す
// numWorkers 個の goroutine で jobs を分担して process(job) を実行し、
// 全部の結果を合計して返す
// ヒント: jobs 用と results 用の channel を作る。jobs は分配し終えたら close する
func SumWithWorkers(jobs []int, numWorkers int, process func(int) int) int {
	jobsCh := make(chan int, len(jobs))
	resultsCh := make(chan int, len(jobs))

	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobsCh {
				resultsCh <- process(j)
			}
		}()
	}

	for _, j := range jobs {
		jobsCh <- j
	}
	close(jobsCh)

	wg.Wait()
	close(resultsCh)

	sum := 0
	for r := range resultsCh {
		sum += r
	}
	return sum
}
