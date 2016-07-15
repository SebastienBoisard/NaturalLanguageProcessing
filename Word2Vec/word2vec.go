package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"os"
)

func readVocabulary(vocabularyFile string) {

	data, err := ioutil.ReadFile(vocabularyFile)

	if err != nil {
		fmt.Println("Vocabulary file not found")
		panic(err)
	}

	w := BuildVocabulary(data)

	fmt.Println(w)
	/*
	   long long a, i = 0;
	   char c;
	   char word[MAX_STRING];
	   FILE *fin = fopen(read_vocab_file, "rb");
	   if (fin == NULL) {
	   printf("Vocabulary file not found\n");
	   exit(1);
	   }
	   for (a = 0; a < vocab_hash_size; a++) vocab_hash[a] = -1;
	   vocab_size = 0;
	   while (1) {
	   ReadWord(word, fin);
	   if (feof(fin)) break;
	   a = AddWordToVocab(word);
	   fscanf(fin, "%lld%c", &vocab[a].cn, &c);
	   i++;
	   }
	   SortVocab();
	   if (debug_mode > 0) {
	   printf("Vocab size: %lld\n", vocab_size);
	   printf("Words in train file: %lld\n", train_words);
	   }
	   fin = fopen(train_file, "rb");
	   if (fin == NULL) {
	   printf("ERROR: training data file not found!\n");
	   exit(1);
	   }
	   fseek(fin, 0, SEEK_END);
	   file_size = ftell(fin);
	   fclose(fin);
	*/
}

const tableSize = 1e8

var vocabSize = 0

const vocabHashSize int = 30000000 // Maximum 30 * 0.7 = 21M words in the vocabulary

var vocabHash [vocabHashSize]int

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

	/*
	   	char word[MAX_STRING];
	     long long a, i;

	     AddWordToVocab((char *)"</s>");
	     while (1) {
	       ReadWord(word, fin);
	       if (feof(fin)) break;
	       train_words++;
	       if ((debug_mode > 1) && (train_words % 100000 == 0)) {
	         printf("%lldK%c", train_words / 1000, 13);
	         fflush(stdout);
	       }
	       i = SearchVocab(word);
	       if (i == -1) {
	         a = AddWordToVocab(word);
	         vocab[a].cn = 1;
	       } else vocab[i].cn++;
	       if (vocab_size > vocab_hash_size * 0.7) ReduceVocab();
	     }
	     SortVocab();
	     if (debug_mode > 0) {
	       printf("Vocab size: %lld\n", vocab_size);
	       printf("Words in train file: %lld\n", train_words);
	     }
	     file_size = ftell(fin);
	     fclose(fin);
	*/
}

func addWordToVocab(word string) int {
	return -1
}

func searchVocab(word string) int {
	return -1
}

func reduceVocab() {
}

func main() {

	BuildVocabulary2("vocabulary.txt")

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
