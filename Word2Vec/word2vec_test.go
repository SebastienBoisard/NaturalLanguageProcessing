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

func TestCreateExpTable(t *testing.T) {

	actualExpTable := createExpTable()

	var tests = []struct {
		expTableIndex int
		wantedValue   float32
	}{
		{0, 0.0024726231566347748},
		{10, 0.00278699607588350773},
		{100, 0.00816257018595933914},
		{999, 0.99749761819839477539},
	}

	for _, test := range tests {
		if actualExpTable[test.expTableIndex] != test.wantedValue {
			t.Errorf("initializeExpTable()[%d] = %f", test.expTableIndex, test.wantedValue)
		}
	}
}

func TestCreateUnigramTable(t *testing.T) {

	corpusFile := "text10.corpus"
	initializeVocabulary()
	learnVocabFromTrainFile(corpusFile)
	initializeNetwork()

	unigramTable := createUnigramTable()

	expectedUnigramTableSize := int(1e8)
	actualUnigramTableSize := len(unigramTable)

	if actualUnigramTableSize != expectedUnigramTableSize {
		t.Error("Expected", expectedUnigramTableSize, "got", actualUnigramTableSize)
	}

	var expectedUnigramTableValue, actualUnigramTableValue int
	expectedUnigramTableValue = 3
	actualUnigramTableValue = unigramTable[6889208]

	if actualUnigramTableValue != expectedUnigramTableValue {
		t.Error("Expected", expectedUnigramTableValue, "got", actualUnigramTableValue)
	}

	expectedUnigramTableValue = 5
	actualUnigramTableValue = unigramTable[10564422]

	if actualUnigramTableValue != expectedUnigramTableValue {
		t.Error("Expected", expectedUnigramTableValue, "got", actualUnigramTableValue)
	}

	expectedUnigramTableValue = 7
	actualUnigramTableValue = unigramTable[15192095]

	if actualUnigramTableValue != expectedUnigramTableValue {
		t.Error("Expected", expectedUnigramTableValue, "got", actualUnigramTableValue)
	}

	expectedUnigramTableValue = 9
	actualUnigramTableValue = unigramTable[16638081]

	if actualUnigramTableValue != expectedUnigramTableValue {
		t.Error("Expected", expectedUnigramTableValue, "got", actualUnigramTableValue)
	}

	expectedUnigramTableValue = 11
	actualUnigramTableValue = unigramTable[18767860]

	if actualUnigramTableValue != expectedUnigramTableValue {
		t.Error("Expected", expectedUnigramTableValue, "got", actualUnigramTableValue)
	}

	expectedUnigramTableValue = 41
	actualUnigramTableValue = unigramTable[37003654]

	if actualUnigramTableValue != expectedUnigramTableValue {
		t.Error("Expected", expectedUnigramTableValue, "got", actualUnigramTableValue)
	}

	expectedUnigramTableValue = 261
	actualUnigramTableValue = unigramTable[76969030]

	if actualUnigramTableValue != expectedUnigramTableValue {
		t.Error("Expected", expectedUnigramTableValue, "got", actualUnigramTableValue)
	}

	expectedUnigramTableValue = 462
	actualUnigramTableValue = unigramTable[95301961]

	if actualUnigramTableValue != expectedUnigramTableValue {
		t.Error("Expected", expectedUnigramTableValue, "got", actualUnigramTableValue)
	}

	expectedUnigramTableValue = 524
	actualUnigramTableValue = unigramTable[99999999]

	if actualUnigramTableValue != expectedUnigramTableValue {
		t.Error("Expected", expectedUnigramTableValue, "got", actualUnigramTableValue)
	}
}
