package util

import (
	"sync"
	"time"
)

const SharedDataLimitLength = 10

/*
Await represents asynchronous wait group
*/
type Await struct{ sync.WaitGroup }

/*
Setup [INTERNAL USE ONLY] set up util package in diarkis
*/
func Setup() {
}

/*
IsPrimitiveDataType returns true if the given value is a primitive.

Data types that returns true are:

	┌─────────┐
	│    uint │
	├─────────┤
	│   uint8 │
	├─────────┤
	│  uint16 │
	├─────────┤
	│  uint32 │
	├─────────┤
	│  uint64 │
	├─────────┤
	│     int │
	├─────────┤
	│    int8 │
	├─────────┤
	│   int16 │
	├─────────┤
	│   int32 │
	├─────────┤
	│   int64 │
	├─────────┤
	│ float32 │
	├─────────┤
	│ float64 │
	├─────────┤
	│  string │
	├─────────┤
	│    bool │
	├─────────┤
	│  []byte │
	└─────────┘
*/
func IsPrimitiveDataType(v interface{}) bool {
	return false
}

/*
Async returns await struct for asynchronous operations:

Usage Example

	await := util.Async(2)
	go func() {
	  asynchronousOperation(func() {
	    await.Done()
	  })
	}
	go func() {
	  asynchronousOperation(func() {
	    await.Done()
	  })
	}
	// this will block until all asynchronous operations are marked by await.Done()
	await.Wait()
	// all asynchronous operations are done now
	finish()
*/
func Async(taskNum int) *Await {
	return nil
}

/*
Parallel executes multiple functions in parallel and calls done when all functions are finished
*/
func Parallel(funcs []func(func(error)), done func(error)) {
}

/*
ForEachParallel executes the given operation function on each item in the list in parallel and calls done when finished
*/
func ForEachParallel(list []interface{}, operation func(interface{}, func(error)), done func(error)) {
}

/*
Waterfall calls multiple functions in the order of the funcs array given and calls done func when finished.
*/
func Waterfall(funcs []func(func(error)), done func(error)) {
}

/*
GenShortID returns a randomly generated HEX encoded ID as a string (8 characters long). Does NOT guarantee uniqueness like UUID v4.
*/
func GenShortID() (string, error) {
	return "", nil
}

/*
GetPublicEndPointMSLB returns the public endpoint address in MS Azure Cloud with Loadbalancer in front
When diarkis process is ready (diarkis.OnReady(callback)), you set this value to server.SetPublicEndPoint(endpoint)
*/
func GetPublicEndPointMSLB() (string, error) {
	return "", nil
}

/*
GetPublicEndPointMS returns the public endpoint address in MS Azure Cloud
When diarkis process is ready (diarkis.OnReady(callback)), you set this value to server.SetPublicEndPoint(endpoint)
*/
func GetPublicEndPointMS() (string, error) {
	return "", nil
}

/*
GetPublicEndPointGCP returns the public endpoint hostname in Google Cloud Computing
When diarkis process is ready (diarkis.OnReady(callback)), you set this value to server.SetPublicEndPoint(endpoint)
*/
func GetPublicEndPointGCP() (string, error) {
	return "", nil
}

/*
GetPublicEndPointAWS returns the public endpoint hostname in Google Cloud Computing
When diarkis process is ready (diarkis.OnReady(callback)), you set this value to server.SetPublicEndPoint(endpoint)
*/
func GetPublicEndPointAWS() (string, error) {
	return "", nil
}

/*
GetPublicEndPointAlibaba returns the public endpoint hostname in Alibaba Cloud
When diarkis process is ready (diarkis.OnReady(callback)), you set this value to server.SetPublicEndPoint(endpoint)
*/
func GetPublicEndPointAlibaba() (string, error) {
	return "", nil
}

/*
GetPublicEndPointTencent returns the public endpoint hostname in Tencent Cloud
When diarkis process is ready (diarkis.OnReady(callback)), you set this value to server.SetPublicEndPoint(endpoint)
*/
func GetPublicEndPointTencent() (string, error) {
	return "", nil
}

/*
GetPublicEndPointLinode returns the public endpoint hostname in Linode
When diarkis process is ready (diarkis.OnReady(callback)), you set this value to server.SetPublicEndPoint(endpoint)
*/
func GetPublicEndPointLinode() (string, error) {
	return "", nil
}

/*
GetPublicEndPointGeneric [INTERNAL USE ONLY] returns the public endpoint hostname by curling ifconfig.io
When diarkis process is ready (diarkis.OnReady(callback)), you set this value to server.SetPublicEndPoint(endpoint)
*/
func GetPublicEndPointGeneric() (string, error) {
	return "", nil
}

/*
IsNullBytes returns true if the given byte array contains only x00
*/
func IsNullBytes(buf []byte) bool {
	return false
}

/*
GetEnv returns a value of an environment variable w/ the given name
all env name must have a prefix of "DIARKIS_"
*/
func GetEnv(name string) string {
	return ""
}

/*
SetEnv sets an environment variable w/ the given name
all env name will have a prefix of "DIARKIS_"
*/
func SetEnv(name string, val string) {
}

/*
WriteToTmp writes a string to a file under /tmp/ - The file with have a prefix of "DIARKIS_"
*/
func WriteToTmp(name string, val string) {
}

/*
ReadFromTmp reads from a file under /tmp/ - The file must have a prefix of "DIARKIS_"
*/
func ReadFromTmp(name string) string {
	return ""
}

/*
DeleteFromTmp removes a file under /tmp/ - The file must have a prefix of "DIARKIS_"
*/
func DeleteFromTmp(name string) {
}

/*
ZuluTimeFormat returns a string of time in UTC Zulu format: RFC 3339
*/
func ZuluTimeFormat(now time.Time) string {
	return ""
}

/*
IndexOf Returns an index of a given element in the given array of strings
*/
func IndexOf(array []string, me string) int {
	return 0
}
func ReverseArray[T any](array []T) []T {
	return nil
}

/*
ReverseBytes Reverses byte array
*/
func ReverseBytes(bytes []byte) []byte {
	return nil
}

/*
AddrToBytes Converts address (address:port) string to byte array
*/
func AddrToBytes(addr string) ([]byte, error) {
	return nil, nil
}

/*
StrConcat concatenates strings
*/
func StrConcat(strlist ...string) string {
	return ""
}

/*
CreateAddressID returns a unique string ID made of UUID v4 and encoded address given.
The length of the returned address ID is always 52.
*/
func CreateAddressID(addr string) (string, error) {
	return "", nil
}

/*
GetAddressFromAddressID returns encoded address as a string from an address ID created by CreateAddressID.
*/
func GetAddressFromAddressID(id string) (string, error) {
	return "", nil
}

/*
NanoSecHex Converts the unix timestamp in nano seconds to a hex string
- Length of the returned string is always 16
*/
func NanoSecHex() string {
	return ""
}

/*
NanoSecRandHex replaces the first 1 byte of unix timestamp in nano seconds to random bytes and converts it to a hex string
- Length of the returned string is always 16
*/
func NanoSecRandHex() string {
	return ""
}

/*
GetID Returns a unique ID with mesh node address encoded in it
- Length of the returned ID is 52
*/
func GetID(nodeAddrList []string) (string, error) {
	return "", nil
}

/*
ParseID Returns unique ID and mesh node address list from an ID created by GetID()
*/
func ParseID(id string) (string, []string, error) {
	return "", nil, nil
}

/*
NowMilliseconds Returns a Unix timestamp in milliseconds
*/
func NowMilliseconds() int64 {
	return 0
}

/*
NowNanoseconds Returns a Unix timestamp in nanoseconds
*/
func NowNanoseconds() int64 {
	return 0
}

/*
NowSeconds Returns a Unix timestamp in seconds
*/
func NowSeconds() int64 {
	return 0
}

/*
RmSpaces Removes all spaces in a string - performs a single allocation,
but may grossly over-allocate if the source string is mainly whitespace
*/
func RmSpaces(str string) string {
	return ""
}

/*
ToFixed Returns a fixed precision of decimal number as a string
*/
func ToFixed(num float64, precision int) string {
	return ""
}

/*
RandomInt Returns a random int between min and max
*/
func RandomInt(min int, max int) int {
	return 0
}

/*
LeadingZero Returns a string with leading 0
*/
func LeadingZero(num int, digit int) string {
	return ""
}

/*
BytesListToBytes converts an array of []byte to a byte array.

[IMPORTANT] each byte array in the list must NOT exceed the size of 255 bytes

Returned byte array format:

Shown below is a single data set and it repeats for as long as the given list.

	+-------------+---------------+-----+
	| size header |     bytes     | ... |
	+-------------+---------------+-----+
	|    1 byte   | variable size | ... |
	+-------------+---------------+-----+
*/
func BytesListToBytes(list [][]byte) []byte {
	return nil
}

/*
BytesToBytesList converts an byte array to a list of byte array.

Input byte array must be the output byte array of BytesListToBytes.
*/
func BytesToBytesList(bytes []byte) [][]byte {
	return nil
}
