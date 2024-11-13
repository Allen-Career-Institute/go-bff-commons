// nolint: gocritic
package utils

import (
	"math"
	"strconv"
	"strings"
)

type Val1 interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64 | ~int | ~int64
}

type Val2 interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64 | ~int | ~int64
}

// GetMin takes two comparable arguments and finds the minimum. In case the types of both the arguments are
// different, it returns the type of the first argument
func GetMin[T Val1, V Val2](val1 T, val2 V) T {
	return T(math.Min(float64(val1), float64(val2)))
}

// GetMax takes two comparable arguments and finds the maximum. In case the types of both the arguments are
// different, it returns the type of the first argument
func GetMax[T Val1, V Val2](val1 T, val2 V) T {
	return T(math.Max(float64(val1), float64(val2)))
}

// GetFormattedNumber returns the number as string in a clean format with at most 1 precision
func GetFormattedNumber(number float32) string {
	formattedNumber := strconv.FormatFloat(float64(number), 'f', 1, 64)
	if strings.HasSuffix(formattedNumber, ".0") {
		return formattedNumber[:strings.LastIndex(formattedNumber, ".")]
	}

	return formattedNumber
}
