package concurrency

import (
	"sync"
	"testing"
	"time"
)

func TestCounter(t *testing.T) {
	c := &Counter{}
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.Inc()
		}()
	}
	wg.Wait()
	if got := c.Value(); got != 100 {
		t.Errorf("Value() = %d, want 100 (run with -race to check for data races)", got)
	}
}

func TestFetchAll(t *testing.T) {
	urls := []string{"a", "b", "c", "d", "e"}
	fetch := func(u string) string {
		time.Sleep(10 * time.Millisecond)
		return "fetched:" + u
	}

	start := time.Now()
	got := FetchAll(urls, fetch)
	elapsed := time.Since(start)

	want := []string{"fetched:a", "fetched:b", "fetched:c", "fetched:d", "fetched:e"}
	if len(got) != len(want) {
		t.Fatalf("len(got) = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("got[%d] = %q, want %q (order must match input)", i, got[i], want[i])
		}
	}

	// 並行に実行されていれば、5件×10msでも合計は100ms未満のはず
	if elapsed > 60*time.Millisecond {
		t.Errorf("elapsed = %v, want < 60ms (are you fetching sequentially instead of concurrently?)", elapsed)
	}
}

func TestSumWithWorkers(t *testing.T) {
	jobs := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	double := func(n int) int { return n * 2 }

	got := SumWithWorkers(jobs, 3, double)
	want := 2 * (10 * 11 / 2) // sum(1..10)*2 = 110
	if got != want {
		t.Errorf("SumWithWorkers = %d, want %d", got, want)
	}
}
