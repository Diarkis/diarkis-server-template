package http

import (
	"net/http"
	"regexp"
	"sync"
)

const Type = "HTTP"
const Ok = 200
const Bad = 400
const NotFound = 404
const Err = 500
const NoService = 503

/*
Route URL route data structure
*/
type Route struct {
	Route      string
	Finder     *regexp.Regexp
	Extractor  *regexp.Regexp
	Handlers   []func(*Response, *Request, *Params, func(error))
	ParamNames []string
}

/*
Params URL parameter data structure
*/
type Params struct{ Values map[string]string }

/*
Response Response data structure
*/
type Response struct {
	Method   string
	URI      string
	Writer   http.ResponseWriter
	Finished bool
	sync.RWMutex
}

/*
Request Request data structure
*/
type Request struct {
	Method   string
	URL      string
	Queries  map[string][]string
	Req      *http.Request
	JSONBody map[string]interface{}
}

/*
Setup Loads configuration file into memory - pass an empty string to load nothing
*/
func Setup(confpath string) {
}

/*
SetupAsAuthServer Sets Up HTTP server as Auth server for TCP/UDP - pass an empty string to load nothing
*/
func SetupAsAuthServer(path string) {
}

/*
SetAllowOrigin sets Access-Control-Allow-Origin response header for CORS
*/
func SetAllowOrigin(origin string) {
}

/*
HandleStaticFiles handle static files
*/
func HandleStaticFiles(uriRoot string, fsPath string) {
}

/*
Online Sets mesh network status to online
*/
func Online() {
}

/*
Offline Sets mesh network status to offline
*/
func Offline() {
}

/*
SetCustomNodeSelector replaces the default node selector with custom selector for Auth
*/
func SetCustomNodeSelector(selector func([]string) string) {
}

/*
Options Registers a handler for an OPTIONS method URL
*/
func Options(uri string, handler func(*Response, *Request, *Params, func(error))) {
}

/*
Get Registers a handler for a GET method URL
*/
func Get(uri string, handler func(*Response, *Request, *Params, func(error))) {
}

/*
Post Registers a handler for a POST method URL
*/
func Post(uri string, handler func(*Response, *Request, *Params, func(error))) {
}

/*
Put Registers a handler for a PUT method URL
*/
func Put(uri string, handler func(*Response, *Request, *Params, func(error))) {
}

/*
Head Registers a handler for a HEAD method URL
*/
func Head(uri string, handler func(*Response, *Request, *Params, func(error))) {
}

/*
Delete Registers a handler for a DELETE method URL
*/
func Delete(uri string, handler func(*Response, *Request, *Params, func(error))) {
}

/*
Command Registers a handler for a given URL
*/
func Command(method string, uri string, handler func(*Response, *Request, *Params, func(error))) {
}

/*
Start Starts HTTP server
*/
func Start(next func(error)) {
}

/*
ShutdownHTTP Stops HTTP server
*/
func ShutdownHTTP() {
}

/*
GetAsString a URL parameter value as a string
*/
func (params *Params) GetAsString(name string) (string, error) {
	return "", nil
}

/*
GetAsInt Returns a URL parameter value as an int
*/
func (params *Params) GetAsInt(name string) (int, error) {
	return 0, nil
}

/*
GetAsInt8 Returns a URL parameter value as an int8
*/
func (params *Params) GetAsInt8(name string) (int8, error) {
	return 0, nil
}

/*
GetAsInt16 Returns a URL parameter value as an int16
*/
func (params *Params) GetAsInt16(name string) (int16, error) {
	return 0, nil
}

/*
GetAsInt32 Returns a URL parameter value as an int32
*/
func (params *Params) GetAsInt32(name string) (int32, error) {
	return 0, nil
}

/*
GetAsInt64 Returns a URL parameter value as an int64
*/
func (params *Params) GetAsInt64(name string) (int64, error) {
	return 0, nil
}

/*
GetAsUint8 Returns a URL parameter value as a uint8
*/
func (params *Params) GetAsUint8(name string) (uint8, error) {
	return 0, nil
}

/*
GetAsUint16 Returns a URL parameter value as a uint16
*/
func (params *Params) GetAsUint16(name string) (uint16, error) {
	return 0, nil
}

/*
GetAsUint32 Returns a URL parameter value as a uint32
*/
func (params *Params) GetAsUint32(name string) (uint32, error) {
	return 0, nil
}

/*
GetAsUint64 Returns a URL parameter value as a uint64
*/
func (params *Params) GetAsUint64(name string) (uint64, error) {
	return 0, nil
}

/*
GetAsFloat64 Returns a URL parameters value as a float64
*/
func (params *Params) GetAsFloat64(name string) (float64, error) {
	return 0, nil
}

/*
GetAsBool Returns a URL parameter value as a boolean
*/
func (params *Params) GetAsBool(name string) (bool, error) {
	return false, nil
}

/*
SetHeader sets a custom response header
*/
func (res *Response) SetHeader(name, value string) {
}

/*
Respond Sends a response packet and HTTP status code
*/
func (res *Response) Respond(data string, status int) {
}

/*
SendBytes sends a response as binary data
*/
func (res *Response) SendBytes(data []byte, status int) {
}

/*
GetPostData Returns postMethod body data by name
*/
func (req *Request) GetPostData(name string) string {
	return ""
}

/*
OnAuthRequest hooks a request handler function to auth request. If invoked before SetupAsAuthServer(),
the hook will be called BEFORE the actual auth handling. If invoked after SetupAsAuthServer(), the hook will  be called AFTER.
*/
func OnAuthRequest(handler func(res *Response, req *Request, params *Params, next func(error))) {
}

/*
GetRealTimeNodeEndpoint returns a WebSocket/TC/UDP server endpoint
Returns a node endpoint that is not taken and online
*/
func GetRealTimeNodeEndpoint(list []string) string {
	return ""
}

/*
GetRoomNumbersByNodeAddress returns the number of rooms on a remote node by its mesh address
*/
func GetRoomNumbersByNodeAddress(addr string) int {
	return 0
}

/*
GetCCUByNodeAddress returns CCU of a remote node by its mesh address
*/
func GetCCUByNodeAddress(addr string) int {
	return 0
}
