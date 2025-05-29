package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

func BuildCharFrequency(file string) map[rune]int {
	charFrequency := make(map[rune]int)
	for _, c := range file {
		if c == 10 {
			continue
		}
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
	nodes := []Node{}
	for k, v := range charFrequency {
		newNode := Node{Weight: v, Char: k, LeftNode: nil, RightNode: nil}
		nodes = append(nodes, newNode)
	}

	tree := BuildHuffmanTree(nodes)
	prefixCodeTable := BuildPrefixCodeTable(tree)
	for k, v := range prefixCodeTable {
		fmt.Printf("%s:%s\n", string(k), strconv.FormatInt(int64(v), 2))
	}
}
