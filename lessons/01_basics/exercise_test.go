package basics

import (
	"math"
	"testing"
)

func almostEqual(a, b float64) bool { return math.Abs(a-b) < 1e-9 }

func TestAdd(t *testing.T) {
	if got := Add(2, 3); got != 5 {
		t.Errorf("Add(2, 3) = %d, want 5", got)
	}
	if got := Add(-4, 1); got != -3 {
		t.Errorf("Add(-4, 1) = %d, want -3", got)
	}
}

func TestRectArea(t *testing.T) {
	if got := RectArea(2.5, 4); !almostEqual(got, 10) {
		t.Errorf("RectArea(2.5, 4) = %v, want 10", got)
	}
}

func TestAverage(t *testing.T) {
	if got := Average(3, 2.0); !almostEqual(got, 2.5) {
		t.Errorf("Average(3, 2.0) = %v, want 2.5", got)
	}
}

func TestGreet(t *testing.T) {
	if got := Greet("Gopher"); got != "Hello, Gopher!" {
		t.Errorf(`Greet("Gopher") = %q, want "Hello, Gopher!"`, got)
	}
}

func TestLength(t *testing.T) {
	if got := Length("golang"); got != 6 {
		t.Errorf(`Length("golang") = %d, want 6`, got)
	}
	if got := Length(""); got != 0 {
		t.Errorf(`Length("") = %d, want 0`, got)
	}
}

func TestCircleArea(t *testing.T) {
	if got := CircleArea(2); !almostEqual(got, 12.56) {
		t.Errorf("CircleArea(2) = %v, want 12.56", got)
	}
}
