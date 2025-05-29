package main

import (
	"fmt"
	"log"
	"os"
)

func BuildCharFrequency(file string) map[rune]int {
	charFrequency := make(map[rune]int)
	for _, c := range file {
		if _, ok := charFrequency[c]; ok {
			charFrequency[c] += 1
		} else {
			charFrequency[c] = 1
		}
	}
	return charFrequency
}
func main() {
	if len(os.Args) < 2 {
		panic("Expected filename argument")
	}
	filename := os.Args[1]
	fmt.Println(filename)
	bytes, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalln("Error reading file", err)
	}
	charFrequency := BuildCharFrequency(string(bytes))
	fmt.Printf("Frequency of X %d\n", charFrequency['X'])
	fmt.Printf("Frequency of t %d\n", charFrequency['t'])
}
