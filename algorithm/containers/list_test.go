/*
Author:    Alexey Osorgin (alexey.osorgin@gmail.com)
Copyright: Alexey Osorgin, 2017

Brief:     Intrusive list's implementation
*/

package containers

import (
	"testing"
)

func TestIfNewListIsEmpty(t *testing.T) {
	list := NewList()
	if list.IsEmpty() == false {
		t.Error("Just initialized list is not empty")
	}
}

func TestPushedItemIsInList(t *testing.T) {
	list := NewList()

	for i := 0; i < 10; i++ {
		list.Push(NewListNode(i))
	}

	i := 9
	for node := list.Next(); node != list; node = node.Next() {
		v := node.Value().(int)
		if v != i {
			t.Errorf("List contained node '%d' instead of '%d'", v, i)
		}
		i--
	}
}

func TestPushedItemInTailIsInList(t *testing.T) {
	list := NewList()

	for i := 0; i < 10; i++ {
		list.PushTail(NewListNode(i))
	}

	i := 0
	for node := list.Next(); node != list; node = node.Next() {
		v := node.Value().(int)
		if v != i {
			t.Errorf("List contained node '%d' instead of '%d'", v, i)
		}
		i++
	}
}

func TestMovedItemIsInSecondList(t *testing.T) {
	list := NewList()

	for i := 0; i < 10; i++ {
		list.Push(NewListNode(i))
	}

	list2 := NewList()
	for i := 0; i < 5; i++ {
		list.Next().MoveTail(list2)
	}

	i := 9
	for node := list2.Next(); node != list2; node = node.Next() {
		v := node.Value().(int)
		if v != i {
			t.Errorf("List contained node '%d' instead of '%d'", v, i)
		}
		i--
	}

	for node := list.Next(); node != list; node = node.Next() {
		v := node.Value().(int)
		if v != i {
			t.Errorf("List contained node '%d' instead of '%d'", v, i)
		}
		i--
	}
}

func TestMovedItemInTailIsInSecondList(t *testing.T) {
	list := NewList()

	for i := 0; i < 10; i++ {
		list.PushTail(NewListNode(i))
	}

	list2 := NewList()
	for i := 0; i < 5; i++ {
		list.Next().Move(list2)
	}

	i := 4
	for node := list2.Next(); node != list2; node = node.Next() {
		v := node.Value().(int)
		if v != i {
			t.Errorf("List contained node '%d' instead of '%d'", v, i)
		}
		i--
	}

	i = 5
	for node := list.Next(); node != list; node = node.Next() {
		v := node.Value().(int)
		if v != i {
			t.Errorf("List contained node '%d' instead of '%d'", v, i)
		}
		i++
	}
}

func TestDeletedItemNotInList(t *testing.T) {
	list := NewList()

	for i := 0; i < 10; i++ {
		list.PushTail(NewListNode(i))
	}

	for i := 0; i < 5; i++ {
		list.Next().Delete()
	}

	i := 5
	for node := list.Next(); node != list; node = node.Next() {
		v := node.Value().(int)
		if v != i {
			t.Errorf("List contained node '%d' instead of '%d'", v, i)
		}
		i++
	}
}
