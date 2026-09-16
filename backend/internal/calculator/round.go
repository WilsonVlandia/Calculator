package calculator

import "math"

// round applies "round half away from zero" at the given number of
// decimal places. math.Round already implements that rule, so rounding
// is done by scaling the value up, rounding to the nearest integer, and
// scaling back down. A negative precision is treated as zero.
func round(value float64, precision int) float64 {
	if precision < 0 {
		precision = 0
	}
	factor := math.Pow(10, float64(precision))
	return math.Round(value*factor) / factor
}
