package diff

import (
	v1 "k8s.io/api/core/v1"
)

type NamespaceDiff struct {
	OnlyInSource []string
	OnlyInTarget []string
}

func IsNamespaceCountEqual(sourceNamespaces, targetNamespaces *v1.NamespaceList) bool {
	return len(sourceNamespaces.Items) == len(targetNamespaces.Items)
}

func GetNamespaceDiffByName(sourceNamespaces, targetNamespaces *v1.NamespaceList) *NamespaceDiff {
	if IsNamespaceCountEqual(sourceNamespaces, targetNamespaces) {
		return nil
	}

	onlyInSource := make([]string, 0)
	onlyInTarget := make([]string, 0)

	for _, sourceNamespace := range sourceNamespaces.Items {
		found := false
		for _, targetNamespace := range targetNamespaces.Items {
			if sourceNamespace.Name == targetNamespace.Name {
				found = true
				break
			}
		}
		if !found {
			onlyInSource = append(onlyInSource, sourceNamespace.Name)
		}
	}

	for _, targetNamespace := range targetNamespaces.Items {
		found := false
		for _, sourceNamespace := range sourceNamespaces.Items {
			if targetNamespace.Name == sourceNamespace.Name {
				found = true
				break
			}
		}
		if !found {
			onlyInTarget = append(onlyInTarget, targetNamespace.Name)
		}
	}

	return &NamespaceDiff{
		OnlyInSource: onlyInSource,
		OnlyInTarget: onlyInTarget,
	}
}
