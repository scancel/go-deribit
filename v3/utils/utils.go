package utils

import (
	"math"
	"time"
)

// ----------
// Utilities
// ----------

func StrPointer(str string) *string {
	return &str
}

func FloatPointer(pFloat float64) *float64 {
	return &pFloat
}

// GetNumber will return NaN or the number, depending on the pointer
func GetNumber(pNumber *float64) float64 {
	var output float64
	output = math.NaN()
	if pNumber != nil {
		output = *pNumber
	}
	return output
}

// HumanReadableTime takes a timestamp in microiseconds and converts it a human-readable time.
func HumanReadableTime(pTimestamp int64) string {
	tm := time.Unix(pTimestamp/1000, 0)
	output := tm.Format("2006-01-02 15:04:05")
	return output
}
