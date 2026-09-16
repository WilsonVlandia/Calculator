// Package calculator implements the arithmetic model of the API. It has
// no dependency on net/http or any other transport concern: it only
// knows how to compute and round numbers.
package calculator

import "math"

// Calculator performs arithmetic operations and rounds every result to
// a fixed decimal precision.
type Calculator struct {
	precision int
}

// New creates a Calculator that rounds every result to the given number
// of decimal places.
func New(precision int) *Calculator {
	return &Calculator{precision: precision}
}

// finalize rounds a computed value and rejects results that are not
// finite (e.g. overflow into +/-Inf, or NaN).
func (c *Calculator) finalize(value float64) (float64, error) {
	if math.IsInf(value, 0) || math.IsNaN(value) {
		return 0, ErrResultNotFinite
	}
	return round(value, c.precision), nil
}

// Add returns a + b.
func (c *Calculator) Add(a, b float64) (float64, error) {
	return c.finalize(a + b)
}

// Subtract returns a - b.
func (c *Calculator) Subtract(a, b float64) (float64, error) {
	return c.finalize(a - b)
}

// Multiply returns a * b.
func (c *Calculator) Multiply(a, b float64) (float64, error) {
	return c.finalize(a * b)
}

// Divide returns a / b. Dividing by zero is a domain error, not a
// panic or an Inf result.
func (c *Calculator) Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, ErrDivisionByZero
	}
	return c.finalize(a / b)
}

// Power returns base raised to exponent. base == 0 with a negative
// exponent is treated as division by zero. base == 0, exponent == 0
// follows the standard math.Pow convention and evaluates to 1.
func (c *Calculator) Power(base, exponent float64) (float64, error) {
	if base == 0 && exponent < 0 {
		return 0, ErrDivisionByZero
	}
	return c.finalize(math.Pow(base, exponent))
}

// Sqrt returns the square root of value. Negative inputs are a domain
// error, not a NaN result.
func (c *Calculator) Sqrt(value float64) (float64, error) {
	if value < 0 {
		return 0, ErrNegativeSqrtInput
	}
	return c.finalize(math.Sqrt(value))
}

// Percentage returns percentage% of value (e.g. Percentage(200, 15)
// returns 15% of 200, i.e. 30), matching the % button of a standard
// calculator.
func (c *Calculator) Percentage(value, percentage float64) (float64, error) {
	return c.finalize(value * (percentage / 100))
}
