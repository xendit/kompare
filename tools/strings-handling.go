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
