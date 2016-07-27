package main

import "flag"

// TODO: add comments to explain the purpose of each of those parameters

var trainFile string
var outputFile string
var layer1Size int
var windowSize int
var occurrenceWordsThreshold float64
var isHierarchicalSoftmaxActivated bool
var numberOfNegativeExamples int
var numberOfThreads int
var numberOfIterations int
var minWordOccurrencesThreshold int
var startingLearningRate float64
var numberOfClasses int
var debugMode int
var binaryMode bool
var saveVocabFile string
var readVocabFile string
var cbowMode bool

// manageParameters parses the command line parameters of Word2Vec
// Warning: a boolean parameter is set to true if it appears in the commande line (e.g. -hs), otherwise it is set to false.
func manageParameters() {

	// train_file
	flag.StringVar(&trainFile, "train_file", "", "Use text data from a file to train the model")

	// output_file
	flag.StringVar(&outputFile, "output_file", "", "Use a file to save the resulting word vectors / word clusters")

	// layer1Size
	flag.IntVar(&layer1Size, "size", 100, "Set size of word vectors; default is 100")

	// window
	flag.IntVar(&windowSize, "window", 5, "Set max skip length between words; default is 5")

	// sample
	flag.Float64Var(&occurrenceWordsThreshold, "sample", 1e-3, "Set threshold for occurrence of words. Those that appear with higher frequency in the training data will be randomly down-sampled; default is 1e-3, useful range is (0, 1e-5)")

	// hs
	flag.BoolVar(&isHierarchicalSoftmaxActivated, "hs", false, "Use Hierarchical Softmax; default is false (not used)")

	// negative
	flag.IntVar(&numberOfNegativeExamples, "negative", 5, "Number of negative examples; default is 5, common values are 3 - 10 (0 = not used)")

	// num_threads
	flag.IntVar(&numberOfThreads, "num_threads", 12, "Number of threads to use (default 12)")

	// iter
	flag.IntVar(&numberOfIterations, "iter", 5, "Run more training iterations (default 5)")

	// min_count
	flag.IntVar(&minWordOccurrencesThreshold, "min-count", 5, "This will discard words that appear less than <int> times; default is 5")

	// alpha
	flag.Float64Var(&startingLearningRate, "alpha", 0.025, "Set the starting learning rate; default is 0.025 for skip-gram and 0.05 for CBOW")

	// classes
	flag.IntVar(&numberOfClasses, "classes", 0, "Output word classes rather than word vectors; default number of classes is 0 (vectors are written)")

	// debug_mode
	flag.IntVar(&debugMode, "debug", 2, "Set the debug mode (default = 2 = more info during training)")

	// binary
	flag.BoolVar(&binaryMode, "binary", false, "Save the resulting vectors in binary moded; default off")

	// save_vocab_file
	flag.StringVar(&saveVocabFile, "save-vocab", "", "The vocabulary will be saved to <file>")

	// read_vocab_file
	flag.StringVar(&readVocabFile, "read-vocab", "", "The vocabulary will be read from <file>, not constructed from the training data")

	// cbow
	flag.BoolVar(&cbowMode, "cbow", false, "Use the continuous bag of words model; default is 1 (use 0 for skip-gram model)")

	flag.Parse()
}
