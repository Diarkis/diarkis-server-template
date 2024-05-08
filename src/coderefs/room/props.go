package room

/*
GetRemoteProperty returns a property of the given room and property name. Returns nil if not found.

	[IMPORTANT] This function works even if the room is not on the same server.
	[IMPORTANT] The fetched property is of interface{} type. Make sure to use util.ToInt, util.ToBytes etc.
	            In order to read the data with the correct data type.
	[IMPORTANT] The callback may be invoked on the server where the room is not held.
	            Do not use functions that require the room to be on the same server.

	[NOTE] This function asynchronous.

Operation example:

	room.GetRemoteProperty(roomID, "counter", func(err error, property interface{}) {
		if err != nil {
			// handle error here...
		}
		// assuming "counter" is an int
		counter, ok := util.ToInt(property)

		if !ok {
			// incorrect data type or "counter" does not exist...
		}
	})

	roomID - Target room ID of the properties.
	name   - Property name to retrieve.
	cb     - A callback to be invoked when the property has been fetched.
*/
func GetRemoteProperty(roomID string, name string, cb func(err error, property interface{})) {
}

/*
GetRemoteProperties returns properties of a room

	[IMPORTANT] This function works even if the room is not on the same server.
	[IMPORTANT] This function does NOT update room's internal TTL.
	[IMPORTANT] The fetched property is of map[string]interface{} type. Make sure to use util.GetAsInt, util.GetAsBytes etc.
	[IMPORTANT] The callback may be invoked on the server where the room is not held.
	            Do not use functions that require the room to be on the same server.

	[NOTE] This function is asynchronous.

Example

	room.GetRemoteProperties(roomID func(err error, properties map[string]interface{}) {
		if err != nil {
			// handle error here...
		}
		// assuming "counter" is an int
		counter, ok := util.GetAsInt(properties, "counter")

		if !ok {
			// incorrect data type of the value does not exist
		}

		message, ok := util.GetAsBytes(properties, "message")

		if !ok {
			// incorrect data type of the value does not exist
		}

		flag, ok    := util.GetAsBool(properties, "bool")

		if !ok {
			// incorrect data type of the value does not exist
		}
	})

Parameters

	roomID - Target room ID of the properties.
	cb     - A callback to be invoked when the properties have been fetched.
*/
func GetRemoteProperties(roomID string, cb func(err error, properties map[string]interface{})) {
}

/*
GetProperty returns a property of the given room and property name. Returns nil if not found.

	[IMPORTANT] This function does NOT work if the room is not on the same server.
	[IMPORTANT] This function does NOT update room's internal TTL.

	[IMPORTANT] The fetched property is of interface{} type. Make sure to use util.ToInt, util.ToBytes etc.

	[NOTE] Uses mutex lock internally.

Example

	property := room.GetProperty(roomID, "counter")
	// assuming "counter" is an int
	counter, ok := util.ToInt(property)

	if !ok {
		// incorrect data type or the value does not exist
	}

Parameters

	roomID - Target room ID of the properties.
	name   - Property name to retrieve.
*/
func GetProperty(roomID string, name string) interface{} {
	return nil
}

/*
GetProperties returns properties of a room

	[IMPORTANT] This function does NOT work if the room is not on the same server.
	[IMPORTANT] This function does NOT update room's internal TTL.

	[IMPORTANT] The fetched property is of map[string]interface{} type. Make sure to use util.GetAsInt, util.GetAsBytes etc.

	[NOTE] Uses mutex lock internally.

Example

	properties := room.GetProperties(roomID)
	counter, ok := util.GetAsInt(properties, "counter")

	if !ok {
		// incorrect data type or the value does not exist
	}

	message, ok := util.GetAsBytes(properties, "message")

	if !ok {
		// incorrect data type or the value does not exist
	}

	flag, ok := util.GetAsBool(properties, "flag")

	if !ok {
		// incorrect data type or the value does not exist
	}

Parameters

	roomID - Target room ID of the properties.
*/
func GetProperties(roomID string) map[string]interface{} {
	return nil
}

/*
UpdateProperties updates or creates properties of a room and returns the updated or created properties.

	[IMPORTANT] This function does NOT work if the room is not on the same server.
	[IMPORTANT] This function DOES update room's internal TTL.
	[IMPORTANT] Room properties do NOT support structs.
	[IMPORTANT] Numeric properties will be float64. Use util.ToInt, ToUint etc. To convert to proper data type.
	[IMPORTANT] Array properties will be []interface{}. Must apply type assertion per element.
	[IMPORTANT] Byte array properties will be converted to base64 encoded string. Use util.ToBytes() to convert it back.
	[IMPORTANT] Operation function is protected by mutex lock, using another mutex lock within the function may cause a deadlock.

	[NOTE] Uses mutex lock internally.

Operation example:

	var updateErr error
	_ := room.UpdateProperties(roomID, func(properties map[string]interface{}) {
		capsule := datacapsule.NewCapsule()
		err := capsule.Import(properties["roomPropertyCapsule"].(map[string]interface{}))
		if err != nil {
			// this is how to propagate an error
			updateErr = err
			return
		}
		counter := capsule.GetAsInt("counter")
		counter++
		capsule.SetAsInt("counter", counter)
		properties["roomPropertyCapsule"] = capsule.Export()
	})
	if updateErr != nil {
		// handle update error here
	}

Use datacapsule.Capsule as property:

	capsule := datacapsule.NewCapsule()
	capsule.SetAsInt8("memberNum", memberNum)
	properties["memberNum"] = capsule.Export() // setting capsule as a property

Parameters

	roomID    - Room ID to update its properties.
	operation - Callback function to perform update operations and returns updated properties as a map.
	            Must return true, when property (properties) is updated and false when there is no update.
*/
func UpdateProperties(roomID string, operation func(properties map[string]interface{}) bool) map[string]interface{} {
	return nil
}

/*
SetProperty assigns a value to the target room associated with the given name. If the property already exists, it will overwrite it.

	[IMPORTANT] This function does NOT work if the room is not on the same server.
	[IMPORTANT] This function does NOT update room's internal TTL.
	[IMPORTANT] Room properties do NOT support structs.
	[IMPORTANT] Numeric properties will be float64. Use util.ToInt, ToUint etc. To convert to proper data type.
	[IMPORTANT] Array properties will be []interface{}. Must apply type assertion per element.
	[IMPORTANT] Byte array properties will be converted to base64 encoded string. Use util.ToBytes() to convert it back.

	[NOTE] If there is a value associated with the name, the value will be replace.
	[NOTE] Uses mutex lock internally.

Parameters

	roomID - Target room ID
	name   - Property name to set the value as
	value  - Property value to be associated with the name
*/
func SetProperty(roomID string, name string, value interface{}) bool {
	return false
}

/*
IncrProperty increments int64 property by given delta and returns the incremented value.

	[IMPORTANT] This function does NOT work if the room is not on the same server.
	[IMPORTANT] This function DOES update room's internal TTL.

If the property does not exist, it creates the property.

	[NOTE] Uses mutex lock internally.

Cases for return value false

	┌──────────────────┬──────────────────────────────────────────────────────┐
	│ Error            │ Reason                                               │
	╞══════════════════╪══════════════════════════════════════════════════════╡
	│ Room not found   │ Either invalid room ID given or room does not exist. │
	│ Property Corrupt │ Target property data type is not int64.              │
	└──────────────────┴──────────────────────────────────────────────────────┘

Parameters

	roomID   - Target room ID.
	propName - Target property name.
	delta    - Delta value to increment by.

Usage Example

	updatedHP, updated := room.IncrProperty(roomID, "HP", damage)
*/
func IncrProperty(roomID string, propName string, delta int64) (int64, bool) {
	return 0, false
}
