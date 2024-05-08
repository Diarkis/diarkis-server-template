package matching

import (
	"io"
)

/*
Test this is used ONLY in tests
*/
func Test(method string, data map[string]interface{}) ([]byte, error) {
	return nil, nil
}

/*
TestDebugDataDump is used ONLY in internal tests.
*/
func TestDebugDataDump() map[string][]*searchItem {
	return nil
}

/*
TTLTest this is used ONLY in tests
*/
func TTLTest(src []int64) []int64 {
	return nil
}

/*
DebugDataDump returns the entire matchmaking data held in memory of the server.

	[IMPORTANT] This is a debug function and must NOT be used in production code at all.

	[NOTE] In order to evaluate key, use CreateKeyByProperties to reproduce the same key with property values that match.
	[NOTE] Uses mutex lock internally.

Example of evaluating the dump data keys:

	dump := matching.DebugDataDump()

	for key, searchItems := range dump {

		// expected key
		expectedKey := matching.CreateKeyByProperties(profileID, tag, expectedProperties)

		// check to see if the key matches the expected key
		if key == expectedKey {
			// good
		} else {
			// bad...
		}

	}
*/
func DebugDataDump() (map[string][]*searchItem, error) {
	return nil, nil
}

/*
DebugDataDumpWriter writes the entire matchmaking data held in memory of the server to io.Writer stream.

	[IMPORTANT] This is a debug function and must NOT be used in production code at all.

	[NOTE] In order to evaluate key, use CreateKeyByProperties to reproduce the same key with property values that match.
	[NOTE] Uses mutex lock internally.
*/
func DebugDataDumpWriter(writer io.Writer) error {
	return nil
}

/*
CreateKeyByProperties generates a key from a tag and search or add properties.

This is meant to be used along with DebugDataDump for test and debugging.

The key generated is the actual search condition data used internally to find matches.
*/
func CreateKeyByProperties(matchingID string, tag string, properties map[string]int) string {
	return ""
}
