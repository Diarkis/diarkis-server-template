package net

const Localhost string = "127.0.0.1"

/*
FindInterfaces Searches for network interfaces and their addresses
*/
func FindInterfaces() {
}

/*
CreateEndPoint Creates server end point if port is not available, it automatically searches an open port
*/
func CreateEndPoint(configs map[string]interface{}, defaultPort string) (string, string) {
	return "", ""
}

/*
CreateFixedEndPoint creates an endpoint with the given address and port
*/
func CreateFixedEndPoint(configs map[string]interface{}, defaultPort string) (string, string) {
	return "", ""
}

/*
GetAddrByInterfaceName Returns an address to bind by network interface name
*/
func GetAddrByInterfaceName(name string) string {
	return ""
}

/*
FindAvailablePort Searches for available port
*/
func FindAvailablePort(addr string, port string) string {
	return ""
}

/*
FindTakenPort searches for a port that is already taken:

This is meant to be used to search for internal communication endpoint
*/
func FindTakenPort(addr string, port string) string {
	return ""
}
