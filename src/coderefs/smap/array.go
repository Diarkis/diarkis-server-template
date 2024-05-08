package smap

import (
	"sync"
)

/*
SyncArray is a thread-safe array.
*/
type SyncArray struct{ sync.RWMutex }

/*
Element is a value struct that is inserted in SyncArray.
Element.ID contains the actual value.
*/
type Element struct {
	ID     string
	Points int64
}

/*
Item represents each element stored in the list with its index.
*/
type Item struct {
	ID     string
	Points int64
	Index  int
}

/*
NewASC creates a new SyncArray in ascending order.
*/
func NewASC() *SyncArray {
	return nil
}

/*
NewDESC creates a new SyncArray in descending order.
*/
func NewDESC() *SyncArray {
	return nil
}

/*
Insert inserts the given element in the ascending or descending order.
Returns the inserted index.
*/
func (sa *SyncArray) Insert(elm *Element) int {
	return 0
}

/*
Remove deletes an element that matches with the value interface{} of the given Element.
Returns the index of the deleted element that holds the given value or -1 if not deleted.
*/
func (sa *SyncArray) Remove(elm *Element) int {
	return 0
}

/*
Search returns the searched result item by the given Element with its index.
*/
func (sa *SyncArray) Search(elm *Element) *Item {
	return nil
}

/*
SearchRange returns a list of items searched by the given Element and given range up to down.
*/
func (sa *SyncArray) SearchRange(elm *Element, up int, down int) []*Item {
	return nil
}

/*
GetAt returns item with its index.
Returns nil if the index is invalid.
*/
func (sa *SyncArray) GetAt(index int) *Item {
	return nil
}

/*
GetRange returns a list of items with in the given range.
*/
func (sa *SyncArray) GetRange(from int, to int) []*Item {
	return nil
}

/*
GetAll returns a list of all items.
*/
func (sa *SyncArray) GetAll() []*Item {
	return nil
}

/*
Length returns the length of the array.
*/
func (sa *SyncArray) Length() int {
	return 0
}

/*
Clear deletes all elements in the array.
*/
func (sa *SyncArray) Clear() {
}
