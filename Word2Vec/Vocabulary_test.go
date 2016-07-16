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

func TestReduceVocab(t *testing.T) {

	initializeVocabulary()

	var position int

	position = addWordToVocab("first")
	vocab[position].frequency = 4

	position = addWordToVocab("second")
	vocab[position].frequency = 2

	position = addWordToVocab("third")
	vocab[position].frequency = 1

	var expectedMinReduce int64
	var expectedVocabSize int

	// Remove words with frequence <= 1 (ie word "third") from the vocabulary

	expectedMinReduce = 2
	expectedVocabSize = 2
	reduceVocab()

	if expectedMinReduce != minReduce {
		t.Error("Expected", expectedMinReduce, "got", minReduce)
	}

	if expectedVocabSize != vocabSize {
		t.Error("Expected", expectedVocabSize, "got", vocabSize)
	}

	// Remove words with frequence <= 2 (ie word "second") from the vocabulary

	expectedMinReduce = 3
	expectedVocabSize = 1
	reduceVocab()

	if expectedMinReduce != minReduce {
		t.Error("Expected", expectedMinReduce, "got", minReduce)
	}

	if expectedVocabSize != vocabSize {
		t.Error("Expected", expectedVocabSize, "got", vocabSize)
	}

	// Remove words with frequence <= 3 (ie no word) from the vocabulary

	expectedMinReduce = 4
	expectedVocabSize = 1
	reduceVocab()

	if expectedMinReduce != minReduce {
		t.Error("Expected", expectedMinReduce, "got", minReduce)
	}

	if expectedVocabSize != vocabSize {
		t.Error("Expected", expectedVocabSize, "got", vocabSize)
	}
}

func TestSearchVocab(t *testing.T) {
	var actualPosition, expectedPosition int

	initializeVocabulary()

	addWordToVocab("first")
	addWordToVocab("second")

	expectedPosition = 0
	actualPosition = searchVocab("first")

	if actualPosition != expectedPosition {
		t.Error("Expected", expectedPosition, "got", actualPosition)
	}

	expectedPosition = 1
	actualPosition = searchVocab("second")

	if actualPosition != expectedPosition {
		t.Error("Expected", expectedPosition, "got", actualPosition)
	}

	expectedPosition = -1
	actualPosition = searchVocab("third")

	if actualPosition != expectedPosition {
		t.Error("Expected", expectedPosition, "got", actualPosition)
	}
}

func TestAddWordToVocab(t *testing.T) {
	var actualPosition, expectedPosition int

	initializeVocabulary()

	expectedPosition = 0
	actualPosition = addWordToVocab("first")

	if actualPosition != expectedPosition {
		t.Error("Expected", expectedPosition, "got", actualPosition)
	}

	expectedPosition = 1
	actualPosition = addWordToVocab("second")

	if actualPosition != expectedPosition {
		t.Error("Expected", expectedPosition, "got", actualPosition)
	}
}

func TestGetWordHash(t *testing.T) {
	var actualHash, expectedHash uint64

	expectedHash = 18516756
	actualHash = getWordHash("hello")

	if actualHash != expectedHash {
		t.Error("Expected", expectedHash, "got", actualHash)
	}

	expectedHash = 22177351
	actualHash = getWordHash("production")

	if actualHash != expectedHash {
		t.Error("Expected", expectedHash, "got", actualHash)
	}

	expectedHash = 26318739
	actualHash = getWordHash("Antidisestablishmentarianism")

	if actualHash != expectedHash {
		t.Error("Expected", expectedHash, "got", actualHash)
	}

	expectedHash = 3568563
	actualHash = getWordHash("antidisestablishmentarianism")

	if actualHash != expectedHash {
		t.Error("Expected", expectedHash, "got", actualHash)
	}

	expectedHash = 97
	actualHash = getWordHash("a")

	if actualHash != expectedHash {
		t.Error("Expected", expectedHash, "got", actualHash)
	}
}
