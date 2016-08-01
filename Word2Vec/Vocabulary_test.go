package main

import (
	"bufio"
	"io"
	"strings"
	"testing"
)

func TestSortVocab(t *testing.T) {

	vocab := initializeVocab()

	var position int

	position = vocab.addWordToVocab("</s>")
	vocab.vocabArray[position].frequency = -1

	position = vocab.addWordToVocab("first")
	vocab.vocabArray[position].frequency = 10

	position = vocab.addWordToVocab("second")
	vocab.vocabArray[position].frequency = 20

	position = vocab.addWordToVocab("third")
	vocab.vocabArray[position].frequency = 40

	vocab.sortVocab()

	expectedVocabSize := 4
	actualVocabSize := vocab.vocabSize

	if expectedVocabSize != actualVocabSize {
		t.Error("Expected", expectedVocabSize, "got", actualVocabSize)
	}

	var expectedWord, actualWord string

	expectedWord = "third"
	actualWord = vocab.vocabArray[1].word

	if expectedWord != actualWord {
		t.Error("Expected", expectedWord, "got", actualWord)
	}

	expectedWord = "second"
	actualWord = vocab.vocabArray[2].word

	if expectedWord != actualWord {
		t.Error("Expected", expectedWord, "got", actualWord)
	}

	expectedWord = "first"
	actualWord = vocab.vocabArray[3].word

	if expectedWord != actualWord {
		t.Error("Expected", expectedWord, "got", actualWord)
	}
}

func TestReduceVocab(t *testing.T) {

	vocab := initializeVocab()

	var position int

	position = vocab.addWordToVocab("first")
	vocab.vocabArray[position].frequency = 4

	position = vocab.addWordToVocab("second")
	vocab.vocabArray[position].frequency = 2

	position = vocab.addWordToVocab("third")
	vocab.vocabArray[position].frequency = 1

	var expectedMinReduce int64
	var expectedVocabSize int

	// Remove words with frequence <= 1 (ie word "third") from the vocabulary

	expectedMinReduce = 2
	expectedVocabSize = 2
	vocab.reduceVocab()

	if expectedMinReduce != vocab.minReduce {
		t.Error("Expected", expectedMinReduce, "got", vocab.minReduce)
	}

	if expectedVocabSize != vocab.vocabSize {
		t.Error("Expected", expectedVocabSize, "got", vocab.vocabSize)
	}

	// Remove words with frequence <= 2 (ie word "second") from the vocabulary

	expectedMinReduce = 3
	expectedVocabSize = 1
	vocab.reduceVocab()

	if expectedMinReduce != vocab.minReduce {
		t.Error("Expected", expectedMinReduce, "got", vocab.minReduce)
	}

	if expectedVocabSize != vocab.vocabSize {
		t.Error("Expected", expectedVocabSize, "got", vocab.vocabSize)
	}

	// Remove words with frequence <= 3 (ie no word) from the vocabulary

	expectedMinReduce = 4
	expectedVocabSize = 1
	vocab.reduceVocab()

	if expectedMinReduce != vocab.minReduce {
		t.Error("Expected", expectedMinReduce, "got", vocab.minReduce)
	}

	if expectedVocabSize != vocab.vocabSize {
		t.Error("Expected", expectedVocabSize, "got", vocab.vocabSize)
	}
}

func TestSearchVocab(t *testing.T) {
	var actualPosition, expectedPosition int

	vocab := initializeVocab()

	vocab.addWordToVocab("first")
	vocab.addWordToVocab("second")

	expectedPosition = 0
	actualPosition = vocab.searchVocab("first")

	if actualPosition != expectedPosition {
		t.Error("Expected", expectedPosition, "got", actualPosition)
	}

	expectedPosition = 1
	actualPosition = vocab.searchVocab("second")

	if actualPosition != expectedPosition {
		t.Error("Expected", expectedPosition, "got", actualPosition)
	}

	expectedPosition = -1
	actualPosition = vocab.searchVocab("third")

	if actualPosition != expectedPosition {
		t.Error("Expected", expectedPosition, "got", actualPosition)
	}
}

func TestAddWordToVocab(t *testing.T) {
	var actualPosition, expectedPosition int

	vocabulary := initializeVocab()

	expectedPosition = 0
	actualPosition = vocabulary.addWordToVocab("first")

	if actualPosition != expectedPosition {
		t.Error("Expected", expectedPosition, "got", actualPosition)
	}

	expectedPosition = 1
	actualPosition = vocabulary.addWordToVocab("second")

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

func TestLearnVocab(t *testing.T) {

	vocab := initializeVocab()

	vocab.learnVocab("vocabulary.txt")

	expectedVocab := []Term{
		{word: "</s>", frequency: 0},
		{word: "the", frequency: 10},
		{word: "of", frequency: 9},
		{word: "and", frequency: 7},
		{word: "a", frequency: 6},
		{word: "that", frequency: 5},
		{word: "is", frequency: 5},
	}

	actualVocab := vocab.vocabArray

	if len(actualVocab) != len(expectedVocab) {
		t.Error("Expected", len(expectedVocab), "got", len(actualVocab))
	}
}

func TestReadWord(t *testing.T) {

	const input = "First second third\nFourth\n"

	var tests = []struct {
		wantedWord  string
		wantedError error
	}{
		{"First", nil},
		{"second", nil},
		{"third", nil},
		{"</s>", nil},
		{"Fourth", nil},
		{"</s>", nil},
		{"", io.EOF},
	}

	reader := bufio.NewReader(strings.NewReader(input))

	for _, test := range tests {

		actualWord, actualErr := readWord(reader)
		if actualErr != test.wantedError || actualWord != test.wantedWord {
			t.Errorf("readWord(\"%s\", %v) = \"%s\", %v", test.wantedWord, test.wantedError, actualWord, actualErr)
		}
	}
}
