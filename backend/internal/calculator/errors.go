package calculator

import "errors"

// ErrDivisionByZero, ErrNegativeSqrtInput and ErrResultNotFinite are the
// sentinel errors for the math domain cases each operation can hit.
// Handlers use errors.Is against these to pick the right HTTP status
// and error code, without the model knowing anything about HTTP.
var (
	// ErrDivisionByZero is returned when an operation would divide by
	// zero (Divide, or Power with a zero base and a negative exponent).
	ErrDivisionByZero = errors.New("division by zero")
	// ErrNegativeSqrtInput is returned when Sqrt is called with a
	// negative value.
	ErrNegativeSqrtInput = errors.New("cannot compute the square root of a negative number")
	// ErrResultNotFinite is returned when a computed result is not a
	// finite number (e.g. it overflows to +/-Inf or is NaN).
	ErrResultNotFinite = errors.New("result is not a finite number")
)
