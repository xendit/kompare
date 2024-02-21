package tools

import (
	"reflect"
	"testing"
)

func TestConvertCamelCaseToSpaces(t *testing.T) {
	type TestCase struct {
		input    string
		expected string
	}

	tests := []TestCase{
		{input: "camelCaseString", expected: "camel Case String"},
		{input: "HTTPResponse", expected: "H T T P Response"},
		{input: "anotherString", expected: "another String"},
	}

	for _, tc := range tests {
		actual := convertCamelCaseToSpaces(tc.input)
		if actual != tc.expected {
			t.Errorf("Expected %s, but got %s", tc.expected, actual)
		}
	}
}

func TestParseCommaSeparateList(t *testing.T) {
	input := "apple,banana,orange"
	expected := []string{"apple", "banana", "orange"}

	actual := ParseCommaSeparateList(input)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %v, but got %v", expected, actual)
	}
}

func TestIsInList(t *testing.T) {
	list := []string{"apple", "banana", "orange"}

	if !IsInList("apple", list) {
		t.Errorf("Expected 'apple' to be in the list")
	}

	if IsInList("grape", list) {
		t.Errorf("Expected 'grape' not to be in the list")
	}
}

func TestAreAnyInLists(t *testing.T) {
	list1 := []string{"apple", "banana", "orange"}
	list2 := []string{"orange", "grape"}

	if !AreAnyInLists(list1, list2) {
		t.Errorf("Expected elements from list1 to be in list2")
	}

	list3 := []string{"pear", "peach"}

	if AreAnyInLists(list1, list3) {
		t.Errorf("Expected elements from list1 not to be in list3")
	}
}

func TestHasCharacter(t *testing.T) {
	str := "hello"

	if !HasCharacter(str, 'h') {
		t.Errorf("Expected 'h' to be in the string")
	}

	if HasCharacter(str, 'z') {
		t.Errorf("Expected 'z' not to be in the string")
	}
}
