package cli

import (
	"fmt"
	"kompare/tools"
	"os"
	"path"

	"github.com/akamensky/argparse"
)

type ArgumentsReceived struct {
	KubeconfigFile, SourceClusterContext, TargetClusterContext, NamespaceName, FiltersForObject, Include, Exclude *string
	VerboseDiffs                                                                                                  *int
	Err                                                                                                           error
}
type ArgumentsReceivedValidated struct {
	KubeconfigFile, SourceClusterContext, TargetClusterContext, NamespaceName, FiltersForObject string
	Include, Exclude                                                                            []string
	VerboseDiffs                                                                                int
	Err                                                                                         error
}

// PaserReader parses the command-line arguments and returns validated arguments.
// It creates a new parser object and defines various flags and options.
// The flags and options include:
//   - 'c' or 'conf' flag for specifying the path to the kubeconfig file (optional).
//   - 's' or 'src' flag for specifying the source cluster's context (optional).
//   - 't' or 'target' flag for specifying the target cluster's context (required).
//   - 'v' or 'verbose' flag for enabling verbose mode to show all diffs (optional).
//   - 'i' or 'include' flag for specifying a list of Kubernetes objects to include (optional).
//   - 'e' or 'exclude' flag for specifying a list of Kubernetes objects to exclude (optional).
//   - 'n' or 'namespace' flag for specifying the namespace to be copied (optional, defaults to 'default').
//   - 'f' or 'filter' flag for specifying what parts of the object to compare (optional).
//
// If an error occurs during parsing, it prints the error and usage information.
// The function returns a struct containing validated arguments.
func PaserReader() ArgumentsReceivedValidated {
	// Create new parser object
	parser := argparse.NewParser("print", "Prints provided string to stdout")
	kubeconfigFile := parser.String("c", "conf", &argparse.Options{Required: false, Help: "Path to the clusters kubeconfig; assume ~/.kube/config if not provided"})
	sourceClusterContext := parser.String("s", "src", &argparse.Options{Required: false, Help: "The Source cluster's context. Origin cluster in the comparison (LHS-left hand side)"})
	targetClusterContext := parser.String("t", "target", &argparse.Options{Required: true, Help: "*The target cluster's context (Required). Cluster used as destination or consequent (RHS - Right hand side)"})
	verboseDiffs := parser.FlagCounter("v", "verbose", &argparse.Options{Help: "-v lists the differences and -vv just shows all the diffs too."})
	IncludeK8sObjects := parser.String("i", "include", &argparse.Options{Help: "List of kubernetes objects names to include, this should be an element or a comma separated list."})
	Excludek8sObjects := parser.String("e", "exclude", &argparse.Options{Help: "List of kubernetes objects to include, this should be an element or a comma separated list."})
	namespaceName := parser.String("n", "namespace", &argparse.Options{Help: "Namespace that needs to be copied. defaults to 'default' namespace. The option also accepts wilcard matching of namespace. E.G.: '*-pci' would match any namespace that ends with -pci. Notice that the '' might be required in some consoles like iterm"})
	filtersForObject := parser.String("f", "filter", &argparse.Options{Help: "Filter what parts of the object I want to compare. must be used together with -i option to apply to that type of objects"})
	err := parser.Parse(os.Args)
	if err != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		fmt.Print(parser.Usage(err))
		return ArgumentsReceivedValidated{
			KubeconfigFile:       "",
			SourceClusterContext: "",
			TargetClusterContext: "",
			NamespaceName:        "",
			FiltersForObject:     "",
			Include:              []string{""},
			Exclude:              []string{""},
			VerboseDiffs:         *verboseDiffs,
			Err:                  err}
	}
	TheArgs := ArgumentsReceived{
		KubeconfigFile:       kubeconfigFile,
		SourceClusterContext: sourceClusterContext,
		TargetClusterContext: targetClusterContext,
		NamespaceName:        namespaceName,
		Include:              IncludeK8sObjects,
		Exclude:              Excludek8sObjects,
		FiltersForObject:     filtersForObject,
		VerboseDiffs:         verboseDiffs,
		Err:                  err}
	ArgumentsReceivedValidated := ValidateParametersFromParserArgs(TheArgs)
	return ArgumentsReceivedValidated
}

// ValidateParametersFromParserArgs validates the arguments received from the command-line parser.
// It extracts the source cluster context, target cluster context, and namespace name from the arguments.
// If the source cluster context is empty, it informs the user that the current kubeconfig context will be used.
// If not empty, it informs the user about the source cluster context being used.
// It sets the kubeconfig file path based on the provided file path or the default ~/.kube/config.
// It validates and parses the include and exclude lists of Kubernetes objects.
// If there are invalid objects in the include or exclude lists, it prints a warning but continues execution.
// It checks for improper usage of filtering options with include and exclude lists and prints warnings accordingly.
// The function returns a struct containing validated arguments for further processing.
func ValidateParametersFromParserArgs(TheArgs ArgumentsReceived) ArgumentsReceivedValidated {
	var strSourceClusterContext, strTargetClusterContext, strNamespaceName string
	strSourceClusterContext = *TheArgs.SourceClusterContext
	strTargetClusterContext = *TheArgs.TargetClusterContext
	if strSourceClusterContext == "" {
		fmt.Println("We will use current kubeconfig context as 'source cluster'.")
	} else {
		fmt.Printf("We will use %s kubeconfig context as 'source cluster' or 'origin cluster'.\n", strSourceClusterContext)
	}
	fmt.Printf("We will use %s kubeconfig context as 'target cluster'.\n", strTargetClusterContext)

	strNamespaceName = *TheArgs.NamespaceName
	configFile := ""
	if *TheArgs.KubeconfigFile != "" {
		configFile = *TheArgs.KubeconfigFile
	} else {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("Error getting the home dir: %v\n", err)
			return ArgumentsReceivedValidated{
				KubeconfigFile: "", SourceClusterContext: "",
				TargetClusterContext: "", NamespaceName: "",
				FiltersForObject: "",
				Include:          []string{""}, Exclude: []string{""},
				VerboseDiffs: *TheArgs.VerboseDiffs, Err: nil}
		}
		configFile = path.Join(homeDir, ".kube", "config")
	}
	invalidInclude, includeStr := ValidateKubernetesObjects(tools.ParseCommaSeparateList(*TheArgs.Include))
	invalidExclude, excludeStr := ValidateKubernetesObjects(tools.ParseCommaSeparateList(*TheArgs.Exclude))

	if invalidInclude != nil {
		fmt.Print("You passed some invalid kubernetes object to incldue as a parameter: ", invalidInclude)
		fmt.Println(". The program will try to execute anyways and ignore this")
	}
	if invalidExclude != nil {
		fmt.Println("You passed some invalid kubernetes object to exclude as a parameter: ", invalidInclude)
		fmt.Println(". The program will try to execute anyways and ignore this")
	}
	if *TheArgs.FiltersForObject != "" && *TheArgs.Exclude != "" {
		fmt.Println("Warning: The -f filtering option was not designed to be used with the -e option.")
		fmt.Println("The program will try to execute anyway, but the output might not be what you expect.")
		fmt.Println("The -f is to be used with one and only one -i include object type at the time.")
	}
	if *TheArgs.FiltersForObject != "" && tools.HasCharacter(*TheArgs.Include, ',') {
		fmt.Println("Warning: The -f filtering option was not designed to be used with multiple -i objects,")
		fmt.Println("The program will try to execute anyway, but the output might not be what you expect.")
		fmt.Println("The -f is to be used with one and only one -i include object type at the time.")
	}
	return ArgumentsReceivedValidated{
		KubeconfigFile:       configFile,
		SourceClusterContext: strSourceClusterContext,
		TargetClusterContext: strTargetClusterContext,
		NamespaceName:        strNamespaceName,
		FiltersForObject:     *TheArgs.FiltersForObject,
		Include:              includeStr,
		Exclude:              excludeStr,
		VerboseDiffs:         *TheArgs.VerboseDiffs,
		Err:                  nil}
}

// ValidateKubernetesObjects validates the given list of Kubernetes object names
// against a list of valid object names and their aliases
// It returns two slices: invalidObjects and validObjects
func ValidateKubernetesObjects(objects []string) ([]string, []string) {
	validObjects := map[string][]string{
		"deployment":         {"deployment", "deployments", "deploy"},
		"ingress":            {"ingress", "ing"},
		"service":            {"service", "svc", "services"},
		"serviceaccount":     {"sa", "serviceaccount", "serviceaccounts"},
		"configmap":          {"configmap", "configmaps", "cm"},
		"secret":             {"secret", "secrets"},
		"namespace":          {"namespace", "ns", "namespaces"},
		"role":               {"role", "roles"},
		"rolebinding":        {"rolebinding", "rolebindings"},
		"clusterrole":        {"clusterrole", "clusterroles"},
		"clusterrolebinding": {"clusterrolebinding", "clusterrolebindings"},
		"crd":                {"crd", "crds", "customresourcedefinition", "customresourcedefinitions"},
		// Add more valid objects and their aliases as needed
	}

	var invalidObjects []string
	var validObjectsStr []string

	for _, obj := range objects {
		found := false
		for standardName, aliases := range validObjects {
			for _, alias := range aliases {
				if obj == alias {
					found = true
					validObjectsStr = append(validObjectsStr, standardName)
					break
				}
			}
			if found {
				break
			}
		}
		if !found {
			invalidObjects = append(invalidObjects, obj)
		}
	}
	return invalidObjects, validObjectsStr
}
