package main

import (
	"log"
	"os"
	"reflect"
	"testing"
)

func compareValues(value1, value2 reflect.Value) bool {

	if value1.Kind() != value2.Kind() {
		// TODO: manage properly this kind of error
		log.Fatal("value1 not the same type than value2")
	}

	switch value1.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value1.Int() == value2.Int()

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return value1.Uint() == value2.Uint()

	case reflect.Float32, reflect.Float64:
		return value1.Float() == value2.Float()

	case reflect.Bool:
		return value1.Bool() == value2.Bool()

	case reflect.String:
		return value1.String() == value2.String()

	default:
		// TODO: manage properly this kind of error
		log.Fatal("unknown type")
	}
	return false
}

func TestManageParameters(t *testing.T) {

	// bin\Word2Vec.exe -train_file text10.txt -output_file vectors.bin -size 200 -window 8 -sample 1e-4 -hs -negative 25 -num_threads 20 -iter 15 -binary
	os.Args = []string{"bin\\Word2Vec.exe", "-train_file", "text10.txt", "-output_file", "vectors.bin", "-size", "200", "-window", "8", "-sample", "1e-4",
		"-hs", "-num_threads", "20", "-iter", "15", "-binary", "-negative", "25"}

	actualParameters := manageParameters()

	// Warning: reflect.ValueOf(...) takes the value of the parameter at the moment, so we have to use the initialization after calling manageParameters().
	var tests = []struct {
		parameterName  string
		parameterValue reflect.Value
		wantedValue    reflect.Value
	}{
		{"-train_file", reflect.ValueOf(actualParameters.trainFile), reflect.ValueOf("text10.txt")},
		{"-output_file", reflect.ValueOf(actualParameters.outputFile), reflect.ValueOf("vectors.bin")},
		{"-cbow", reflect.ValueOf(actualParameters.cbowMode), reflect.ValueOf(false)},
		{"-hs", reflect.ValueOf(actualParameters.isHierarchicalSoftmaxActivated), reflect.ValueOf(true)},
		{"-size", reflect.ValueOf(actualParameters.layer1Size), reflect.ValueOf(200)},
		{"-window", reflect.ValueOf(actualParameters.windowSize), reflect.ValueOf(8)},
		{"-sample", reflect.ValueOf(actualParameters.occurrenceWordsThreshold), reflect.ValueOf(1e-4)},
		{"-negative", reflect.ValueOf(actualParameters.numberOfNegativeExamples), reflect.ValueOf(25)},
		{"-num_threads", reflect.ValueOf(actualParameters.numberOfThreads), reflect.ValueOf(20)},
		{"-iter", reflect.ValueOf(actualParameters.numberOfIterations), reflect.ValueOf(15)},
		{"-binary", reflect.ValueOf(actualParameters.binaryMode), reflect.ValueOf(true)},
	}

	for _, test := range tests {
		if compareValues(test.parameterValue, test.wantedValue) == false {
			t.Errorf("manageParameters(\"%s %v\") = %v", test.parameterName, test.wantedValue, test.parameterValue)
		}
	}
}
