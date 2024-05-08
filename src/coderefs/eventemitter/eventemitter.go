package eventemitter

import (
	"sync"
)

/*
EventEmitter represents the event emitter.
*/
type EventEmitter struct{ sync.RWMutex }

/*
New creates a new event emitter
*/
func New(name string) *EventEmitter {
	return nil
}

/*
On assigns a callback to an event
*/
func (emitter *EventEmitter) On(eventName string, callback func(interface{})) {
}

/*
Emit triggers an event and invokes the callbacks
*/
func (emitter *EventEmitter) Emit(eventName string, params interface{}) {
}

/*
Remove deletes a callback assigned to the given event.

	[IMPORTANT] The callback passed is the callback to be removed.
*/
func (emitter *EventEmitter) Remove(eventName string, callback func(interface{})) bool {
	return false
}
