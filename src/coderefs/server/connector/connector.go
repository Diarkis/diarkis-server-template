package connector

/*
SetPublicEndPoint assigns a public endpoint for the connector
*/
func SetPublicEndPoint(addr string) {
}

/*
Setup [INTERNAL USE ONLY] is used to expose connector REST APIs for external server process
to communicate with Diarkis process in order to use Diarkis' auto-scaling and other features
*/
func Setup(path string) int64 {
	return 0
}

/*
Initialize sets up find Diarkis server process endpoint
*/
func Initialize() string {
	return ""
}

/*
GetOpenPort get an open port for client server to use from Diarkis server
*/
func GetOpenPort() (int, error) {
	return 0, nil
}

/*
GetAddress get the public endpoint address from Diarkis server for client server to use
*/
func GetAddress() string {
	return ""
}

/*
Ready informs Diarkis server that client server is ready to receive clients
*/
func Ready() int {
	return 0
}

/*
Health informs Diarkis server that client server is health
*/
func Health(ccu int) int {
	return 0
}

/*
Allocate informs Diarkis server that client server no longer able to receive clients
*/
func Allocate() int {
	return 0
}

/*
Shutdown informs Diarkis server that client server is shutting down
This instructs Diarkis server to also shutdown
*/
func Shutdown() int {
	return 0
}

/*
SendToDiarkis sends a communication message from client server to Diarkis server
*/
func SendToDiarkis(endpoint string, mode string, value int) (int, []byte) {
	return 0, nil
}
