package compare

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-test/deep"

	"kompare/cli"
	"kompare/tools"

	v1 "k8s.io/api/apps/v1"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	batchv1 "k8s.io/api/batch/v1"
	Corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	RbacV1 "k8s.io/api/rbac/v1"
	apiextensionv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/client-go/kubernetes"
)

type TypeAssertionFunc func(interface{}) (bool, interface{})

// typeAssertions is a map containing type assertion functions for various Kubernetes resource lists.
// Each entry in the map consists of a string key representing the type of Kubernetes resource list
// and a corresponding TypeAssertionFunc, which is a function type.
var typeAssertions = map[string]TypeAssertionFunc{
	"*Corev1.NamespaceList": func(obj interface{}) (bool, interface{}) {
		val, ok := obj.(*Corev1.NamespaceList)
		return ok, val
	},
	"*v1.DeploymentList": func(obj interface{}) (bool, interface{}) {
		val, ok := obj.(*v1.DeploymentList)
		return ok, val
	},
	"*autoscalingv1.HorizontalPodAutoscalerList": func(obj interface{}) (bool, interface{}) {
		val, ok := obj.(*autoscalingv1.HorizontalPodAutoscalerList)
		return ok, val
	},
	"*batchv1.CronJobList": func(obj interface{}) (bool, interface{}) {
		val, ok := obj.(*batchv1.CronJobList)
		return ok, val
	},
	"*apiextensionv1.CustomResourceDefinitionList": func(obj interface{}) (bool, interface{}) {
		val, ok := obj.(*apiextensionv1.CustomResourceDefinitionList)
		return ok, val
	},
	"*networkingv1.IngressList": func(obj interface{}) (bool, interface{}) {
		val, ok := obj.(*networkingv1.IngressList)
		return ok, val
	},
	"*Corev1.ServiceList": func(obj interface{}) (bool, interface{}) {
		val, ok := obj.(*Corev1.ServiceList)
		return ok, val
	},
	"*Corev1.ConfigMapList": func(obj interface{}) (bool, interface{}) {
		val, ok := obj.(*Corev1.ConfigMapList)
		return ok, val
	},
	"*Corev1.SecretList": func(obj interface{}) (bool, interface{}) {
		val, ok := obj.(*Corev1.SecretList)
		return ok, val
	},
	"*Corev1.ServiceAccountList": func(obj interface{}) (bool, interface{}) {
		val, ok := obj.(*Corev1.ServiceAccountList)
		return ok, val
	},
	"*RbacV1.RoleList": func(obj interface{}) (bool, interface{}) {
		val, ok := obj.(*RbacV1.RoleList)
		return ok, val
	},
	"*RbacV1.RoleBindingList": func(obj interface{}) (bool, interface{}) {
		val, ok := obj.(*RbacV1.RoleBindingList)
		return ok, val
	},
	"*RbacV1.ClusterRoleList": func(obj interface{}) (bool, interface{}) {
		val, ok := obj.(*RbacV1.ClusterRoleList)
		return ok, val
	},
	"*RbacV1.ClusterRoleBindingList": func(obj interface{}) (bool, interface{}) {
		val, ok := obj.(*RbacV1.ClusterRoleBindingList)
		return ok, val
	},
}

func assertType(typeName string, obj interface{}) (bool, interface{}) {
	if assertionFunc, ok := typeAssertions[typeName]; ok {
		return assertionFunc(obj)
	} else {
		return false, nil
	}
}

func GetTypeInfo(obj interface{}) (string, interface{}) {
	typeName := "unknown"
	var objValue interface{}

	for t, assertionFunc := range typeAssertions {
		if success, value := assertionFunc(obj); success {
			typeName = t
			objValue = value
			break
		}
	}

	return typeName, objValue
}

// GenericCountListElements counts the number of elements in a generic list object.
// It checks the type of the object and extracts its value.
// If the object has an "Items" field and it's of slice type, it returns the length of the slice.
// If the object does not meet the criteria for counting elements, it returns 0.
func GenericCountListElements(obj interface{}) int {
	// Check the type of the object
	_, objValue := GetTypeInfo(obj)

	// Check if it's a list type with an "Items" field
	if hasItemsField(objValue) {
		// Get the "Items" field of the list
		itemsField := reflect.ValueOf(objValue).Elem().FieldByName("Items")

		// Check if the field is a slice
		if itemsField.Kind() == reflect.Slice {
			// Return the length of the slice
			return itemsField.Len()
		}
	}

	return 0
}

// hasItemsField checks if an object has an "Items" field.
// It takes an interface{} as input and inspects its type and value using reflection.
// If the object is a pointer to a struct and the struct has a field named "Items",
// it returns true if the "Items" field exists in the value, otherwise returns false.
func hasItemsField(obj interface{}) bool {
	objType := reflect.TypeOf(obj)
	objValue := reflect.ValueOf(obj)

	// Check if it's a struct
	if objType.Kind() == reflect.Ptr && objType.Elem().Kind() == reflect.Struct {
		// Dereference the pointer to get the struct type
		structType := objType.Elem()

		// Check if the struct has a field named "Items"
		if _, found := structType.FieldByName("Items"); found {
			// Check if the "Items" field exists in the value
			return objValue.Elem().FieldByName("Items").IsValid()
		}
	}

	return false
}

// CompareNumbersGenericOutput prints a comparison result between two numbers along with a description of what they represent.
// It takes two integers (number1 and number2) representing counts from different clusters and an interface{} (what) as input.
// The interface{} parameter is expected to describe the type of objects being compared.
// It formats and prints a comparison message indicating the count of objects represented by 'what'
// in the source cluster (number1) compared to the count in the target cluster (number2).
func CompareNumbersGenericOutput(number1, number2 int, what interface{}) {
	fmt.Printf("The number of %s in the source cluster is %d and there are %d in the target cluster.\n",
		tools.ConvertTypeStringToHumanReadable(what), number1, number2)
}

// IterateGenericSimpleDiff iterates over two generic interfaces representing lists of objects
// and identifies the differences between them.
// It takes sourceInterface and targetInterface as input, which are interfaces representing lists of objects.
// It calculates the lengths of both interfaces and compares them.
// If the lengths are not equal, it prints a notice indicating the discrepancy.
// It currently returns empty slices for differences since the comparison logic is not implemented.
// Further implementation is required to compare the actual objects and identify differences.
// The function returns two empty string slices for differences.
func IterateGenericSimpleDiff(sourceInterface, targetInterface interface{}) ([]string, []string) {
	lenSourceInterface := GenericCountListElements(sourceInterface)
	lenTargetInterface := GenericCountListElements(targetInterface)
	if lenSourceInterface != lenTargetInterface {
		// var onlyInSource, onlyInTarget []string

		fmt.Printf("NOTICE: not equal number of %v!!!\n", tools.ConvertTypeStringToHumanReadable(sourceInterface))
		// need to compare interfaces by finding each items name like in the static IterateDeploymentsSimpleDiff function
	}

	return nil, nil
}

// CompareByName compares objects in two interfaces by their names.
// It takes two interface{} parameters (firstInterface and secondInterface) representing lists of objects,
// and a string message_heading as input.
// It extracts the "Items" field from both interfaces and iterates through the items in the first interface.
// For each item in the first interface, it checks if the item is present in the second interface.
// If an item is not present in the second interface, it generates and prints a message using the message_heading
// and appends the item's name to the diffNameList.
// The function returns a slice containing the names of the items that are present in the first interface
// but not in the second interface.
func CompareByName(firstInterface, secondInterface interface{}, message_heading string) []string {
	var diffNameList []string

	// Extract the "Items" field from the first and second interfaces
	firstItems := reflect.ValueOf(firstInterface).Elem().FieldByName("Items")
	secondItems := reflect.ValueOf(secondInterface).Elem().FieldByName("Items")

	// Check if both fields are slices
	if firstItems.Kind() == reflect.Slice && secondItems.Kind() == reflect.Slice {
		// Loop through the items in the first interface
		for i := 0; i < firstItems.Len(); i++ {
			item := firstItems.Index(i).Interface()

			// Check if the item is not present in the second interface
			if !containsItem(item, secondItems) {
				fmt.Printf(generateMessage(message_heading, tools.ConvertTypeStringToHumanReadable(item), getName(item)))
				diffNameList = append(diffNameList, getName(item))
			}
		}
	}
	return diffNameList
}

// generateMessage generates a generic message using a template, object type, and item name.
// It takes three string parameters: template, objectType, and ItemName.
// It formats and returns a message by inserting the objectType and ItemName into the template.
func generateMessage(template, objectType, ItemName string) string {
	return fmt.Sprintf(template, objectType, ItemName)
}

// containsItem checks if an item is present in the second interface.
// It takes an item interface{} and a reflect.Value (secondItems) representing the second interface.
// It loops through the items in the second interface and compares each item's name with the name of the provided item.
// If a matching item is found, it returns true; otherwise, it returns false.
func containsItem(item interface{}, secondItems reflect.Value) bool {
	// Loop through the items in the second interface
	for i := 0; i < secondItems.Len(); i++ {
		secondItem := secondItems.Index(i).Interface()

		// Compare items by name
		if getName(item) == getName(secondItem) {
			return true
		}
	}

	return false
}

// getName retrieves the name of an item assuming it has a "Name" field.
// It takes an item interface{} as input and extracts the value of the "Name" field.
// If the field is valid and of type string, it returns the string value of the field.
// If the field is invalid or not of type string, it returns an empty string.
func getName(item interface{}) string {
	nameField := reflect.ValueOf(item).FieldByName("Name")
	if nameField.IsValid() && nameField.Kind() == reflect.String {
		return nameField.String()
	}
	return ""
}

// getNestedFieldValue retrieves the value of a nested field within a structure using reflection.
// It takes a reflect.Value (obj) representing the structure and a slice of strings (fieldNames) representing the nested field names.
// It iterates through each field name in the fieldNames slice and accesses the corresponding nested field in the structure.
// If the structure contains pointers, it dereferences them to access the nested fields.
// If a nested field is not found or is invalid, it returns an error indicating the missing field.
// Otherwise, it returns the reflect.Value of the nested field.
func getNestedFieldValue(obj reflect.Value, fieldNames []string) (reflect.Value, error) {
	for _, fieldName := range fieldNames {
		// Dereference pointers in the nested structure
		if obj.Kind() == reflect.Ptr {
			obj = obj.Elem()
		}

		// Access the nested field
		obj = obj.FieldByName(fieldName)

		// Check if the field is valid
		if !obj.IsValid() {
			return reflect.Value{}, fmt.Errorf("Field %s not found", fieldName)
		}
	}

	return obj, nil
}

// DeepCompare performs a deep comparison between two interfaces representing lists of objects.
// It compares the objects based on specified criteria and returns a list of differences along with their names and namespaces.
// It takes sourceInterface and targetInterface as input interfaces and DiffCriteria as a slice of strings representing comparison criteria.
// It iterates over the 'Items' fields of both sourceInterface and targetInterface and compares each item's 'Name' field.
// If the 'Name' fields match, it compares the specified DiffCriteria fields of the objects using the deep.Equal function from the 'github.com/go-test/deep' package.
// It constructs DiffWithName structs containing the object name, namespace, difference details, and property name for each difference found.
// The function returns a slice of DiffWithName containing the differences between the source and target interfaces based on the specified criteria.
func DeepCompare(sourceInterface, targetInterface interface{}, DiffCriteria []string) ([]DiffWithName, error) {
	var tmpDiff DiffWithName
	var diffSourceTarget []DiffWithName
	// Get type information for source and target
	_, sourceObject := GetTypeInfo(sourceInterface)
	_, targetObject := GetTypeInfo(targetInterface)
	sourceItemsField := reflect.ValueOf(sourceObject).Elem().FieldByName("Items")
	targetItemsField := reflect.ValueOf(targetObject).Elem().FieldByName("Items")
	// Check if 'Items' is a slice in both source and target objects
	if sourceItemsField.Kind() == reflect.Slice && targetItemsField.Kind() == reflect.Slice {
		// Iterate over sourceItems
		for i := 0; i < sourceItemsField.Len(); i++ {
			sourceItem := sourceItemsField.Index(i).Interface()

			// Iterate over targetItems
			for j := 0; j < targetItemsField.Len(); j++ {
				targetItem := targetItemsField.Index(j).Interface()
				// Compare 'Name' fields
				sourceName, _ := getNestedFieldValue(reflect.ValueOf(sourceItem), []string{"Name"})
				targetName, _ := getNestedFieldValue(reflect.ValueOf(targetItem), []string{"Name"})
				sourceNamespace, _ := getNestedFieldValue(reflect.ValueOf(sourceItem), []string{"Namespace"})
				if sourceName.String() == targetName.String() {
					for _, v := range DiffCriteria {
						sourceDiffCriteriaField, err := getNestedFieldValue(reflect.ValueOf(sourceItem), strings.Split(v, "."))
						if err != nil {
							fmt.Printf("Error accessing field: %v\n", err)
							continue
						}
						targetDiffCriteriaField, err := getNestedFieldValue(reflect.ValueOf(targetItem), strings.Split(v, "."))
						if err != nil {
							fmt.Printf("Error accessing field: %v\n", err)
							continue
						}
						xdiff := deep.Equal(sourceDiffCriteriaField.Interface(), targetDiffCriteriaField.Interface())
						tmpDiff.Name = targetName.String()
						tmpDiff.Namespace = sourceNamespace.String()
						tmpDiff.Diff = xdiff
						tmpDiff.PropertyName = v
						diffSourceTarget = append(diffSourceTarget, tmpDiff)
					}
				}
			}
		}
	} else {
		fmt.Println("'Items' field is not a slice in source or target object.")
	}

	return diffSourceTarget, nil
}

// ShowResourceComparison compares two sets of resources from different clusters and identifies differences based on specified criteria.
// It takes sourceResource and targetResource as input interfaces representing lists of resources from different clusters,
// and diffCriteria as a slice of strings representing comparison criteria.
// It calculates the lengths of sourceResource and targetResource and compares them.
// If the lengths are different, it prints a message indicating the discrepancy and performs a number comparison.
// It then compares the resources in both clusters using the CompareByName function and prints the differences.
// It also performs a deep comparison of resources based on the specified diffCriteria using the DeepCompare function.
// The function returns a slice of DiffWithName containing the differences between the source and target resources,
// along with any error encountered during the comparison.
func ShowResourceComparison(sourceResource, targetResource interface{}, diffCriteria []string, args cli.ArgumentsReceivedValidated) ([]DiffWithName, error) {
	var TheDiff []DiffWithName
	lensourceResource := GenericCountListElements(sourceResource)
	lentargetResource := GenericCountListElements(targetResource)
	resourceType := tools.ConvertTypeStringToHumanReadable(sourceResource)

	messageheading := "* These two cluster do not have the same number of " + resourceType + ", please check it manually! *"
	lenMessageheading := len(messageheading)
	if args.VerboseDiffs != 0 {
		if lentargetResource != lensourceResource {

			fmt.Println(strings.Repeat("*", lenMessageheading))
			fmt.Println(messageheading)
			fmt.Println(strings.Repeat("*", lenMessageheading))
			CompareNumbersGenericOutput(lensourceResource, lentargetResource, targetResource)
		}
		fmt.Println(strings.Repeat("*", lenMessageheading))
		sourceMessageTemplate := "- First cluster has %s: %s, but it's not in the second cluster\n"
		resultStringsSvT := CompareByName(sourceResource, targetResource, sourceMessageTemplate)
		if len(resultStringsSvT) > 0 {
			fmt.Println(strings.Repeat("*", lenMessageheading))
		} else {
			fmt.Println("Done compering source cluster versus target cluster's ", resourceType)
		}
		targetmessageTemplate := "- Second cluster has %s: %s, but it's not in the first cluster\n"
		resultStringsTvS := CompareByName(targetResource, sourceResource, targetmessageTemplate)
		if len(resultStringsTvS) > 0 {
			fmt.Println(strings.Repeat("*", lenMessageheading))
		} else {
			fmt.Println("Done compering target cluster versus source cluster's ", resourceType)
		}
		TheDiff, _ = DeepCompare(targetResource, sourceResource, diffCriteria)
		return TheDiff, nil
	}
	if lentargetResource != lensourceResource {
		fmt.Println(strings.Repeat("*", lenMessageheading))
		fmt.Println(messageheading)
		fmt.Println(strings.Repeat("*", lenMessageheading))
		CompareNumbersGenericOutput(lensourceResource, lentargetResource, targetResource)
	}
	TheDiff, _ = DeepCompare(targetResource, sourceResource, diffCriteria)
	return TheDiff, nil
}

// Merge two DiffWithName structs
func mergeDiffs(diff1, diff2 DiffWithName) DiffWithName {
	mergedDiff := diff1

	// Merge the fields from diff2 into mergedDiff
	if diff2.MessageHeading != "" {
		mergedDiff.MessageHeading = diff2.MessageHeading
	}
	if diff2.SourceMessage != "" {
		mergedDiff.SourceMessage = diff2.SourceMessage
	}
	if diff2.TargetMessage != "" {
		mergedDiff.TargetMessage = diff2.TargetMessage
	}
	if len(diff2.Diff) > 0 {
		mergedDiff.Diff = append(mergedDiff.Diff, diff2.Diff...)
	}

	return mergedDiff
}

// FormatDiffHumanReadable formats differences in a human-readable format.
// It takes a slice of DiffWithName (differences) containing information about differences between resources.
// It iterates over each difference and constructs a human-readable format containing details such as resource type, object name, namespace, and specific differences.
// If a property name, object name, namespace, or differences exist for a particular difference, it includes them in the formatted output.
// The function returns a string containing the human-readable formatted differences.
func FormatDiffHumanReadable(differences []DiffWithName) string {
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
		}
	}
	return formattedDiff.String()
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

func startsWithMapPattern(input string) (string, string, bool) {
	// Define the prefix pattern
	prefix := "map["

	// Check if the input string starts with the prefix
	if strings.HasPrefix(input, prefix) {
		// Find the index of the closing bracket "]"
		closingBracketIndex := strings.Index(input, "]")

		// Check if the closing bracket is followed by a colon ":"
		if closingBracketIndex != -1 && closingBracketIndex < len(input)-1 && input[closingBracketIndex+1] == ':' {
			// Extract the substring between the brackets and after the colon
			key := input[len(prefix):closingBracketIndex]
			value := input[closingBracketIndex+2:] // Skipping the ":"
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

// CompareVerboseVSNonVerbose compares two sets of namespaces from different clusters based on specified criteria.
// It takes sourceNameSpacesList and targetNameSpacesList as input interfaces representing lists of namespaces from different clusters,
// diffCriteria as a slice of strings representing comparison criteria, and boolverboseDiffs as a pointer to a boolean indicating whether to display verbose differences.
// If args is the aguments passed like for instance VerboseDiffs for the level of verbosity desired.
// The function returns a slice of DiffWithName containing the differences between the source and target namespaces,
// along with any error encountered during the comparison.
func CompareVerboseVSNonVerbose(sourceNameSpacesList, targetNameSpacesList interface{}, diffCriteria []string, args cli.ArgumentsReceivedValidated) ([]DiffWithName, error) {
	if args.VerboseDiffs != 0 {
		if args.VerboseDiffs > 1 {
			TheDiff, err := ShowResourceComparison(sourceNameSpacesList, targetNameSpacesList, diffCriteria, args)
			fmt.Println(FormatDiffHumanReadable(TheDiff))
			return TheDiff, err
		} else if args.VerboseDiffs == 1 {
			return ShowResourceComparison(sourceNameSpacesList, targetNameSpacesList, diffCriteria, args)
		}

	}
	// sumary goes here.
	return ShowResourceComparison(sourceNameSpacesList, targetNameSpacesList, diffCriteria, args)
}

// GenericCompareResources compares resources between two Kubernetes clusters based on specified criteria.
// It takes clientsetToSource and clientsetToTarget as pointers to kubernetes.Clientset representing connections to the source and target clusters,
// namespaceName as a string specifying the namespace to compare resources in,
// resourceGetter as a function to retrieve the list of resources from a clientset and namespace,
// diffCriteria as a slice of strings representing comparison criteria,
// and boolverboseDiffs as a pointer to a boolean indicating whether to display verbose differences.
// It retrieves the list of resources from both the source and target clusters using the resourceGetter function.
// It then performs a comparison between the resources from both clusters using the CompareVerboseVSNonVerbose function.
// The function returns a slice of DiffWithName containing the differences between the source and target resources,
// along with any error encountered during the comparison.
func GenericCompareResources(clientsetToSource, clientsetToTarget *kubernetes.Clientset, namespaceName string, resourceGetter func(*kubernetes.Clientset, string) (interface{}, error), diffCriteria []string, args cli.ArgumentsReceivedValidated) ([]DiffWithName, error) {
	var TheDiff []DiffWithName

	sourceResources, err := resourceGetter(clientsetToSource, namespaceName)
	if err != nil {
		return TheDiff, fmt.Errorf("error getting source resources: %v", err)
	}

	targetResources, err := resourceGetter(clientsetToTarget, namespaceName)
	if err != nil {
		return TheDiff, fmt.Errorf("error getting target resources: %v", err)
	}

	// Type assertion to convert the interface{} to a slice of the specific type
	sourceSlice, ok := sourceResources.([]interface{})
	if !ok {
		return TheDiff, fmt.Errorf("unexpected type for source resources")
	}

	targetSlice, ok := targetResources.([]interface{})
	if !ok {
		return TheDiff, fmt.Errorf("unexpected type for target resources")
	}

	return CompareVerboseVSNonVerbose(sourceSlice, targetSlice, diffCriteria, args)
}
