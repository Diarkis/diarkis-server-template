package http

import (
	"sync"
)

/*
ResponseData represents an HTTP response JSON data.
*/
type ResponseData struct{ sync.RWMutex }

/*
NewResponseData creates a new HTTP ResponseData.
*/
func NewResponseData() *ResponseData {
	return nil
}

/*
Add adds a key and value to be encoded into JSON.
The value data type supported is only primitive data types.
*/
func (rd *ResponseData) Add(key string, value interface{}) bool {
	return false
}

/*
JSON encodes added keys and values into JSON.
*/
func (rd *ResponseData) JSON() string {
	return ""
}
