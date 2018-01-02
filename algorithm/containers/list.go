/*
Author:    Alexey Osorgin (alexey.osorgin@gmail.com)
Copyright: Alexey Osorgin, 2017

Brief:     Intrusive list's implementation
*/

package containers

// ListNode is used to store node into list
type ListNode interface {
	Prev() ListNode
	Next() ListNode
	Value() interface{}

	IsEmpty() bool
	Push(node ListNode)
	PushTail(node ListNode)
	Move(head ListNode)
	MoveTail(head ListNode)
	Delete()
}

// Init is used o initalize list

type binaryList struct {
	prevNode ListNode
	nextNode ListNode
	value    interface{}
}

func NewList() ListNode {
	node := binaryList{}
	node.init()
	return &node
}

func NewListNode(value interface{}) ListNode {
	node := binaryList{value: value}
	node.init()
	return &node
}

func (l *binaryList) Value() interface{} {
	return l.value
}

func (l *binaryList) Prev() ListNode {
	return l.prevNode
}

func (l *binaryList) setPrev(node ListNode) {
	l.prevNode = node
}

func (l *binaryList) Next() ListNode {
	return l.nextNode
}

func (l *binaryList) setNext(node ListNode) {
	l.nextNode = node
}

func (l *binaryList) init() {
	l.setPrev(l)
	l.setNext(l)
}

// IsEmpty is used to check if list contains no elements
func (l *binaryList) IsEmpty() bool {
	return l.Prev() == l && l.Next() == l
}

func (l *binaryList) Push(n ListNode) {
	node := n.(*binaryList)
	node.setPrev(l)
	node.setNext(l.Next())
	l.Next().(*binaryList).setPrev(n)
	l.setNext(n)
}

func (l *binaryList) PushTail(n ListNode) {
	node := n.(*binaryList)
	node.setNext(l)
	node.setPrev(l.Prev())
	l.Prev().(*binaryList).setNext(node)
	l.setPrev(node)
}

func (l *binaryList) Move(head ListNode) {
	l.Delete()
	head.Push(l)
}

func (l *binaryList) MoveTail(head ListNode) {
	l.Delete()
	head.PushTail(l)
}

func (l *binaryList) Delete() {
	l.Prev().(*binaryList).setNext(l.Next())
	l.Next().(*binaryList).setPrev(l.Prev())
	l.init()
}
