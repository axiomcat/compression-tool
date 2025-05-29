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

func bitStringToByte(s string) byte {
	val, err := strconv.ParseInt(s, 2, 8)
	if err != nil {
		panic(err)
	}
	return byte(val)
}

func main() {
	if len(os.Args) < 2 {
		panic("Expected filename argument")
	}
	filename := os.Args[1]
	bytes, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalln("Error reading file", err)
	}
	fmt.Println("Building char frequency")
	charFrequency := BuildCharFrequency(string(bytes))
	nodes := []Node{}
	for k, v := range charFrequency {
		newNode := Node{Weight: v, Char: k, LeftNode: nil, RightNode: nil}
		nodes = append(nodes, newNode)
	}
	fmt.Println("Finished char frequency")
	fmt.Println("Building tree")
	tree := BuildHuffmanTree(nodes)
	fmt.Println("Finished tree")
	fmt.Println("Building CodeTable")
	prefixCodeTable := BuildPrefixCodeTable(tree)
	fmt.Println("Finished CodeTable")

	fmt.Println("Start encoding")
	compressedFile := []byte{}
	for k, v := range prefixCodeTable {
		headerTableEntry := fmt.Sprintf("%x:%s,", k, v)
		compressedFile = append(compressedFile, headerTableEntry...)
	}

	fmt.Println("\tFinished header encoding")

	compressedFile = append(compressedFile, 'µ')

	encodedFile := ""
	inputLen := len(bytes)
	fmt.Printf("Input len %d\n", inputLen)
	for i, c := range string(bytes) {
		if i%100000 == 0 {
			fmt.Printf("%d/%d\n", i, inputLen)
		}
		encodedFile += prefixCodeTable[c]
	}
	fmt.Println("\tFinished building full string")
	l := len(encodedFile)
	for i := 0; i < l; i += 7 {
		end := min(l, i+7)
		bitString := encodedFile[i:end]
		val, err := strconv.ParseInt(bitString, 2, 8)
		if err != nil {
			panic(err)
		}
		compressedFile = append(compressedFile, byte(val))
	}

	os.WriteFile("compressed.txt", compressedFile, 0644)
}
