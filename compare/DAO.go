package compare

type DiffWithName struct {
	Name           string
	Namespace      string
	Diff           []string
	PropertyName   string
	MessageHeading string
	SourceMessage  string
	TargetMessage  string
}
