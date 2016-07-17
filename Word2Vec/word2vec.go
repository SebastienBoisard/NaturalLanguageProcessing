package main

import (
	"flag"
	"fmt"
	"math"
)

const tableSize = 1e8

func createUnigramTable(vocab []Term) [tableSize]int {

	const power float64 = 0.75

	var unigramTable [tableSize]int

	var trainWordsPow float64
	trainWordsPow = 0.0

	for a := 0; a < vocabSize; a++ {
		trainWordsPow += math.Pow(float64(vocab[a].frequency), float64(power))
	}

	i := 0

	d1 := math.Pow(float64(vocab[i].frequency), power) / trainWordsPow

	for a := 0; a < tableSize; a++ {
		unigramTable[a] = i
		if float64(a)/float64(tableSize) > d1 {
			i++
			d1 += math.Pow(float64(vocab[i].frequency), power) / trainWordsPow
		}

		if i >= vocabSize {
			i = vocabSize - 1
		}
	}

	return unigramTable
}

const maxCodeLength = 40

// createBinaryTree creates binary Huffman tree using the word counts
// Frequent words will have short unique binary codes
func createBinaryTree() {

	var code [maxCodeLength]byte
	var point [maxCodeLength]int

	binary := make([]byte, vocabSize*2+1)
	parentNode := make([]int, vocabSize*2+1)

	count := make([]int64, vocabSize*2+1)

	for a := 0; a < vocabSize; a++ {
		count[a] = vocab[a].frequency
	}

	for a := vocabSize; a < vocabSize*2; a++ {
		count[a] = 1e15
	}

	pos1 := vocabSize - 1
	pos2 := vocabSize

	var min1i, min2i int

	// Following algorithm constructs the Huffman tree by adding one node at a time
	for a := 0; a < vocabSize-1; a++ {
		// First, find two smallest nodes 'min1, min2'
		if pos1 >= 0 {
			if count[pos1] < count[pos2] {
				min1i = pos1
				pos1--
			} else {
				min1i = pos2
				pos2++
			}
		} else {
			min1i = pos2
			pos2++
		}
		if pos1 >= 0 {
			if count[pos1] < count[pos2] {
				min2i = pos1
				pos1--
			} else {
				min2i = pos2
				pos2++
			}
		} else {
			min2i = pos2
			pos2++
		}
		count[vocabSize+a] = count[min1i] + count[min2i]
		parentNode[min1i] = vocabSize + a
		parentNode[min2i] = vocabSize + a
		binary[min2i] = 1
	}
	// Now assign binary code to each vocabulary word
	for a := 0; a < vocabSize; a++ {
		b := a
		i := 0
		for {
			code[i] = binary[b]
			point[i] = b
			i++
			b = parentNode[b]
			if b == vocabSize*2-2 {
				break
			}
		}
		vocab[a].codelen = byte(i)
		vocab[a].point[0] = vocabSize - 2
		for b := 0; b < i; b++ {
			vocab[a].code[i-b-1] = code[b]
			vocab[a].point[i-b] = point[b] - vocabSize
		}
	}
	//free(count);
	//free(binary);
	//free(parent_node);
}

func main() {

	// train_file
	trainFile := flag.String("train_file", "", "Use text data from a file to train the model")

	// output_file
	outputFile := flag.String("output_file", "", "Use a file to save the resulting word vectors / word clusters")

	// layer1_size
	wordVectorsSize := flag.Int("size", 100, "Set size of word vectors; default is 100")

	// window
	windowSize := flag.Int("window", 5, "Set max skip length between words; default is 5")

	// sample
	occurrenceWordsThreshold := flag.Float64("sample", 1e-3, "Set threshold for occurrence of words. Those that appear with higher frequency in the training data will be randomly down-sampled; default is 1e-3, useful range is (0, 1e-5)")

	// hs
	isHierarchicalSoftmaxActivated := flag.Bool("hs", false, "Use Hierarchical Softmax; default is false (not used)")

	// negative
	numberOfNegativeExamples := flag.Int("negative", 5, "Number of negative examples; default is 5, common values are 3 - 10 (0 = not used)")

	// num_threads
	numberOfThreads := flag.Int("num_threads", 12, "Number of threads to use (default 12)")

	// iter
	numberOfIterations := flag.Int("iter", 5, "Run more training iterations (default 5)")

	// min_count
	minWordOccurrencesThreshold := flag.Int("min-count", 5, "This will discard words that appear less than <int> times; default is 5")

	// alpha
	startingLearningRate := flag.Float64("alpha", 0.025, "Set the starting learning rate; default is 0.025 for skip-gram and 0.05 for CBOW")

	// classes
	numberOfClasses := flag.Int("classes", 0, "Output word classes rather than word vectors; default number of classes is 0 (vectors are written)")

	// debug_mode
	debugMode := flag.Int("debug", 2, "Set the debug mode (default = 2 = more info during training)")

	// binary
	binaryMode := flag.Bool("binary", false, "Save the resulting vectors in binary moded; default off")

	// save_vocab_file
	saveVocabFile := flag.String("save-vocab", "", "The vocabulary will be saved to <file>")

	// read_vocab_file
	readVocabFile := flag.String("read-vocab", "", "The vocabulary will be read from <file>, not constructed from the training data")

	// cbow
	cbowMode := flag.Bool("cbow", true, "Use the continuous bag of words model; default is 1 (use 0 for skip-gram model)")

	flag.Parse()

	fmt.Println("train_file:", *trainFile)
	fmt.Println("output_file:", *outputFile)
	fmt.Println("word_vectors_size:", *wordVectorsSize)
	fmt.Println("window_size:", *windowSize)
	fmt.Println("sample:", *occurrenceWordsThreshold)
	fmt.Println("hs:", *isHierarchicalSoftmaxActivated)
	fmt.Println("negative:", *numberOfNegativeExamples)
	fmt.Println("num_threads:", *numberOfThreads)
	fmt.Println("iter:", *numberOfIterations)
	fmt.Println("min-count:", *minWordOccurrencesThreshold)

	fmt.Println("alpha:", *startingLearningRate)
	fmt.Println("classes:", *numberOfClasses)
	fmt.Println("debug:", *debugMode)
	fmt.Println("binary:", *binaryMode)
	fmt.Println("save-vocab:", *saveVocabFile)
	fmt.Println("read-vocab:", *readVocabFile)
	fmt.Println("cbow:", *cbowMode)

}
