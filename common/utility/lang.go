package utility

// ? :
func If[T any](condition bool, valueIfTrue T, valueIfFalse T) T {
	if condition {
		return valueIfTrue
	} else {
		return valueIfFalse
	}
}
