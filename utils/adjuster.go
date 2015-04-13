package utils

// An Adjuster is a function which takes a floating point value, most likely
// between 0 and 1, and returns another floating point value, again between 0
// and 1 (though this is not a requirement).
type Adjuster func(float64) float64

// Adder returns an Adjuster which adds the value with to the returned Adjusters
// argument.
func Adder(with float64) Adjuster {
	return func(i float64) float64 {
		return i + with
	}
}

// Multiplier returns an Adjuster which multiplies a given value with ratio.
func Multiplier(ratio float64) Adjuster {
	return func(i float64) float64 {
		return i * ratio
	}
}
