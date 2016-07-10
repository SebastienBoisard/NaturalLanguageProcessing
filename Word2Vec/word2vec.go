package main

import (
	"flag"
	"fmt"
)

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

/*
printf("WORD VECTOR estimation toolkit v 0.1c\n\n");
printf("Options:\n");
printf("Parameters for training:\n");













printf("\nExamples:\n");
printf("./word2vec -train data.txt -output vec.txt -size 200 -window 5 -sample 1e-4 -negative 5 -hs 0 -binary 0 -cbow 1 -iter 3\n\n");
return 0;
}


output_file[0] = 0;
save_vocab_file[0] = 0;
read_vocab_file[0] = 0;
if (cbow) alpha = 0.05;


*/
