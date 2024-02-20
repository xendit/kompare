package tools

import (
	"encoding/json"
	"fmt"
	"strings"

	"kompare/DAO"
)

// FormatDiffHumanReadable formats differences in a human-readable format.
// It takes a slice of DiffWithName (differences) containing information about differences between resources.
// It iterates over each difference and constructs a human-readable format containing details such as resource type, object name, namespace, and specific differences.
// If a property name, object name, namespace, or differences exist for a particular difference, it includes them in the formatted output.
// The function returns a string containing the human-readable formatted differences.
func FormatDiffHumanReadable(differences []DAO.DiffWithName) string {
	var formattedDiff strings.Builder
	for _, diff := range differences {
		if len(diff.Diff) != 0 {
			if diff.PropertyName != "" {
				formattedDiff.WriteString(fmt.Sprintf("Kubernetes resource definition type: %s\n", diff.PropertyName))
			}
			if diff.Name != "" {
				formattedDiff.WriteString(fmt.Sprintf("Object Name: %s\n", diff.Name))
			}
			if diff.Namespace != "" {
				formattedDiff.WriteString(fmt.Sprintf("Namespace: %s\n", diff.Namespace))
			}
			formattedDiff.WriteString("Differences:\n")
			if len(diff.Diff) > 0 {
				for _, d := range diff.Diff {
					key, value, result := startsWithMapPattern(d)
					if result {
						x, y, z := ExtractSubstrings(value)
						leftMultiline := strings.Contains(x, "\n")
						rightMultiline := strings.Contains(z, "\n")
						if leftMultiline || rightMultiline {
							formattedDiff.WriteString(fmt.Sprintf("- %s:\n", key))
							if isJSONCompatible(x) {
								prettyX, err := prettifyJSON(x)
								if err != nil {
									// Handle error
									formattedDiff.WriteString(fmt.Sprintf("Error prettifying left side JSON: %v\n", err))
								} else {
									formattedDiff.WriteString(fmt.Sprintf("%s\n", prettyX))
								}
							} else {
								formattedDiff.WriteString(fmt.Sprintf("%s:\n", x))
							}
							formattedDiff.WriteString(fmt.Sprintf(" %s\n", y))
							if isJSONCompatible(z) {
								prettyZ, err := prettifyJSON(z)
								if err != nil {
									// Handle error
									formattedDiff.WriteString(fmt.Sprintf("Error prettifying right side JSON: %v\n", err))
								} else {
									formattedDiff.WriteString(fmt.Sprintf("%s\n", prettyZ))
								}
							} else {
								formattedDiff.WriteString(fmt.Sprintf("%s:\n", z))
							}
						} else {
							formattedDiff.WriteString(fmt.Sprintf("- %s: %s %s %s\n", key, x, y, z))
						}
					} else {
						formattedDiff.WriteString(fmt.Sprintf("- %v\n", d))
					}
				}
				formattedDiff.WriteString("\n")
			}
		} else {
			fmt.Println("No differences found.")
		}
	}
	return formattedDiff.String()
}

func startsWithMapPattern(input string) (string, string, bool) {
	// Define the prefix pattern
	prefix := "map["
	// Check if the input string starts with the prefix
	if strings.Contains(input, prefix) {
		// Find the index of the closing bracket "]"
		closingBracketIndex := strings.Index(input, ":")

		// Check if the closing bracket is followed by a colon ":"
		if closingBracketIndex != -1 && closingBracketIndex < len(input)-1 && input[closingBracketIndex+1] == ':' {
			// Extract the substring between the brackets and after the colon
			key := strings.TrimSpace(input[len(prefix):closingBracketIndex])
			value := strings.TrimSpace(input[closingBracketIndex+1:]) // Trimming spaces
			return key, value, true
		}
	}

	return "", "", false
}

func ExtractSubstrings(value string) (string, string, string) {
	parts := strings.Split(value, "!= ")

	if len(parts) < 2 {
		return "", "", ""
	}

	leftSide := parts[0]
	rightSide := parts[1]

	// Split operatorAndRightSide to get the operator and right side
	leftSide = strings.TrimSpace(leftSide)
	rightSide = strings.TrimSpace(rightSide)
	return leftSide, "!=", rightSide
}

func isJSONCompatible(jsonStr string) bool {
	var data interface{}
	err := json.Unmarshal([]byte(jsonStr), &data)
	return err == nil
}

func prettifyJSON(jsonStr string) (string, error) {
	var data interface{}
	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		return "", err
	}

	prettyJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}
	return string(prettyJSON), nil
}
