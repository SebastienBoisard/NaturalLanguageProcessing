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

// vocabMaxSize can be changed
var vocabMaxSize = 1000

const maxString = 100

var vocab []Term

var vocabHash []int

var vocabSize int

var minReduce int64

var minCount = int64(5)

func (term Term) String() string {
	return fmt.Sprintf("{word=%s; frequency=%d; point=%v; code=%v; codelen=%d}\n", term.word, term.frequency, term.point[:term.codelen+1], term.code[:term.codelen], term.codelen)
}

func initializeVocabulary() {
	vocab = make([]Term, 0, vocabMaxSize)
	vocabHash = make([]int, vocabHashSize)

	for a := 0; a < vocabHashSize; a++ {
		vocabHash[a] = -1
	}

	vocabSize = 0
	minReduce = 1
}

func learnVocabFromTrainFile(trainFileName string) {

	for a := 0; a < vocabHashSize; a++ {
		vocabHash[a] = -1
	}

	vocabSize = 0

	file, err := os.Open(trainFileName)

	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(file)

	addWordToVocab("</s>")

	counter := 0

	for {
		word, err := readWord(reader)

		if err != nil {
			fmt.Println("learnVocabFromTrainFile [", counter, "] break err=", err)
			break
		}

		// fmt.Println("learnVocabFromTrainFile [", counter, "] word=", word)

		// if (feof(fin)) break;

		trainWords++

		pos := searchVocab(word)
		// fmt.Println("learnVocabFromTrainFile [", counter, "] pos1=", pos)

		if pos == -1 {
			pos = addWordToVocab(word)
			// fmt.Println("learnVocabFromTrainFile [", counter, "] pos2=", pos)
		}

		vocab[pos].frequency++

		if float32(vocabSize) > float32(vocabHashSize)*0.7 {
			reduceVocab()
		}

		counter++

	}
	/*
		scanner := bufio.NewScanner(file)

		// Set the Split method to ScanWords.
		scanner.Split(bufio.ScanWords)

		counter := 0
		// Scan all words from the file.
		for scanner.Scan() {
			word := scanner.Text()

			fmt.Println("learnVocabFromTrainFile [", counter, "] word=", word)
			//fmt.Println(word)

			pos := searchVocab(word)

			fmt.Println("learnVocabFromTrainFile [", counter, "] pos1=", pos)

			if pos == -1 {
				pos = addWordToVocab(word)
				fmt.Println("learnVocabFromTrainFile [", counter, "] pos2=", pos)
			}

			vocab[pos].frequency++

			if float32(vocabSize) > float32(vocabHashSize)*0.7 {
				reduceVocab()
			}
			//fmt.Print("vocab[", pos, "]=", vocab[pos])
			counter++
		}
	*/
	//	vocab = vocab[:vocabSize]

	sortVocab()

	// fmt.Println("learnVocabFromTrainFile vocabSize=", vocabSize)
	// fmt.Println("learnVocabFromTrainFile trainWords=", trainWords)
}

// reduceVocab reduces the vocabulary by removing infrequent words
func reduceVocab() {
	// fmt.Println("reduceVocab BEGIN")

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

	// fmt.Println("reduceVocab END")
}

// searchVocab returns position of a word in the vocabulary; if the word is not found, returns -1
// 0 is for '\n'
func searchVocab(word string) int {
	hash := getWordHash(word)
	for {
		if vocabHash[hash] == -1 {
			return -1
		}

		if word == vocab[vocabHash[hash]].word {
			// fmt.Println("searchVocab word=", word, "hash=", hash, "vocabHash=", vocabHash[hash])
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

	vocab = append(vocab, Term{word: word, frequency: 0})
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
	if s[i].frequency == s[j].frequency {
		return s[i].word < s[j].word
	}
	return s[i].frequency > s[j].frequency
}

// SortVocab sorts the vocabulary by frequency using word counts
func sortVocab() {
	// fmt.Println("sortVocab BEGIN")
	// Sort the vocabulary and keep </s> at the first position

	// for index, currentVocab := range vocab {
	// 	fmt.Println("SortVocab before sort vocab[", index, "].word=", currentVocab.word, "frequency=", currentVocab.frequency)
	// }

	// TODO: manage the first vocab </s>
	sort.Sort(Terms(vocab[1:]))

	// for index, currentVocab := range vocab {
	// 	fmt.Println("SortVocab after sort vocab[", index, "].word=", currentVocab.word, "frequency=", currentVocab.frequency)
	// }

	for a := 0; a < vocabHashSize; a++ {
		vocabHash[a] = -1
	}

	for index, currentVocab := range vocab[1:] {
		if currentVocab.frequency < minCount {
			// fmt.Println("SortVocab delete vocab[", index, "].word=", currentVocab.word)
			vocab = vocab[:index+1]
			break
		}
	}

	trainWords = 0

	for index, currentVocab := range vocab {
		hash := getWordHash(currentVocab.word)
		for vocabHash[hash] != -1 {
			hash = (hash + 1) % uint64(vocabHashSize)
		}
		vocabHash[hash] = index
		trainWords += int(currentVocab.frequency)
		// fmt.Println("SortVocab keep vocab[", index, "].word=", currentVocab.word, "with new hash=", hash)
	}

	vocabSize = len(vocab)

	// fmt.Println("sortVocab END")

	/*
		//	size := vocabSize
		trainWords = 0
		for a := len(vocab)-1; a > 0; a-- {
			fmt.Println("sortVocab vocab=", currentVocab.word)
			//	for a := 0; a < size; a++ {
			// Words occuring less than minCount times will be discarded from the vocab
			if vocab[a].frequency < minCount && a != 0 {
				fmt.Print("------- begin remove vocab:", vocab[a])
				vocab = append(vocab[:a], vocab[a+1:]...)
				vocabSize--
				fmt.Println(vocab)
				fmt.Println("------- end remove vocab")
			} else {
				// Hash will be re-computed, as after the sorting it is not actual
				hash := getWordHash(vocab[a].word)
				for vocabHash[hash] != -1 {
					hash = (hash + 1) % uint64(vocabHashSize)
				}
				vocabHash[hash] = a + 1
				trainWords += vocab[a].frequency
				// TODO: add vocabHash tests in Vocabulary_test.go
			}
		}
	*/

	/*
	   vocab = (struct vocab_word *)realloc(vocab, (vocab_size + 1) * sizeof(struct vocab_word));
	   // Allocate memory for the binary tree construction
	   for (a = 0; a < vocab_size; a++) {
	     vocab[a].code = (char *)calloc(MAX_CODE_LENGTH, sizeof(char));
	     vocab[a].point = (int *)calloc(MAX_CODE_LENGTH, sizeof(int));
	   }
	*/
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
func readWordIndex(reader *bufio.Reader) int {
	word, _ := readWord(reader)
	//  if (feof(fin)) return -1;
	// fmt.Println("readWordIndex word=", word)
	return searchVocab(word)
}
