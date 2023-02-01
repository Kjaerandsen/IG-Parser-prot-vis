package shared

import "strings"

/*
Replaces/escapes selected symbols in as far as relevant for export (e.g., quotation marks).
*/
func EscapeSymbolsForExport(rawValue string) string {
	// Replace quotation marks with single quotes
	return strings.ReplaceAll(rawValue, "\"", "'")
}

/*
Aggregates values in a given array that exceed a given threshold value.
Returns sum, if it exceeds a given default value. Otherwise it returns the given default value.
*/
func AggregateIfGreaterThan(arr []int, threshold int, defaultValue int) int {
	sum := 0
	for i := 0; i < len(arr); i++ {
		if arr[i] > threshold {
			// Only aggregate values above threshold
			sum += arr[i]
		}
	}
	if sum > defaultValue {
		// Only return sum if greater than default value
		return sum
	} else {
		return defaultValue
	}
}

/*
Return the maximum value for all values in a given array.
If none of the values is higher, it returns a given default value.
*/
func FindMaxValue(arr []int, defaultValue int) int {
	max := 0
	for i := 0; i < len(arr); i++ {
		if arr[i] > max {
			max += arr[i]
		}
	}
	if max > defaultValue {
		return max
	} else {
		return defaultValue
	}
}