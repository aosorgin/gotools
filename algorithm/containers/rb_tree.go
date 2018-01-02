/*
Author:    Alexey Osorgin (alexey.osorgin@gmail.com)
Copyright: Alexey Osorgin, 2017

Brief:     RB-tree's implementation
*/

package containers

// KeyType is used as rb-tree node's key type
type KeyType int

type node struct {
	key    KeyType
	isRed  bool
	parent *node
	left   *node
	right  *node
}

func (n *node) grandparent() *node {
	if n.parent != nil {
		return n.parent.parent
	}
	return nil
}

func (n *node) uncle() *node {
	grand := n.grandparent()
	if grand == nil {
		return nil
	}

	if grand.left == n {
		return grand.right
	}

	return grand.left
}

func (n *node) subling() *node {
	if n.parent.left == n {
		return n.parent.right
	}
	return n.parent.left
}

func (n *node) rotateLeft() {
	pivot := n.right
	if n.parent != nil {
		if n.parent.left == n {
			n.parent.left = pivot
		} else {
			n.parent.right = pivot
		}
	}

	pivot.parent = n.parent
	n.right = pivot.left
	if n.right != nil {
		n.right.parent = n
	}
	pivot.left = n
	n.parent = pivot
}

func (n *node) rotateRight() {
	pivot := n.left
	if n.parent != nil {
		if n.parent.left == n {
			n.parent.left = pivot
		} else {
			n.parent.right = pivot
		}
	}

	pivot.parent = n.parent
	n.left = pivot.right
	if n.left != nil {
		n.left.parent = n
	}
	pivot.right = n
	n.parent = pivot
}

func (n *node) parentToInsert(key KeyType) *node {
	cur, parent := n, n.parent
	for cur != nil {
		if key == cur.key {
			return nil
		}

		parent = cur
		if key < cur.key {
			cur = cur.left
		} else {
			cur = cur.right
		}
	}

	return parent
}

// outerLeft finds out minimum node where node.key >= key
func (n *node) outerLeft(key KeyType) *node {
	cur := n
	var candidate *node
	for cur != nil {
		if key < cur.key {
			cur = cur.left
		} else {
			candidate = cur
			cur = cur.right
		}
	}

	return candidate
}

// outerRight finds out minimum node where node.key > key
func (n *node) outerRight(key KeyType) *node {
	cur := n
	var candidate *node
	for cur != nil {
		if key < cur.key {
			cur = cur.left
		} else {
			if key > cur.key {
				candidate = cur
			}
			cur = cur.right
		}
	}

	return candidate
}

func (n *node) lookup(key KeyType) *node {
	node := n.outerLeft(key)
	if node != nil && node.key == key {
		return node
	}
	return nil
}

func optimizeAfterInsert(n *node) *node {
	if n.parent == nil {
		n.isRed = false
		return nil
	}

	if n.parent.isRed == false {
		return nil
	}

	u := n.uncle()
	if u != nil && u.isRed == true {
		n.parent.isRed = false
		u.isRed = false
		g := n.grandparent()
		g.isRed = true
		return g
	}

	g := n.grandparent()
	if n.parent == g.left && n == n.parent.right {
		n.parent.rotateLeft()
		n = n.left
	} else if n.parent == g.right && n == n.parent.left {
		n.parent.rotateRight()
		n = n.right
	}

	g = n.grandparent()
	n.parent.isRed = false
	g.isRed = true

	if n == n.parent.left && n.parent == g.left {
		g.rotateRight()
	} else {
		g.rotateLeft()
	}

	return nil
}

func (n *node) insert(key KeyType) bool {
	parent := n.parentToInsert(key)
	if parent == nil {
		return false
	}

	node := &node{
		key:    key,
		parent: parent,
		isRed:  true,
	}

	if key < parent.key {
		parent.left = node
	} else {
		parent.right = node
	}

	for node != nil {
		node = optimizeAfterInsert(node)
	}

	return true
}

func exclude(n *node) {
	n.parent = nil
	n.left = nil
	n.right = nil
}

func bugOn(condition bool) {
	if condition == true {
		panic("BUG!!!")
	}
}

func deleteOptimized(n *node) (parent *node) {
	/* Check if node has at least one ont null child. This child has to be the red one */

	var child *node
	if n.left != nil {
		child = n.left
	} else if n.right != nil {
		child = n.right
	}

	/* if node is root just make the child as root if it exists */

	if n.parent == nil {
		if child != nil {
			child.parent = nil
			child.isRed = false
			return child
		}
		return nil
	}

	parent = n.parent

	if child != nil {
		if n.parent.left == n {
			n.parent.left = child
		} else {
			n.parent.right = child
		}
		child.parent = n.parent
		child.isRed = false
		return
	}

	if n.parent.left == n {
		optimizedAfterDeletion(n)
		n.parent.left = nil
	} else {
		optimizedAfterDeletion(n)
		n.parent.right = nil
	}
	return
}

func optimizedAfterDeletion(n *node) {
	if n.parent == nil {
		n.isRed = false
		return
	}

	/* if sibling is a right child of parent */

	if n.parent.left == n {
		if n.isRed == true {
			return
		}

		sibling := n.parent.right

		/* if sibling is black and it has at least one red child */
		if sibling.isRed == false {
			if sibling.right != nil && sibling.right.isRed == true {
				sibling.right.isRed = false
				n.parent.rotateLeft()
			} else if sibling.left != nil && sibling.left.isRed == true {
				sibling.isRed = false
				sibling.left.isRed = false
				sibling.rotateRight()
				n.parent.rotateLeft()
			} else {
				sibling.isRed = true
				if n.parent.isRed == true {
					n.parent.isRed = false
				} else {
					optimizedAfterDeletion(n.parent)
				}
			}
		} else {
			/* if sibling is black */
			bugOn(sibling.left == nil || sibling.left.isRed == true)
			bugOn(sibling.right == nil || sibling.right.isRed == true)
			sibling.left.isRed = true
			n.parent.rotateLeft()
		}
	} else {
		if n.isRed == true {
			return
		}

		sibling := n.parent.left

		/* if sibling is black and it has at least one red child */
		if sibling.isRed == false {
			if sibling.left != nil {
				bugOn(sibling.left.isRed == false)
				sibling.left.isRed = false
				n.parent.rotateRight()
			} else if sibling.right != nil {
				bugOn(sibling.right.isRed == false)
				sibling.isRed = false
				sibling.left.isRed = false
				sibling.rotateLeft()
				n.parent.rotateRight()
			} else {
				sibling.isRed = true
				if n.parent.isRed == true {
					n.parent.isRed = false
				} else {
					optimizedAfterDeletion(n.parent)
				}
			}
		} else {
			/* if sibling is black */
			bugOn(sibling.left == nil || sibling.left.isRed == true)
			bugOn(sibling.right == nil || sibling.right.isRed == true)
			sibling.right.isRed = true
			n.parent.rotateRight()
		}
	}
}

func (n *node) delete(key KeyType) (isDeleted bool, root *node) {
	node := n.lookup(key)
	if node == nil {
		return false, n
	}

	isDeleted = true

	if node.left != nil && node.right != nil {
		successor := node.right.outerRight(key)
		/* TODO: exchange node with successor */
		node.key = successor.key
		deleteOptimized(successor)
		exclude(successor)
		root = n
	} else {
		root = deleteOptimized(node)
		for root != nil && root.parent != nil {
			root = root.parent
		}
		exclude(node)
	}
	if root != nil && root.isRed == true {
		root.isRed = false
	}
	return
}

// Set is used to store keys in rb-tree
type Set struct {
	root *node
}

// Lookup find the node where node.key == key
func (s *Set) Lookup(key KeyType) bool {
	n := s.root.outerLeft(key)
	if n == nil || key != n.key {
		return false
	}

	return true
}

// Insert node to rb-tree
func (s *Set) Insert(key KeyType) bool {
	if s.root == nil {
		s.root = &node{
			key: key,
		}
		return true
	}

	if s.root.insert(key) == false {
		return false
	}

	for s.root.parent != nil {
		s.root = s.root.parent
	}

	return true
}

// Delete is remove node from rb-tree
func (s *Set) Delete(key KeyType) (result bool) {
	if s.root == nil {
		return false
	}

	result, s.root = s.root.delete(key)
	return
}
