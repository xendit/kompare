package diff

import (
	"kompare/kubernetes/dao"
)

func IsNamespaceCountEqual(sourceNamespaces, targetNamespaces []dao.Namespace) bool {
	return len(sourceNamespaces) == len(targetNamespaces)
}

func GetNamespaceDiffByName(sourceNamespaces, targetNamespaces []dao.Namespace) ([]dao.Namespace, []dao.Namespace) {
	if IsNamespaceCountEqual(sourceNamespaces, targetNamespaces) {
		return nil, nil
	}

	onlyInSource := make([]dao.Namespace, 0)
	onlyInTarget := make([]dao.Namespace, 0)

	for _, sourceNamespace := range sourceNamespaces {
		found := false
		for _, targetNamespace := range targetNamespaces {
			if sourceNamespace.Name == targetNamespace.Name {
				found = true
				break
			}
		}
		if !found {
			onlyInSource = append(onlyInSource, sourceNamespace)
		}
	}

	for _, targetNamespace := range targetNamespaces {
		found := false
		for _, sourceNamespace := range sourceNamespaces {
			if targetNamespace.Name == sourceNamespace.Name {
				found = true
				break
			}
		}
		if !found {
			onlyInTarget = append(onlyInTarget, targetNamespace)
		}
	}

	return onlyInSource, onlyInTarget
}
