package main

import "testing"

func TestVocabularyBuilding(t *testing.T) {
	expected := []string{"le", "petit", "poucet"}
	actual := BuildVocabulary([]byte("le petit poucet"))[:]
	if actual == nil {
		t.Error("Test failed")
	}

	if testEq(expected, actual) == false {
		t.Error("Test failed")
	}
}

func testEq(a, b []string) bool {

	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
