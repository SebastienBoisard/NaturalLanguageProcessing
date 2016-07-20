package main

import "testing"

func TestCreateBinaryTree(t *testing.T) {

	initializeVocabulary()

	learnVocabFromTrainFile("vocabulary.txt")

	createBinaryTree()

	expectedVocab := []Term{
		{word: "</s>", frequency: 0, point: [maxCodeLength]int{5, 3, -7}, code: [maxCodeLength]byte{0, 1}, codelen: 2},
		{word: "the", frequency: 10, point: [maxCodeLength]int{5, 3, -6}, code: [maxCodeLength]byte{0, 0}, codelen: 2},
		{word: "of", frequency: 9, point: [maxCodeLength]int{5, 4, 2, -5}, code: [maxCodeLength]byte{1, 1, 0}, codelen: 3},
		{word: "and", frequency: 7, point: [maxCodeLength]int{5, 4, 1, -4}, code: [maxCodeLength]byte{1, 0, 1}, codelen: 3},
		{word: "a", frequency: 6, point: [maxCodeLength]int{5, 4, 1, -3}, code: [maxCodeLength]byte{1, 0, 0}, codelen: 3},
		{word: "that", frequency: 5, point: [maxCodeLength]int{5, 4, 2, 0, -2}, code: [maxCodeLength]byte{1, 1, 1, 1}, codelen: 4},
		{word: "is", frequency: 5, point: [maxCodeLength]int{5, 4, 2, 0, -1}, code: [maxCodeLength]byte{1, 1, 1, 0}, codelen: 4},
	}

	actualVocab := vocab

	if len(actualVocab) != len(expectedVocab) {
		t.Error("Expected", len(expectedVocab), "got", len(actualVocab))
	}

	for index := range expectedVocab {

		expectedWord := expectedVocab[index].word
		actualWord := actualVocab[index].word

		if actualWord != expectedWord {
			t.Error("Expected", expectedWord, "got", actualWord)
		}

		expectedFrequency := expectedVocab[index].frequency
		actualFrequency := actualVocab[index].frequency

		if actualFrequency != expectedFrequency {
			t.Error("Expected", expectedFrequency, "got", actualFrequency)
		}

		expectedPoint := expectedVocab[index].point
		actualPoint := actualVocab[index].point

		if actualPoint != expectedPoint {
			t.Error("Expected", expectedPoint, "got", actualPoint)
		}

		expectedCode := expectedVocab[index].code
		actualCode := actualVocab[index].code

		if actualCode != expectedCode {
			t.Error("Expected", expectedCode, "got", actualCode)
		}

		expectedCodelen := expectedVocab[index].codelen
		actualCodelen := actualVocab[index].codelen

		if actualCodelen != expectedCodelen {
			t.Error("Expected", expectedCodelen, "got", actualCodelen)
		}

	}
}

func TestInitializeExpTable(t *testing.T) {

	initializeExpTable()

	var expectedExp, actualExp float64

	expectedExp = 0.0024726231566347748
	actualExp = expTable[0]

	if expectedExp != actualExp {
		t.Error("Expected", expectedExp, "got", actualExp)
	}

	expectedExp = 0.00278699607588350773
	actualExp = expTable[10]

	if expectedExp != actualExp {
		t.Error("Expected", expectedExp, "got", actualExp)
	}

	expectedExp = 0.00816257018595933914
	actualExp = expTable[100]

	if expectedExp != actualExp {
		t.Error("Expected", expectedExp, "got", actualExp)
	}

	expectedExp = 0.99749761819839477539
	actualExp = expTable[999]

	if expectedExp != actualExp {
		t.Error("Expected", expectedExp, "got", actualExp)
	}
}
