package diarkisexec

/*
SetMeshCommandHandler assigns a callback to the given mesh command handler to the given cmd.

	[IMPORTANT] This function must be used BEFORE invoking diarkisexec.StartDiarkis()
	[IMPORTANT] You may NOT assign a callback that is used by built-in internal mesh command handlers.
	            Diarkis internally uses cmd from 0 to 10000.
	[IMPORTANT] You may NOT assign multiple callbacks to the same cmd.
*/
func SetMeshCommandHandler(cmd uint16, handler MeshCommandHandler) {
}

/*
SendMeshCommand sends a mesh command of the given cmd to the given internal mesh address.

The server of the given address will invoke the command handler of the given cmd.

	[IMPORTANT] This function must be used AFTER invoking diarkisexec.StartDiarkis()
	[IMPORTANT] The command handler must be defined using diarkisexec.SetMeshCommandHandler(cmd uint16, handler MeshCommandHandler).

Parameters

	cmd       - Mesh command ID.
	addresses - Internal mesh address of the server to send the command to.
	data      - Data to be sent to the handler.
	            The valid data type of data must be the following two types only:
	            1. map[string]interface{}
	            2. struct

The diagram below show how command is sent to multiple servers when you give multiple addresses:

	                                                           ┌──────────┐
	                                                     ┌────▶︎│ Server D │
	                                                     │     └──────────┘
	                               ┌──────────┐ <2> Send │
	                           ┌──▶︎│ Server B │──────────┤     ┌──────────┐
	                           │   └──────────┘          └────▶︎│ Server E │
	┌──────────┐ <1> Send      │                               └──────────┘
	│ Server A │ ──────────────┤                               ┌──────────┐
	└──────────┘               │   ┌──────────┐          ┌────▶︎│ Server F │
	                           └──▶︎│ Server C │──────────┤     └──────────┘
	                               └──────────┘          │
	                                                     │     ┌──────────┐
	                                                     └────▶︎│ Server G │
	                                                           └──────────┘
*/
func SendMeshCommand(cmd uint16, addresses []string, data interface{}, reliable bool) error {
	return nil
}

/*
SendMeshRequest sends a mesh command of the given cmd to the given internal mesh address and expects a response back from the server.

The server of the given address will invoke the command handler of the given cmd and sends back a response.

	[IMPORTANT] This function must be used AFTER invoking diarkisexec.StartDiarkis()
	[IMPORTANT] The command handler must be defined using diarkisexec.SetMeshCommandHandler(cmd uint16, handler MeshCommandHandler).

Error Cases

	┌───────────────────┬─────────────────────────────────────────────────────────────────────────┐
	│ Error             │ Reason                                                                  │
	╞═══════════════════╪═════════════════════════════════════════════════════════════════════════╡
	│ Invalid data type │ Input data type must be either a struct or map[string]interface{}.      │
	│ Handler error     │ Handler function of the request returned an error.                      │
	│ Network error     │ Mesh network error. Failed to send or receive server-to-server message. │
	╘═══════════════════╧═════════════════════════════════════════════════════════════════════════╛

The diagram below show how SendMeshRequest works:

	┌──────────┐ <1> Send request  ┌──────────┐
	│ Server A │ ─────────────────▶︎│ Server B │
	│          │ ◀︎──────────────── │          │
	└──────────┘ <2> Send Response └──────────┘

Parameters

	cmd       - Mesh command ID.
	addresses - Internal mesh address of the server to send the command to.
	data      - Data to be sent to the handler.
	            The valid data type of data must be the following two types only:
	            1. map[string]interface{}
	            2. struct
	callback  - The callback to be invoked when the response is received from the remote server.
*/
func SendMeshRequest(cmd uint16, address string, data interface{}, callback func(err error, response map[string]interface{})) {
}

/*
CreateReturnBytes converts the given map[string]interface{} into a byte error for MeshCommandHandler return value.
*/
func CreateReturnBytes(data map[string]interface{}) ([]byte, error) {
	return nil, nil
}

/*
GetServerType returns the server type of itself.

	[IMPORTANT] This function must be used AFTER invoking diarkisexec.StartDiarkis()

	[NOTE] Uses mutex lock internally.
	[NOTE] The default server types: HTTP, UDP, and TCP.
	[NOTE] Server type can be customized by using DIARKIS_SERVER_TYPE=$(server_type) env.
*/
func GetServerType() string {
	return ""
}

/*
GetServerTypeByMeshAddress returns the server type of a remote server by its internal mesh address.

	[IMPORTANT] This function must be used AFTER invoking diarkisexec.StartDiarkis()

	[NOTE] Uses mutex lock internally.
	[NOTE] The default server types: HTTP, UDP, and TCP.
	[NOTE] Server type can be customized by using DIARKIS_SERVER_TYPE=$(server_type) env.
*/
func GetServerTypeByMeshAddress(meshAddress string) string {
	return ""
}

/*
GetServerRole returns the server role of itself.

	[IMPORTANT] This function must be used AFTER invoking diarkisexec.StartDiarkis()

	[NOTE] Uses mutex lock internally.
	[NOTE] Server roles are automatically determined by server's network protocol.
	       The valid server roles are: HTTP, UDP, and TCP.
*/
func GetServerRole() string {
	return ""
}

/*
IsServerRoleHTTP returns true if the server role is HTTP.

	[IMPORTANT] This function does not work correctly before Diarkis server is ready.

	[NOTE]      Server role is decided according to the network protocol the server uses and it does not change unlike server type.
*/
func IsServerRoleHTTP() bool {
	return false
}

/*
IsServerRoleUDP returns true if the server role is UDP.

	[IMPORTANT] This function does not work correctly before Diarkis server is ready.

	[NOTE]      Server role is decided according to the network protocol the server uses and it does not change unlike server type.
*/
func IsServerRoleUDP() bool {
	return false
}

/*
IsServerRoleTCP returns true if the server role is UDP.

	[IMPORTANT] This function does not work correctly before Diarkis server is ready.

	[NOTE]      Server role is decided according to the network protocol the server uses and it does not change unlike server type.
*/
func IsServerRoleTCP() bool {
	return false
}

/*
GetServerRoleByMeshAddress returns the server role of a remote server by its internal mesh address.

	[IMPORTANT] This function must be used AFTER invoking diarkisexec.StartDiarkis()

	[NOTE] Uses mutex lock internally.
	[NOTE] Server roles are automatically determined by server's network protocol.
	       The valid server roles are: HTTP, UDP, and TCP.
*/
func GetServerRoleByMeshAddress(meshAddress string) string {
	return ""
}

/*
GetPublicEndpoint returns the server's public endpoint.

	[IMPORTANT] This function must be used AFTER invoking diarkisexec.StartDiarkis()

	[NOTE] Uses mutex lock internally.
	[NOTE] This function is available for UDP and TCP servers only.
*/
func GetPublicEndpoint() string {
	return ""
}

/*
GetPublicEndpointByMeshAddress returns the remote server's public endpoint by the remote server's internal mesh address.

	[IMPORTANT] This function must be used AFTER invoking diarkisexec.StartDiarkis()

	[NOTE] Uses mutex lock internally.
	[NOTE] This function is available for UDP and TCP servers only.
*/
func GetPublicEndpointByMeshAddress(meshAddress string) string {
	return ""
}

/*
GetMeshAddress returns the server's internal mesh address.

	[IMPORTANT] This function must be used AFTER invoking diarkisexec.StartDiarkis()

	[NOTE] Uses mutex lock internally.
*/
func GetMeshAddress() string {
	return ""
}

/*
GetMeshAddressByServerType returns a randomly selected internal mesh address of the given server type.

	[IMPORTANT] This function must be used AFTER invoking diarkisexec.StartDiarkis()

	[NOTE] Uses mutex lock internally.
	[NOTE] The default server types: HTTP, UDP, and TCP.
	[NOTE] Server type can be customized by using DIARKIS_SERVER_TYPE=$(server_type) env.
*/
func GetMeshAddressByServerType(serverType string) string {
	return ""
}

/*
GetMeshAddressesByServerType returns an array of internal mesh addresses of the given server type.

	[IMPORTANT] This function must be used AFTER invoking diarkisexec.StartDiarkis()

	[NOTE] Uses mutex lock internally.
	[NOTE] The default server types: HTTP, UDP, and TCP.
	[NOTE] Server type can be customized by using DIARKIS_SERVER_TYPE=$(server_type) env.
*/
func GetMeshAddressesByServerType(serverType string) []string {
	return nil
}

/*
GetMeshAddressByServerRole returns a randomly selected internal mesh address of the given server role.

	[IMPORTANT] This function must be used AFTER invoking diarkisexec.StartDiarkis()

	[NOTE] Uses mutex lock internally.
	[NOTE] Server roles are automatically determined by server's network protocol.
	       The valid server roles are: HTTP, UDP, and TCP.
*/
func GetMeshAddressByServerRole(serverRole string) string {
	return ""
}

/*
GetMeshAddressesByServerRole returns an array of internal mesh addresses of the given server role.

	[IMPORTANT] This function must be used AFTER invoking diarkisexec.StartDiarkis()

	[NOTE] Uses mutex lock internally.
	[NOTE] Server roles are automatically determined by server's network protocol.
	       The valid server roles are: HTTP, UDP, and TCP.
*/
func GetMeshAddressesByServerRole(serverRole string) []string {
	return nil
}

/*
IsOffline returns true if the server is marked to be shutdown.

	[IMPORTANT] This function must be used AFTER invoking diarkisexec.StartDiarkis()

	[NOTE] Uses mutex lock internally.
*/
func IsOffline() bool {
	return false
}

/*
IsOfflineByMeshAddress returns true if the remote server of the given mesh address is marked to be shutdown.

	[IMPORTANT] This function must be used AFTER invoking diarkisexec.StartDiarkis()

	[NOTE] Uses mutex lock internally.
*/
func IsOfflineByMeshAddress(meshAddress string) bool {
	return false
}

/*
IsTaken returns true if the server is in taken state.

	[IMPORTANT] This function must be used AFTER invoking diarkisexec.StartDiarkis()

	[NOTE] Uses mutex lock internally.
*/
func IsTaken() bool {
	return false
}

/*
IsTakenByMeshAddress returns true if the remote server of the given mesh address is in taken state.

	[IMPORTANT] This function must be used AFTER invoking diarkisexec.StartDiarkis()

	[NOTE] Uses mutex lock internally.
*/
func IsTakenByMeshAddress(meshAddress string) bool {
	return false
}

/*
SetSharedData updates a shared data to be propagated to all server nodes in the cluster.

The propagation of the shared data may take some time.

	[IMPORTANT] The number of shared data keys you may store is limited to 10 keys.
	[IMPORTANT] Updated value may suffer from race condition.
	            If multiple server nodes attempt to update the same key,
	            The value of the key maybe overwritten.
*/
func SetSharedData(key string, value int16) bool {
	return false
}

/*
RemoveSharedData removes the given shared key
and propagates the removal to all server nodes in the cluster.

The propagation of the shared data may take some time.

	[IMPORTANT] Updated value may suffer from race condition.
	            If multiple server nodes attempt to update the same key,
	            The value of the key maybe overwritten.
*/
func RemoveSharedData(key string) bool {
	return false
}

/*
GetSharedData returns synchronized shared data by its key.
*/
func GetSharedData(key string) (int16, bool) {
	return 0, false
}

/*
SetOnSharedDataUpdate assigns a callback to be invoked when a shared data is updated.

	[NOTE] To remove the assigned callback use RemoveOnSharedDataUpdate(cb func(key string, value int16)).
*/
func SetOnSharedDataUpdate(cb func(key string, value int16)) {
}

/*
SetOnSharedDataRemove assigns a callback to be invoked when a shared data is deleted.

	[NOTE] To remove the assigned callback, use RemoveOnSharedDataRemove(cb func(key string)).
*/
func SetOnSharedDataRemove(cb func(key string)) {
}

/*
RemoveOnSharedDataUpdate removes the given callback that has been assigned by SetOnSharedDataUpdate.

	[IMPORTANT] This function uses mutex lock internally.
*/
func RemoveOnSharedDataUpdate(cb func(key string, value int16)) {
}

/*
RemoveOnSharedDataRemove removes the given callback that has been assigned by SetOnSharedDataRemove.

	[IMPORTANT] This function uses mutex lock internally.
*/
func RemoveOnSharedDataRemove(cb func(key string)) {
}
