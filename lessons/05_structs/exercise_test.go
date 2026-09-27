package structs

import (
	"slices"
	"testing"
)

func TestRect(t *testing.T) {
	r := Rect{Width: 3, Height: 4}
	if got := r.Area(); got != 12 {
		t.Errorf("Area() = %v, want 12", got)
	}
	if got := r.Perimeter(); got != 14 {
		t.Errorf("Perimeter() = %v, want 14", got)
	}
}

func TestAccount(t *testing.T) {
	a := NewAccount("alice")
	if a == nil {
		t.Fatal("NewAccount returned nil")
	}
	if a.Owner != "alice" || a.Balance != 0 {
		t.Fatalf("NewAccount(\"alice\") = %+v, want Owner=alice Balance=0", *a)
	}

	a.Deposit(1000)
	if a.Balance != 1000 {
		t.Errorf("after Deposit(1000), Balance = %d, want 1000", a.Balance)
	}

	if err := a.Withdraw(300); err != nil {
		t.Errorf("Withdraw(300) returned error: %v", err)
	}
	if a.Balance != 700 {
		t.Errorf("after Withdraw(300), Balance = %d, want 700", a.Balance)
	}

	if err := a.Withdraw(5000); err == nil {
		t.Errorf("Withdraw(5000) with balance 700 should return an error")
	}
	if a.Balance != 700 {
		t.Errorf("failed Withdraw must not change Balance, got %d", a.Balance)
	}
}

func TestTodoList(t *testing.T) {
	var l TodoList
	l.Add("buy milk")
	l.Add("learn go")
	l.Add("sleep")

	if len(l.Tasks) != 3 {
		t.Fatalf("after 3 Adds, len(Tasks) = %d, want 3", len(l.Tasks))
	}
	if err := l.Complete(1); err != nil {
		t.Errorf("Complete(1) returned error: %v", err)
	}
	if !l.Tasks[1].Done {
		t.Errorf("Tasks[1].Done = false after Complete(1)")
	}
	if got, want := l.Pending(), []string{"buy milk", "sleep"}; !slices.Equal(got, want) {
		t.Errorf("Pending() = %v, want %v", got, want)
	}
	for _, bad := range []int{-1, 3, 100} {
		if err := l.Complete(bad); err == nil {
			t.Errorf("Complete(%d) should return an error", bad)
		}
	}
}
