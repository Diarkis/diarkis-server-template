package timed

/*
Number represents a number (int) that changes over time either increment or decrement as time passes
*/
type Number struct{}

/*
State represents a state (int) that changes over time either forward or backward as time passes
*/
type State struct{}

/*
Event represents an event that occurs at certain interval and executes pre-registered operations
*/
type Event struct{}

/*
NewNumber creates an instance of *timed.Number.

Parameters

	start    - Initial value to start from.
	min      - Minimum allowed value.
	max      - Maximum allowed value.
	interval - Interval in seconds for the value to be updated.
*/
func NewNumber(start, min, max, step int, interval int64) (*Number, error) {
	return nil, nil
}

/*
NewState creates an instance of *timed.State.

Parameters

	states   - An array of available states.
	           The order of the array will be the order to state change.
	start    - Starting index of the states array.
	interval - Interval in seconds for the state to change.
*/
func NewState(states []int, start int, interval int64) (*State, error) {
	return nil, nil
}

/*
NewEvent creates an instance of *timed.Event

Parameters

	interval - Interval in seconds for the event to be triggered.
*/
func NewEvent(interval int64) (*Event, error) {
	return nil, nil
}

/*
SerializeNumber serializes an instance of timed.Number into a string
*/
func SerializeNumber(num *Number) string {
	return ""
}

/*
DeserializeNumber deserializes a serialized string of an instance of timed.Number
*/
func DeserializeNumber(str string) (*Number, error) {
	return nil, nil
}

/*
SerializeState serializes a state.
*/
func SerializeState(state *State) string {
	return ""
}

/*
DeserializeState deserializes a serialized state.
*/
func DeserializeState(str string) (*State, error) {
	return nil, nil
}

/*
Once executes the given callback function after the given time (in milliseconds) once.
*/
func Once(wait int64, callback func(interface{}), options interface{}) {
}
