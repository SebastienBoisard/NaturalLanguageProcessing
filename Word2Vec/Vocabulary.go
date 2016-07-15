package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Term ...
type Term struct {
	frequency int64
}

// BuildVocabulary2 builds the vocabulary list
func BuildVocabulary2(fileName string) {

	file, err := os.Open(fileName)

	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(file)

	// Set the Split method to ScanWords.
	scanner.Split(bufio.ScanWords)

	wordMap := make(map[string]*Term)

	// Scan all words from the file.
	for scanner.Scan() {
		word := scanner.Text()
		fmt.Println(word)

		term, ok := wordMap[word]
		if ok == true {
			term.frequency++
		} else {
			wordMap[word] = &Term{frequency: 1}
		}
	}

	for word, term := range wordMap {
		fmt.Println("word=", word, " frequency=", term.frequency)
	}

	file.Close()
}

// BuildVocabulary builds the vocabulary list
func BuildVocabulary(data []byte) []string {

	w := strings.FieldsFunc(string(data), func(r rune) bool {
		switch r {
		case ',', '.', ' ':
			return true
		}
		return false
	})

	return w
}
