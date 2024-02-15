---
  /* uint8 */
  proto.$(PROPERTY_NAME_UPPER) = uint8(bytes[offset])
  offset++
---
  /* int8 */
  proto.$(PROPERTY_NAME_UPPER) = int8(bytes[offset])
  offset++
---
  /* uint16 */
  proto.$(PROPERTY_NAME_UPPER) = binary.BigEndian.Uint16(bytes[offset:offset + 2])
  offset += 2
---
  /* int16 */
  proto.$(PROPERTY_NAME_UPPER) = int16(binary.BigEndian.Uint16(bytes[offset:offset + 2]))
  offset += 2
---
  /* uint32 */
  proto.$(PROPERTY_NAME_UPPER) = binary.BigEndian.Uint32(bytes[offset:offset + 4])
  offset += 4
---
  /* int32 */
  proto.$(PROPERTY_NAME_UPPER) = int32(binary.BigEndian.Uint32(bytes[offset:offset + 4]))
  offset += 4
---
  /* uint64 */
  proto.$(PROPERTY_NAME_UPPER) = binary.BigEndian.Uint64(bytes[offset:offset + 8])
  offset += 8
---
  /* int64 */
  proto.$(PROPERTY_NAME_UPPER) = int64(binary.BigEndian.Uint64(bytes[offset:offset + 8]))
  offset += 8
---
  /* float32 */
  $(PROPERTY_NAME_LOWER)Bytes := util.ReverseBytes(bytes[offset:offset + 4])
  $(PROPERTY_NAME_LOWER)Bits := binary.BigEndian.Uint32($(PROPERTY_NAME_LOWER)Bytes)
  offset += 4
  proto.$(PROPERTY_NAME_UPPER) = math.Float32frombits($(PROPERTY_NAME_LOWER)Bits)
---
  /* float64 */
  $(PROPERTY_NAME_LOWER)Bytes := util.ReverseBytes(bytes[offset:offset + 8])
  $(PROPERTY_NAME_LOWER)Bits := binary.BigEndian.Uint64($(PROPERTY_NAME_LOWER)Bytes)
  offset += 8
  proto.$(PROPERTY_NAME_UPPER) = math.Float64frombits($(PROPERTY_NAME_LOWER)Bits)
---
  /* $(CUSTOM_DATA_TYPE) */
  $(PROPERTY_NAME_LOWER)Size := int(binary.BigEndian.Uint16((bytes[offset:offset + 2])))
  if $(PROPERTY_NAME_LOWER)Size + offset > len(bytes) {
    return errors.New("UnpackError")
  }
  offset += 2
  $(PROPERTY_NAME_LOWER)Bytes := bytes[offset:offset + $(PROPERTY_NAME_LOWER)Size]
  proto.$(PROPERTY_NAME_UPPER) = &$(CUSTOM_DATA_TYPE){ Ver: $(VER), Cmd: $(CMD) }
  proto.$(PROPERTY_NAME_UPPER).Unpack($(PROPERTY_NAME_LOWER)Bytes)
  offset += $(PROPERTY_NAME_LOWER)Size
---
  /* bool */
  $(PROPERTY_NAME_LOWER)Byte := bytes[offset]
  if $(PROPERTY_NAME_LOWER)Byte == 0x01 {
    proto.$(PROPERTY_NAME_UPPER) = true
  } else {
    proto.$(PROPERTY_NAME_UPPER) = false
  }
  offset++
---
  /* string */
  $(PROPERTY_NAME_LOWER)Size := int(binary.BigEndian.Uint16(bytes[offset:offset + 2]))
  if $(PROPERTY_NAME_LOWER)Size + offset > len(bytes) {
    return errors.New("UnpackError")
  }
  offset += 2
  proto.$(PROPERTY_NAME_UPPER) = string(bytes[offset:offset + $(PROPERTY_NAME_LOWER)Size])
  offset += $(PROPERTY_NAME_LOWER)Size
---
  /* []byte */
  $(PROPERTY_NAME_LOWER)Size := int(binary.BigEndian.Uint16(bytes[offset:offset + 2]))
  if $(PROPERTY_NAME_LOWER)Size + offset > len(bytes) {
    return errors.New("UnpackError")
  }
  offset += 2
  proto.$(PROPERTY_NAME_UPPER) = bytes[offset:offset + $(PROPERTY_NAME_LOWER)Size]
  offset += $(PROPERTY_NAME_LOWER)Size
---
  /* []uint8 */
  $(PROPERTY_NAME_LOWER)Size := int(binary.BigEndian.Uint16(bytes[offset:offset + 2]))
  if $(PROPERTY_NAME_LOWER)Size + offset > len(bytes) {
    return errors.New("UnpackError")
  }
  offset += 2
  for i := 0; i < $(PROPERTY_NAME_LOWER)Size; i++ {
    proto.$(PROPERTY_NAME_UPPER) = append(proto.$(PROPERTY_NAME_UPPER), uint8(bytes[offset]))
    offset++
  }
---
  /* []int8 */
  $(PROPERTY_NAME_LOWER)Size := int(binary.BigEndian.Uint16(bytes[offset:offset + 2]))
  if $(PROPERTY_NAME_LOWER)Size + offset > len(bytes) {
    return errors.New("UnpackError")
  }
  offset += 2
  for i := 0; i < $(PROPERTY_NAME_LOWER)Size; i++ {
    proto.$(PROPERTY_NAME_UPPER) = append(proto.$(PROPERTY_NAME_UPPER), int8(bytes[offset]))
    offset++
  }
---
  /* []uint16 */
  $(PROPERTY_NAME_LOWER)Size := int(binary.BigEndian.Uint16(bytes[offset:offset + 2]))
  if $(PROPERTY_NAME_LOWER)Size + offset > len(bytes) {
    return errors.New("UnpackError")
  }
  offset += 2
  for i := 0; i < $(PROPERTY_NAME_LOWER)Size; i++ {
    proto.$(PROPERTY_NAME_UPPER) = append(proto.$(PROPERTY_NAME_UPPER), binary.BigEndian.Uint16(bytes[offset:offset + 2]))
    offset += 2
  }
---
  /* []int16 */
  $(PROPERTY_NAME_LOWER)Size := int(binary.BigEndian.Uint16(bytes[offset:offset + 2]))
  if $(PROPERTY_NAME_LOWER)Size + offset > len(bytes) {
    return errors.New("UnpackError")
  }
  offset += 2
  for i := 0; i < $(PROPERTY_NAME_LOWER)Size; i++ {
    proto.$(PROPERTY_NAME_UPPER) = append(proto.$(PROPERTY_NAME_UPPER), int16(binary.BigEndian.Uint16(bytes[offset:offset + 2])))
    offset += 2
  }
---
  /* []uint32 */
  $(PROPERTY_NAME_LOWER)Size := int(binary.BigEndian.Uint16(bytes[offset:offset + 2]))
  if $(PROPERTY_NAME_LOWER)Size + offset > len(bytes) {
    return errors.New("UnpackError")
  }
  offset += 2
  for i := 0; i < $(PROPERTY_NAME_LOWER)Size; i++ {
    proto.$(PROPERTY_NAME_UPPER) = append(proto.$(PROPERTY_NAME_UPPER), binary.BigEndian.Uint32(bytes[offset:offset + 4]))
    offset += 4
  }
---
  /* []int32 */
  $(PROPERTY_NAME_LOWER)Size := int(binary.BigEndian.Uint16(bytes[offset:offset + 2]))
  if $(PROPERTY_NAME_LOWER)Size + offset > len(bytes) {
    return errors.New("UnpackError")
  }
  offset += 2
  for i := 0; i < $(PROPERTY_NAME_LOWER)Size; i++ {
    proto.$(PROPERTY_NAME_UPPER) = append(proto.$(PROPERTY_NAME_UPPER), int32(binary.BigEndian.Uint32(bytes[offset:offset + 4])))
    offset += 4
  }
---
  /* []uint64 */
  $(PROPERTY_NAME_LOWER)Size := int(binary.BigEndian.Uint16(bytes[offset:offset + 2]))
  if $(PROPERTY_NAME_LOWER)Size + offset > len(bytes) {
    return errors.New("UnpackError")
  }
  offset += 2
  for i := 0; i < $(PROPERTY_NAME_LOWER)Size; i++ {
    proto.$(PROPERTY_NAME_UPPER) = append(proto.$(PROPERTY_NAME_UPPER), binary.BigEndian.Uint64(bytes[offset:offset + 8]))
    offset += 8
  }
---
  /* []int64 */
  $(PROPERTY_NAME_LOWER)Size := int(binary.BigEndian.Uint16(bytes[offset:offset + 2]))
  if $(PROPERTY_NAME_LOWER)Size + offset > len(bytes) {
    return errors.New("UnpackError")
  }
  offset += 2
  for i := 0; i < $(PROPERTY_NAME_LOWER)Size; i++ {
    proto.$(PROPERTY_NAME_UPPER) = append(proto.$(PROPERTY_NAME_UPPER), int64(binary.BigEndian.Uint64(bytes[offset:offset + 8])))
    offset += 8
  }
---
  /* []float32 */
  $(PROPERTY_NAME_LOWER)Size := int(binary.BigEndian.Uint16(bytes[offset:offset + 2]))
  if $(PROPERTY_NAME_LOWER)Size + offset > len(bytes) {
    return errors.New("UnpackError")
  }
  offset += 2
  for i := 0; i < $(PROPERTY_NAME_LOWER)Size; i++ {
    b := util.ReverseBytes(bytes[offset:offset + 4])
    bits := binary.BigEndian.Uint32(b)
    offset += 4
    proto.$(PROPERTY_NAME_UPPER) = append(proto.$(PROPERTY_NAME_UPPER), math.Float32frombits(bits))
  }
---
  /* []float64 */
  $(PROPERTY_NAME_LOWER)Size := int(binary.BigEndian.Uint16(bytes[offset:offset + 2]))
  if $(PROPERTY_NAME_LOWER)Size + offset > len(bytes) {
    return errors.New("UnpackError")
  }
  offset += 2
  for i := 0; i < $(PROPERTY_NAME_LOWER)Size; i++ {
    b := util.ReverseBytes(bytes[offset:offset + 8])
    bits := binary.BigEndian.Uint64(b)
    offset += 8
    proto.$(PROPERTY_NAME_UPPER) = append(proto.$(PROPERTY_NAME_UPPER), math.Float64frombits(bits))
  }
---
  /* []bool */
  $(PROPERTY_NAME_LOWER)Size := int(binary.BigEndian.Uint16(bytes[offset:offset + 2]))
  if $(PROPERTY_NAME_LOWER)Size + offset > len(bytes) {
    return errors.New("UnpackError")
  }
  offset += 2
  for i := 0; i < $(PROPERTY_NAME_LOWER)Size; i++ {
    v := true
    if bytes[offset] == uint8(0x02) {
      v = false
    }
    proto.$(PROPERTY_NAME_UPPER) = append(proto.$(PROPERTY_NAME_UPPER), v)
    offset++
  }
---
  /* []string */
  $(PROPERTY_NAME_LOWER)Size := int(binary.BigEndian.Uint16(bytes[offset:offset + 2]))
  if $(PROPERTY_NAME_LOWER)Size + offset > len(bytes) {
    return errors.New("UnpackError")
  }
  offset += 2
  for i := 0; i < $(PROPERTY_NAME_LOWER)Size; i++ {
    size := int(binary.BigEndian.Uint16((bytes[offset:offset + 2])))
    offset += 2
    proto.$(PROPERTY_NAME_UPPER) = append(proto.$(PROPERTY_NAME_UPPER), string(bytes[offset:offset + size]))
    offset += size
  }
---
  /* [][]byte */
  $(PROPERTY_NAME_LOWER)Size := int(binary.BigEndian.Uint16(bytes[offset:offset + 2]))
  if $(PROPERTY_NAME_LOWER)Size + offset > len(bytes) {
    return errors.New("UnpackError")
  }
  offset += 2
  for i := 0; i < $(PROPERTY_NAME_LOWER)Size; i++ {
    size := int(binary.BigEndian.Uint16((bytes[offset:offset + 2])))
    offset += 2
    proto.$(PROPERTY_NAME_UPPER) = append(proto.$(PROPERTY_NAME_UPPER), bytes[offset:offset + size])
    offset += size
  }
---
  /* []$(CUSTOM_DATA_TYPE) */
  $(PROPERTY_NAME_LOWER)Length := int(binary.BigEndian.Uint16((bytes[offset:offset + 2])))
  offset += 2
  proto.$(PROPERTY_NAME_UPPER) = make([]*$(CUSTOM_DATA_TYPE), $(PROPERTY_NAME_LOWER)Length)
  for i := 0; i < $(PROPERTY_NAME_LOWER)Length; i++ {
    $(PROPERTY_NAME_LOWER)Size := int(binary.BigEndian.Uint16((bytes[offset:offset + 2])))
    if $(PROPERTY_NAME_LOWER)Size + offset > len(bytes) {
      return errors.New("UnpackError")
    }
    offset += 2
    $(PROPERTY_NAME_LOWER)Bytes := bytes[offset:offset + $(PROPERTY_NAME_LOWER)Size]
    item := &$(CUSTOM_DATA_TYPE){ Ver: $(VER), Cmd: $(CMD) }
    item.Unpack($(PROPERTY_NAME_LOWER)Bytes)
    proto.$(PROPERTY_NAME_UPPER)[i] = item
    offset += $(PROPERTY_NAME_LOWER)Size
  }
