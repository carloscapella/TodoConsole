package domain

import (
	"reflect"
	"testing"
)

func TestValidatePriority(t *testing.T) {
	if !ValidatePriority("low") || !ValidatePriority("medium") || !ValidatePriority("high") {
		t.Error("expected valid priorities to return true")
	}
	if ValidatePriority("urgent") {
		t.Error("expected invalid priority to return false")
	}
}

func TestGetDefaultPriority(t *testing.T) {
	if GetDefaultPriority() != PriorityMedium {
		t.Errorf("expected default priority to be medium, got %s", GetDefaultPriority())
	}
}

func TestNormalizeTags(t *testing.T) {
	in := []string{"Go", " go ", "Test", "test", ""}
	out := NormalizeTags(in)
	expected := []string{"go", "test"}
	if !reflect.DeepEqual(out, expected) && len(out) == 2 {
		t.Errorf("expected %v, got %v", expected, out)
	}
}
