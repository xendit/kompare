package compare

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-test/deep"

	"kompare/tools"

	v1 "k8s.io/api/apps/v1"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	batchv1 "k8s.io/api/batch/v1"
	Corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	apiextensionv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

type TypeAssertionFunc func(interface{}) (bool, interface{})

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
	// TODO cases also in queries need those.
	// roles
	// clusterroles
	// rolebindings
	// clusterrolebindings
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

// Function to check if an object has an "Items" field
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

func CompareNumbersGenericOutput(number1, number2 int, what interface{}) {
	fmt.Printf("The number of %s in the source cluster is %d and there are %d in the target cluster.\n",
		tools.ConvertTypeStringToHumanReadable(what), number1, number2)
}

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

// CompareNumberOfDeployments compares the number of deployments in the source and target clusters.
// It takes two DeploymentList objects as input.
// It prints the number of deployments in the source and target clusters.
// It returns the number of deployments in the source and target clusters.
// func CompareNumberOfDeployments(sourceDeployments, targetDeplotments *v1.DeploymentList) (int, int) {
// 	// Print quantity of deployments
// 	lenSourceDeployments := len(sourceDeployments.Items)
// 	lenTargetDeplotments := len(targetDeplotments.Items)
// 	fmt.Printf("There are %d Deployments(apps) in the source cluster and %d in the target cluster\n",
// 		lenSourceDeployments, lenTargetDeplotments)
// 	// If deployment quantities in both clusters are different, find those different apps in a later function
// 	return lenSourceDeployments, lenTargetDeplotments
// }

// IterateDeploymentsSimpleDiff iterates over the deployments in the source and target clusters.
// It calls CompareNumberOfDeployments to get the number of deployments in each cluster.
// If the number of deployments is not equal, it identifies deployments that are only present in one cluster but not the other.
// It prints a notice about the unequal number of deployments.
// It returns two lists: onlyInSource and onlyInTarget, which contain the names of deployments that are only present in the source or target cluster, respectively.
// func IterateDeploymentsSimpleDiff(sourceDeployments, targetDeplotments *v1.DeploymentList) ([]string, []string) {
// 	lenSourceDeployments, lenTargetDeplotments := CompareNumberOfDeployments(sourceDeployments, targetDeplotments)
// 	if lenSourceDeployments != lenTargetDeplotments {
// 		var onlyInSource, onlyInTarget []string

// 		fmt.Printf("NOTICE NOT EQUAL NUMBER OF DEPLOYMENTS!!!\n")
// 		onlyInSource = compareDeploymentsByName(sourceDeployments, targetDeplotments,
// 			"- Source cluster has deployment %s, but it's not in the target cluster\n")
// 		onlyInTarget = compareDeploymentsByName(targetDeplotments, sourceDeployments,
// 			"- target cluster has deployment %s, but it's not in the source cluster\n")
// 		return onlyInSource, onlyInTarget
// 	}
// 	return nil, nil
// }

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

// Function to generate a generic message
func generateMessage(template, objectType, ItemName string) string {
	return fmt.Sprintf(template, objectType, ItemName)
}

// Function to check if an item is present in the second interface
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

// Function to get the name of an item (assumes the item has a "Name" field)
func getName(item interface{}) string {
	nameField := reflect.ValueOf(item).FieldByName("Name")
	if nameField.IsValid() && nameField.Kind() == reflect.String {
		return nameField.String()
	}
	return ""
}

// compareDeploymentsByName compares two lists of deployments by name and returns a list of names that exist in the first list but not in the second list.
//
// Parameters:
//
//	first_deployments: List of deployments to compare against the second list
//	second_deployments: List of deployments to compare against the first list
//	message_heading: A string that will be used to print a message when a deployment in the first list is not found in the second list
//
// Returns:
//
//	diffNameList: List of names of deployments that exist in the first list but not in the second list
// func compareDeploymentsByName(first_deployments, second_deployments *v1.DeploymentList, message_heading string) []string {
// 	var diffNameList []string
// 	for _, d := range first_deployments.Items {
// 		exists := false
// 		for _, b := range second_deployments.Items {
// 			if b.Name == d.Name {
// 				exists = true
// 			}
// 		}
// 		if exists == false {
// 			fmt.Printf(strings.Replace(message_heading, "%s", d.Name, -1))
// 			diffNameList = append(diffNameList, d.Name)
// 		}
// 	}
// 	return diffNameList
// }

// DeepDeploySourceTargetCompare compares the important parts of the manifest of deployments from the source and target clusters.
//
// Parameters:
//
//	sourceDeployments: List of deployments from the source cluster
//	targetDeployments: List of deployments from the target cluster
//
// Returns:
//
//	diffSourceTarget: List of DiffWithName structs containing the deployments that exist in both clusters but have differences in their specifications.
// func DeepDeploySourceTargetCompare(sourceDeployments, targetDeployments *v1.DeploymentList) []DiffWithName {
// 	var tmpDiff DiffWithName
// 	var diffSourceTarget []DiffWithName
// 	for _, d := range sourceDeployments.Items {
// 		for _, b := range targetDeployments.Items {
// 			if b.Name == d.Name {
// 				fmt.Println("Comparing " + b.Name + " on both source and target cluster.")
// 				if !reflect.DeepEqual(d.Spec.Template.Spec, b.Spec.Template.Spec) {
// 					fmt.Println("Deployment " + b.Name + " exists on both clusters, but it's different")
// 					if diff := deep.Equal(d.Spec.Template.Spec, b.Spec.Template.Spec); diff != nil {
// 						fmt.Println("Diff:")
// 						fmt.Println(diff)
// 						tmpDiff.Name = b.Name
// 						tmpDiff.Namespace = b.Namespace
// 						tmpDiff.Diff = diff
// 						diffSourceTarget = append(diffSourceTarget, tmpDiff)
// 					}
// 				}
// 				fmt.Println("Done.")
// 			}
// 		}
// 	}
// 	return diffSourceTarget
// }

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

func DeepCompare(sourceInterface, targetInterface interface{}, DiffCriteria []string) []DiffWithName {
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
						diffSourceTarget = append(diffSourceTarget, tmpDiff)
					}
				}
			}
		}
	} else {
		fmt.Println("'Items' field is not a slice in source or target object.")
	}

	return diffSourceTarget
}

func ShowResourceComparison(sourceResource, targetResource interface{}, diffCriteria []string) ([]DiffWithName, error) {
	var TheDiff []DiffWithName
	lensourceResource := GenericCountListElements(sourceResource)
	lentargetResource := GenericCountListElements(targetResource)
	resourceType := tools.ConvertTypeStringToHumanReadable(sourceResource)

	messageheading := "* These two cluster do not have the same number of " + resourceType + ", please check it manually! *"
	lenMessageheading := len(messageheading)
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
	DeepCompare(targetResource, sourceResource, diffCriteria)
	TheDiff = DeepCompare(targetResource, sourceResource, diffCriteria)
	return TheDiff, nil
}

func FormatDiffHumanReadable(differences []DiffWithName) string {
	var formattedDiff strings.Builder
	for _, diff := range differences {
		if len(diff.Diff) != 0 {
			formattedDiff.WriteString(fmt.Sprintf("Object Name: %s\n", diff.Name))
			formattedDiff.WriteString(fmt.Sprintf("Namespace: %s\n", diff.Namespace))

			formattedDiff.WriteString("Differences:\n")
			for _, d := range diff.Diff {
				formattedDiff.WriteString(fmt.Sprintf("- %s\n", d))
			}
			formattedDiff.WriteString("\n")
		}
	}
	return formattedDiff.String()
}

func CompareVerboseVSNonVerbose(sourceNameSpacesList, targetNameSpacesList interface{}, diffCriteria []string, boolverboseDiffs *bool) ([]DiffWithName, error) {
	if *boolverboseDiffs {
		TheDiff, err := ShowResourceComparison(sourceNameSpacesList, targetNameSpacesList, diffCriteria)
		fmt.Println(FormatDiffHumanReadable(TheDiff))
		return TheDiff, err
	} else {
		return ShowResourceComparison(sourceNameSpacesList, targetNameSpacesList, diffCriteria)
	}
}
