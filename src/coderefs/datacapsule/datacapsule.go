package datacapsule

/*
Capsule represents data map of interface{} and it enforces its properties' data type despite the properties being interface{} preventing
unexpected errors and behaviors because of the interface{} properties.

	For example Capsule can be used as a property value of Room. When you set Capsule as Room property,
	use capsule.Export() to convert it into map[string]interface{}
	And when you get Capsule from Room property, use capsule.Import(roomPropertyCapsuleMap).
*/
type Capsule struct{}

/*
NewCapsule returns a new empty Capsule instance.
*/
func NewCapsule() *Capsule {
	return nil
}

/*
Pack converts Capsule into a byte array to be transported etc.

Example of using datacapsule to create a return value for mesh HandleCommand and SendRequest:

	// A mesh command handler with a callback for mesh.SendRequest
	mesh.HandleCommand(someCmdID, func(requestData map[string]interface{}) ([]byte, error) {
	  cp := datacapsule.NewCapsule()
	  cp.SetAsString("name", "Bob")
	  cp.SetAsBytes("byteMessate", []byte("Hello"))

	  resp := make(map[string]interface{})

	  packed, err := cp.Pack()
	  if err != nil {
	    // handle error here...
	  }

	  resp["packed"] = packed

	  return mesh.CreateReturnBytes(resp)
	})
*/
func Pack(c *Capsule) ([]byte, error) {
	return nil, nil
}

/*
Unpack converts packed Capsule back to Capsule.

Example of datacapsule usage with mesh HandleCommand and SendRequest

	// The callback that handles data returned by mesh.HandleCommand
	mesh.SendRequest(someCmdID, targetNodeAddress, func(err error, resp map[string]interface{}) {
	  packed := mesh.GetBytes(resp, "packed")
	  cp, err := datacapsule.Unpack(packed)

	  if err != nil {
	    // handle error...
	  }

	  name, err := cp.GetAsString("name")

	  if err != nil {
	    // handle error...
	  }

	  byteMessage, err := cp.GetAsBytes("byteMessage")

	  if err != nil {
	    // handle error...
	  }

	})
*/
func Unpack(src []byte) (*Capsule, error) {
	return nil, nil
}

/*
Export returns its internal map data.
*/
func (c *Capsule) Export() map[string]interface{} {
	return nil
}

/*
Import replaces its internal map data with imported data map.
*/
func (c *Capsule) Import(src map[string]interface{}) error {
	return nil
}

/*
SetAsInt8 sets an int8 value to the internal map with name.
*/
func (c *Capsule) SetAsInt8(name string, val int8) {
}

/*
SetAsUint8 sets a uint8 value to the internal map with name.
*/
func (c *Capsule) SetAsUint8(name string, val uint8) {
}

/*
SetAsInt16 sets an int16 value to the internal map with name.
*/
func (c *Capsule) SetAsInt16(name string, val int16) {
}

/*
SetAsUint16 sets a uint16 value to the internal map with name.
*/
func (c *Capsule) SetAsUint16(name string, val uint16) {
}

/*
SetAsInt32 sets an int32 value to the internal map with name.
*/
func (c *Capsule) SetAsInt32(name string, val int32) {
}

/*
SetAsUint32 sets a uint32 value to the internal map with name.
*/
func (c *Capsule) SetAsUint32(name string, val uint32) {
}

/*
SetAsInt64 sets an int64 value to the internal map with name.
*/
func (c *Capsule) SetAsInt64(name string, val int64) {
}

/*
SetAsUint64 sets a uint64 value to the internal map with name.
*/
func (c *Capsule) SetAsUint64(name string, val uint64) {
}

/*
SetAsFloat64 sets a float64 value to the internal map with name.
*/
func (c *Capsule) SetAsFloat64(name string, val float64) {
}

/*
SetAsBool sets a bool value to the internal map with name.
*/
func (c *Capsule) SetAsBool(name string, val bool) {
}

/*
SetAsString sets a string value to the internal map with name.
*/
func (c *Capsule) SetAsString(name string, val string) {
}

/*
SetAsCapsule sets a Capsule value to the internal map with name.
*/
func (c *Capsule) SetAsCapsule(name string, val *Capsule) {
}

/*
SetAsArray sets an array of Capsules to the internal map with name.
*/
func (c *Capsule) SetAsArray(name string, val []*Capsule) {
}

/*
SetAsMap sets a map of Capsules with string keys to the internal map with name.
*/
func (c *Capsule) SetAsMap(name string, val map[string]*Capsule) {
}

/*
SetAsBytes sets a byte array to the internal map with name.
*/
func (c *Capsule) SetAsBytes(name string, val []byte) {
}

/*
GetAsInt8 returns the value of the name given.
*/
func (c *Capsule) GetAsInt8(name string) (int8, error) {
	return 0, nil
}

/*
GetAsUint8 returns the value of the name given.
*/
func (c *Capsule) GetAsUint8(name string) (uint8, error) {
	return 0, nil
}

/*
GetAsInt16 returns the value of the name given.
*/
func (c *Capsule) GetAsInt16(name string) (int16, error) {
	return 0, nil
}

/*
GetAsUint16 returns the value of the name given.
*/
func (c *Capsule) GetAsUint16(name string) (uint16, error) {
	return 0, nil
}

/*
GetAsInt32 returns the value of the name given.
*/
func (c *Capsule) GetAsInt32(name string) (int32, error) {
	return 0, nil
}

/*
GetAsUint32 returns the value of the name given.
*/
func (c *Capsule) GetAsUint32(name string) (uint32, error) {
	return 0, nil
}

/*
GetAsInt64 returns the value of the name given.
*/
func (c *Capsule) GetAsInt64(name string) (int64, error) {
	return 0, nil
}

/*
GetAsUint64 returns the value of the name given.
*/
func (c *Capsule) GetAsUint64(name string) (uint64, error) {
	return 0, nil
}

/*
GetAsFloat64 returns the value of the name given.
*/
func (c *Capsule) GetAsFloat64(name string) (float64, error) {
	return 0, nil
}

/*
GetAsBool returns the value of the name given.
*/
func (c *Capsule) GetAsBool(name string) (bool, error) {
	return false, nil
}

/*
GetAsString returns the value of the name given.
*/
func (c *Capsule) GetAsString(name string) (string, error) {
	return "", nil
}

/*
GetAsCapsule returns the value of the name given.
*/
func (c *Capsule) GetAsCapsule(name string) (*Capsule, error) {
	return nil, nil
}

/*
GetAsArray returns the value of the name given.
*/
func (c *Capsule) GetAsArray(name string) ([]*Capsule, error) {
	return nil, nil
}

/*
GetAsMap returns the value of the name given.
*/
func (c *Capsule) GetAsMap(name string) (map[string]*Capsule, error) {
	return nil, nil
}

/*
GetAsBytes returns the value of the name given.
*/
func (c *Capsule) GetAsBytes(name string) ([]byte, error) {
	return nil, nil
}
