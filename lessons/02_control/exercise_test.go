package control

import "testing"

func TestIsEven(t *testing.T) {
	cases := map[int]bool{0: true, 1: false, 2: true, 7: false, -4: true}
	for n, want := range cases {
		if got := IsEven(n); got != want {
			t.Errorf("IsEven(%d) = %v, want %v", n, got, want)
		}
	}
}

func TestMax(t *testing.T) {
	if got := Max(3, 8); got != 8 {
		t.Errorf("Max(3, 8) = %d, want 8", got)
	}
	if got := Max(-1, -5); got != -1 {
		t.Errorf("Max(-1, -5) = %d, want -1", got)
	}
}

func TestSumTo(t *testing.T) {
	cases := map[int]int{1: 1, 4: 10, 10: 55, 100: 5050}
	for n, want := range cases {
		if got := SumTo(n); got != want {
			t.Errorf("SumTo(%d) = %d, want %d", n, got, want)
		}
	}
}

func TestGrade(t *testing.T) {
	cases := map[int]string{100: "A", 80: "A", 79: "B", 60: "B", 45: "C", 40: "C", 39: "D", 0: "D"}
	for score, want := range cases {
		if got := Grade(score); got != want {
			t.Errorf("Grade(%d) = %q, want %q", score, got, want)
		}
	}
}

func TestFizzBuzz(t *testing.T) {
	cases := map[int]string{1: "1", 3: "Fizz", 5: "Buzz", 9: "Fizz", 10: "Buzz", 15: "FizzBuzz", 30: "FizzBuzz", 22: "22"}
	for n, want := range cases {
		if got := FizzBuzz(n); got != want {
			t.Errorf("FizzBuzz(%d) = %q, want %q", n, got, want)
		}
	}
}

func TestCountPrimes(t *testing.T) {
	cases := map[int]int{1: 0, 2: 1, 10: 4, 30: 10}
	for n, want := range cases {
		if got := CountPrimes(n); got != want {
			t.Errorf("CountPrimes(%d) = %d, want %d", n, got, want)
		}
	}
}
