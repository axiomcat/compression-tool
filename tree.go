package main

import (
	"container/heap"
	"fmt"
	"slices"
)

type Node struct {
	Weight    int
	Char      rune
	LeftNode  *Node
	RightNode *Node
}

func compareNodes(nodeA Node, nodeB Node) int {
	return nodeA.Weight - nodeB.Weight

}

func printNodes(nodes []Node) {
	for _, n := range nodes {
		if n.Char == 10 {
			continue
		}
		fmt.Printf("%s:%d|", string(n.Char), n.Weight)
	}
	fmt.Println()
}

func printTree(node *Node) {
	if node == nil {
		return
	}
	printTree(node.LeftNode)
	fmt.Printf("%d:%s\n", node.Weight, string(node.Char))
	printTree(node.RightNode)
}

func BuildHuffmanTree(nodes []Node) Node {
	slices.SortFunc(nodes, compareNodes)
	usePQ := true
	if usePQ {
		pq := make(PriorityQueue, len(nodes))
		i := 0

		for _, n := range nodes {
			pq[i] = &Item{
				value:    n,
				priority: n.Weight,
				index:    -i,
			}
			i++
		}
		heap.Init(&pq)
		for pq.Len() > 1 {
			tempL := heap.Pop(&pq).(*Item)
			tempR := heap.Pop(&pq).(*Item)
			newRoot := Node{Weight: tempL.value.Weight + tempR.value.Weight, LeftNode: &tempL.value, RightNode: &tempR.value}
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

	for len(nodes) > 1 {
		tempL := nodes[0]
		tempR := nodes[1]
		newRoot := Node{Weight: tempL.Weight + tempR.Weight, LeftNode: &tempL, RightNode: &tempR}
		nodes = slices.Delete(nodes, 0, 2)
		nodes = append(nodes, newRoot)
		slices.SortFunc(nodes, compareNodes)
	}
	return nodes[0]
}

func BuildPrefixCodeTable(tree Node) map[rune]string {
	prefixCodeTable := make(map[rune]string)
	buildPrefixCodeTableAux(&tree, prefixCodeTable, "")
	// for k, v := range prefixCodeTable {
	// 	fmt.Printf("%s:%s\n", string(k), v)
	// }
	return prefixCodeTable
}

func buildPrefixCodeTableAux(tree *Node, table map[rune]string, currPrefix string) {
	// Is a leaf
	if tree.Char != 0 {
		table[tree.Char] = currPrefix
		return
	}
	buildPrefixCodeTableAux(tree.LeftNode, table, currPrefix+"0")
	buildPrefixCodeTableAux(tree.RightNode, table, currPrefix+"1")
}
