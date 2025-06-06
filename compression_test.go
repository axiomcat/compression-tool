package main

import (
	"math/rand"
	"os"
	"reflect"
	"testing"
)

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

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!@$%&*)*&)_-=[];',./?><:{}'")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func TestEncoderRandom(t *testing.T) {
	s := RandStringRunes(rand.Intn(10000))
	file := []byte(s)
	encodedFile := encode(file)
	decodedFile := decode(encodedFile)
	if !reflect.DeepEqual(file, decodedFile) {
		t.Errorf("File and decoded file are not the same: Expected %s but got %s", string(file), string(decodedFile))
	}
}

var FileTest = []struct {
	filename string
}{
	{"tests/test1.in"},
	{"tests/test2.in"},
	{"tests/test3.in"},
	{"tests/book.txt"},
}

func TestEncoderFiles(t *testing.T) {
	for _, testData := range FileTest {
		file, err := os.ReadFile(testData.filename)
		if err != nil {
			t.Error("Error reading input file", err)
		}
		encodedFile := encode(file)
		decodedFile := decode(encodedFile)
		if !reflect.DeepEqual(file, decodedFile) {
			t.Errorf("File and decoded file are not the same: Expected %s but got %s", string(file), string(decodedFile))
		}
	}
}
