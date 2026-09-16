package calculator

import (
	"errors"
	"math"
	"testing"
)

func TestAdd(t *testing.T) {
	calc := New(2)

	tests := []struct {
		name     string
		a, b     float64
		expected float64
	}{
		{"positive numbers", 2, 3, 5},
		{"negative numbers", -2, -3, -5},
		{"mixed signs", -5, 8, 3},
		{"floating point noise", 0.1, 0.2, 0.3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := calc.Add(tt.a, tt.b)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("Add(%v, %v) = %v, want %v", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestSubtract(t *testing.T) {
	calc := New(2)

	tests := []struct {
		name     string
		a, b     float64
		expected float64
	}{
		{"positive result", 5, 3, 2},
		{"negative result", 3, 5, -2},
		{"subtract zero", 7, 0, 7},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := calc.Subtract(tt.a, tt.b)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("Subtract(%v, %v) = %v, want %v", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestMultiply(t *testing.T) {
	calc := New(2)

	tests := []struct {
		name     string
		a, b     float64
		expected float64
	}{
		{"positive numbers", 4, 5, 20},
		{"multiply by zero", 100, 0, 0},
		{"negative times positive", -3, 4, -12},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := calc.Multiply(tt.a, tt.b)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("Multiply(%v, %v) = %v, want %v", tt.a, tt.b, result, tt.expected)
			}
		})
	}

	t.Run("overflow is reported as a domain error, not Inf", func(t *testing.T) {
		_, err := calc.Multiply(math.MaxFloat64, math.MaxFloat64)
		if !errors.Is(err, ErrResultNotFinite) {
			t.Fatalf("expected ErrResultNotFinite, got %v", err)
		}
	})
}

func TestDivide(t *testing.T) {
	calc := New(2)

	tests := []struct {
		name     string
		a, b     float64
		expected float64
	}{
		{"exact division", 10, 2, 5},
		{"fractional result", 10, 3, 3.33},
		{"negative divisor", 10, -2, -5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := calc.Divide(tt.a, tt.b)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("Divide(%v, %v) = %v, want %v", tt.a, tt.b, result, tt.expected)
			}
		})
	}

	t.Run("division by zero is a domain error", func(t *testing.T) {
		_, err := calc.Divide(10, 0)
		if !errors.Is(err, ErrDivisionByZero) {
			t.Fatalf("expected ErrDivisionByZero, got %v", err)
		}
	})
}

func TestPower(t *testing.T) {
	calc := New(2)

	tests := []struct {
		name     string
		base     float64
		exponent float64
		expected float64
	}{
		{"positive exponent", 2, 3, 8},
		{"exponent zero", 5, 0, 1},
		{"zero base, zero exponent convention", 0, 0, 1},
		{"negative exponent", 2, -1, 0.5},
		{"fractional exponent", 4, 0.5, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := calc.Power(tt.base, tt.exponent)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("Power(%v, %v) = %v, want %v", tt.base, tt.exponent, result, tt.expected)
			}
		})
	}

	t.Run("zero base with negative exponent is a domain error", func(t *testing.T) {
		_, err := calc.Power(0, -1)
		if !errors.Is(err, ErrDivisionByZero) {
			t.Fatalf("expected ErrDivisionByZero, got %v", err)
		}
	})

	t.Run("overflow is reported as a domain error, not Inf", func(t *testing.T) {
		_, err := calc.Power(10, 1000)
		if !errors.Is(err, ErrResultNotFinite) {
			t.Fatalf("expected ErrResultNotFinite, got %v", err)
		}
	})
}

func TestSqrt(t *testing.T) {
	calc := New(2)

	tests := []struct {
		name     string
		value    float64
		expected float64
	}{
		{"perfect square", 4, 2},
		{"zero", 0, 0},
		{"non-perfect square", 2, 1.41},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := calc.Sqrt(tt.value)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("Sqrt(%v) = %v, want %v", tt.value, result, tt.expected)
			}
		})
	}

	t.Run("negative input is a domain error", func(t *testing.T) {
		_, err := calc.Sqrt(-4)
		if !errors.Is(err, ErrNegativeSqrtInput) {
			t.Fatalf("expected ErrNegativeSqrtInput, got %v", err)
		}
	})
}

func TestPercentage(t *testing.T) {
	calc := New(2)

	tests := []struct {
		name       string
		value      float64
		percentage float64
		expected   float64
	}{
		{"15 percent of 200", 200, 15, 30},
		{"zero percentage", 100, 0, 0},
		{"percentage over 100", 50, 200, 100},
		{"negative value", -100, 10, -10},
		{"negative percentage", 100, -10, -10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := calc.Percentage(tt.value, tt.percentage)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("Percentage(%v, %v) = %v, want %v", tt.value, tt.percentage, result, tt.expected)
			}
		})
	}
}

func TestRoundingPrecisionIsConfigurable(t *testing.T) {
	tests := []struct {
		name      string
		precision int
		expected  float64
	}{
		{"zero decimals", 0, 3},
		{"two decimals", 2, 3.33},
		{"four decimals", 4, 3.3333},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calc := New(tt.precision)
			result, err := calc.Divide(10, 3)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("Divide(10, 3) with precision %d = %v, want %v", tt.precision, result, tt.expected)
			}
		})
	}
}

func TestRoundHalfAwayFromZero(t *testing.T) {
	calc := New(0)

	result, err := calc.Add(2.5, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != 3 {
		t.Errorf("expected round-half-away-from-zero to round 2.5 to 3, got %v", result)
	}

	result, err = calc.Add(-2.5, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != -3 {
		t.Errorf("expected round-half-away-from-zero to round -2.5 to -3, got %v", result)
	}
}
