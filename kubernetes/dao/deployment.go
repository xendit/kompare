package dao

// TODO: to decide which fields to take
type Deployment struct {
	Name      string
	Namespace string
	Specs     []DeploymentSpec
}

type DeploymentSpec struct {
	Name string
}
