package timed

/*
Incrementer instance will increment by step at each interval
*/
func (num *Number) Incrementer() {
}

/*
Decrementer instance will decrement by step at each interval
*/
func (num *Number) Decrementer() {
}

/*
Start starts the instance
*/
func (num *Number) Start() bool {
	return false
}

/*
GetCurrentValue returns the current value
*/
func (num *Number) GetCurrentValue() int {
	return 0
}

/*
Incr increments the value by the given value
*/
func (num *Number) Incr(val int) {
}

/*
Decr decrements the value by the given value
*/
func (num *Number) Decr(val int) {
}

/*
Set sets the current value to be the given value
*/
func (num *Number) Set(val int) {
}

/*
ToMax sets the current value to be max
*/
func (num *Number) ToMax() {
}

/*
ToMin sets the current value to be min
*/
func (num *Number) ToMin() {
}

/*
GetProperties returns all properties:
start, min, max, step, interval
*/
func (num *Number) GetProperties() (int, int, int, int, int64) {
	return 0, 0, 0, 0, 0
}

/*
IsIncrementer returns true if the instance is incrementer.
*/
func (num *Number) IsIncrementer() bool {
	return false
}

/*
IsDecrementer returns true if the instance is decrementer.
*/
func (num *Number) IsDecrementer() bool {
	return false
}
