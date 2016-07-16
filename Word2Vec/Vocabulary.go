package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Term ...
type Term struct {
	frequency int64
	word      string
}

// BuildVocabulary2 builds the vocabulary list
func BuildVocabulary2(fileName string) {

	file, err := os.Open(fileName)

	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(file)

	// Set the Split method to ScanWords.
	scanner.Split(bufio.ScanWords)

	wordMap := make(map[string]*Term)

	// Scan all words from the file.
	for scanner.Scan() {
		word := scanner.Text()
		fmt.Println(word)

		term, ok := wordMap[word]
		if ok == true {
			term.frequency++
		} else {
			wordMap[word] = &Term{frequency: 1}
		}
	}

	for word, term := range wordMap {
		fmt.Println("word=", word, " frequency=", term.frequency)
	}

	file.Close()
}

// BuildVocabulary builds the vocabulary list
func BuildVocabulary(data []byte) []string {

	w := strings.FieldsFunc(string(data), func(r rune) bool {
		switch r {
		case ',', '.', ' ':
			return true
		}
		return false
	})

	return w
}

const vocabHashSize int = 30000000 // Maximum 30 * 0.7 = 21M words in the vocabulary

// vocabMaxSize can be changed
var vocabMaxSize = 1000

var vocab []Term
var vocabHash []int

func initializeVocabulary() {
	vocab = make([]Term, vocabMaxSize)
	vocabHash = make([]int, vocabHashSize)
}

func learnVocabFromTrainFile(trainFileName string, vocab []Term) {

	for a := 0; a < vocabHashSize; a++ {
		vocabHash[a] = -1
	}

	file, err := os.Open(trainFileName)

	if err != nil {
		panic(err)
	}

	vocabSize = 0

	addWordToVocab("</s>")

	scanner := bufio.NewScanner(file)

	// Set the Split method to ScanWords.
	scanner.Split(bufio.ScanWords)

	// Scan all words from the file.
	for scanner.Scan() {
		word := scanner.Text()
		fmt.Println(word)

		pos := searchVocab(word)
		if pos == -1 {
			pos = addWordToVocab(word)
			vocab[pos].frequency = 1
		} else {
			vocab[pos].frequency++
		}

		if float32(vocabSize) > float32(vocabHashSize)*0.7 {
			reduceVocab()
		}
	}
}

// searchVocab returns position of a word in the vocabulary; if the word is not found, returns -1
func searchVocab(word string) int {
	hash := getWordHash(word)
	for {
		if vocabHash[hash] == -1 {
			return -1
		}
		if word == vocab[vocabHash[hash]].word {
			return vocabHash[hash]
		}
		hash = (hash + 1) % uint64(vocabHashSize)
	}
}

// addWordToVocab adds a word to the vocabulary
func addWordToVocab(word string) int {

	wordLength := len(word) + 1

	if wordLength > maxString {
		wordLength = maxString
	}

	vocab[vocabSize].word = word
	vocab[vocabSize].frequency = 0
	vocabSize++

	hash := getWordHash(word)

	for vocabHash[hash] != -1 {
		hash = (hash + 1) % uint64(vocabHashSize)
	}

	vocabHash[hash] = vocabSize - 1

	return vocabSize - 1
}

// getWordHash returns hash value of a word
func getWordHash(word string) uint64 {
	var hash uint64
	for a := 0; a < len(word); a++ {
		hash = hash*257 + uint64(word[a])
	}
	hash = hash % uint64(vocabHashSize)
	return hash
}

var minReduce int64 = 1

// reduceVocab reduces the vocabulary by removing infrequent tokens
func reduceVocab() {

	for a := 0; a < vocabSize; a++ {
		if vocab[a].frequency <= minReduce {
			//append(vocab[:a], vocab[a+1:])
		}
	}

	vocabSize = len(vocab)

	for a := 0; a < vocabHashSize; a++ {
		vocabHash[a] = -1
	}
	for a := 0; a < vocabSize; a++ {
		// Hash will be re-computed, as it is not actual
		hash := getWordHash(vocab[a].word)
		for vocabHash[hash] != -1 {
			hash = (hash + 1) % uint64(vocabHashSize)
		}
		vocabHash[hash] = a
	}
	minReduce++
}