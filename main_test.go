package main

import (
	"kompare/connect"
	"kompare/query"
	"reflect"
	"testing"
)

func TestConnectNow(t *testing.T) {
	// Test case 1; Make sure we get a client set to kubernetes.
	configFile := "/Users/abel.guzman/.kube/config"
	clientsetToSource, err := connect.ConnectNow(&configFile)
	expectedType := "*kubernetes.Clientset"
	if err != nil {
		t.Errorf("Error Connecting to cluster %v", err)
	}
	if reflect.ValueOf(clientsetToSource).Type().String() != expectedType {
		t.Errorf("Test case 1; TestConnectNow(configFile) returned %s, expected %s", reflect.TypeOf(clientsetToSource).Kind().String(), expectedType)
	}
}

func TestListNameSpaces(t *testing.T) { // Test case 2;
	// Test case 2; list namespaces.
	configFile := "/Users/abel.guzman/.kube/config"
	clientsetToSource, err := connect.ConnectNow(&configFile)
	nameSpacesList, err := query.ListNameSpaces(clientsetToSource)
	if err != nil {
		t.Errorf("Error getting namespace list: %v\n", err)
	}
	expectedType := "*v1.NamespaceList"
	if reflect.ValueOf(nameSpacesList).Type().String() != expectedType {
		t.Errorf("Test case 2; TestConnectNow(configFile) returned %s, expected %s", reflect.TypeOf(clientsetToSource).Kind().String(), expectedType)
	}
	if len(nameSpacesList.Items) <= 0 {
		t.Errorf("Test case 2; TestConnectNow(configFile) returned %v, expected > 0 namespaces", nameSpacesList.Items)
	}
}
