package timed

/*
Circular defines the state to be a circular - meaning the change of state does not stop
*/
func (state *State) Circular() {
}

/*
NonCircular defines the state to be a non-circular - meaning the change of state stops at the either beginning or the end
*/
func (state *State) NonCircular() {
}

/*
Start starts the state instance
*/
func (state *State) Start() {
}

/*
Forward the change direction of the state will be forward
*/
func (state *State) Forward() {
}

/*
Backward the change direction of the state will be backward
*/
func (state *State) Backward() {
}

/*
GetCurrentState returns the current state
*/
func (state *State) GetCurrentState() int {
	return 0
}

/*
Next moves the state forward by the given value
*/
func (state *State) Next(step int) {
}

/*
Back moves the state backwards by the given value
*/
func (state *State) Back(step int) {
}

/*
GetProperties returns the properties: states, interval
*/
func (state *State) GetProperties() ([]int, int64) {
	return nil, 0
}

/*
ToStart moves the state to the first state of the states array
*/
func (state *State) ToStart() {
}

/*
ToEnd moves the state to the end state of the states array
*/
func (state *State) ToEnd() {
}
