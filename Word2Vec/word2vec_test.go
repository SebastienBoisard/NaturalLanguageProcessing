package main

import "testing"

func TestCreateBinaryTree(t *testing.T) {

	actualVocab := []Term{
		{word: "</s>", frequency: 1},
		{word: "the", frequency: 10},
		{word: "of", frequency: 9},
		{word: "and", frequency: 7},
		{word: "a", frequency: 6},
		{word: "is", frequency: 5},
		{word: "that", frequency: 5},
	}

	createBinaryTree(actualVocab)

	expectedVocab := []Term{
		{word: "</s>", frequency: 1, point: [maxCodeLength]int{5, 3, -7}, code: [maxCodeLength]byte{0, 1}, codelen: 2},
		{word: "the", frequency: 10, point: [maxCodeLength]int{5, 3, -6}, code: [maxCodeLength]byte{0, 0}, codelen: 2},
		{word: "of", frequency: 9, point: [maxCodeLength]int{5, 4, 2, -5}, code: [maxCodeLength]byte{1, 1, 0}, codelen: 3},
		{word: "and", frequency: 7, point: [maxCodeLength]int{5, 4, 1, -4}, code: [maxCodeLength]byte{1, 0, 1}, codelen: 3},
		{word: "a", frequency: 6, point: [maxCodeLength]int{5, 4, 1, -3}, code: [maxCodeLength]byte{1, 0, 0}, codelen: 3},
		{word: "is", frequency: 5, point: [maxCodeLength]int{5, 4, 2, 0, -2}, code: [maxCodeLength]byte{1, 1, 1, 1}, codelen: 4},
		{word: "that", frequency: 5, point: [maxCodeLength]int{5, 4, 2, 0, -1}, code: [maxCodeLength]byte{1, 1, 1}, codelen: 4},
	}

	for index := range expectedVocab {

		if actualVocab[index].point != expectedVocab[index].point {
			t.Errorf("createBinaryTree()[%d].point(%v) = %v", index, expectedVocab[index].point, actualVocab[index].point)
		}

		if actualVocab[index].code != expectedVocab[index].code {
			t.Errorf("createBinaryTree()[%d].code(%v) = %v", index, expectedVocab[index].code, actualVocab[index].code)
		}

		if actualVocab[index].codelen != expectedVocab[index].codelen {
			t.Errorf("createBinaryTree()[%d].codelen(%d) = %d", index, expectedVocab[index].codelen, actualVocab[index].codelen)
		}

	}
}

func TestCreateExpTable(t *testing.T) {

	actualExpTable := createExpTable()

	var tests = []struct {
		expTableIndex int
		wantedValue   float32
	}{
		{0, 0.0024726231566347748},
		{10, 0.00278699607588350773},
		{100, 0.00816257018595933914},
		{999, 0.99749761819839477539},
	}

	for _, test := range tests {
		if actualExpTable[test.expTableIndex] != test.wantedValue {
			t.Errorf("createExpTable()[%d] = %f", test.expTableIndex, test.wantedValue)
		}
	}
}

func TestCreateUnigramTable(t *testing.T) {

	actualVocab := []Term{
		{word: "scale", frequency: 5, point: [maxCodeLength]int{523, 521, 517, 510, 498, 478, 443, 383, 293, 15, 8, 11, -24}, code: [maxCodeLength]byte{0, 0, 1, 1, 1, 1, 0, 0, 1, 1, 1}, codelen: 11},
		{word: "seem", frequency: 5, point: [maxCodeLength]int{523, 521, 517, 510, 498, 478, 443, 383, 293, 158, 11, -23}, code: [maxCodeLength]byte{0, 0, 1, 1, 1, 1, 0, 0, 1, 1, 0}, codelen: 11},
		{word: "services", frequency: 5, point: [maxCodeLength]int{523, 521, 517, 510, 498, 478, 443, 383, 293, 158, 10, -22}, code: [maxCodeLength]byte{0, 0, 1, 1, 1, 1, 0, 0, 1, 0, 1}, codelen: 11},
		{word: "several", frequency: 5, point: [maxCodeLength]int{523, 521, 517, 510, 498, 478, 443, 383, 293, 158, 10, -21}, code: [maxCodeLength]byte{0, 0, 1, 1, 1, 1, 0, 0, 1, 0, 0}, codelen: 11},
		{word: "sovereignty", frequency: 5, point: [maxCodeLength]int{523, 521, 517, 510, 498, 478, 443, 383, 293, 157, 9, -20}, code: [maxCodeLength]byte{0, 0, 1, 1, 1, 1, 0, 0, 0, 1, 1}, codelen: 11},
		{word: "statistical", frequency: 5, point: [maxCodeLength]int{523, 521, 517, 510, 498, 478, 443, 383, 293, 157, 9, -19}, code: [maxCodeLength]byte{0, 0, 1, 1, 1, 1, 0, 0, 0, 1, 0}, codelen: 11},
		{word: "stirner", frequency: 5, point: [maxCodeLength]int{523, 521, 517, 510, 498, 478, 443, 383, 293, 157, 8, -18}, code: [maxCodeLength]byte{0, 0, 1, 1, 1, 1, 0, 0, 0, 0, 1}, codelen: 11},
		{word: "study", frequency: 5, point: [maxCodeLength]int{523, 521, 517, 510, 498, 478, 443, 383, 293, 157, 8, -17}, code: [maxCodeLength]byte{0, 0, 1, 1, 1, 1, 0, 0, 0, 0, 0}, codelen: 11},
		{word: "syndicalist", frequency: 5, point: [maxCodeLength]int{523, 521, 517, 510, 498, 478, 442, 382, 292, 156, 7, -16}, code: [maxCodeLength]byte{0, 0, 1, 1, 1, 0, 1, 1, 1, 1, 1}, codelen: 11},
		{word: "terms", frequency: 5, point: [maxCodeLength]int{523, 521, 517, 510, 498, 478, 442, 382, 292, 156, 7, -15}, code: [maxCodeLength]byte{0, 0, 1, 1, 1, 0, 1, 1, 1, 1, 0}, codelen: 11},
		{word: "theorists", frequency: 5, point: [maxCodeLength]int{523, 521, 517, 510, 498, 478, 442, 382, 292, 156, 6, -14}, code: [maxCodeLength]byte{0, 0, 1, 1, 1, 0, 1, 1, 1, 0, 1}, codelen: 11},
		{word: "things", frequency: 5, point: [maxCodeLength]int{523, 521, 517, 510, 498, 478, 442, 382, 292, 156, 6, -13}, code: [maxCodeLength]byte{0, 0, 1, 1, 1, 0, 1, 1, 1, 0, 0}, codelen: 11},
		{word: "throughout", frequency: 5, point: [maxCodeLength]int{523, 521, 517, 510, 498, 478, 442, 382, 292, 155, 5, -12}, code: [maxCodeLength]byte{0, 0, 1, 1, 1, 0, 1, 1, 0, 1, 1}, codelen: 11},
		{word: "trade", frequency: 5, point: [maxCodeLength]int{523, 521, 517, 510, 498, 478, 442, 382, 292, 155, 5, -11}, code: [maxCodeLength]byte{0, 0, 1, 1, 1, 0, 1, 1, 0, 1, 0}, codelen: 11},
		{word: "tropical", frequency: 5, point: [maxCodeLength]int{523, 521, 517, 510, 498, 478, 442, 382, 292, 155, 4, -10}, code: [maxCodeLength]byte{0, 0, 1, 1, 1, 0, 1, 1, 0, 0, 1}, codelen: 11},
		{word: "understand", frequency: 5, point: [maxCodeLength]int{523, 521, 517, 510, 498, 478, 442, 382, 292, 155, 4, -9}, code: [maxCodeLength]byte{0, 0, 1, 1, 1, 0, 1, 1, 0, 0, 0}, codelen: 11},
		{word: "unusual", frequency: 5, point: [maxCodeLength]int{523, 521, 517, 510, 498, 478, 442, 382, 291, 154, 3, -8}, code: [maxCodeLength]byte{0, 0, 1, 1, 1, 0, 1, 0, 1, 1, 1}, codelen: 11},
		{word: "upon", frequency: 5, point: [maxCodeLength]int{523, 521, 517, 510, 498, 478, 442, 382, 291, 154, 3, -7}, code: [maxCodeLength]byte{0, 0, 1, 1, 1, 0, 1, 0, 1, 1, 0}, codelen: 11},
		{word: "values", frequency: 5, point: [maxCodeLength]int{523, 521, 517, 510, 498, 478, 442, 382, 291, 154, 2, -6}, code: [maxCodeLength]byte{0, 0, 1, 1, 1, 0, 1, 0, 1, 0, 1}, codelen: 11},
		{word: "voting", frequency: 5, point: [maxCodeLength]int{523, 521, 517, 510, 498, 478, 442, 382, 291, 154, 2, -5}, code: [maxCodeLength]byte{0, 0, 1, 1, 1, 0, 1, 0, 1, 0, 0}, codelen: 11},
		{word: "wealth", frequency: 5, point: [maxCodeLength]int{523, 521, 517, 510, 498, 478, 442, 382, 291, 153, 1, -4}, code: [maxCodeLength]byte{0, 0, 1, 1, 1, 0, 1, 0, 0, 1, 1}, codelen: 11},
		{word: "western", frequency: 5, point: [maxCodeLength]int{523, 521, 517, 510, 498, 478, 442, 382, 291, 153, 1, -3}, code: [maxCodeLength]byte{0, 0, 1, 1, 1, 0, 1, 0, 0, 1, 0}, codelen: 11},
		{word: "william", frequency: 5, point: [maxCodeLength]int{523, 521, 517, 510, 498, 478, 442, 382, 291, 153, 0, -2}, code: [maxCodeLength]byte{0, 0, 1, 1, 1, 0, 1, 0, 0, 0, 1}, codelen: 11},
		{word: "zayed", frequency: 5, point: [maxCodeLength]int{523, 521, 517, 510, 498, 478, 442, 382, 291, 153, 0, -1}, code: [maxCodeLength]byte{0, 0, 1, 1, 1, 0, 1, 0, 0, 0, 0}, codelen: 11},
	}

	actualUnigramTable := createUnigramTable(actualVocab)

	expectedUnigramTableSize := int(1e8)
	actualUnigramTableSize := len(actualUnigramTable)

	if actualUnigramTableSize != expectedUnigramTableSize {
		t.Error("createUnigramTable().len(%d) = %d", expectedUnigramTableSize, actualUnigramTableSize)
	}

	var tests = []struct {
		unigramTableIndex int
		wantedValue       int
	}{
		{6889208, 1},
		{10564422, 2},
		{15192095, 3},
		{16638081, 3},
		{18767860, 4},
		{37003654, 8},
		{76969030, 18},
		{95301961, 22},
		{99999999, 23},
	}

	for _, test := range tests {
		if actualUnigramTable[test.unigramTableIndex] != test.wantedValue {
			t.Errorf("createUnigramTable()[%d] = %d", test.unigramTableIndex, test.wantedValue)
		}
	}
}
