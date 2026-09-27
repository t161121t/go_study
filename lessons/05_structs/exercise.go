package structs

// 課題1: 長方形
type Rect struct {
	Width  float64
	Height float64
}

// 面積を返す
func (r Rect) Area() float64 {
	// TODO
	return 0
}

// 周の長さを返す
func (r Rect) Perimeter() float64 {
	// TODO
	return 0
}

// 課題2: 銀行口座
type Account struct {
	Owner   string
	Balance int
}

// 残高 0 の口座を作って返す
func NewAccount(owner string) *Account {
	// TODO
	return nil
}

// amount を入金する
// ヒント: 残高を書き換えるのでポインタレシーバ (*Account) になっている
func (a *Account) Deposit(amount int) {
	// TODO
}

// amount を出金する。残高が足りなければ残高は変えずにエラーを返す
// ヒント: import "errors" して errors.New("...") でエラーを作る
func (a *Account) Withdraw(amount int) error {
	// TODO
	return nil
}

// 課題3: TODO リスト
type Task struct {
	Title string
	Done  bool
}

type TodoList struct {
	Tasks []Task
}

// タイトル title の未完了タスクを末尾に追加する
func (l *TodoList) Add(title string) {
	// TODO
}

// index 番目のタスクを完了にする。範囲外ならエラーを返す
func (l *TodoList) Complete(index int) error {
	// TODO
	return nil
}

// 未完了タスクのタイトル一覧を返す（追加した順）
func (l *TodoList) Pending() []string {
	// TODO
	return nil
}
