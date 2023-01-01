package main

import (
	"fmt"
)

// Node is a node of a tree.
type Node struct {
	Key   int
	Left  *Node
	Right *Node
}

// BST is a binary search tree.
type BST struct {
	Root *Node
}

// insert insert a node to tree.
func (b *BST) insert(key int) {
	if b.Root == nil {
		b.Root = &Node{
			Key:   key,
			Left:  nil,
			Right: nil,
		}
	} else {
		recursiveInsert(b.Root, &Node{
			Key:   key,
			Left:  nil,
			Right: nil,
		})
	}
}

// recursiveInsert insert a new node to targetNode recursively.
func recursiveInsert(targetNode *Node, newNode *Node) {
	// if a newNode is smaller than targetNode, insert a newNode to left child node.
	// if a newNode is a bigger than targetNode, insert a newNode to right childe node.
	if newNode.Key < targetNode.Key {
		if targetNode.Left == nil {
			targetNode.Left = newNode
		} else {
			recursiveInsert(targetNode.Left, newNode)
		}
	} else {
		if targetNode.Right == nil {
			targetNode.Right = newNode
		} else {
			recursiveInsert(targetNode.Right, newNode)
		}
	}
}

// remove remove a key from tree.
func (b *BST) remove(key int) {
	recursiveRemove(b.Root, key)
}

// recursiveRemove remove a key from tree recursively.
func recursiveRemove(targetNode *Node, key int) *Node {
	if targetNode == nil {
		return nil
	}

	if key < targetNode.Key {
		targetNode.Left = recursiveRemove(targetNode.Left, key)
		return targetNode
	}

	if key > targetNode.Key {
		targetNode.Right = recursiveRemove(targetNode.Right, key)
		return targetNode
	}

	if targetNode.Left == nil && targetNode.Right == nil {
		targetNode = nil
		return nil
	}

	if targetNode.Left == nil {
		targetNode = targetNode.Right
		return targetNode
	}

	if targetNode.Right == nil {
		targetNode = targetNode.Left
		return targetNode
	}

	leftNodeOfMostRightNode := targetNode.Right

	for {
		if leftNodeOfMostRightNode != nil && leftNodeOfMostRightNode.Left != nil {
			leftNodeOfMostRightNode = leftNodeOfMostRightNode.Left
		} else {
			break
		}
	}

	targetNode.Key = leftNodeOfMostRightNode.Key
	targetNode.Right = recursiveRemove(targetNode.Right, targetNode.Key)
	return targetNode
}

// search search a key from tree.
func (b *BST) search(key int) bool {
	result := recursiveSearch(b.Root, key)

	return result
}

// recursiveSearch search a key from tree recursively.
func recursiveSearch(targetNode *Node, key int) bool {
	if targetNode == nil {
		return false
	}

	if key < targetNode.Key {
		return recursiveSearch(targetNode.Left, key)
	}

	if key > targetNode.Key {
		return recursiveSearch(targetNode.Right, key)
	}

	// targetNode == key
	return true
}

// depth-first search
// inOrderTraverse traverse tree by in-order.
func (b *BST) inOrderTraverse() {
	recursiveInOrderTraverse(b.Root)
}

// recursiveInOrderTraverse traverse tree by in-order recursively.
func recursiveInOrderTraverse(n *Node) {
	if n != nil {
		recursiveInOrderTraverse(n.Left)
		fmt.Printf("%d\n", n.Key)
		recursiveInOrderTraverse(n.Right)
	}
}

// depth-first search
// preOrderTraverse traverse by pre-order.
func (b *BST) preOrderTraverse() {
	recursivePreOrderTraverse(b.Root)
}

// recursivePreOrderTraverse traverse by pre-order recursively.
func recursivePreOrderTraverse(n *Node) {
	if n != nil {
		fmt.Printf("%d\n", n.Key)
		recursivePreOrderTraverse(n.Left)
		recursivePreOrderTraverse(n.Right)
	}
}

// depth-first search
// postOrderTraverse traverse by post-order.
func (b *BST) postOrderTraverse() {
	recursivePostOrderTraverse(b.Root)
}

// recursivePostOrderTraverse traverse by post-order recursively.
func recursivePostOrderTraverse(n *Node) {
	if n != nil {
		recursivePostOrderTraverse(n.Left)
		recursivePostOrderTraverse(n.Right)
		fmt.Printf("%v\n", n.Key)
	}
}

// breadth-first search
// levelOrderTraverse traverse by level-order.
func (b *BST) levelOrderTraverse() {
	if b != nil {
		queue := []*Node{b.Root}

		for len(queue) > 0 {
			currentNode := queue[0]
			fmt.Printf("%d ", currentNode.Key)

			queue = queue[1:]

			if currentNode.Left != nil {
				queue = append(queue, currentNode.Left)
			}

			if currentNode.Right != nil {
				queue = append(queue, currentNode.Right)
			}
		}
	}
}

func main() {
	tree := &BST{}

	tree.insert(10)
	tree.insert(2)
	tree.insert(3)
	tree.insert(3)
	tree.insert(3)
	tree.insert(15)
	tree.insert(14)
	tree.insert(18)
	tree.insert(16)
	tree.insert(16)

	tree.remove(3)
	tree.remove(10)
	tree.remove(16)

	fmt.Println(tree.search(10))
	fmt.Println(tree.search(19))

	// Traverse
	tree.inOrderTraverse()
	tree.preOrderTraverse()
	tree.postOrderTraverse()
	tree.levelOrderTraverse()

	fmt.Printf("%#v\n", tree)
}
