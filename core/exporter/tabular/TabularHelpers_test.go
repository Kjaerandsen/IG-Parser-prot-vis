package tabular

import (
	"testing"
)

/*
Tests basic addition of element to array
*/
func TestAddingElement(t *testing.T) {

	arr := []string{"One", "Two"}

	arr = addElementIfNotExisting("Two", arr)

	if len(arr) != 2 {
		t.Error("Element should not have been added")
	}

	arr = addElementIfNotExisting("Three", arr)

	if len(arr) != 3 {
		t.Error("Element should have been added")
	}

}

/*
Test moving to last position, with and without addition during moving.
*/
func TestMoveElementToFirstPosition(t *testing.T) {

	arr := []string{"One", "Two", "Three"}

	arr = moveElementToFirstPosition("Three", arr, true)

	if len(arr) != 3 {
		t.Error("Element should not have been added")
	}

	if arr[0] != "Three" {
		t.Error("Element has not been moved to first position")
	}

	arr = moveElementToFirstPosition("Fourth", arr, true)

	if len(arr) != 4 {
		t.Error("Element should have been added")
	}

	if arr[0] != "Fourth" {
		t.Error("Element has not been moved to first position")
	}

	arr = moveElementToFirstPosition("Fifth", arr, false)

	if len(arr) != 4 {
		t.Error("Element should have been added")
	}

}

/*
Test moving to last position, with and without addition during moving.
*/
func TestMoveElementToLastPosition(t *testing.T) {

	arr := []string{"One", "Two", "Three"}

	arr = moveElementToLastPosition("Two", arr, true)

	if len(arr) != 3 {
		t.Error("Element should not have been added")
	}

	if arr[2] != "Two" {
		t.Error("Element has not been moved to last position")
	}

	arr = moveElementToLastPosition("Fourth", arr, true)

	if len(arr) != 4 {
		t.Error("Element should have been added")
	}

	if arr[3] != "Fourth" {
		t.Error("Element has not been moved to last position")
	}

	arr = moveElementToLastPosition("Two", arr, true)

	if len(arr) != 4 {
		t.Error("Element should have been added")
	}

	if arr[3] != "Two" {
		t.Error("Element has not been moved to last position")
	}

	arr = moveElementToLastPosition("Fifth", arr, false)

	if len(arr) != 4 {
		t.Error("Element should not have been added")
	}

	if arr[0] != "One" || arr[1] != "Three" || arr[2] != "Fourth" || arr[3] != "Two" {
		t.Error("Array has change, even though it shouldn't have")
	}

}

/*
Tests the cleaning of input string from line breaks
*/
func TestCleanInput(t *testing.T) {

	input := "Program Manager\n has objectives \r\n and we have| variable input \n\n    to clean.\n\n\n"

	input = CleanInput(input, "|")

	expectedOutput := "Program Manager  has objectives   and we have variable input       to clean.   "

	if input != expectedOutput {
		t.Fatal("Preprocessing did not work according to expectation.")
	}
}

/*
Tests basic processing of text for output generation
*/
func TestPerformOutputSpecificAdjustments(t *testing.T) {

	input := "\" Example txt with \" quotation mark embedded \" randomly ..."

	// Google-specific
	output := performOutputSpecificAdjustments(input, OUTPUT_TYPE_GOOGLE_SHEETS)
	if output != "'' Example txt with ' quotation mark embedded ' randomly ..." {
		t.Fatal("Test error: Substitution for characters during export resulted in : " + output)
	}

	// no Google-specific consideration
	output = performOutputSpecificAdjustments(input, "")
	if output != "' Example txt with ' quotation mark embedded ' randomly ..." {
		t.Fatal("Test error: Substitution for characters during export resulted in : " + output)
	}

}
