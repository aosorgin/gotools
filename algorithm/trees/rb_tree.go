package trees

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

// OuterLeft finds out node where node.key <= key
func (n *node) outerLeft(key KeyType) *node {
	cur, parent := n, n.parent
	for cur != nil {
		if key == cur.key {
			return cur
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

// Set is used to store keys in rb-tree
type Set struct {
	root *node
}

// Lookup find the node where node.key == key
func (s *Set) Lookup(key KeyType) bool {
	n := s.root.outerLeft(key)
	if key == n.key {
		return true
	}

	return false
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

// Insert node to rb-tree
func (s *Set) Insert(key KeyType) bool {
	if s.root == nil {
		s.root = &node{
			key: key,
		}
		return true
	}

	parent := s.root.outerLeft(key)
	if parent.key == key {
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

	for s.root.parent != nil {
		s.root = s.root.parent
	}

	return true
}
