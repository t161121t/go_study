package functions

import "testing"

func TestMinMax(t *testing.T) {
	if lo, hi := MinMax(7, 3); lo != 3 || hi != 7 {
		t.Errorf("MinMax(7, 3) = %d, %d, want 3, 7", lo, hi)
	}
	if lo, hi := MinMax(2, 9); lo != 2 || hi != 9 {
		t.Errorf("MinMax(2, 9) = %d, %d, want 2, 9", lo, hi)
	}
}

func TestDivide(t *testing.T) {
	got, err := Divide(10, 3)
	if err != nil || got != 3 {
		t.Errorf("Divide(10, 3) = %d, %v, want 3, nil", got, err)
	}
	got, err = Divide(1, 0)
	if err == nil {
		t.Errorf("Divide(1, 0) should return an error, got nil")
	}
	if got != 0 {
		t.Errorf("Divide(1, 0) value = %d, want 0", got)
	}
}

func TestValidateAge(t *testing.T) {
	for _, ok := range []int{0, 20, 150} {
		if err := ValidateAge(ok); err != nil {
			t.Errorf("ValidateAge(%d) = %v, want nil", ok, err)
		}
	}
	for _, ng := range []int{-1, 151, 999} {
		if err := ValidateAge(ng); err == nil {
			t.Errorf("ValidateAge(%d) should return an error", ng)
		}
	}
}

func TestAverage(t *testing.T) {
	got, err := Average(1, 2, 3, 4)
	if err != nil || got != 2.5 {
		t.Errorf("Average(1, 2, 3, 4) = %v, %v, want 2.5, nil", got, err)
	}
	if _, err := Average(); err == nil {
		t.Errorf("Average() should return an error")
	}
}

func TestApplyN(t *testing.T) {
	double := func(v int) int { return v * 2 }
	if got := ApplyN(double, 1, 3); got != 8 {
		t.Errorf("ApplyN(double, 1, 3) = %d, want 8", got)
	}
	if got := ApplyN(double, 5, 0); got != 5 {
		t.Errorf("ApplyN(double, 5, 0) = %d, want 5", got)
	}
}

func TestCounter(t *testing.T) {
	a := Counter()
	for want := 1; want <= 3; want++ {
		if got := a(); got != want {
			t.Errorf("counter a: call %d = %d, want %d", want, got, want)
		}
	}
	b := Counter() // 別のカウンターは独立している
	if got := b(); got != 1 {
		t.Errorf("new counter b first call = %d, want 1", got)
	}
}
