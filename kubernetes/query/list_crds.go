package query

import (
	"context"
	"kompare/kubernetes"

	apiextensionv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ListCustomResourceDefinitions(client *kubernetes.Client) (*apiextensionv1.CustomResourceDefinitionList, error) {
	result, err := client.APIExtensionClientset.ApiextensionsV1().CustomResourceDefinitions().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return result, nil
}
