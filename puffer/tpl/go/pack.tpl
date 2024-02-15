---
  /* uint8 */
  $(PROPERTY_NAME_LOWER)Bytes := make([]byte, 1)
  $(PROPERTY_NAME_LOWER)Bytes[0] = proto.$(PROPERTY_NAME_UPPER)
  bytes = append(bytes, $(PROPERTY_NAME_LOWER)Bytes...)
---
  /* uint16 */
  $(PROPERTY_NAME_LOWER)Bytes := make([]byte, 2)
  binary.BigEndian.PutUint16($(PROPERTY_NAME_LOWER)Bytes, uint16(proto.$(PROPERTY_NAME_UPPER)))
  bytes = append(bytes, $(PROPERTY_NAME_LOWER)Bytes...)
---
  /* uint32 */
  $(PROPERTY_NAME_LOWER)Bytes := make([]byte, 4)
  binary.BigEndian.PutUint32($(PROPERTY_NAME_LOWER)Bytes, uint32(proto.$(PROPERTY_NAME_UPPER)))
  bytes = append(bytes, $(PROPERTY_NAME_LOWER)Bytes...)
---
  /* uint64 */
  $(PROPERTY_NAME_LOWER)Bytes := make([]byte, 8)
  binary.BigEndian.PutUint64($(PROPERTY_NAME_LOWER)Bytes, uint64(proto.$(PROPERTY_NAME_UPPER)))
  bytes = append(bytes, $(PROPERTY_NAME_LOWER)Bytes...)
---
  /* int8 */
  $(PROPERTY_NAME_LOWER)Bytes := make([]byte, 1)
  $(PROPERTY_NAME_LOWER)Bytes[0] = byte(proto.$(PROPERTY_NAME_UPPER))
  bytes = append(bytes, $(PROPERTY_NAME_LOWER)Bytes...)
---
  /* int16 */
  $(PROPERTY_NAME_LOWER)Bytes := make([]byte, 2)
  binary.BigEndian.PutUint16($(PROPERTY_NAME_LOWER)Bytes, uint16(proto.$(PROPERTY_NAME_UPPER)))
  bytes = append(bytes, $(PROPERTY_NAME_LOWER)Bytes...)
---
  /* int32 */
  $(PROPERTY_NAME_LOWER)Bytes := make([]byte, 4)
  binary.BigEndian.PutUint32($(PROPERTY_NAME_LOWER)Bytes, uint32(proto.$(PROPERTY_NAME_UPPER)))
  bytes = append(bytes, $(PROPERTY_NAME_LOWER)Bytes...)
---
  /* int64 */
  $(PROPERTY_NAME_LOWER)Bytes := make([]byte, 8)
  binary.BigEndian.PutUint64($(PROPERTY_NAME_LOWER)Bytes, uint64(proto.$(PROPERTY_NAME_UPPER)))
  bytes = append(bytes, $(PROPERTY_NAME_LOWER)Bytes...)
---
  /* float32 */
  $(PROPERTY_NAME_LOWER)Bytes := make([]byte, 4)
  $(PROPERTY_NAME_LOWER)Bits := math.Float32bits(proto.$(PROPERTY_NAME_UPPER))
  $(PROPERTY_NAME_LOWER)Bytes[0] = byte($(PROPERTY_NAME_LOWER)Bits >> 24)
  $(PROPERTY_NAME_LOWER)Bytes[1] = byte($(PROPERTY_NAME_LOWER)Bits >> 16)
  $(PROPERTY_NAME_LOWER)Bytes[2] = byte($(PROPERTY_NAME_LOWER)Bits >> 8)
  $(PROPERTY_NAME_LOWER)Bytes[3] = byte($(PROPERTY_NAME_LOWER)Bits)
  $(PROPERTY_NAME_LOWER)Bytes = util.ReverseBytes($(PROPERTY_NAME_LOWER)Bytes)
  bytes = append(bytes, $(PROPERTY_NAME_LOWER)Bytes...)
---
  /* float64 */
  $(PROPERTY_NAME_LOWER)Bytes := make([]byte, 8)
  $(PROPERTY_NAME_LOWER)Bits := math.Float64bits(proto.$(PROPERTY_NAME_UPPER))
  $(PROPERTY_NAME_LOWER)Bytes[0] = byte($(PROPERTY_NAME_LOWER)Bits >> 56)
  $(PROPERTY_NAME_LOWER)Bytes[1] = byte($(PROPERTY_NAME_LOWER)Bits >> 48)
  $(PROPERTY_NAME_LOWER)Bytes[2] = byte($(PROPERTY_NAME_LOWER)Bits >> 40)
  $(PROPERTY_NAME_LOWER)Bytes[3] = byte($(PROPERTY_NAME_LOWER)Bits >> 32)
  $(PROPERTY_NAME_LOWER)Bytes[4] = byte($(PROPERTY_NAME_LOWER)Bits >> 24)
  $(PROPERTY_NAME_LOWER)Bytes[5] = byte($(PROPERTY_NAME_LOWER)Bits >> 24)
  $(PROPERTY_NAME_LOWER)Bytes[6] = byte($(PROPERTY_NAME_LOWER)Bits >> 16)
  $(PROPERTY_NAME_LOWER)Bytes[7] = byte($(PROPERTY_NAME_LOWER)Bits >> 8)
  $(PROPERTY_NAME_LOWER)Bytes = util.ReverseBytes($(PROPERTY_NAME_LOWER)Bytes)
  bytes = append(bytes, $(PROPERTY_NAME_LOWER)Bytes...)
---
  /* $(CUSTOM_DATA_TYPE) */
  $(PROPERTY_NAME_LOWER)SizeBytes := make([]byte, 2)
  $(PROPERTY_NAME_LOWER)Packed := proto.$(PROPERTY_NAME_UPPER).Pack()
  binary.BigEndian.PutUint16($(PROPERTY_NAME_LOWER)SizeBytes, uint16(len($(PROPERTY_NAME_LOWER)Packed)))
  bytes = append(bytes, $(PROPERTY_NAME_LOWER)SizeBytes...)
  bytes = append(bytes, $(PROPERTY_NAME_LOWER)Packed...)
---
  /* bool */
  $(PROPERTY_NAME_LOWER)Bytes := make([]byte, 1)
  if proto.$(PROPERTY_NAME_UPPER) {
    $(PROPERTY_NAME_LOWER)Bytes[0] = uint8(1) /* true */
  } else {
    $(PROPERTY_NAME_LOWER)Bytes[0] = uint8(2) /* false */
  }
  bytes = append(bytes, $(PROPERTY_NAME_LOWER)Bytes...)
---
  /* string */
  $(PROPERTY_NAME_LOWER)SizeBytes := make([]byte, 2)
  binary.BigEndian.PutUint16($(PROPERTY_NAME_LOWER)SizeBytes, uint16(len(proto.$(PROPERTY_NAME_UPPER))))
  bytes = append(bytes, $(PROPERTY_NAME_LOWER)SizeBytes...)
  bytes = append(bytes, []byte(proto.$(PROPERTY_NAME_UPPER))...)
---
  /* []byte */
  $(PROPERTY_NAME_LOWER)SizeBytes := make([]byte, 2)
  binary.BigEndian.PutUint16($(PROPERTY_NAME_LOWER)SizeBytes, uint16(len(proto.$(PROPERTY_NAME_UPPER))))
  bytes = append(bytes, $(PROPERTY_NAME_LOWER)SizeBytes...)
  bytes = append(bytes, proto.$(PROPERTY_NAME_UPPER)...)
---
  /* []uint8 */
  $(PROPERTY_NAME_LOWER)SizeBytes := make([]byte, 2)
  binary.BigEndian.PutUint16($(PROPERTY_NAME_LOWER)SizeBytes, uint16(len(proto.$(PROPERTY_NAME_UPPER))))
  bytes = append(bytes, $(PROPERTY_NAME_LOWER)SizeBytes...)
  for i := 0; i < len(proto.$(PROPERTY_NAME_UPPER)); i++ {
    b := make([]byte, 1)
    b[0] = proto.$(PROPERTY_NAME_UPPER)[i]
    bytes = append(bytes, b...)
  }
---
  /* []uint16 */
  $(PROPERTY_NAME_LOWER)SizeBytes := make([]byte, 2)
  binary.BigEndian.PutUint16($(PROPERTY_NAME_LOWER)SizeBytes, uint16(len(proto.$(PROPERTY_NAME_UPPER))))
  bytes = append(bytes, $(PROPERTY_NAME_LOWER)SizeBytes...)
  for i := 0; i < len(proto.$(PROPERTY_NAME_UPPER)); i++ {
    b := make([]byte, 2)
    binary.BigEndian.PutUint16(b, uint16(proto.$(PROPERTY_NAME_UPPER)[i]))
    bytes = append(bytes, b...)
  }
---
  /* []uint32 */
  $(PROPERTY_NAME_LOWER)SizeBytes := make([]byte, 2)
  binary.BigEndian.PutUint16($(PROPERTY_NAME_LOWER)SizeBytes, uint16(len(proto.$(PROPERTY_NAME_UPPER))))
  bytes = append(bytes, $(PROPERTY_NAME_LOWER)SizeBytes...)
  for i := 0; i < len(proto.$(PROPERTY_NAME_UPPER)); i++ {
    b := make([]byte, 4)
    binary.BigEndian.PutUint32(b, uint32(proto.$(PROPERTY_NAME_UPPER)[i]))
    bytes = append(bytes, b...)
  }
---
  /* []uint64 */
  $(PROPERTY_NAME_LOWER)SizeBytes := make([]byte, 2)
  binary.BigEndian.PutUint16($(PROPERTY_NAME_LOWER)SizeBytes, uint16(len(proto.$(PROPERTY_NAME_UPPER))))
  bytes = append(bytes, $(PROPERTY_NAME_LOWER)SizeBytes...)
  for i := 0; i < len(proto.$(PROPERTY_NAME_UPPER)); i++ {
    b := make([]byte, 8)
    binary.BigEndian.PutUint64(b, uint64(proto.$(PROPERTY_NAME_UPPER)[i]))
    bytes = append(bytes, b...)
  }
---
  /* []int8 */
  $(PROPERTY_NAME_LOWER)SizeBytes := make([]byte, 2)
  binary.BigEndian.PutUint16($(PROPERTY_NAME_LOWER)SizeBytes, uint16(len(proto.$(PROPERTY_NAME_UPPER))))
  bytes = append(bytes, $(PROPERTY_NAME_LOWER)SizeBytes...)
  for i := 0; i < len(proto.$(PROPERTY_NAME_UPPER)); i++ {
    b := make([]byte, 1)
    b[0] = uint8(proto.$(PROPERTY_NAME_UPPER)[i])
    bytes = append(bytes, b...)
  }
---
  /* []int16 */
  $(PROPERTY_NAME_LOWER)SizeBytes := make([]byte, 2)
  binary.BigEndian.PutUint16($(PROPERTY_NAME_LOWER)SizeBytes, uint16(len(proto.$(PROPERTY_NAME_UPPER))))
  bytes = append(bytes, $(PROPERTY_NAME_LOWER)SizeBytes...)
  for i := 0; i < len(proto.$(PROPERTY_NAME_UPPER)); i++ {
    b := make([]byte, 2)
    binary.BigEndian.PutUint16(b, uint16(proto.$(PROPERTY_NAME_UPPER)[i]))
    bytes = append(bytes, b...)
  }
---
  /* []int32 */
  $(PROPERTY_NAME_LOWER)SizeBytes := make([]byte, 2)
  binary.BigEndian.PutUint16($(PROPERTY_NAME_LOWER)SizeBytes, uint16(len(proto.$(PROPERTY_NAME_UPPER))))
  bytes = append(bytes, $(PROPERTY_NAME_LOWER)SizeBytes...)
  for i := 0; i < len(proto.$(PROPERTY_NAME_UPPER)); i++ {
    b := make([]byte, 4)
    binary.BigEndian.PutUint32(b, uint32(proto.$(PROPERTY_NAME_UPPER)[i]))
    bytes = append(bytes, b...)
  }
---
  /* []int64 */
  $(PROPERTY_NAME_LOWER)SizeBytes := make([]byte, 2)
  binary.BigEndian.PutUint16($(PROPERTY_NAME_LOWER)SizeBytes, uint16(len(proto.$(PROPERTY_NAME_UPPER))))
  bytes = append(bytes, $(PROPERTY_NAME_LOWER)SizeBytes...)
  for i := 0; i < len(proto.$(PROPERTY_NAME_UPPER)); i++ {
    b := make([]byte, 8)
    binary.BigEndian.PutUint64(b, uint64(proto.$(PROPERTY_NAME_UPPER)[i]))
    bytes = append(bytes, b...)
  }
---
  /* []float32 */
  $(PROPERTY_NAME_LOWER)SizeBytes := make([]byte, 2)
  binary.BigEndian.PutUint16($(PROPERTY_NAME_LOWER)SizeBytes, uint16(len(proto.$(PROPERTY_NAME_UPPER))))
  bytes = append(bytes, $(PROPERTY_NAME_LOWER)SizeBytes...)
  for i := 0; i < len(proto.$(PROPERTY_NAME_UPPER)); i++ {
    b := make([]byte, 4)
    bits := math.Float32bits(proto.$(PROPERTY_NAME_UPPER)[i])
    b[0] = byte(bits >> 24)
    b[1] = byte(bits >> 16)
    b[2] = byte(bits >> 8)
    b[3] = byte(bits)
    b = util.ReverseBytes(b)
    bytes = append(bytes, b...)
  }
---
  /* []float64 */
  $(PROPERTY_NAME_LOWER)SizeBytes := make([]byte, 2)
  binary.BigEndian.PutUint16($(PROPERTY_NAME_LOWER)SizeBytes, uint16(len(proto.$(PROPERTY_NAME_UPPER))))
  bytes = append(bytes, $(PROPERTY_NAME_LOWER)SizeBytes...)
  for i := 0; i < len(proto.$(PROPERTY_NAME_UPPER)); i++ {
    b := make([]byte, 8)
    bits := math.Float64bits(proto.$(PROPERTY_NAME_UPPER)[i])
    b[0] = byte(bits >> 56)
    b[1] = byte(bits >> 48)
    b[2] = byte(bits >> 40)
    b[3] = byte(bits >> 32)
    b[4] = byte(bits >> 24)
    b[5] = byte(bits >> 16)
    b[6] = byte(bits >> 8)
    b[7] = byte(bits)
    b = util.ReverseBytes(b)
    bytes = append(bytes, b...)
  }
---
  /* []bool */
  $(PROPERTY_NAME_LOWER)SizeBytes := make([]byte, 2)
  binary.BigEndian.PutUint16($(PROPERTY_NAME_LOWER)SizeBytes, uint16(len(proto.$(PROPERTY_NAME_UPPER))))
  bytes = append(bytes, $(PROPERTY_NAME_LOWER)SizeBytes...)
  for i := 0; i < len(proto.$(PROPERTY_NAME_UPPER)); i++ {
    b := make([]byte, 1)
    if proto.$(PROPERTY_NAME_UPPER)[i] {
      b[0] = 0x01
    } else {
      b[0] = 0x02
    }
    bytes = append(bytes, b...)
  }
---
  /* []string */
  $(PROPERTY_NAME_LOWER)SizeBytes := make([]byte, 2)
  binary.BigEndian.PutUint16($(PROPERTY_NAME_LOWER)SizeBytes, uint16(len(proto.$(PROPERTY_NAME_UPPER))))
  bytes = append(bytes, $(PROPERTY_NAME_LOWER)SizeBytes...)
  for i := 0; i < len(proto.$(PROPERTY_NAME_UPPER)); i++ {
    b := make([]byte, 2)
    binary.BigEndian.PutUint16(b, uint16(len(proto.$(PROPERTY_NAME_UPPER)[i])))
    bytes = append(bytes, b...)
    bytes = append(bytes, []byte(proto.$(PROPERTY_NAME_UPPER)[i])...)
  }
---
  /* [][]byte */
  $(PROPERTY_NAME_LOWER)SizeBytes := make([]byte, 2)
  binary.BigEndian.PutUint16($(PROPERTY_NAME_LOWER)SizeBytes, uint16(len(proto.$(PROPERTY_NAME_UPPER))))
  bytes = append(bytes, $(PROPERTY_NAME_LOWER)SizeBytes...)
  for i := 0; i < len(proto.$(PROPERTY_NAME_UPPER)); i++ {
    b := make([]byte, 2)
    binary.BigEndian.PutUint16(b, uint16(len(proto.$(PROPERTY_NAME_UPPER)[i])))
    bytes = append(bytes, b...)
    bytes = append(bytes, proto.$(PROPERTY_NAME_UPPER)[i]...)
  }
---
  /* []$(CUSTOM_DATA_TYPE) */
  $(PROPERTY_NAME_LOWER)LengthBytes := make([]byte, 2)
  $(PROPERTY_NAME_LOWER)Length := len(proto.$(PROPERTY_NAME_UPPER))
  binary.BigEndian.PutUint16($(PROPERTY_NAME_LOWER)LengthBytes, uint16($(PROPERTY_NAME_LOWER)Length))
  bytes = append(bytes, $(PROPERTY_NAME_LOWER)LengthBytes...)
  for i := 0; i < $(PROPERTY_NAME_LOWER)Length; i++ {
    $(PROPERTY_NAME_LOWER)SizeBytes := make([]byte, 2)
    $(PROPERTY_NAME_LOWER)Packed := proto.$(PROPERTY_NAME_UPPER)[i].Pack()
    binary.BigEndian.PutUint16($(PROPERTY_NAME_LOWER)SizeBytes, uint16(len($(PROPERTY_NAME_LOWER)Packed)))
    bytes = append(bytes, $(PROPERTY_NAME_LOWER)SizeBytes...)
    bytes = append(bytes, $(PROPERTY_NAME_LOWER)Packed...)
  }
