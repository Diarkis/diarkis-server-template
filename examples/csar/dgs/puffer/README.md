Puffer allows you to define and generate data protocols for easy byte-serialization.

## Usage

### command

`make gen` - generate packet definition files
`make clean` - remove generated files

## Code Example

```
func Expose() { // add command using puffer handler
    ...
	// puffer version sample
	diarkisexec.SetServerCommandHandler(custom.EchoVer, custom.EchoCmd, echoPufferCmd) // sample puffer echo command
    ...
}

func echoPufferCmd(ver uint8, cmd uint16, payload []byte, userData *user.User, next func(error)) {
	logger.Debug("Hello puffer command has received %#v from the client SID:%s - UID:%s", payload, userData.SID, userData.ID)
	// unpack []byte to struct
	echoData := custom.NewEcho()
	err := echoData.Unpack(payload) // You can unpack []byte to go struct

	if err != nil {
		logger.Error("Failed to unpack echo data: %v", err)
		userData.ServerRespond(nil, ver, cmd, server.Err, true)
		next(nil)
		return
	}

	logger.Debug("Unpacked echo data: %#v", echoData)
    userData.ServerRespond(echoData.Pack(), ver, cmd, server.Ok, true) // You can get []byte by using Pack. ( echoData.Pack equals payload in this example.)
	// move on to the next command handler if there is any
	next(nil)
}

```

## Data Types Used in JSON definition files

Table of property data types.

                                  Corresponding Data Type

| JSON Data Type Dictation |        Go |       C# |                                             C++ |
| -----------------------: | --------: | -------: | ----------------------------------------------: |
|                       u8 |     uint8 |     byte |                                         uint8_t |
|                      u16 |    uint16 |   ushort |                                        uint16_t |
|                      u32 |    uint32 |     uint |                                        uint32_t |
|                      u64 |    uint64 |    ulong |                                        uint64_t |
|                       i8 |      int8 |    sbyte |                                          int8_t |
|                      i16 |     int16 |    short |                                         int16_t |
|                      i32 |     int32 |      int |                                         int32_t |
|                      i64 |     int64 |     long |                                         int64_t |
|                      f32 |   float32 |    float |                                           float |
|                      f64 |   float64 |   double |                                          double |
|                     bool |   boolean |     bool |                                            bool |
|                   string |    string |   string |                              Diarkis::StdString |
|                      b[] |    []byte |   byte[] |                     Diarkis::StdVector<uint8_t> |
|                     u8[] |   []uint8 |   byte[] |                     Diarkis::StdVector<uint8_t> |
|                    u16[] |  []uint16 | ushort[] |                    Diarkis::StdVector<uint16_t> |
|                    u32[] |  []uint32 |   uint[] |                    Diarkis::StdVector<uint32_t> |
|                    u64[] |  []uint64 |  ulong[] |                    Diarkis::StdVector<uint64_t> |
|                     i8[] |    []int8 |  sbyte[] |                      Diarkis::StdVector<int8_t> |
|                    i16[] |   []int16 |  short[] |                     Diarkis::StdVector<int16_t> |
|                    i32[] |   []int32 |    int[] |                     Diarkis::StdVector<int32_t> |
|                    i64[] |   []int64 |   long[] |                     Diarkis::StdVector<int64_t> |
|                    f32[] | []float32 |  float[] |                       Diarkis::StdVector<float> |
|                    f64[] | []float64 | double[] |                      Diarkis::StdVector<double> |
|                   bool[] |    []bool |   bool[] |                        Diarkis::StdVector<bool> |
|                 string[] |  []string | string[] |          Diarkis::StdVector<Diarkis::StdString> |
|                    b[][] |  [][]byte | byte[][] | Diarkis::StdVector<Diarkis::StdVector<uint8_t>> |

## Custom Data Types

You may use protocols that you define in JSON files as custom data types.

Example:

```
{
  "States": {
    "Flag":      "bool",
    "mode":      "u8",
    "timestamp": "i64"
  },
  "UserData": {
    "ID": "string",
    "States": "States"
  }
}
```

## directory architecture

`make gen` will generate the following directories.
go - go code
cs - c# code
cpp - C++ code

The code is generated for each runtime, and can be used on the client for cs and cpp, and can be used on the client and server in common for go.
`json_definitions` contains puffer definitions.
`json_definitions/samples` is a sample definition that is not used in the code right now.
We recommend that you delete it when actually using it.

## Note

The current version of puffer has a problem with overwriting the same command name even if the package name is different.
This will be fixed in the next version.
