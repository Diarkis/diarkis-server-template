package http

/*
OnEndPointAPIResponse assigns a callback on /endpoint/type/:type/user/:user API to add extra data to the response JSON.
*/
func OnEndPointAPIResponse(cb func(string, string, string, string) (string, interface{})) {
}

/*
ForceEndPointAPIV1Response forces the response format of /endpoint/type/:type/user/:user to be the same as when the request contains "ResponseVersion:v1" header.

This is only for the purpose of backward compatibility support.
*/
func ForceEndPointAPIV1Response() {
}
