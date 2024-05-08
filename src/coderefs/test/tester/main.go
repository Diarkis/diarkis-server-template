package tester

/*
Tests represents a set of unit tests
*/
type Tests struct{}

/*
Test represents a unit test
*/
type Test struct{}

/*
Wait keeps the process running until the given flag becomes true
*/
func Wait(flag bool) {
}

/*
NewTests creates a new Tests instance
*/
func NewTests(name string) *Tests {
	return nil
}

/*
RunTests executes multiple *Tests as a single continuous test.
*/
func RunTests(list []*Tests) {
}

/*
Define creates a unit test. If you pass an error to the callback, the unit test fails.

	tests.Define("Can run some actions", func(next func(error)) {

	  // we run some actions here
	  // by calling next(), we tell the unit test container to move to the next test
	  // if you pass an error to next(), the test fails and unit test terminates immediately
	  next(nil)

	})
*/
func (ts *Tests) Define(label string, logic func(callback func(err error))) {
}

/*
OnEnd assigns a callback to be called at the end of all tests.

This is useful when you have to clean up such as database etc. after your tests
*/
func (ts *Tests) OnEnd(cb func()) {
}

/*
Count returns the number of tests.
*/
func (ts *Tests) Count() int {
	return 0
}

/*
Run executes all unit tests in the order of the tests defined.

If a unit test fails, the entire Tests fail.
*/
func (ts *Tests) Run(finished func(err error)) {
}

/*
Async executes a given callback function asynchronously.

You must use this function instead of goroutine in your tests.
*/
func (ts *Tests) Async(cb func()) {
}

/*
GetCurrentTest returns the current test.
*/
func (ts *Tests) GetCurrentTest() *Test {
	return nil
}

/*
Assert evaluates the two given values and if they are evaluated to be different values, the test fails immediately.
*/
func (ts *Tests) Assert(given interface{}, expected interface{}) {
}

/*
AssertNotEqual evaluates the two given values and if they are evaluated to be the same values, the test fails immediately.
*/
func (ts *Tests) AssertNotEqual(given interface{}, expected interface{}) {
}

/*
DebugLogging outputs debug logging during the tests if the test process is started with --debug parameter.
*/
func (ts *Tests) DebugLogging(d ...interface{}) {
}
