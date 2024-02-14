package DAO

import (
	"testing"
)

func TestDiffWithName(t *testing.T) {
	// Initialize a DiffWithName instance
	diff := DiffWithName{
		Name:           "TestName",
		Namespace:      "TestNamespace",
		Diff:           []string{"diff1", "diff2"},
		PropertyName:   "TestProperty",
		MessageHeading: "TestHeading",
		SourceMessage:  "SourceMessage",
		TargetMessage:  "TargetMessage",
	}

	// Assert that the fields are set correctly
	if diff.Name != "TestName" {
		t.Errorf("Expected Name to be 'TestName', got %s", diff.Name)
	}
	if diff.Namespace != "TestNamespace" {
		t.Errorf("Expected Namespace to be 'TestNamespace', got %s", diff.Namespace)
	}
	if len(diff.Diff) != 2 || diff.Diff[0] != "diff1" || diff.Diff[1] != "diff2" {
		t.Errorf("Expected Diff to be ['diff1', 'diff2'], got %v", diff.Diff)
	}
	if diff.PropertyName != "TestProperty" {
		t.Errorf("Expected PropertyName to be 'TestProperty', got %s", diff.PropertyName)
	}
	if diff.MessageHeading != "TestHeading" {
		t.Errorf("Expected MessageHeading to be 'TestHeading', got %s", diff.MessageHeading)
	}
	if diff.SourceMessage != "SourceMessage" {
		t.Errorf("Expected SourceMessage to be 'SourceMessage', got %s", diff.SourceMessage)
	}
	if diff.TargetMessage != "TargetMessage" {
		t.Errorf("Expected TargetMessage to be 'TargetMessage', got %s", diff.TargetMessage)
	}
}
