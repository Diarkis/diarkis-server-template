package field

/*
 []propagationItem

*/

/*
Pack encodes DispatchData struct to a byte array to be delivered over the command.
*/
func (proto *entity) Pack() []byte {
	return nil
}

/*
Unpack decodes the command payload byte array to DispatchData struct.
*/
func (proto *entity) Unpack(bytes []byte) error {
	return nil
}
func (proto *entityList) Pack() []byte {
	return nil
}
func (proto *entityList) Unpack(bytes []byte) error {
	return nil
}

/*
GetPackedSize returns the size of the returned value of Pack()
*/
func (proto *dispatchData) GetPackedSize() int {
	return 0
}
func (proto *dispatchData) Pack(bytes []byte) []byte {
	return nil
}
func (proto *dispatchData) Unpack(bytes []byte) error {
	return nil
}
func (proto *propagationData) Pack() []byte {
	return nil
}
func (proto *propagationData) Unpack(bytes []byte) error {
	return nil
}
func (proto *propagationItem) Pack() []byte {
	return nil
}
func (proto *propagationItem) Unpack(bytes []byte) error {
	return nil
}
func (proto *propagationItemList) GetPackedSize() int {
	return 0
}
func (proto *propagationItemList) Pack(bytes []byte) []byte {
	return nil
}
func (proto *propagationItemList) Unpack(bytes []byte) error {
	return nil
}
