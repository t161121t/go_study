package functions

// 課題1: 最小値と最大値を両方返す
// 例: MinMax(7, 3) → 3, 7
func MinMax(a, b int) (int, int) {
	// TODO
	return 0, 0
}

// 課題2: 割り算。b が 0 ならエラーを返す
// 成功時: 商と nil を返す
// 失敗時: 0 とエラーを返す（エラーメッセージは自由）
func Divide(a, b int) (int, error) {
	// TODO
	return 0, nil
}

// 課題3: 年齢を検証する
// 0〜150 の範囲なら nil、それ以外ならエラーを返す
func ValidateAge(age int) error {
	// TODO
	return nil
}

// 課題4: 可変長引数の平均を返す。引数が0個ならエラーを返す
// 例: Average(1, 2, 3, 4) → 2.5, nil
func Average(nums ...int) (float64, error) {
	// TODO
	return 0, nil
}

// 課題5: 関数 f を x に n 回適用した結果を返す
// 例: ApplyN(func(v int) int { return v * 2 }, 1, 3) → 8 （1→2→4→8）
func ApplyN(f func(int) int, x int, n int) int {
	// TODO
	return 0
}

// 課題6: 呼ぶたびに 1, 2, 3, ... を返す関数を返す（クロージャ）
// 例:
//
//	next := Counter()
//	next() → 1
//	next() → 2
//
// ヒント: 関数の外側の変数を、返す関数の中から書き換える
func Counter() func() int {
	// TODO
	return func() int { return 0 }
}
