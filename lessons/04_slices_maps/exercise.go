package slicesmaps

// 課題1: スライスの合計を返す
func Sum(nums []int) int {
	// TODO
	return 0
}

// 課題2: スライスの最大値を返す。空ならエラー... ではなく、ここでは ok=false を返す
// 例: MaxOf([]int{3, 9, 2}) → 9, true
//
//	MaxOf([]int{})        → 0, false
func MaxOf(nums []int) (int, bool) {
	// TODO
	return 0, false
}

// 課題3: 偶数だけを取り出した新しいスライスを返す（元の順番を保つ）
// 例: Evens([]int{1, 2, 3, 4}) → []int{2, 4}
func Evens(nums []int) []int {
	// TODO
	return nil
}

// 課題4: 順番を逆にした「新しい」スライスを返す（引数のスライスは変更しない）
// 例: Reverse([]int{1, 2, 3}) → []int{3, 2, 1}
func Reverse(nums []int) []int {
	// TODO
	return nil
}

// 課題5: 文章に含まれる各単語の出現回数を数える
// 例: WordCount("go is fun go") → map[string]int{"go": 2, "is": 1, "fun": 1}
// ヒント: strings.Fields
func WordCount(text string) map[string]int {
	// TODO
	return nil
}

// 課題6: 重複を取り除いたスライスを返す（最初に出てきた順番を保つ）
// 例: Unique([]string{"a", "b", "a", "c", "b"}) → []string{"a", "b", "c"}
// ヒント: 「もう見たか」をマップで記録する
func Unique(items []string) []string {
	// TODO
	return nil
}
