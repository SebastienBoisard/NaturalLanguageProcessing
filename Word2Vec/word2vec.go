package main

import (
	"bufio"
	"log"
	"fmt"
	"math"
	"os"
)

const tableSize = 1e8

var table []int

const maxCodeLength = 40

const maxSentenceLength = 1000

const maxExp = 6

const expTableSize = 1000

var expTable [expTableSize + 1]float32

var startingAlpha float32

var fileSize int

var wordCountActual = 0

var alpha = float32(0.025)

var trainWords int

var syn0, syn1, syn1neg []float32

var isEndFile = false

func createUnigramTable() []int {

	fmt.Println("createUnigramTable BEGIN")

	const power float64 = 0.75

	unigramTable := make([]int, tableSize)

	var trainWordsPow float64

	for a := 0; a < vocabSize; a++ {
		trainWordsPow += math.Pow(float64(vocab[a].frequency), power)
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

// createBinaryTree creates binary Huffman tree using the word counts
// Frequent words will have short unique binary codes
func createBinaryTree() {

	fmt.Println("createBinaryTree BEGIN")

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

func initializeNetwork() {

	//long long a, b;
	//unsigned long long nextRandom = 1;
	syn0 = make([]float32, vocabSize*layer1Size)

	if isHierarchicalSoftmaxActivated == true {
		syn1 = make([]float32, vocabSize*layer1Size)
	}

	var nextRandom uint64 = 1

	if numberOfNegativeExamples > 0 {

		syn1neg = make([]float32, vocabSize*layer1Size)
	}


	for a := 0; a < vocabSize; a++ {
		for b := 0; b < layer1Size; b++ {
			nextRandom = nextRandom*25214903917 + 11

			// fmt.Printf("InitNet next_random=%d\n", nextRandom)

			// var f float32

			// f = float32(nextRandom & 0xFFFF)
			// fmt.Printf("InitNet f1=%.20f\n", f)

			// f = f / float32(65536)
			// fmt.Printf("InitNet f2=%.20f\n", f)

			// f = f - 0.5
			// fmt.Printf("InitNet f3=%.20f\n", f)

			// f = f / float32(layer1Size)
			// fmt.Printf("InitNet f4=%.20f\n", f)


			// fmt.Printf("InitNet v=%.20f\n", (((float32(nextRandom & 0xFFFF)) / float32(65536)) - float32(0.5))  / float32(layer1Size) )

			syn0[a*layer1Size+b] = (((float32(nextRandom & 0xFFFF)) / float32(65536)) - float32(0.5)) / float32(layer1Size)



				fmt.Printf("InitNet syn0a[%d]=%.20f\n", a * layer1Size + b, syn0[a * layer1Size + b])
				// fmt.Printf("InitNet syn0b[%d]=%.20f\n", a * layer1Size + b, f);

		}
	}

	createBinaryTree()
}

func trainModelThread(id int) {

	fmt.Println("trainModelThread[", id, "] BEGIN")

	//  long long a, b, d, cw, word, last_word, sentenceLength = 0, sentencePosition = 0;
	//long long wordCount = 0, lastWordCount = 0, sen[maxSentenceLength + 1];
	//long long l1, l2, c, target, label, local_iter = iter;
	//real f, g;
	//	var cw int
	//	var g float32

	var label int
	var target int

	nextRandom := uint64(id)

	var wordCount, lastWordCount int
	var sen [maxSentenceLength + 1]int

	sentencePosition := 0
	localIter := numberOfIterations

	neu1 := make([]float32, layer1Size)
	neu1e := make([]float32, layer1Size)

	isEndFile = false
	fi, err := os.Open(trainFile)
	if err != nil {
		log.Fatal(err)
	}
	defer fi.Close()

	filePosition := int64(fileSize / numberOfThreads * id)
	fi.Seek(filePosition, 0)

	fmt.Println("trainModelThread[", id, "] numberOfThreads=", numberOfThreads, "filePosition=", filePosition)
	fmt.Printf("trainModelThread[ %d ] alpha= %.20f\n", id, alpha)

	reader := bufio.NewReader(fi)

	sentenceLength := 0
	counter2 := 0
	for {
		counter2++

		fmt.Println("trainModelThread[", id, "][", counter2, "] wordCount=", wordCount, "lastWordCount=", lastWordCount)

		// for c := 0; c < layer1Size; c++ {
		// 	fmt.Printf("trainModelThread[ %d ][ %d ] new1[ %d ]= %.20f\n", id, counter2, c, neu1[c])
		// }
		idx_neu1 := 1
		fmt.Printf("trainModelThread[ %d ][ %d ] 0 new1[ %d ]= %.20f\n", id, counter2, idx_neu1, neu1[idx_neu1])

		if wordCount-lastWordCount > 10000 {
			wordCountActual += wordCount - lastWordCount
			lastWordCount = wordCount

			fmt.Printf("%cAlpha: %f  Progress: %.2f%%\n", 13, alpha, float32(wordCountActual)/float32(numberOfIterations*trainWords+1)*100)
			alpha = float32(startingAlpha * (1 - float32(wordCountActual)/float32(numberOfIterations*trainWords+1)))

			if alpha < startingAlpha*0.0001 {
				alpha = startingAlpha * 0.0001
			}
		}

		// fmt.Println("trainModelThread[", id, "][", counter2, "] sentenceLength=", sentenceLength)

		if sentenceLength == 0 {

			counter := 0
			for {
				counter++
				word := readWordIndex(reader)

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
				if occurrenceWordsThreshold > 0 {

					// fmt.Println("trainModelThread[", id, "] vocab[word].frequency=", vocab[word].frequency, "occurrenceWordsThreshold=", occurrenceWordsThreshold, "trainWords=", trainWords)


					ran := float32(math.Sqrt(float64(float32(vocab[word].frequency)/(float32(occurrenceWordsThreshold)*float32(trainWords)))) + 1.0) * 
						               (float32(occurrenceWordsThreshold) * float32(trainWords)) / float32(vocab[word].frequency)
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

		if isEndFile == true || wordCount > trainWords/numberOfThreads {
			wordCountActual += wordCount - lastWordCount
			localIter--
			if localIter == 0 {
				fmt.Println("trainModelThread[", id, "][", counter2, "] locaIter==0 so break")
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
		for c := 0; c < layer1Size; c++ {
			neu1[c] = 0
		}
		for c := 0; c < layer1Size; c++ {
			neu1e[c] = 0
		}

		nextRandom = nextRandom*25214903917 + 11
		b := int(nextRandom % uint64(windowSize))

		// fmt.Println("trainModelThread[", id, "][", counter2, "] nextRandom=", nextRandom, "b=", b)

		if cbowMode == true {
			//train the cbow architecture

			// fmt.Println("trainModelThread[", id, "][", counter2, "] cbowMode on")

			// in -> hidden
			cw := 0

			for a := int(b); a < windowSize*2+1-int(b); a++ {
				if a != windowSize {
					c := sentencePosition - windowSize + a
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

					fmt.Printf("trainModelThread[ %d ][ %d ] 1a new1[ %d ]= %.20f\n", id, counter2, idx_neu1, neu1[idx_neu1])

					fmt.Printf("trainModelThread[ %d ][ %d ] lastWord=%d syn0[%d]=%.20f\n", id, counter2, lastWord, 
						idx_neu1+lastWord*layer1Size, syn0[idx_neu1+lastWord*layer1Size])

					for c := 0; c < layer1Size; c++ {
						neu1[c] += syn0[c+lastWord*layer1Size]
						// fmt.Printf("trainModelThread[ %d ][ %d ] 1b new1[ %d ]= %.20f\n", id, counter2, c, neu1[c])
					}
					fmt.Printf("trainModelThread[ %d ][ %d ] 1b new1[ %d ]= %.20f\n", id, counter2, idx_neu1, neu1[idx_neu1])
					cw++
				}
			}

			if cw > 0 {

				// fmt.Println("trainModelThread[", id, "][", counter2, "] cw > 0")

				for c := 0; c < layer1Size; c++ {
					neu1[c] /= float32(cw)
					fmt.Printf("trainModelThread[ %d ][ %d ] 2 new1[ %d ]= %.20f\n", id, counter2, c, neu1[c])
				}
				fmt.Printf("trainModelThread[ %d ][ %d ] 2 new1[ %d ]= %.20f\n", id, counter2, idx_neu1, neu1[idx_neu1])

				if isHierarchicalSoftmaxActivated == true {

					// fmt.Println("trainModelThread[", id, "][", counter2, "] isHierarchicalSoftmaxActivated")

					for d := 0; d < int(vocab[word].codelen); d++ {
						f := float32(0.0)
						l2 := vocab[word].point[d] * layer1Size
						// Propagate hidden -> output
						for c := 0; c < layer1Size; c++ {
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

						fmt.Printf("trainModelThread[ %d ][ %d ] hs f=%.20f\n", id, counter2, f)

						// 'g' is the gradient multiplied by the learning rate
						g := float32(1.0-float32(vocab[word].code[d])-f) * alpha

						fmt.Printf("trainModelThread[ %d ][ %d ] hs g=%.20f\n", id, counter2, g)

						// Propagate errors output -> hidden
						for c := 0; c < layer1Size; c++ {
							neu1e[c] += g * syn1[c+l2]
							// fmt.Printf("trainModelThread[ %d ][ %d ] hs neu1e[ %d ]= %.20f\n", id, counter2, c, neu1e[c])
						}
						// Learn weights hidden -> output
						for c := 0; c < layer1Size; c++ {
							syn1[c+l2] += g * neu1[c]
							// fmt.Printf("trainModelThread[ %d ][ %d ] hs syn1[ %d ]= %.20f\n", id, counter2, c+l2, syn1[c+l2])
						}
					}
				}

				// NEGATIVE SAMPLING
				if numberOfNegativeExamples > 0 {

					// fmt.Println("trainModelThread[", id, "][", counter2, "] numberOfNegativeExamples > 0")

					var label int
					var g float32
					for d := 0; d < numberOfNegativeExamples+1; d++ {
						if d == 0 {
							target = word
							label = 1
						} else {
							nextRandom = nextRandom*25214903917 + 11
							// fmt.Println("trainModelThread[", id, "] nextRandom=", nextRandom, "(nextRandom>>16)%tableSize=", (nextRandom>>16)%tableSize)
							target = table[(nextRandom>>16)%tableSize]
							if target == 0 {
								target = int(nextRandom%uint64(vocabSize-1)) + 1
							}
							if target == word {
								continue
							}
							label = 0
						}
						l2 := target * layer1Size
						fmt.Printf("trainModelThread[ %d ][ %d ] negative l2=%d\n", id, counter2, l2)

						f := float32(0.0)
						for c := 0; c < layer1Size; c++ {
							f += neu1[c] * syn1neg[c+l2]
							// fmt.Printf("trainModelThread[ %d ][ %d ] negative f[%d]=%.20f  neu1[%d]=%.20f  syn1neg[%d]=%.20f   neu1[%d]*syn1neg[%d]=%.20f\n", id, counter2, c, f, c, neu1[c], c+l2, syn1neg[c+l2], c, c+l2, neu1[c] * syn1neg[c+l2])
						}

						fmt.Printf("trainModelThread[ %d ][ %d ] negative f=%.20f\n", id, counter2, f)

						if f > maxExp {
							g = float32(label-1) * alpha
							fmt.Printf("trainModelThread[ %d ][ %d ] negative g1=%.20f\n", id, counter2, g)
						} else {
							if f < -maxExp {
								g = float32(label-0) * alpha
								fmt.Printf("trainModelThread[ %d ][ %d ] negative g2=%.20f\n", id, counter2, g)
							} else {
								exp_idx := int((f+maxExp)*(expTableSize/maxExp/2.0))
								g = (float32(label) - expTable[exp_idx]) * alpha


								//g = (label - expTable[(int)((f + MAX_EXP) * (EXP_TABLE_SIZE / MAX_EXP / 2))]) * alpha;

								// fmt.Printf("trainModelThread[ %d ][ %d ] negative maxEp=%d  expTableSize=%d  f=%.20f  exp_idx=%d  expTable[exp_idx]=%.20f\n", 
									// id, counter2, maxExp, expTableSize, f, exp_idx, expTable[exp_idx])

								// fmt.Printf("trainModelThread[ %d ][ %d ] negative f+maxEp=%f  expTableSize/maxExp/2=%d  (f+maxExp)*(expTableSize/maxExp/2)=%f\n", 
									// id, counter2, (f+maxExp), (expTableSize/maxExp/2), ((f+maxExp)*(expTableSize/maxExp/2)))

								// fmt.Printf("trainModelThread[ %d ][ %d ] negative label=%d f=%.20f maxExp=%d expTableSize=%d alpha=%.20f g3=%.20f\n", 
									// id, counter2, label, f, maxExp, expTableSize, alpha, g)
								// fmt.Printf("trainModelThread[ %d ][ %d ] negative g3=%.20f\n", id, counter2, g)
							}
						}
						
						
						// fmt.Printf("trainModelThread[ %d ][ %d ] negative neu1e[ %d ]= %.20f\n", id, counter2, 0, neu1e[0])
						for c := 0; c < layer1Size; c++ {
							neu1e[c] += g * syn1neg[c+l2]
          				// fmt.Printf("trainModelThread[ %d ][ %d ] negative neu1e[ %d ]= %.20f\n", id, counter2, c, neu1e[c])
						}

						// fmt.Printf("trainModelThread[ %d ][ %d ] negative l2=%d g=%.20f syn1neg[ %d ]= %.20f\n", id, counter2, l2, g, l2, syn1neg[l2])
						for c := 0; c < layer1Size; c++ {
							syn1neg[c+l2] += g * neu1[c]
          				// fmt.Printf("trainModelThread[ %d ][ %d ] negative syn1neg[%d]=%.20f  neu1[%d]=%.20f\n", id, counter2, c+l2, syn1neg[c+l2], c, neu1[c])
						}
					}
				}

				// hidden -> in
				for a := int(b); a < int(windowSize)*2+1-int(b); a++ {
					if a != int(windowSize) {
						c := sentencePosition - windowSize + a
						fmt.Printf("trainModelThread[ %d ][ %d ] hidden a=%d c=%d\n", id, counter2, a, c)
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
						for c := 0; c < layer1Size; c++ {
							syn0[c+lastWord*layer1Size] += neu1e[c]
						}
					}
				}
			}
		} else {
			//train skip-gram

			fmt.Println("trainModelThread[", id, "] skip-gram")

			for a := b; a < windowSize*2+1-b; a++ {
				if a != windowSize {
					c := sentencePosition - windowSize + a
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
					l1 := lastWord * layer1Size
					for c := 0; c < layer1Size; c++ {
						neu1e[c] = 0
					}
					// HIERARCHICAL SOFTMAX
					if isHierarchicalSoftmaxActivated == true {
						for d := 0; d < int(vocab[word].codelen); d++ {
							f := float32(0.0)
							l2 := vocab[word].point[d] * layer1Size
							// Propagate hidden -> output
							for c := 0; c < layer1Size; c++ {
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
							g := (1 - float32(vocab[word].code[d]) - f) * alpha
							// Propagate errors output -> hidden
							for c := 0; c < layer1Size; c++ {
								neu1e[c] += g * syn1[c+l2]
							}
							// Learn weights hidden -> output
							for c := 0; c < layer1Size; c++ {
								syn1[c+l2] += g * syn0[c+l1]
							}
						}
					}

					// NEGATIVE SAMPLING
					if numberOfNegativeExamples > 0 {
						for d := 0; d < numberOfNegativeExamples+1; d++ {
							if d == 0 {
								target = word
								label = 1
							} else {
								nextRandom = nextRandom*25214903917 + 11
								target = table[(nextRandom>>16)%tableSize]
								if target == 0 {
									target = int(nextRandom%uint64(vocabSize-1)) + 1
								}
								if target == word {
									continue
								}
								label = 0
							}
							var g float32
							l2 := target * layer1Size
							f := float32(0.0)
							for c := 0; c < layer1Size; c++ {
								f += syn0[c+l1] * syn1neg[c+l2]
							}
							if f > maxExp {
								g = float32(label-1) * alpha
							} else {
								if f < -maxExp {
									g = float32(label-0) * alpha
								} else {
									g = float32(float32(label)-expTable[int((f+maxExp)*(expTableSize/maxExp/2))]) * alpha
								}
							}
							for c := 0; c < layer1Size; c++ {
								neu1e[c] += g * syn1neg[c+l2]
							}
							for c := 0; c < layer1Size; c++ {
								syn1neg[c+l2] += g * syn0[c+l1]
							}
						}
					}

					// Learn weights input -> hidden
					for c := 0; c < layer1Size; c++ {
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

func trainModel() {
	fmt.Println("trainModel BEGIN")

	//long a, b, c, d;
	//FILE *fo;
	//pthread_t *pt = (pthread_t *)malloc(num_threads * sizeof(pthread_t));
	//printf("Starting training using file %s\n", train_file);

	startingAlpha = float32(startingLearningRate)

	learnVocabFromTrainFile(trainFile)

	initializeNetwork()

	if numberOfNegativeExamples > 0 {
		table = createUnigramTable()
	}

	//start = clock();

	trainModelThread(0)

	//for (a = 0; a < num_threads; a++) pthread_create(&pt[a], NULL, TrainModelThread, (void *)a);
	//for (a = 0; a < num_threads; a++) pthread_join(pt[a], NULL);

	//fo = fopen(output_file, "wb");
	file, err := os.Create(outputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	if numberOfClasses == 0 {
		// Save the word vectors
		fmt.Fprintf(file, "%d %d\n", vocabSize, layer1Size)
		for a := 0; a < vocabSize; a++ {
			fmt.Fprintf(file, "%s ", vocab[a].word)
			if binaryMode == true {
				for b := 0; b < layer1Size; b++ {
					//fwrite(&syn0[a * layer1Size + b], sizeof(real), 1, fo);
				}
			} else {
				for b := 0; b < layer1Size; b++ {
					fmt.Fprintf(file, "%f ", syn0[a*layer1Size+b])
				}
			}
			fmt.Fprintf(file, "\n")
		}
	} else {
		// Run K-means on the word vectors
		// int clcn = classes, iter = 10, closeid;

		clcn := numberOfClasses
		numberOfIterations = 10
		var closeid int

		var closev, x float32
		cl := make([]int, vocabSize)
		centcn := make([]int, numberOfClasses)
		cent := make([]float32, numberOfClasses*layer1Size)

		for a := 0; a < vocabSize; a++ {
			cl[a] = a % clcn
		}

		for a := 0; a < numberOfIterations; a++ {
			for b := 0; b < clcn*layer1Size; b++ {
				cent[b] = 0
			}
			for b := 0; b < clcn; b++ {
				centcn[b] = 1
			}
			for c := 0; c < vocabSize; c++ {
				for d := 0; d < layer1Size; d++ {
					cent[layer1Size*cl[c]+d] += syn0[c*layer1Size+d]
				}
				centcn[cl[c]]++
			}
			for b := 0; b < clcn; b++ {
				closev = 0
				for c := 0; c < layer1Size; c++ {
					cent[layer1Size*b+c] = cent[layer1Size*b+c] / float32(centcn[b])
					closev += cent[layer1Size*b+c] * cent[layer1Size*b+c]
				}
				closev = float32(math.Sqrt(float64(closev)))
				for c := 0; c < layer1Size; c++ {
					cent[layer1Size*b+c] /= closev
				}
			}
			for c := 0; c < vocabSize; c++ {
				closev = -10
				closeid = 0
				for d := 0; d < clcn; d++ {
					x = 0
					for b := 0; b < layer1Size; b++ {
						x += cent[layer1Size*d+b] * syn0[c*layer1Size+b]
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
		for a := 0; a < vocabSize; a++ {
			fmt.Fprintf(file, "%s %d\n", vocab[a].word, cl[a])
		}
	}
}

func initializeExpTable() {
	for i := 0; i < expTableSize; i++ {
		expTable[i] = float32(math.Exp(float64((float32(i)/float32(expTableSize)*2.0-1.0) * float32(maxExp)))) // Precompute the exp() table
		// fmt.Printf("expTable1[%d]=%.20f\n", i, expTable[i])
		expTable[i] = expTable[i] / (expTable[i] + 1)                         // Precompute f(x) = x / (x + 1)
	}
}

func main() {

	// Word2Vec -train_file text10.txt -output_file vectors2.bin -cbow  -size 200 -window 8 -negative 25 -sample 1e-4 -num_threads 1 -iter 1 > a2.txt
   // ./word2vec -train text10.txt -output vectors1.bin -cbow 1 -size 200 -window 8 -negative 25 -hs 0 -sample 1e-4 -threads 1 -binary 0 -iter 1 > a1.txt
	manageParameters()
	if (cbowMode == true) {
		alpha = 0.05
	}

	initializeExpTable()
	initializeVocabulary()
	trainModel()

}
