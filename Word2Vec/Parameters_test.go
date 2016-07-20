package main

import (
	"os"
	"testing"
)

func TestManageParameters(t *testing.T) {

	os.Args = []string{"bin\\Word2Vec.exe", "-train_file", "text10.txt", "-output_file", "vectors.bin", "-size", "200", "-window", "8", "-sample", "1e-4",
		"-hs", "-negative", "25", "-num_threads", "20", "-iter", "15", "-binary"}

	expectedTrainFile := "text10.txt"
	actualTrainFile := &trainFile

	expectedOutputFile := "vectors.bin"
	actualOutputFile := &outputFile

	expectedCbowMode := true
	actualCbowMode := &cbowMode

	expectedLayer1Size := 200
	actualLayer1Size := &layer1Size

	expectedWindowSize := 8
	actualWindowSize := &windowSize

	expectedOccurrenceWordsThreshold := 1e-4
	actualOccurrenceWordsThreshold := &occurrenceWordsThreshold

	expectedIsHierarchicalSoftmaxActivated := true
	actualIsHierarchicalSoftmaxActivated := &isHierarchicalSoftmaxActivated

	expectedNumberOfNegativeExamples := 25
	actualNumberOfNegativeExamples := &numberOfNegativeExamples

	expectedNumberOfThreads := 20
	actualNumberOfThreads := &numberOfThreads

	expectedNumberOfIterations := 15
	actualNumberOfIterations := &numberOfIterations

	expectedBinaryMode := true
	actualBinaryMode := &binaryMode

	manageParameters()

	if expectedTrainFile != *actualTrainFile {
		t.Error("Expected", expectedTrainFile, "got", *actualTrainFile)
	}

	if expectedCbowMode != *actualCbowMode {
		t.Error("Expected", expectedCbowMode, "got", *actualCbowMode)
	}

	if expectedOutputFile != *actualOutputFile {
		t.Error("Expected", expectedOutputFile, "got", *actualOutputFile)
	}

	if expectedLayer1Size != *actualLayer1Size {
		t.Error("Expected", expectedLayer1Size, "got", *actualLayer1Size)
	}

	if expectedWindowSize != *actualWindowSize {
		t.Error("Expected", expectedWindowSize, "got", *actualWindowSize)
	}

	if expectedNumberOfNegativeExamples != *actualNumberOfNegativeExamples {
		t.Error("Expected", expectedNumberOfNegativeExamples, "got", *actualNumberOfNegativeExamples)
	}

	if expectedIsHierarchicalSoftmaxActivated != *actualIsHierarchicalSoftmaxActivated {
		t.Error("Expected", expectedIsHierarchicalSoftmaxActivated, "got", *actualIsHierarchicalSoftmaxActivated)
	}

	if expectedOccurrenceWordsThreshold != *actualOccurrenceWordsThreshold {
		t.Error("Expected", expectedOccurrenceWordsThreshold, "got", *actualOccurrenceWordsThreshold)
	}

	if expectedNumberOfThreads != *actualNumberOfThreads {
		t.Error("Expected", expectedNumberOfThreads, "got", *actualNumberOfThreads)
	}

	if expectedNumberOfIterations != *actualNumberOfIterations {
		t.Error("Expected", expectedNumberOfIterations, "got", *actualNumberOfIterations)
	}

	if expectedBinaryMode != *actualBinaryMode {
		t.Error("Expected", expectedBinaryMode, "got", *actualBinaryMode)
	}
}
