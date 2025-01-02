package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	len   int
	front *ListItem
	back  *ListItem
}

func (list *list) Len() int {
	return list.len
}

func (list *list) Front() *ListItem {
	return list.front
}

func (list *list) Back() *ListItem {
	return list.back
}

func (list *list) PushFront(v interface{}) *ListItem {
	item := &ListItem{Value: v, Next: list.front}

	if list.front != nil {
		list.front.Prev = item
	} else {
		list.back = item // If list is empty, the new element becomes both front and back
	}

	list.front = item
	list.len++

	return item
}

func (list *list) PushBack(v interface{}) *ListItem {
	item := &ListItem{Value: v, Prev: list.back}

	if list.back != nil {
		list.back.Next = item
	} else {
		list.front = item // If list is empty, the new element becomes both front and back
	}

	list.back = item
	list.len++

	return item
}

func (list *list) Remove(i *ListItem) {
	if prev := i.Prev; prev != nil {
		prev.Next = i.Next
	} else {
		list.front = i.Next
	}

	if next := i.Next; next != nil {
		next.Prev = i.Prev
	} else {
		list.back = i.Prev
	}

	list.len--
}

func (list *list) MoveToFront(i *ListItem) {
	list.Remove(i)
	list.PushFront(i.Value)
}

func NewList() List {
	return new(list)
}
