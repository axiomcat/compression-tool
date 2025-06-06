package main

import (
	"container/heap"
	"slices"
	"unicode/utf8"
)

type Node struct {
	Weight    int
	Char      rune
	LeftNode  *Node
	RightNode *Node
	IsLeaf    bool
}

func compareNodes(nodeA Node, nodeB Node) int {
	return nodeA.Weight - nodeB.Weight

}

func BuildHuffmanTree(nodes []Node) Node {
	slices.SortFunc(nodes, compareNodes)
	pq := make(PriorityQueue, len(nodes))
	i := 0

	for _, n := range nodes {
		pq[i] = &Item{
			value:    n,
			priority: -n.Weight,
			index:    i,
		}
		i++
	}
	heap.Init(&pq)
	for pq.Len() > 1 {
		tempL := heap.Pop(&pq).(*Item)
		tempR := heap.Pop(&pq).(*Item)
		newRoot := Node{Weight: tempL.value.Weight + tempR.value.Weight, LeftNode: &tempL.value, RightNode: &tempR.value, IsLeaf: false}
		newItem := &Item{
			value:    newRoot,
			priority: 1,
		}
		heap.Push(&pq, newItem)
		pq.Update(newItem, newItem.value, -newItem.value.Weight)
	}
	treeRoot := heap.Pop(&pq).(*Item)
	return treeRoot.value
}

func BuildPrefixCodeTable(tree Node) map[rune]string {
	prefixCodeTable := make(map[rune]string)
	buildPrefixCodeTableAux(&tree, prefixCodeTable, "")
	return prefixCodeTable
}

func buildPrefixCodeTableAux(tree *Node, table map[rune]string, currPrefix string) {
	// Is a leaf
	if tree.IsLeaf {
		table[tree.Char] = currPrefix
		return
	}
	buildPrefixCodeTableAux(tree.LeftNode, table, currPrefix+"0")
	buildPrefixCodeTableAux(tree.RightNode, table, currPrefix+"1")
}

func BuildHeaderTree(tree *Node, path string) string {
	if tree.Char != 0 {
		return "1" + string(tree.Char)
	}
	if tree == nil {
		return ""
	}
	l := BuildHeaderTree(tree.LeftNode, path+"0")
	r := BuildHeaderTree(tree.RightNode, path+"1")
	return "0" + l + r
}

func BuildTreeFromHeader(header string, currentPrefix string, prefixCodeTable map[string]rune) string {
	headerValue := header[0]
	if headerValue == '0' {
		// Remove 0 from the header
		header = header[1:]
		header = BuildTreeFromHeader(header, currentPrefix+"0", prefixCodeTable)
		header = BuildTreeFromHeader(header, currentPrefix+"1", prefixCodeTable)
		return header
	} else {
		// Get the char from the header
		char, w := utf8.DecodeRuneInString(header[1:])
		// remove 1 and the character from the array
		header = header[1+w:]
		prefixCodeTable[currentPrefix] = char
		return header
	}
}
