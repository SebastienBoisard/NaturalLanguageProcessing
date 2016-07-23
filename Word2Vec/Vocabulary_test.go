package main

import (
	"bufio"
	"os"
	"testing"
)

func TestSortVocab(t *testing.T) {

	initializeVocabulary()

	var position int

	position = addWordToVocab("</s>")
	vocab[position].frequency = -1

	position = addWordToVocab("first")
	vocab[position].frequency = 10

	position = addWordToVocab("second")
	vocab[position].frequency = 20

	position = addWordToVocab("third")
	vocab[position].frequency = 40

	sortVocab()

	expectedVocabSize := 4
	actualVocabSize := vocabSize

	if expectedVocabSize != actualVocabSize {
		t.Error("Expected", expectedVocabSize, "got", actualVocabSize)
	}

	var expectedWord, actualWord string

	expectedWord = "third"
	actualWord = vocab[1].word

	if expectedWord != actualWord {
		t.Error("Expected", expectedWord, "got", actualWord)
	}

	expectedWord = "second"
	actualWord = vocab[2].word

	if expectedWord != actualWord {
		t.Error("Expected", expectedWord, "got", actualWord)
	}

	expectedWord = "first"
	actualWord = vocab[3].word

	if expectedWord != actualWord {
		t.Error("Expected", expectedWord, "got", actualWord)
	}
}

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

func TestLearnVocabFromTrainFile(t *testing.T) {

	initializeVocabulary()

	learnVocabFromTrainFile("vocabulary.txt")

	expectedVocab := []Term{
		{word: "</s>", frequency: 0},
		{word: "the", frequency: 10},
		{word: "of", frequency: 9},
		{word: "and", frequency: 7},
		{word: "a", frequency: 6},
		{word: "that", frequency: 5},
		{word: "is", frequency: 5},
	}

	actualVocab := vocab

	if len(actualVocab) != len(expectedVocab) {
		t.Error("Expected", len(expectedVocab), "got", len(actualVocab))
	}

	//	fmt.Println(vocab[:vocabSize])

	//createBinaryTree()

	//	fmt.Println(vocab[:vocabSize])
}

func TestReadWord(t *testing.T) {
	testFile := "text11.txt"
	f, err := os.Open(testFile)
	if err != nil {
		t.Error(testFile, " error", err)
	}
	defer f.Close()

	reader := bufio.NewReader(f)

	var expectedWord, actualWord string

	expectedWord = "First"
	actualWord, err = readWord(reader)

	if err != nil {
		t.Error("Expected nil")
	}

	if actualWord != expectedWord {
		t.Error("Expected", expectedWord, "got", actualWord)
	}

	expectedWord = "second"
	actualWord, err = readWord(reader)

	if err != nil {
		t.Error("Expected nil")
	}

	if actualWord != expectedWord {
		t.Error("Expected", expectedWord, "got", actualWord)
	}

	expectedWord = "third"
	actualWord, err = readWord(reader)

	if err != nil {
		t.Error("Expected nil")
	}

	if actualWord != expectedWord {
		t.Error("Expected", expectedWord, "got", actualWord)
	}

	expectedWord = "</s>"
	actualWord, err = readWord(reader)

	if err != nil {
		t.Error("Expected nil")
	}

	if actualWord != expectedWord {
		t.Error("Expected", expectedWord, "got", actualWord)
	}

	expectedWord = "Fourth"
	actualWord, err = readWord(reader)

	if err != nil {
		t.Error("Expected nil")
	}

	if actualWord != expectedWord {
		t.Error("Expected", expectedWord, "got", actualWord)
	}

	expectedWord = "</s>"
	actualWord, err = readWord(reader)

	if err != nil {
		t.Error("Expected nil")
	}

	if actualWord != expectedWord {
		t.Error("Expected", expectedWord, "got", actualWord)
	}
}
