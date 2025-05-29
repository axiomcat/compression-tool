package main

import "testing"

var frequencyTest = []struct {
	in   string
	char rune
	out  int
}{
	{"aaabbc", 'a', 3},
	{"aaabbc", 'b', 2},
	{"aaabbc", 'c', 1},
}

func TestFrequency(t *testing.T) {
	for _, testData := range frequencyTest {
		charFrequency := BuildCharFrequency(testData.in)
		currentCharFreq := charFrequency[testData.char]
		if testData.out != currentCharFreq {
			t.Errorf("Expected frequency %d of char %s in string %s, found %d", testData.out, string(testData.char), testData.in, currentCharFreq)
		}
	}
}
