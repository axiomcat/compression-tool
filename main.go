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

	currentPrefix := ""
	for _, c := range string(file) {
		prefixStr := prefixCodeTable[c]
		currentPrefix += prefixStr
		if len(currentPrefix) > 7 {
			bitString := currentPrefix[:7]
			byteString := bitStringToByte(bitString)
			encodedFile = append(encodedFile, byteString)
			currentPrefix = currentPrefix[7:]
		}
	}
	// We alway have text remaining in currentPrefix
	bitString := currentPrefix
	byteString := bitStringToByte(bitString)
	encodedFile = append(encodedFile, byteString)

	// Count number of padding zeroes to add it to the last byte
	zeroPos, zeroCount := 0, 0
	for zeroPos < len(bitString) && bitString[zeroPos] == '0' {
		zeroCount += 1
		zeroPos += 1
	}
	// If all bits are zero, we dont' want to count one of them
	if zeroCount == len(bitString) {
		zeroCount -= 1
	}
	lastByteInfo := byte(zeroCount)
	encodedFile = append(encodedFile, lastByteInfo)

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
	bitStringRune := make(map[string]rune)
	remainingHeader := BuildTreeFromHeader(header, "", bitStringRune)
	if remainingHeader != "" {
		log.Fatalln("Can't continue because header was not completly processed, remaining:", remainingHeader)
	}
	fullBitString := ""
	for i, b := range encodedText[:len(encodedText)-1] {
		bitString := strconv.FormatInt(int64(b), 2)
		// Pad left with zeroes
		if i < len(encodedText)-2 {
			bitString = fmt.Sprintf("%07s", bitString)
		} else {
			// The last byte of the encoded file contains how many zeroes should we
			// use to pad the last byte of the original text, because we can't default
			// to a bitstring on length 7
			lastByteInfo := encodedText[len(encodedText)-1]
			numberOfPaddingZeroesAtEnd := int(lastByteInfo)
			if numberOfPaddingZeroesAtEnd > 0 {
				bitString = fmt.Sprintf("%0*s", numberOfPaddingZeroesAtEnd+len(bitString), bitString)
			}
		}
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
		log.Fatalln("Error reading input file", err)
	}

	var result []byte
	if process == "encode" {
		result = encode(file)
		reducedSize := 100 - (100 * float32(len(result)) / float32(len(file)))
		fmt.Printf("Filesize from %d to %d bytes. A reduction of %.2f\n", len(file), len(result), reducedSize)
	} else {
		result = decode(file)
	}

	if outputFile != "" {
		os.WriteFile(outputFile, result, 0644)
	} else {
		fmt.Println(string(result))
	}
}
