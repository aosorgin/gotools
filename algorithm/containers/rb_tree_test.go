/*
Author:    Alexey Osorgin (alexey.osorgin@gmail.com)
Copyright: Alexey Osorgin, 2017

Brief:     Tests for rb-tree's implementation
*/

package containers

import (
	"fmt"
	"testing"
)

func TestInsertedNodesAreFouneded(t *testing.T) {
	var s Set
	for i := 0; i < 100000; i++ {
		if s.Insert(KeyType(i)) == false {
			t.Error(fmt.Errorf("failed fo insert node '%d'", i))
		}
	}

	for i := 0; i < 100000; i++ {
		if s.Lookup(KeyType(i)) == false {
			t.Error(fmt.Errorf("failed fo insert node '%d'", i))
		}
	}
}

func TestDeletedNodesAreNotFounded(t *testing.T) {
	var s Set
	maxNodes := 10
	for i := 0; i < maxNodes; i++ {
		if s.Insert(KeyType(i)) == false {
			t.Error(fmt.Errorf("failed to insert node '%d'", i))
		}
	}

	for i := 0; i < maxNodes; i++ {
		if s.Delete(KeyType(i)) == false {
			t.Error(fmt.Errorf("failed to remove node '%d'", i))
		}

		for j := 0; j <= i; j++ {
			if s.Lookup(KeyType(j)) == true {
				t.Error(fmt.Errorf("deleted node '%d' is found", j))
			}
		}

		for j := i + 1; j < maxNodes; j++ {
			if s.Lookup(KeyType(j)) == false {
				t.Error(fmt.Errorf("failed fo find node '%d'", j))
			}
		}
	}
}

func TestDoNotFoundUnknownNodes(t *testing.T) {
	var s Set
	for i := 0; i < 1000; i++ {
		if s.Insert(KeyType(i)) == false {
			t.Errorf("failed fo insert node '%d'", i)
		}
	}

	for i := 1001; i < 200; i++ {
		if s.Lookup(KeyType(i)) == true {
			t.Errorf("uninsterted node '%d' is found", i)
		}
	}
}

func getNodePath(n *node) string {
	nodeID := func(x *node) string {
		res := fmt.Sprintf("%d", x.key)
		if x.isRed {
			res += "(r)"
		} else {
			res += "(b)"
		}
		return res
	}

	res := nodeID(n)
	for n.parent != nil {
		n = n.parent
		res += nodeID(n)
	}
	return res
}

func getTreeDepth(t *testing.T, n *node, depth int, maxDepth *int) {
	if n == nil {
		return
	}

	nextDepth := depth
	if !n.isRed {
		nextDepth++
	}
	if n.left != nil {
		getTreeDepth(t, n.left, nextDepth, maxDepth)
	}
	if n.right != nil {
		getTreeDepth(t, n.right, nextDepth, maxDepth)
	}

	if n.left == nil && n.right == nil {
		if *maxDepth == 0 {
			*maxDepth = nextDepth
		} else if *maxDepth != nextDepth {
			printTree(n)
			t.Fatalf("There is node '%s' with depth '%d'. Max depth id '%d'", getNodePath(n), depth+1, *maxDepth)
		}
	}
}

func printChildren(n *node) {
	if n.left != nil {
		fmt.Printf("key: %d, side: left, parent: %d, red: %t\n", n.left.key, n.left.parent.key, n.left.isRed)
		printChildren(n.left)
	}
	if n.right != nil {
		fmt.Printf("key: %d, side: right, parent: %d, red: %t\n", n.right.key, n.right.parent.key, n.right.isRed)
		printChildren(n.right)
	}
}

func printTree(n *node) {
	if n == nil {
		return
	}

	for n.parent != nil {
		n = n.parent
	}
	fmt.Printf("key: %d, red: %t\n", n.key, n.isRed)
	printChildren(n)
}

func TestAllPathsHaveTheSameBlackNodesCount(t *testing.T) {
	max := 5000
	var s Set
	for i := 0; i < max; i++ {
		if s.Insert(KeyType(i)) == false {
			t.Error(fmt.Errorf("failed fo insert node '%d'", i))
		}
		var maxDepth int
		getTreeDepth(t, s.root, 0, &maxDepth)
	}

	/*for i := 0; i < max; i++ {
		if s.Delete(KeyType(i)) == false {
			t.Error(fmt.Errorf("failed fo remove node '%d'", i))
		}
		var maxDepth int
		getTreeDepth(t, s.root, 0, &maxDepth)
	}*/
}

func TestAllPathsHaveTheSameBlackNodesCountFixed(t *testing.T) {
	max := 12

	var s Set
	for i := 0; i < max; i++ {
		if s.Insert(KeyType(i)) == false {
			t.Error(fmt.Errorf("failed fo insert node '%d'", i))
		}
	}

	printTree(s.root)

	for i := 0; i < max; i++ {
		fmt.Printf("\nRemoving %d...\n\n", i)
		if s.Delete(KeyType(i)) == false {
			t.Error(fmt.Errorf("failed fo remove node '%d'", i))
		}
		printTree(s.root)
		var maxDepth int
		getTreeDepth(t, s.root, 0, &maxDepth)
	}
}
