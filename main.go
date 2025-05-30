package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

const HEADER_END = '#'

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

func bitStringToByte(s string) byte {
	val, err := strconv.ParseInt(s, 2, 8)
	if err != nil {
		panic(err)
	}
	return byte(val)
}

func encode(file []byte) []byte {

	charFrequency := BuildCharFrequency(string(file))

	nodes := []Node{}
	for k, v := range charFrequency {
		newNode := Node{Weight: v, Char: k, LeftNode: nil, RightNode: nil}
		nodes = append(nodes, newNode)
	}

	tree := BuildHuffmanTree(nodes)
	prefixCodeTable := BuildPrefixCodeTable(tree)

	// Build header
	headerString := BuildHeaderTree(&tree, "")
	encodedFile := []byte(headerString)
	encodedFile = append(encodedFile, HEADER_END)

	fullBitString := ""
	for _, c := range string(file) {
		fullBitString += prefixCodeTable[c]
	}
	fmt.Println(fullBitString)
	l := len(fullBitString)
	for i := 0; i < l; i += 7 {
		end := min(l, i+7)
		bitString := fullBitString[i:end]
		byteString := bitStringToByte(bitString)
		fmt.Printf("Transforming %s to byte %d (%s)\n", bitString, byteString, string(byteString))
		encodedFile = append(encodedFile, byteString)
	}
	return encodedFile
}

func decode(file []byte) []byte {
	headerEndPos := 0
	for i, c := range file {
		if c == HEADER_END {
			headerEndPos = i
			break
		}
	}
	header := string(file[0:headerEndPos])
	encodedText := file[headerEndPos+1:]
	fmt.Println(string(encodedText))
	bitStringRune := make(map[string]rune)
	remainingHeader := BuildTreeFromHeader(header, "", bitStringRune)
	if remainingHeader != "" {
		log.Fatalln("Can't continue because header was not completly processed, remaining:", remainingHeader)
	}
	fullBitString := ""
	for i, b := range encodedText {
		bitString := strconv.FormatInt(int64(b), 2)
		// Pad left with zeroes
		if i < len(encodedText)-1 {
			bitString = fmt.Sprintf("%07s", bitString)
		}
		fmt.Printf("Got byte: %d which is %s and equivalent to %s\n", b, string(b), bitString)
		fullBitString += bitString
	}
	decodedText := []byte{}
	l, r := 0, 0
	for l < len(fullBitString) && r < len(fullBitString) {
		currBitString := fullBitString[l : r+1]
		if runeValue, ok := bitStringRune[currBitString]; ok {
			decodedText = append(decodedText, byte(runeValue))
			l = r + 1
			r = l
		} else {
			r += 1
		}
	}
	return decodedText
}

func main() {
	var process string
	var filename string
	var outputFile string
	flag.StringVar(&process, "p", "encode", "Set process to 'encode' or 'decode'")
	flag.StringVar(&filename, "f", "", "Path to the file to code or decode")
	flag.StringVar(&outputFile, "o", "", "Path to the output file")
	flag.CommandLine.Parse(os.Args[1:])

	if filename == "" {
		log.Fatalln("Please define an input file with the flag -f")
	}

	file, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalln("Erro reading input file", err)
	}

	var result []byte
	if process == "encode" {
		result = encode(file)
	} else {
		result = decode(file)
		fmt.Println(result)
	}

	if outputFile != "" {
		os.WriteFile(outputFile, result, 0644)
	} else {
		fmt.Println(string(result))
	}
}
