package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"unicode/utf8"
)

const HEADER_END = 'ӧ'

func BuildCharFrequency(file string) map[rune]int {
	charFrequency := make(map[rune]int)
	for i := 0; i < len(file); {
		r, w := utf8.DecodeRuneInString(file[i:])
		if _, ok := charFrequency[r]; ok {
			charFrequency[r] += 1
		} else {
			charFrequency[r] = 1
		}

		i += w
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
		newNode := Node{Weight: v, Char: k, LeftNode: nil, RightNode: nil, IsLeaf: true}
		nodes = append(nodes, newNode)
	}

	tree := BuildHuffmanTree(nodes)
	prefixCodeTable := BuildPrefixCodeTable(tree)

	// Build header
	headerString := BuildHeaderTree(&tree, "")
	encodedFile := []byte(headerString)
	encodedFile = utf8.AppendRune(encodedFile, HEADER_END)

	currentPrefix := ""
	for i := 0; i < len(file); {
		c, w := utf8.DecodeRune(file[i:])
		prefixStr := prefixCodeTable[c]
		currentPrefix += prefixStr
		for len(currentPrefix) > 7 {
			bitString := currentPrefix[:7]
			byteString := bitStringToByte(bitString)
			encodedFile = append(encodedFile, byteString)
			currentPrefix = currentPrefix[7:]
		}
		i += w
	}
	bitString := currentPrefix

	// We alway have text remaining in currentPrefix
	for len(currentPrefix) > 0 {
		limit := min(len(currentPrefix), 7)
		bitString = currentPrefix[:limit]
		byteString := bitStringToByte(bitString)
		encodedFile = append(encodedFile, byteString)
		currentPrefix = currentPrefix[limit:]
	}

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
	headerStart, headerEnd := 0, 0
	for i := 0; i < len(file); {
		c, w := utf8.DecodeRune(file[i:])
		if c == HEADER_END {
			headerStart = i
			headerEnd = i + w
			break
		}
		i += w
	}
	header := string(file[0:headerStart])
	encodedText := file[headerEnd:]
	bitStringRune := make(map[string]rune)
	remainingHeader := BuildTreeFromHeader(header, "", bitStringRune)
	if remainingHeader != "" {
		log.Println(header)
		log.Fatalln("Can't continue because header was not completly processed, remaining:", remainingHeader)
	}
	currentBitString := ""

	decodedText := []byte{}
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
		currentBitString += bitString

		l, r := 0, 0
		for l < len(currentBitString) && r < len(currentBitString) {
			bitStringToDecode := currentBitString[l : r+1]
			if runeValue, ok := bitStringRune[bitStringToDecode]; ok {
				decodedText = utf8.AppendRune(decodedText, runeValue)
				l = r + 1
				r = l
			} else {
				r += 1
			}
		}
		currentBitString = currentBitString[l:]
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
		fmt.Printf("Filesize from %d to %d bytes. A reduction of %.2f%%\n", len(file), len(result), reducedSize)
	} else {
		result = decode(file)
	}

	if outputFile != "" {
		os.WriteFile(outputFile, result, 0644)
	} else {
		fmt.Println(string(result))
	}
}
