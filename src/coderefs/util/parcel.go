package util

import (
	"sync"
)

/*
Parcel represents a map of primitive data with string keys for safe data transportation such as mesh etc.
*/
type Parcel struct{ sync.RWMutex }

/*
NewParcel creates a new parcel instance.
*/
func NewParcel() *Parcel {
	return nil
}

/*
Add stores a new value or replace an existing value along with its key.
The value MUST be of a primitive data type.
*/
func (p *Parcel) Add(key string, value interface{}) bool {
	return false
}

/*
Remove deletes a stored key along with its value.
*/
func (p *Parcel) Remove(key string) bool {
	return false
}

/*
Export returns a copy of stored keys and values.
*/
func (p *Parcel) Export() map[string]interface{} {
	return nil
}
