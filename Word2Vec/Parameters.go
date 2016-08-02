package main

import "flag"

// TODO: add comments to explain the purpose of each of those parameters

// Parameters contains all the command-line parameters of Word2Vec
type Parameters struct {
	trainFile                      string
	outputFile                     string
	layer1Size                     int
	windowSize                     int
	occurrenceWordsThreshold       float64
	isHierarchicalSoftmaxActivated bool
	numberOfNegativeExamples       int
	numberOfThreads                int
	numberOfIterations             int
	minWordOccurrencesThreshold    int
	startingLearningRate           float32
	numberOfClasses                int
	debugMode                      int
	binaryMode                     bool
	saveVocabFile                  string
	readVocabFile                  string
	cbowMode                       bool
}

// manageParameters parses the command line parameters of Word2Vec
// Warning: a boolean parameter is set to true if it appears in the commande line (e.g. -hs), otherwise it is set to false.
func manageParameters() Parameters {

	parameters := Parameters{}

	// train_file
	flag.StringVar(&(parameters.trainFile), "train_file", "", "Use text data from a file to train the model")

	// output_file
	flag.StringVar(&parameters.outputFile, "output_file", "", "Use a file to save the resulting word vectors / word clusters")

	// layer1Size
	flag.IntVar(&parameters.layer1Size, "size", 100, "Set size of word vectors; default is 100")

	// window
	flag.IntVar(&parameters.windowSize, "window", 5, "Set max skip length between words; default is 5")

	// sample
	flag.Float64Var(&parameters.occurrenceWordsThreshold, "sample", 1e-3, "Set threshold for occurrence of words. Those that appear with higher frequency in the training data will be randomly down-sampled; default is 1e-3, useful range is (0, 1e-5)")

	// hs
	flag.BoolVar(&parameters.isHierarchicalSoftmaxActivated, "hs", false, "Use Hierarchical Softmax; default is false (not used)")

	// negative
	flag.IntVar(&parameters.numberOfNegativeExamples, "negative", 5, "Number of negative examples; default is 5, common values are 3 - 10 (0 = not used)")

	// num_threads
	flag.IntVar(&parameters.numberOfThreads, "num_threads", 12, "Number of threads to use (default 12)")

	// iter
	flag.IntVar(&parameters.numberOfIterations, "iter", 5, "Run more training iterations (default 5)")

	// min_count
	flag.IntVar(&parameters.minWordOccurrencesThreshold, "min-count", 5, "This will discard words that appear less than <int> times; default is 5")

	// alpha
	var startingLearningRate64 float64
	flag.Float64Var(&startingLearningRate64, "alpha", 0.025, "Set the starting learning rate; default is 0.025 for skip-gram and 0.05 for CBOW")
	parameters.startingLearningRate = float32(startingLearningRate64)

	// classes
	flag.IntVar(&parameters.numberOfClasses, "classes", 0, "Output word classes rather than word vectors; default number of classes is 0 (vectors are written)")

	// debug_mode
	flag.IntVar(&parameters.debugMode, "debug", 2, "Set the debug mode (default = 2 = more info during training)")

	// binary
	flag.BoolVar(&parameters.binaryMode, "binary", false, "Save the resulting vectors in binary moded; default off")

	// save_vocab_file
	flag.StringVar(&parameters.saveVocabFile, "save-vocab", "", "The vocabulary will be saved to <file>")

	// read_vocab_file
	flag.StringVar(&parameters.readVocabFile, "read-vocab", "", "The vocabulary will be read from <file>, not constructed from the training data")

	// cbow
	flag.BoolVar(&parameters.cbowMode, "cbow", false, "Use the continuous bag of words model; default is 1 (use 0 for skip-gram model)")

	flag.Parse()

	return parameters
}
