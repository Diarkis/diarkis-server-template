package timed

import (
	"sync"
)

/*
Queue represents the data structure of a function queue
*/
type Queue struct{ sync.RWMutex }

/*
NewQueue creates a new timed queue with interval.

	[IMPORTANT] If interval is 0, queue will not start and Add() will execute operation immediately
*/
func NewQueue(interval int64) *Queue {
	return nil
}

/*
Start starts the queue
*/
func (q *Queue) Start() {
}

/*
Add adds a new operation into queue
*/
func (q *Queue) Add(operation func()) {
}

/*
Stop stops queue, but keeps the items in the queue
*/
func (q *Queue) Stop() {
}

/*
Reset stops queue and discards all operations in the queue
*/
func (q *Queue) Reset() {
}
