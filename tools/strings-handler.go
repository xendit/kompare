package tools

import (
	"reflect"
	"strings"
	"unicode"
)

func ConvertTypeStringToHumanReadable(what interface{}) string {
	objType := reflect.TypeOf(what)

	// If it's a pointer, get the element type
	if objType.Kind() == reflect.Ptr {
		objType = objType.Elem()
	}

	// Get the package and type names
	aString := objType.Name()
	// pkgPath := objType.PkgPath()
	dotIndex := strings.Index(aString, ".")
	aString = aString[dotIndex+1:]
	aString = convertCamelCaseToSpaces(aString)
	// Transform "List" to "in the list"
	aString = strings.TrimSuffix(aString, " List")
	aString += " in the list"
	return aString
}

func convertCamelCaseToSpaces(s string) string {
	var result string
	for i, char := range s {
		if i > 0 && unicode.IsUpper(char) {
			result += " "
		}
		result += string(char)
	}
	return result
}

func ParseCommaSeparateList(s string) []string {
	if s == "" {
		return []string{}
	}

	return strings.Split(s, ",")
}

// IsInList checks if a string is present in a list of strings
func IsInList(str string, list []string) bool {
	for _, s := range list {
		if str == s {
			return true
		}
	}
	return false
}

// AreAnyInLists checks if any element from the first list is present in the second list
func AreAnyInLists(list1, list2 []string) bool {
	for _, s1 := range list1 {
		for _, s2 := range list2 {
			if s1 == s2 {
				return true
			}
		}
	}
	return false
}

// HasCharacter checks if a string contains a specific character
func HasCharacter(str string, char byte) bool {
	for i := 0; i < len(str); i++ {
		if str[i] == char {
			return true
		}
	}
	return false
}
