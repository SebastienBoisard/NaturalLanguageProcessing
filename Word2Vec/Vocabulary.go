package main

import (
	"bufio"
	"os"
	"sort"
)

// Term ...
type Term struct {
	frequency int64
	word      string
}

const vocabHashSize int = 30000000 // Maximum 30 * 0.7 = 21M words in the vocabulary

// vocabMaxSize can be changed
var vocabMaxSize = 1000

const maxString = 100

var vocab []Term

var vocabHash []int

var vocabSize int

var minReduce int64

var minCount = int64(5)

var trainWords int64

func initializeVocabulary() {
	vocab = make([]Term, vocabMaxSize)
	vocabHash = make([]int, vocabHashSize)

	for a := 0; a < vocabHashSize; a++ {
		vocabHash[a] = -1
	}

	vocabSize = 0
	minReduce = 1
}

func learnVocabFromTrainFile(trainFileName string, vocab []Term) {

	for a := 0; a < vocabHashSize; a++ {
		vocabHash[a] = -1
	}

	vocabSize = 0

	file, err := os.Open(trainFileName)

	if err != nil {
		panic(err)
	}

	addWordToVocab("</s>")

	scanner := bufio.NewScanner(file)

	// Set the Split method to ScanWords.
	scanner.Split(bufio.ScanWords)

	// Scan all words from the file.
	for scanner.Scan() {
		word := scanner.Text()
		//fmt.Println(word)

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

// reduceVocab reduces the vocabulary by removing infrequent words
func reduceVocab() {

	nbWordRemoved := 0
	for a := 0; a < vocabSize; a++ {
		if vocab[a].frequency <= minReduce {
			vocab = append(vocab[:a], vocab[a+1:]...)
			nbWordRemoved++
		}
	}

	vocabSize = vocabSize - nbWordRemoved

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
// addWordToVocab returns the position of the new word in the vocabulary
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
// To get the same results than the original version, the hash value must be
// a uint64 type.
func getWordHash(word string) uint64 {
	var hash uint64
	for a := 0; a < len(word); a++ {
		hash = hash*257 + uint64(word[a])
	}
	// TODO: verify that it's useful to apply a modulo to a uint64
	hash = hash % uint64(vocabHashSize)

	return hash
}

// Terms is...
type Terms []Term

func (s Terms) Len() int {
	return len(s)
}
func (s Terms) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s Terms) Less(i, j int) bool {
	return s[i].frequency < s[j].frequency
}

// SortVocab sorts the vocabulary by frequency using word counts
func sortVocab() {
	// Sort the vocabulary and keep </s> at the first position

	// TODO: manage the first vocab </s>
	sort.Sort(Terms(vocab[0:vocabSize]))

	for a := 0; a < vocabHashSize; a++ {
		vocabHash[a] = -1
	}

	size := vocabSize
	trainWords = 0
	for a := 0; a < size; a++ {
		// Words occuring less than minCount times will be discarded from the vocab
		if vocab[a].frequency < minCount && a != 0 {
			vocab = append(vocab[:a], vocab[a+1:]...)
			vocabSize--
		} else {
			// Hash will be re-computed, as after the sorting it is not actual
			hash := getWordHash(vocab[a].word)
			for vocabHash[hash] != -1 {
				hash = (hash + 1) % uint64(vocabHashSize)
			}
			vocabHash[hash] = a
			trainWords += vocab[a].frequency
			// TODO: add vocabHash tests in Vocabulary_test.go
		}
	}

	/*
	   vocab = (struct vocab_word *)realloc(vocab, (vocab_size + 1) * sizeof(struct vocab_word));
	   // Allocate memory for the binary tree construction
	   for (a = 0; a < vocab_size; a++) {
	     vocab[a].code = (char *)calloc(MAX_CODE_LENGTH, sizeof(char));
	     vocab[a].point = (int *)calloc(MAX_CODE_LENGTH, sizeof(int));
	   }
	*/
}
