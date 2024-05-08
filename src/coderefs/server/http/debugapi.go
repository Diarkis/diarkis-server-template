package http

/*
ExposeDebugAPI exposes debug APIs

--------------------------------------------------

Get the list of available server addresses.

The addresses are public addresses.

	GET /debug/server/list/type/:serverType

serverType - TCP, UDP or custom server type.

--------------------------------------------------

Get the endpoint of the server by its public address.

The API allows you to choose which server to create a new user on.

	GET /debug/endpoint/type/:serverType/user/:uid/address/:addr

	POST /debug/endpoint/type/:serverType/user/:uid/address/:addr

serverType - TCP, UDP or custom server type.

uid        - User ID

addr       - Server address and port: format 0.0.0.0:8888

# POST Request

▶︎ Request Body format:

Request body will follow name=value format and separated by an & if there will be more than one name and value set.

	[IMPORTANT] Currently user data does NOT automatically URL decode request body data if URL encoded.

▶︎ Example With Form:

	foo=foo&bar=1234567&flag=true

The above will be translated to the following as user property data. Each value is stored as a string.

	foo  = "foo"
	bar  = "123456"
	flag = "true"

▶︎ Example With JSON (Content-Type:application/json header must be provided):

If Content-Type:application/json header is provided, The API will read the request body as JSON.

	{ "foo": "foo", "bar": 123456, "flag": true }

The above JSON request body will be translated to the following as user property data:

All values will be interface{} and they must be type asserted before using them.

	foo  = "foo"
	bar  = 123456
	flag = true

--------------------------------------------------
*/
func ExposeDebugAPI() {
}
