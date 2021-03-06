package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"math"
	"os"
)

var parameters Parameters

const tableSize = 1e8

var table []int

const maxCodeLength = 40

const maxSentenceLength = 1000

const maxExp = 6

const expTableSize = 1000

var expTable []float32

var fileSize int

var wordCountActual = 0

var learningRate = float32(0.025)

var trainWords int

var syn0, syn1, syn1neg []float32

var isEndFile = false

func createUnigramTable(vocabArray []Term) []int {

	const power float64 = 0.75

	unigramTable := make([]int, tableSize)

	var trainWordsPow float64

	vocabSize := len(vocabArray)

	for a := 0; a < vocabSize; a++ {
		trainWordsPow += math.Pow(float64(vocabArray[a].frequency), power)
	}

	i := 0

	d1 := math.Pow(float64(vocabArray[i].frequency), power) / trainWordsPow
	for a := 0; a < tableSize; a++ {
		unigramTable[a] = i
		if float64(a)/float64(tableSize) > d1 {
			i++
			d1 += math.Pow(float64(vocabArray[i].frequency), power) / trainWordsPow
		}

		if i >= vocabSize {
			i = vocabSize - 1
		}
	}

	return unigramTable
}

// createBinaryTree creates binary Huffman tree using the word counts
// Frequent words will have short unique binary codes
func createBinaryTree(vocabArray []Term) {

	var code [maxCodeLength]byte
	var point [maxCodeLength]int

	binary := make([]byte, len(vocabArray)*2+1)
	parentNode := make([]int, len(vocabArray)*2+1)

	count := make([]int64, len(vocabArray)*2+1)

	for a := 0; a < len(vocabArray); a++ {
		count[a] = vocabArray[a].frequency
	}

	for a := len(vocabArray); a < len(vocabArray)*2; a++ {
		count[a] = 1e15
	}

	pos1 := len(vocabArray) - 1
	pos2 := len(vocabArray)

	var min1i, min2i int

	// Following algorithm constructs the Huffman tree by adding one node at a time
	for a := 0; a < len(vocabArray)-1; a++ {
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
		count[len(vocabArray)+a] = count[min1i] + count[min2i]
		parentNode[min1i] = len(vocabArray) + a
		parentNode[min2i] = len(vocabArray) + a
		binary[min2i] = 1
	}
	// Now assign binary code to each vocabulary word
	for a := 0; a < len(vocabArray); a++ {
		b := a
		i := 0
		for {
			code[i] = binary[b]
			point[i] = b
			i++
			b = parentNode[b]
			if b == len(vocabArray)*2-2 {
				break
			}
		}
		vocabArray[a].codelen = byte(i)
		vocabArray[a].point[0] = len(vocabArray) - 2
		for b := 0; b < i; b++ {
			vocabArray[a].code[i-b-1] = code[b]
			vocabArray[a].point[i-b] = point[b] - len(vocabArray)
		}
	}
}

func initializeNetwork(vocabSize int, layer1Size int) {

	syn0 = make([]float32, vocabSize*layer1Size)

	if parameters.isHierarchicalSoftmaxActivated == true {
		syn1 = make([]float32, vocabSize*layer1Size)
	}

	var nextRandom uint64 = 1

	if parameters.numberOfNegativeExamples > 0 {

		syn1neg = make([]float32, vocabSize*layer1Size)
	}

	for a := 0; a < vocabSize; a++ {
		for b := 0; b < layer1Size; b++ {
			nextRandom = nextRandom*25214903917 + 11

			syn0[a*layer1Size+b] = (((float32(nextRandom & 0xFFFF)) / float32(65536)) - float32(0.5)) / float32(layer1Size)
		}
	}

}

func trainModelThread(id int, vocab Vocab) {

	var label int
	var target int

	nextRandom := uint64(id)

	var wordCount, lastWordCount int
	var sen [maxSentenceLength + 1]int

	sentencePosition := 0
	localIter := parameters.numberOfIterations

	neu1 := make([]float32, parameters.layer1Size)
	neu1e := make([]float32, parameters.layer1Size)

	isEndFile = false
	fi, err := os.Open(parameters.trainFile)
	if err != nil {
		log.Fatal(err)
	}
	defer fi.Close()

	filePosition := int64(fileSize / parameters.numberOfThreads * id)
	fi.Seek(filePosition, 0)

	// fmt.Println("trainModelThread[", id, "] numberOfThreads=", numberOfThreads, "filePosition=", filePosition)
	// fmt.Printf("trainModelThread[ %d ] learningRate= %.20f\n", id, learningRate)

	reader := bufio.NewReader(fi)

	sentenceLength := 0
	counter2 := 0
	for {
		counter2++

		// fmt.Println("trainModelThread[", id, "][", counter2, "] wordCount=", wordCount, "lastWordCount=", lastWordCount)

		// for c := 0; c < layer1Size; c++ {
		// 	fmt.Printf("trainModelThread[ %d ][ %d ] new1[ %d ]= %.20f\n", id, counter2, c, neu1[c])
		// }
		// idxNeu1 := 1
		// fmt.Printf("trainModelThread[ %d ][ %d ] 0 new1[ %d ]= %.20f\n", id, counter2, idxNeu1, neu1[idxNeu1])

		if wordCount-lastWordCount > 10000 {
			wordCountActual += wordCount - lastWordCount
			lastWordCount = wordCount

			fmt.Printf("%clearningRate: %f  Progress: %.2f%%\n", 13, learningRate, float32(wordCountActual)/float32(parameters.numberOfIterations*trainWords+1)*100)
			learningRate = float32(parameters.startingLearningRate * (1 - float32(wordCountActual)/float32(parameters.numberOfIterations*trainWords+1)))

			if learningRate < parameters.startingLearningRate*0.0001 {
				learningRate = parameters.startingLearningRate * 0.0001
			}
		}

		// fmt.Println("trainModelThread[", id, "][", counter2, "] sentenceLength=", sentenceLength)

		if sentenceLength == 0 {

			counter := 0
			for {
				counter++
				word := vocab.readWordIndex(reader)

				// if word > 0 {
				// 	fmt.Println("trainModelThread[", id, "][", counter2, "][", counter, "] word=", vocab[word].word, " (", word, ")")
				// } else {
				// 	fmt.Println("trainModelThread[", id, "][", counter2, "][", counter, "] word_id=", word)
				// }

				if isEndFile == true {
					break
				}

				// Test if the word was found
				if word == -1 {
					continue
				}

				wordCount++

				// Test if the word is '\n'
				if word == 0 {
					break
				}
				// The subsampling randomly discards frequent words while keeping the ranking same
				if parameters.occurrenceWordsThreshold > 0 {

					// fmt.Println("trainModelThread[", id, "] vocab[word].frequency=", vocab[word].frequency, "occurrenceWordsThreshold=", occurrenceWordsThreshold, "trainWords=", trainWords)

					ran := float32(math.Sqrt(float64(float32(vocab.vocabArray[word].frequency)/(float32(parameters.occurrenceWordsThreshold)*float32(trainWords))))+1.0) *
						(float32(parameters.occurrenceWordsThreshold) * float32(trainWords)) / float32(vocab.vocabArray[word].frequency)
					nextRandom = nextRandom*25214903917 + 11

					// fmt.Println("trainModelThread[", id, "] ran=", ran)
					// fmt.Println("trainModelThread[", id, "] nextRandom=", nextRandom)

					if ran < float32(nextRandom&0xFFFF)/float32(65536.0) {
						continue
					}
				}
				sen[sentenceLength] = word
				sentenceLength++
				if sentenceLength >= maxSentenceLength {
					break
				}
			}

			sentencePosition = 0

			// for i, v := range sen[:sentenceLength] {
			// 	fmt.Println("trainModelThread[", id, "][", counter2, "] sen[", i, "]=", v)
			// }
		}

		if isEndFile == true || wordCount > trainWords/parameters.numberOfThreads {
			wordCountActual += wordCount - lastWordCount
			localIter--
			if localIter == 0 {
				// fmt.Println("trainModelThread[", id, "][", counter2, "] locaIter==0 so break")
				break
			}
			wordCount = 0
			lastWordCount = 0
			sentenceLength = 0
			//fi.fseek(file_size / (long long)num_threads * (long long)id, 0)
			continue
		}

		word := sen[sentencePosition]
		// fmt.Println("trainModelThread[", id, "][", counter2, "] word=", word, "sentencePosition=", sentencePosition)

		if word == -1 {
			continue
		}
		for c := 0; c < parameters.layer1Size; c++ {
			neu1[c] = 0
		}
		for c := 0; c < parameters.layer1Size; c++ {
			neu1e[c] = 0
		}

		nextRandom = nextRandom*25214903917 + 11
		b := int(nextRandom % uint64(parameters.windowSize))

		// fmt.Println("trainModelThread[", id, "][", counter2, "] nextRandom=", nextRandom, "b=", b)

		if parameters.cbowMode == true {
			//train the cbow architecture

			// fmt.Println("trainModelThread[", id, "][", counter2, "] cbowMode on")

			// in -> hidden
			cw := 0

			for a := int(b); a < parameters.windowSize*2+1-int(b); a++ {
				if a != parameters.windowSize {
					c := sentencePosition - parameters.windowSize + a
					if c < 0 {
						continue
					}
					if c >= sentenceLength {
						continue
					}
					lastWord := sen[c]
					if lastWord == -1 {
						continue
					}

					// fmt.Printf("trainModelThread[ %d ][ %d ] 1a new1[ %d ]= %.20f\n", id, counter2, idxNeu1, neu1[idxNeu1])

					// fmt.Printf("trainModelThread[ %d ][ %d ] lastWord=%d syn0[%d]=%.20f\n", id, counter2, lastWord,
					// idxNeu1+lastWord*layer1Size, syn0[idxNeu1+lastWord*layer1Size])

					for c := 0; c < parameters.layer1Size; c++ {
						neu1[c] += syn0[c+lastWord*parameters.layer1Size]
						// fmt.Printf("trainModelThread[ %d ][ %d ] 1b new1[ %d ]= %.20f\n", id, counter2, c, neu1[c])
					}
					// fmt.Printf("trainModelThread[ %d ][ %d ] 1b new1[ %d ]= %.20f\n", id, counter2, idxNeu1, neu1[idxNeu1])
					cw++
				}
			}

			if cw > 0 {

				// fmt.Println("trainModelThread[", id, "][", counter2, "] cw > 0")

				for c := 0; c < parameters.layer1Size; c++ {
					neu1[c] /= float32(cw)
					// fmt.Printf("trainModelThread[ %d ][ %d ] 2 new1[ %d ]= %.20f\n", id, counter2, c, neu1[c])
				}
				// fmt.Printf("trainModelThread[ %d ][ %d ] 2 new1[ %d ]= %.20f\n", id, counter2, idxNeu1, neu1[idxNeu1])

				if parameters.isHierarchicalSoftmaxActivated == true {

					// fmt.Println("trainModelThread[", id, "][", counter2, "] isHierarchicalSoftmaxActivated")

					for d := 0; d < int(vocab.vocabArray[word].codelen); d++ {
						f := float32(0.0)
						l2 := vocab.vocabArray[word].point[d] * parameters.layer1Size
						// Propagate hidden -> output
						for c := 0; c < parameters.layer1Size; c++ {
							f += neu1[c] * syn1[c+l2]
						}
						if f <= -maxExp {
							continue
						} else {
							if f >= maxExp {
								continue
							} else {
								f = expTable[int((f+maxExp)*(expTableSize/maxExp/2))]
							}
						}

						// fmt.Printf("trainModelThread[ %d ][ %d ] hs f=%.20f\n", id, counter2, f)

						// 'g' is the gradient multiplied by the learning rate
						g := float32(1.0-float32(vocab.vocabArray[word].code[d])-f) * learningRate

						// fmt.Printf("trainModelThread[ %d ][ %d ] hs g=%.20f\n", id, counter2, g)

						// Propagate errors output -> hidden
						for c := 0; c < parameters.layer1Size; c++ {
							neu1e[c] += g * syn1[c+l2]
							// fmt.Printf("trainModelThread[ %d ][ %d ] hs neu1e[ %d ]= %.20f\n", id, counter2, c, neu1e[c])
						}
						// Learn weights hidden -> output
						for c := 0; c < parameters.layer1Size; c++ {
							syn1[c+l2] += g * neu1[c]
							// fmt.Printf("trainModelThread[ %d ][ %d ] hs syn1[ %d ]= %.20f\n", id, counter2, c+l2, syn1[c+l2])
						}
					}
				}

				// NEGATIVE SAMPLING
				if parameters.numberOfNegativeExamples > 0 {

					// fmt.Println("trainModelThread[", id, "][", counter2, "] numberOfNegativeExamples > 0")

					var label int
					var g float32
					for d := 0; d < parameters.numberOfNegativeExamples+1; d++ {
						if d == 0 {
							target = word
							label = 1
						} else {
							nextRandom = nextRandom*25214903917 + 11
							// fmt.Println("trainModelThread[", id, "] nextRandom=", nextRandom, "(nextRandom>>16)%tableSize=", (nextRandom>>16)%tableSize)
							target = table[(nextRandom>>16)%tableSize]
							if target == 0 {
								target = int(nextRandom%uint64(vocab.vocabSize-1)) + 1
							}
							if target == word {
								continue
							}
							label = 0
						}
						l2 := target * parameters.layer1Size
						// fmt.Printf("trainModelThread[ %d ][ %d ] negative l2=%d\n", id, counter2, l2)

						f := float32(0.0)
						for c := 0; c < parameters.layer1Size; c++ {
							f += neu1[c] * syn1neg[c+l2]
							// fmt.Printf("trainModelThread[ %d ][ %d ] negative f[%d]=%.20f  neu1[%d]=%.20f  syn1neg[%d]=%.20f   neu1[%d]*syn1neg[%d]=%.20f\n", id, counter2, c, f, c, neu1[c], c+l2, syn1neg[c+l2], c, c+l2, neu1[c] * syn1neg[c+l2])
						}

						// fmt.Printf("trainModelThread[ %d ][ %d ] negative f=%.20f\n", id, counter2, f)

						if f > maxExp {
							g = float32(label-1) * learningRate
							// fmt.Printf("trainModelThread[ %d ][ %d ] negative g1=%.20f\n", id, counter2, g)
						} else {
							if f < -maxExp {
								g = float32(label-0) * learningRate
								// fmt.Printf("trainModelThread[ %d ][ %d ] negative g2=%.20f\n", id, counter2, g)
							} else {
								expIdx := int((f + maxExp) * (expTableSize / maxExp / 2.0))
								g = (float32(label) - expTable[expIdx]) * learningRate

								//g = (label - expTable[(int)((f + MAX_EXP) * (EXP_TABLE_SIZE / MAX_EXP / 2))]) * learningRate;

								// fmt.Printf("trainModelThread[ %d ][ %d ] negative maxEp=%d  expTableSize=%d  f=%.20f  expIdx=%d  expTable[expIdx]=%.20f\n",
								// id, counter2, maxExp, expTableSize, f, expIdx, expTable[expIdx])

								// fmt.Printf("trainModelThread[ %d ][ %d ] negative f+maxEp=%f  expTableSize/maxExp/2=%d  (f+maxExp)*(expTableSize/maxExp/2)=%f\n",
								// id, counter2, (f+maxExp), (expTableSize/maxExp/2), ((f+maxExp)*(expTableSize/maxExp/2)))

								// fmt.Printf("trainModelThread[ %d ][ %d ] negative label=%d f=%.20f maxExp=%d expTableSize=%d learningRate=%.20f g3=%.20f\n",
								// id, counter2, label, f, maxExp, expTableSize, learningRate, g)
								// fmt.Printf("trainModelThread[ %d ][ %d ] negative g3=%.20f\n", id, counter2, g)
							}
						}

						// fmt.Printf("trainModelThread[ %d ][ %d ] negative neu1e[ %d ]= %.20f\n", id, counter2, 0, neu1e[0])
						for c := 0; c < parameters.layer1Size; c++ {
							neu1e[c] += g * syn1neg[c+l2]
							// fmt.Printf("trainModelThread[ %d ][ %d ] negative neu1e[ %d ]= %.20f\n", id, counter2, c, neu1e[c])
						}

						// fmt.Printf("trainModelThread[ %d ][ %d ] negative l2=%d g=%.20f syn1neg[ %d ]= %.20f\n", id, counter2, l2, g, l2, syn1neg[l2])
						for c := 0; c < parameters.layer1Size; c++ {
							syn1neg[c+l2] += g * neu1[c]
							// fmt.Printf("trainModelThread[ %d ][ %d ] negative syn1neg[%d]=%.20f  neu1[%d]=%.20f\n", id, counter2, c+l2, syn1neg[c+l2], c, neu1[c])
						}
					}
				}

				// hidden -> in
				for a := int(b); a < int(parameters.windowSize)*2+1-int(b); a++ {
					if a != int(parameters.windowSize) {
						c := sentencePosition - parameters.windowSize + a
						// fmt.Printf("trainModelThread[ %d ][ %d ] hidden a=%d c=%d\n", id, counter2, a, c)
						if c < 0 {
							continue
						}
						if c >= sentenceLength {
							continue
						}
						lastWord := sen[c]
						if lastWord == -1 {
							continue
						}
						for c := 0; c < parameters.layer1Size; c++ {
							syn0[c+lastWord*parameters.layer1Size] += neu1e[c]
						}
					}
				}
			}
		} else {
			//train skip-gram

			// fmt.Println("trainModelThread[", id, "] skip-gram")

			for a := b; a < parameters.windowSize*2+1-b; a++ {
				if a != parameters.windowSize {
					c := sentencePosition - parameters.windowSize + a
					if c < 0 {
						continue
					}
					if c >= sentenceLength {
						continue
					}
					lastWord := sen[c]
					if lastWord == -1 {
						continue
					}
					l1 := lastWord * parameters.layer1Size
					for c := 0; c < parameters.layer1Size; c++ {
						neu1e[c] = 0
					}
					// HIERARCHICAL SOFTMAX
					if parameters.isHierarchicalSoftmaxActivated == true {
						for d := 0; d < int(vocab.vocabArray[word].codelen); d++ {
							f := float32(0.0)
							l2 := vocab.vocabArray[word].point[d] * parameters.layer1Size
							// Propagate hidden -> output
							for c := 0; c < parameters.layer1Size; c++ {
								f += syn0[c+l1] * syn1[c+l2]
							}

							if f <= -maxExp {
								continue
							} else {
								if f >= maxExp {
									continue
								} else {
									f = expTable[(int)((f+maxExp)*(expTableSize/maxExp/2))]
								}
							}

							// 'g' is the gradient multiplied by the learning rate
							g := (1 - float32(vocab.vocabArray[word].code[d]) - f) * learningRate
							// Propagate errors output -> hidden
							for c := 0; c < parameters.layer1Size; c++ {
								neu1e[c] += g * syn1[c+l2]
							}
							// Learn weights hidden -> output
							for c := 0; c < parameters.layer1Size; c++ {
								syn1[c+l2] += g * syn0[c+l1]
							}
						}
					}

					// NEGATIVE SAMPLING
					if parameters.numberOfNegativeExamples > 0 {
						for d := 0; d < parameters.numberOfNegativeExamples+1; d++ {
							if d == 0 {
								target = word
								label = 1
							} else {
								nextRandom = nextRandom*25214903917 + 11
								target = table[(nextRandom>>16)%tableSize]
								if target == 0 {
									target = int(nextRandom%uint64(vocab.vocabSize-1)) + 1
								}
								if target == word {
									continue
								}
								label = 0
							}
							var g float32
							l2 := target * parameters.layer1Size
							f := float32(0.0)
							for c := 0; c < parameters.layer1Size; c++ {
								f += syn0[c+l1] * syn1neg[c+l2]
							}
							if f > maxExp {
								g = float32(label-1) * learningRate
							} else {
								if f < -maxExp {
									g = float32(label-0) * learningRate
								} else {
									g = float32(float32(label)-expTable[int((f+maxExp)*(expTableSize/maxExp/2))]) * learningRate
								}
							}
							for c := 0; c < parameters.layer1Size; c++ {
								neu1e[c] += g * syn1neg[c+l2]
							}
							for c := 0; c < parameters.layer1Size; c++ {
								syn1neg[c+l2] += g * syn0[c+l1]
							}
						}
					}

					// Learn weights input -> hidden
					for c := 0; c < parameters.layer1Size; c++ {
						syn0[c+l1] += neu1e[c]
					}
				}
			}
		}

		sentencePosition++
		if sentencePosition >= sentenceLength {
			sentenceLength = 0
			continue
		}
	}
}

func trainModel(vocab Vocab) {

	trainModelThread(0, vocab)

	//for (a = 0; a < num_threads; a++) pthread_create(&pt[a], NULL, TrainModelThread, (void *)a);
	//for (a = 0; a < num_threads; a++) pthread_join(pt[a], NULL);
}

func saveData(outputFileName string, vocab Vocab) {
	file, err := os.Create(outputFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	if parameters.numberOfClasses == 0 {
		saveWordsAsVectors(file, vocab, syn0)
	} else {
		saveWordsAsClasses(file, vocab, syn0)
	}
}

func saveWordsAsVectors(output io.Writer, vocab Vocab, syn0 []float32) {
	// Save the word vectors
	fmt.Fprintf(output, "%d %d\n", vocab.vocabSize, parameters.layer1Size)
	for a := 0; a < vocab.vocabSize; a++ {
		fmt.Fprintf(output, "%s ", vocab.vocabArray[a].word)
		if parameters.binaryMode == true {
			for b := 0; b < parameters.layer1Size; b++ {
				binary.Write(output, binary.LittleEndian, syn0[a*parameters.layer1Size+b])
			}
		} else {
			for b := 0; b < parameters.layer1Size; b++ {
				fmt.Fprintf(output, "%f ", syn0[a*parameters.layer1Size+b])
			}
		}
		fmt.Fprintf(output, "\n")
	}
}

func saveWordsAsClasses(output io.Writer, vocab Vocab, syn0 []float32) {
	// Run K-means on the word vectors
	// int clcn = classes, iter = 10, closeid;

	clcn := parameters.numberOfClasses
	numberOfIterations := 10
	var closeid int

	var closev, x float32
	cl := make([]int, vocab.vocabSize)
	centcn := make([]int, parameters.numberOfClasses)
	cent := make([]float32, parameters.numberOfClasses*parameters.layer1Size)

	for a := 0; a < vocab.vocabSize; a++ {
		cl[a] = a % clcn
	}

	for a := 0; a < numberOfIterations; a++ {
		for b := 0; b < clcn*parameters.layer1Size; b++ {
			cent[b] = 0
		}
		for b := 0; b < clcn; b++ {
			centcn[b] = 1
		}
		for c := 0; c < vocab.vocabSize; c++ {
			for d := 0; d < parameters.layer1Size; d++ {
				cent[parameters.layer1Size*cl[c]+d] += syn0[c*parameters.layer1Size+d]
			}
			centcn[cl[c]]++
		}
		for b := 0; b < clcn; b++ {
			closev = 0
			for c := 0; c < parameters.layer1Size; c++ {
				cent[parameters.layer1Size*b+c] = cent[parameters.layer1Size*b+c] / float32(centcn[b])
				closev += cent[parameters.layer1Size*b+c] * cent[parameters.layer1Size*b+c]
			}
			closev = float32(math.Sqrt(float64(closev)))
			for c := 0; c < parameters.layer1Size; c++ {
				cent[parameters.layer1Size*b+c] /= closev
			}
		}
		for c := 0; c < vocab.vocabSize; c++ {
			closev = -10
			closeid = 0
			for d := 0; d < clcn; d++ {
				x = 0
				for b := 0; b < parameters.layer1Size; b++ {
					x += cent[parameters.layer1Size*d+b] * syn0[c*parameters.layer1Size+b]
				}
				if x > closev {
					closev = x
					closeid = d
				}
			}
			cl[c] = closeid
		}
	}

	// Save the K-means classes
	for a := 0; a < vocab.vocabSize; a++ {
		fmt.Fprintf(output, "%s %d\n", vocab.vocabArray[a].word, cl[a])
	}
}

func createExpTable() []float32 {
	expTable := make([]float32, expTableSize+1)
	for i := 0; i < expTableSize; i++ {
		expTable[i] = float32(math.Exp(float64((float32(i)/float32(expTableSize)*2.0 - 1.0) * float32(maxExp)))) // Precompute the exp() table
		expTable[i] = expTable[i] / (expTable[i] + 1)                                                            // Precompute f(x) = x / (x + 1)
	}
	return expTable
}

func main() {

	// Go version: Word2Vec -train_file text10.txt -output_file vectors2.bin -cbow  -size 200 -window 8 -negative 25 -sample 1e-4 -num_threads 1 -iter 1 > a2.txt
	// C version:  ./word2vec -train text10.txt -output vectors1.bin -cbow 1 -size 200 -window 8 -negative 25 -hs 0 -sample 1e-4 -threads 1 -binary 0 -iter 1 > a1.txt
	parameters = manageParameters()

	if parameters.cbowMode == true {
		learningRate = 0.05
	}

	expTable = createExpTable()

	vocab := initializeVocab()

	vocab.learnVocab(parameters.trainFile)

	initializeNetwork(vocab.vocabSize, parameters.layer1Size)

	createBinaryTree(vocab.vocabArray)

	if parameters.numberOfNegativeExamples > 0 {
		table = createUnigramTable(vocab.vocabArray)
	}

	trainModel(vocab)

	saveData(parameters.outputFile, vocab)
}
