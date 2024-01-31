package compare

import (
	"fmt"
	"kompare/query"
)

// Compare CRDs using generic functions from module "compare"
func CompareCRDs(targetContext, configFile string, boolverboseDiffs *bool) ([]DiffWithName, error) {
	var TheDiff []DiffWithName
	sourceCRDs, err := query.ListCRDs("", configFile)
	if err != nil {
		fmt.Printf("Error getting CRDs list: %v\n", err)
		return TheDiff, err
	}
	targetCRDs, err := query.ListCRDs(targetContext, configFile)
	if err != nil {
		fmt.Printf("Error getting CRDs list: %v\n", err)
		return TheDiff, err
	}
	diffCriteria := []string{"Spec", "Name"}
	return CompareVerboseVSNonVerbose(sourceCRDs, targetCRDs, diffCriteria, boolverboseDiffs)
}
