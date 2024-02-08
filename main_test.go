package main

import (
	"testing"

	corev1 "k8s.io/api/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestFilterNamespaces(t *testing.T) {
	// Create a sample namespace list
	namespaces := &corev1.NamespaceList{
		Items: []corev1.Namespace{
			{ObjectMeta: metav1.ObjectMeta{Name: "test"}},
			{ObjectMeta: metav1.ObjectMeta{Name: "example"}},
			{ObjectMeta: metav1.ObjectMeta{Name: "hello"}},
			{ObjectMeta: metav1.ObjectMeta{Name: "world"}},
		},
	}

	// Test case when there are matching namespaces
	result := filterNamespaces(namespaces, "*l*")
	expected := &corev1.NamespaceList{
		Items: []corev1.Namespace{
			{ObjectMeta: metav1.ObjectMeta{Name: "hello"}},
			{ObjectMeta: metav1.ObjectMeta{Name: "world"}},
			{ObjectMeta: metav1.ObjectMeta{Name: "example"}},
		},
	}
	if !equalNamespaceLists(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}

	// Test case when there are no matching namespaces
	result = filterNamespaces(namespaces, "abc")
	expected = &corev1.NamespaceList{}
	if !equalNamespaceLists(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}

	// Test case when the input namespace list is empty
	emptyNamespaces := &corev1.NamespaceList{}
	result = filterNamespaces(emptyNamespaces, "*")
	expected = &corev1.NamespaceList{}
	if !equalNamespaceLists(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

// Function to compare NamespaceLists
func equalNamespaceLists(a, b *corev1.NamespaceList) bool {
	if len(a.Items) != len(b.Items) {
		return false
	}

	namespacesA := make(map[string]int)
	namespacesB := make(map[string]int)

	for _, ns := range a.Items {
		namespacesA[ns.Name]++
	}

	for _, ns := range b.Items {
		namespacesB[ns.Name]++
	}

	for name, countA := range namespacesA {
		countB := namespacesB[name]
		if countA != countB {
			return false
		}
	}

	return true
}

func TestMatchWildcard(t *testing.T) {
	// Test case when the string matches the wildcard pattern
	result := matchWildcard("hello", "h*llo")
	if !result {
		t.Error("Expected true, but got false")
	}

	// Test case when the string does not match the wildcard pattern
	result = matchWildcard("world", "h*llo")
	if result {
		t.Error("Expected false, but got true")
	}

	// Test case when there is an error in pattern matching
	result = matchWildcard("hello", "[")
	if result {
		t.Error("Expected false, but got true")
	}
}

func TestDetectNamespacePattern(t *testing.T) {
	// Test case for empty pattern
	result := DetectNamespacePattern("")
	if result != "empty" {
		t.Errorf("Expected 'empty' but got '%s'", result)
	}

	// Test case for wildcard pattern
	result = DetectNamespacePattern("*.example.com")
	if result != "wildcard" {
		t.Errorf("Expected 'wildcard' but got '%s'", result)
	}

	// Test case for specific pattern
	result = DetectNamespacePattern("example.com")
	if result != "specific" {
		t.Errorf("Expected 'specific' but got '%s'", result)
	}
}
