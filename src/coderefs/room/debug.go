package room

/*
DataDump implies a payload format for the dump room data command
*/
type DataDump struct {
	ID         string                 `json:"RoomID"`
	MaxMembers int                    `json:"MaxMembers"`
	AllowEmpty bool                   `json:"AllowEmpty"`
	MemberIDs  map[string]string      `json:"MemberIDs"`
	Properties map[string]interface{} `json:"Properties"`
}
