package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// Term ...
type Term struct {
	frequency int64
	word      string
	point     [maxCodeLength]int
	code      [maxCodeLength]byte
	codelen   byte
}

const vocabHashSize int = 30000000 // Maximum 30 * 0.7 = 21M words in the vocabulary
const maxString = 100

// vocabMaxSize can be changed
var vocabMaxSize = 1000

//var vocab []Term

//var vocabHash []int

//var vocabSize int

//var minReduce int64

var minCount = int64(5)

func (term Term) String() string {
	return fmt.Sprintf("{word=%s; frequency=%d; point=%v; code=%v; codelen=%d}\n", term.word, term.frequency, term.point[:term.codelen+1], term.code[:term.codelen], term.codelen)
}

// Vocab contains all the data of the corpus
type Vocab struct {
	vocabArray []Term
	vocabSize  int
	vocabHash  []int
	minReduce  int64
}

func initializeVocab() Vocab {

	vocab := Vocab{}

	vocab.vocabArray = make([]Term, 0, vocabMaxSize)

	vocab.vocabSize = 0

	vocab.vocabHash = make([]int, vocabHashSize)

	for a := 0; a < vocabHashSize; a++ {
		vocab.vocabHash[a] = -1
	}

	vocab.minReduce = 1

	return vocab
}

func (vocab *Vocab) learnVocab(trainFileName string) {

	file, err := os.Open(trainFileName)

	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(file)

	vocab.addWordToVocab("</s>")

	counter := 0

	for {
		word, err := readWord(reader)

		if err != nil {
			break
		}

		trainWords++

		pos := vocab.searchVocab(word)

		if pos == -1 {
			pos = vocab.addWordToVocab(word)
		}

		vocab.vocabArray[pos].frequency++

		if float32(vocab.vocabSize) > float32(vocabHashSize)*0.7 {
			vocab.reduceVocab()
		}

		counter++
	}

	vocab.sortVocab()
}

// reduceVocab reduces the vocabulary by removing infrequent words
func (vocab *Vocab) reduceVocab() {

	nbWordRemoved := 0
	for a := 0; a < vocab.vocabSize; a++ {
		if vocab.vocabArray[a].frequency <= vocab.minReduce {
			vocab.vocabArray = append(vocab.vocabArray[:a], vocab.vocabArray[a+1:]...)
			nbWordRemoved++
		}
	}

	vocab.vocabSize = vocab.vocabSize - nbWordRemoved

	for a := 0; a < vocabHashSize; a++ {
		vocab.vocabHash[a] = -1
	}
	for a := 0; a < vocab.vocabSize; a++ {
		// Hash will be re-computed, as it is not actual
		hash := getWordHash(vocab.vocabArray[a].word)
		for vocab.vocabHash[hash] != -1 {
			hash = (hash + 1) % uint64(vocabHashSize)
		}
		vocab.vocabHash[hash] = a
	}
	vocab.minReduce++
}

// searchVocab returns position of a word in the vocabulary; if the word is not found, returns -1
// 0 is for '\n'
func (vocab *Vocab) searchVocab(word string) int {
	hash := getWordHash(word)
	for {
		if vocab.vocabHash[hash] == -1 {
			return -1
		}

		if word == vocab.vocabArray[vocab.vocabHash[hash]].word {
			return vocab.vocabHash[hash]
		}

		hash = (hash + 1) % uint64(vocabHashSize)
	}
}

// addWordToVocab adds a word to the vocabulary
// addWordToVocab returns the position of the new word in the vocabulary
func (vocab *Vocab) addWordToVocab(word string) int {

	vocab.vocabArray = append(vocab.vocabArray, Term{word: word, frequency: 0})
	vocab.vocabSize++

	hash := getWordHash(word)

	for vocab.vocabHash[hash] != -1 {
		hash = (hash + 1) % uint64(vocabHashSize)
	}

	vocab.vocabHash[hash] = vocab.vocabSize - 1

	return vocab.vocabSize - 1
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
	if s[i].frequency == s[j].frequency {
		return s[i].word < s[j].word
	}
	return s[i].frequency > s[j].frequency
}

// SortVocab sorts the vocabulary by frequency using word counts
func (vocab *Vocab) sortVocab() {
	// fmt.Println("sortVocab BEGIN")
	// Sort the vocabulary and keep </s> at the first position

	// for index, currentVocab := range vocab {
	// 	fmt.Println("SortVocab before sort vocab[", index, "].word=", currentVocab.word, "frequency=", currentVocab.frequency)
	// }

	// TODO: manage the first vocab </s>
	sort.Sort(Terms(vocab.vocabArray[1:]))

	// for index, currentVocab := range vocab {
	// 	fmt.Println("SortVocab after sort vocab[", index, "].word=", currentVocab.word, "frequency=", currentVocab.frequency)
	// }

	for a := 0; a < vocabHashSize; a++ {
		vocab.vocabHash[a] = -1
	}

	for index, currentVocab := range vocab.vocabArray[1:] {
		if currentVocab.frequency < minCount {
			// fmt.Println("SortVocab delete vocab[", index, "].word=", currentVocab.word)
			vocab.vocabArray = vocab.vocabArray[:index+1]
			break
		}
	}

	trainWords = 0

	for index, currentVocab := range vocab.vocabArray {
		hash := getWordHash(currentVocab.word)
		for vocab.vocabHash[hash] != -1 {
			hash = (hash + 1) % uint64(vocabHashSize)
		}
		vocab.vocabHash[hash] = index
		trainWords += int(currentVocab.frequency)
		// fmt.Println("SortVocab keep vocab[", index, "].word=", currentVocab.word, "with new hash=", hash)
	}

	vocab.vocabSize = len(vocab.vocabArray)
}

// Reads a single word from a file, assuming space + tab + EOL to be word boundaries
func readWord(reader *bufio.Reader) (string, error) {
	var word string

	for {
		ch, err := reader.ReadByte()
		if err != nil {
			isEndFile = true
			return word, err
		}
		if ch == 13 || ch == ',' || ch == '.' {
			continue
		}
		if ch == ' ' || ch == '\t' || ch == '\n' {
			if len(word) > 0 {
				if ch == '\n' {
					reader.UnreadByte()
				}
				break
			}
			if ch == '\n' {
				word = "</s>"
				break
			}
			continue
		}
		word += string(ch)
	}
	return word, nil
}

// Reads a word and returns its index in the vocabulary
func (vocab *Vocab) readWordIndex(reader *bufio.Reader) int {
	word, _ := readWord(reader)
	//  if (feof(fin)) return -1;
	// fmt.Println("readWordIndex word=", word)
	return vocab.searchVocab(word)
}
