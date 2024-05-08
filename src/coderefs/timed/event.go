package timed

/*
Start starts the timed event
*/
func (event *Event) Start() {
}

/*
Stop stops the timed event
*/
func (event *Event) Stop() {
}

/*
OnTick registers a callback function to be executed at every event interval (tick)
*/
func (event *Event) OnTick(callback func()) {
}

/*
OnError registers a callback function to be executed if callbacks panics at event tick
*/
func (event *Event) OnError(callback func(error)) {
}
