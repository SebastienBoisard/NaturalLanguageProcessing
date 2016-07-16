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

/*
func TestLearnVocabFromTrainFile(t *testing.T) {
	initializeVocabulary();

	learnVocabFromTrainFile("vocabulary.txt", vocab) {
}*/

func TestGetWordHash(t *testing.T) {
	var actualHash, expectedHash uint64

	expectedHash = 18516756
	actualHash = getWordHash("hello")

	if actualHash != expectedHash {
		t.Error("Expected ", expectedHash, "got", actualHash)
	}

	expectedHash = 22177351
	actualHash = getWordHash("production")

	if actualHash != expectedHash {
		t.Error("Expected ", expectedHash, "got", actualHash)
	}

	expectedHash = 26318739
	actualHash = getWordHash("Antidisestablishmentarianism")

	if actualHash != expectedHash {
		t.Error("Expected ", expectedHash, "got", actualHash)
	}

	expectedHash = 3568563
	actualHash = getWordHash("antidisestablishmentarianism")

	if actualHash != expectedHash {
		t.Error("Expected ", expectedHash, "got", actualHash)
	}

	expectedHash = 97
	actualHash = getWordHash("a")

	if actualHash != expectedHash {
		t.Error("Expected ", expectedHash, "got", actualHash)
	}
}
