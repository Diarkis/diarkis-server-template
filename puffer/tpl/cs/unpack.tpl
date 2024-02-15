---
		/* byte */
		$(PROPERTY_NAME_UPPER) = (byte)bytes[offset];
		offset++;
---
		/* sbyte */
		$(PROPERTY_NAME_UPPER) = (sbyte)bytes[offset];
		offset++;
---
		/* ushort */
		byte[] $(PROPERTY_NAME_LOWER)Bytes = DiarkisPacket.Slice(bytes, offset, 2);
		Array.Reverse($(PROPERTY_NAME_LOWER)Bytes);
		$(PROPERTY_NAME_UPPER) = BitConverter.ToUInt16($(PROPERTY_NAME_LOWER)Bytes);
		offset += 2;
---
		/* short */
		byte[] $(PROPERTY_NAME_LOWER)Bytes = DiarkisPacket.Slice(bytes, offset, 2);
		Array.Reverse($(PROPERTY_NAME_LOWER)Bytes);
		$(PROPERTY_NAME_UPPER) = (short)BitConverter.ToUInt16($(PROPERTY_NAME_LOWER)Bytes);
		offset += 2;
---
		/* uint */
		byte[] $(PROPERTY_NAME_LOWER)Bytes = DiarkisPacket.Slice(bytes, offset, 4);
		Array.Reverse($(PROPERTY_NAME_LOWER)Bytes);
		$(PROPERTY_NAME_UPPER) = BitConverter.ToUInt32($(PROPERTY_NAME_LOWER)Bytes);
		offset += 4;
---
		/* int */
		byte[] $(PROPERTY_NAME_LOWER)Bytes = DiarkisPacket.Slice(bytes, offset, 4);
		Array.Reverse($(PROPERTY_NAME_LOWER)Bytes);
		$(PROPERTY_NAME_UPPER) = (int)BitConverter.ToUInt32($(PROPERTY_NAME_LOWER)Bytes);
		offset += 4;
---
		/* ulong */
		byte[] $(PROPERTY_NAME_LOWER)Bytes = DiarkisPacket.Slice(bytes, offset, 8);
		Array.Reverse($(PROPERTY_NAME_LOWER)Bytes);
		$(PROPERTY_NAME_UPPER) = BitConverter.ToUInt64($(PROPERTY_NAME_LOWER)Bytes);
		offset += 8;
---
		/* long */
		byte[] $(PROPERTY_NAME_LOWER)Bytes = DiarkisPacket.Slice(bytes, offset, 8);
		Array.Reverse($(PROPERTY_NAME_LOWER)Bytes);
		$(PROPERTY_NAME_UPPER) = (long)BitConverter.ToUInt64($(PROPERTY_NAME_LOWER)Bytes);
		offset += 8;
---
		/* float */
		byte[] $(PROPERTY_NAME_LOWER)Bytes = DiarkisPacket.Slice(bytes, offset, 4);
		Array.Reverse($(PROPERTY_NAME_LOWER)Bytes);
		$(PROPERTY_NAME_UPPER) = BitConverter.ToSingle($(PROPERTY_NAME_LOWER)Bytes);
		offset += 4;
---
		/* double */
		byte[] $(PROPERTY_NAME_LOWER)Bytes = DiarkisPacket.Slice(bytes, offset, 8);
		Array.Reverse($(PROPERTY_NAME_LOWER)Bytes);
		$(PROPERTY_NAME_UPPER) = BitConverter.ToDouble($(PROPERTY_NAME_LOWER)Bytes);
		offset += 8;
---
		/* $(CUSTOM_DATA_TYPE) */
		byte[] $(PROPERTY_NAME_LOWER)SizeBytes = DiarkisPacket.Slice(bytes, offset, 2);
		Array.Reverse($(PROPERTY_NAME_LOWER)SizeBytes);
		int $(PROPERTY_NAME_LOWER)Size = (int)BitConverter.ToUInt16($(PROPERTY_NAME_LOWER)SizeBytes);
		if ($(PROPERTY_NAME_LOWER)Size + offset > bytes.Length)
    {
      return false;
    }
		offset += 2;
		byte[] $(PROPERTY_NAME_LOWER)Bytes = DiarkisPacket.Slice(bytes, offset, $(PROPERTY_NAME_LOWER)Size);
		$(PROPERTY_NAME_UPPER) = new Diarkis.$(CUSTOM_DATA_TYPE)();
		$(PROPERTY_NAME_UPPER).Unpack($(PROPERTY_NAME_LOWER)Bytes);
		offset += $(PROPERTY_NAME_LOWER)Size;
---
		/* bool */
		$(PROPERTY_NAME_UPPER) = (byte)bytes[offset] == 0x01 ? true : false;;
		offset++;
---
		/* string */
		byte[] $(PROPERTY_NAME_LOWER)SizeBytes = DiarkisPacket.Slice(bytes, offset, 2);
		Array.Reverse($(PROPERTY_NAME_LOWER)SizeBytes);
		int $(PROPERTY_NAME_LOWER)Size = (int)BitConverter.ToUInt16($(PROPERTY_NAME_LOWER)SizeBytes);
		if ($(PROPERTY_NAME_LOWER)Size + offset > bytes.Length)
    {
      return false;
    }
		offset += 2;
		$(PROPERTY_NAME_UPPER) = Encoding.UTF8.GetString(DiarkisPacket.Slice(bytes, offset, $(PROPERTY_NAME_LOWER)Size));
		offset += $(PROPERTY_NAME_LOWER)Size;
---
		/* byte[] */
		byte[] $(PROPERTY_NAME_LOWER)SizeBytes = DiarkisPacket.Slice(bytes, offset, 2);
		Array.Reverse($(PROPERTY_NAME_LOWER)SizeBytes);
		int $(PROPERTY_NAME_LOWER)Size = (int)BitConverter.ToUInt16($(PROPERTY_NAME_LOWER)SizeBytes);
		if ($(PROPERTY_NAME_LOWER)Size + offset > bytes.Length)
    {
      return false;
    }
		offset += 2;
		$(PROPERTY_NAME_UPPER) = DiarkisPacket.Slice(bytes, offset, $(PROPERTY_NAME_LOWER)Size);
		offset += $(PROPERTY_NAME_LOWER)Size;
---
		/* sbyte[] */
		byte[] $(PROPERTY_NAME_LOWER)SizeBytes = DiarkisPacket.Slice(bytes, offset, 2);
		Array.Reverse($(PROPERTY_NAME_LOWER)SizeBytes);
		int $(PROPERTY_NAME_LOWER)Size = (int)BitConverter.ToUInt16($(PROPERTY_NAME_LOWER)SizeBytes);
		if ($(PROPERTY_NAME_LOWER)Size + offset > bytes.Length)
    {
      return false;
    }
		offset += 2;
		$(PROPERTY_NAME_UPPER) = new sbyte[$(PROPERTY_NAME_LOWER)Size];
		for (int i = 0; i < $(PROPERTY_NAME_LOWER)Size; i++)
		{
			$(PROPERTY_NAME_UPPER)[i] = (sbyte)bytes[offset];
			offset++;
		}
---
		/* ushort[] */
		byte[] $(PROPERTY_NAME_LOWER)SizeBytes = DiarkisPacket.Slice(bytes, offset, 2);
		Array.Reverse($(PROPERTY_NAME_LOWER)SizeBytes);
		int $(PROPERTY_NAME_LOWER)Size = (int)BitConverter.ToUInt16($(PROPERTY_NAME_LOWER)SizeBytes);
		if ($(PROPERTY_NAME_LOWER)Size + offset > bytes.Length)
    {
      return false;
    }
		offset += 2;
		$(PROPERTY_NAME_UPPER) = new ushort[$(PROPERTY_NAME_LOWER)Size];
		for (int i = 0; i < $(PROPERTY_NAME_LOWER)Size; i++)
		{
		  byte[] $(PROPERTY_NAME_LOWER)ItemBytes = DiarkisPacket.Slice(bytes, offset, 2);
			Array.Reverse($(PROPERTY_NAME_LOWER)ItemBytes);
			$(PROPERTY_NAME_UPPER)[i] = BitConverter.ToUInt16($(PROPERTY_NAME_LOWER)ItemBytes);
			offset += 2;
		}
---
		/* short[] */
		byte[] $(PROPERTY_NAME_LOWER)SizeBytes = DiarkisPacket.Slice(bytes, offset, 2);
		Array.Reverse($(PROPERTY_NAME_LOWER)SizeBytes);
		int $(PROPERTY_NAME_LOWER)Size = (int)BitConverter.ToUInt16($(PROPERTY_NAME_LOWER)SizeBytes);
		if ($(PROPERTY_NAME_LOWER)Size + offset > bytes.Length)
    {
      return false;
    }
		offset += 2;
		$(PROPERTY_NAME_UPPER) = new short[$(PROPERTY_NAME_LOWER)Size];
		for (int i = 0; i < $(PROPERTY_NAME_LOWER)Size; i++)
		{
		  byte[] $(PROPERTY_NAME_LOWER)ItemBytes = DiarkisPacket.Slice(bytes, offset, 2);
			Array.Reverse($(PROPERTY_NAME_LOWER)ItemBytes);
			$(PROPERTY_NAME_UPPER)[i] = (short)BitConverter.ToUInt16($(PROPERTY_NAME_LOWER)ItemBytes);
			offset += 2;
		}
---
		/* uint[] */
		byte[] $(PROPERTY_NAME_LOWER)SizeBytes = DiarkisPacket.Slice(bytes, offset, 2);
		Array.Reverse($(PROPERTY_NAME_LOWER)SizeBytes);
		int $(PROPERTY_NAME_LOWER)Size = (int)BitConverter.ToUInt16($(PROPERTY_NAME_LOWER)SizeBytes);
		if ($(PROPERTY_NAME_LOWER)Size + offset > bytes.Length)
    {
      return false;
    }
		offset += 2;
		$(PROPERTY_NAME_UPPER) = new uint[$(PROPERTY_NAME_LOWER)Size];
		for (int i = 0; i < $(PROPERTY_NAME_LOWER)Size; i++)
		{
		  byte[] $(PROPERTY_NAME_LOWER)ItemBytes = DiarkisPacket.Slice(bytes, offset, 4);
			Array.Reverse($(PROPERTY_NAME_LOWER)ItemBytes);
			$(PROPERTY_NAME_UPPER)[i] = BitConverter.ToUInt32($(PROPERTY_NAME_LOWER)ItemBytes);
			offset += 4;
		}
---
		/* int[] */
		byte[] $(PROPERTY_NAME_LOWER)SizeBytes = DiarkisPacket.Slice(bytes, offset, 2);
		Array.Reverse($(PROPERTY_NAME_LOWER)SizeBytes);
		int $(PROPERTY_NAME_LOWER)Size = (int)BitConverter.ToUInt16($(PROPERTY_NAME_LOWER)SizeBytes);
		if ($(PROPERTY_NAME_LOWER)Size + offset > bytes.Length)
    {
      return false;
    }
		offset += 2;
		$(PROPERTY_NAME_UPPER) = new int[$(PROPERTY_NAME_LOWER)Size];
		for (int i = 0; i < $(PROPERTY_NAME_LOWER)Size; i++)
		{
		  byte[] $(PROPERTY_NAME_LOWER)ItemBytes = DiarkisPacket.Slice(bytes, offset, 4);
			Array.Reverse($(PROPERTY_NAME_LOWER)ItemBytes);
			$(PROPERTY_NAME_UPPER)[i] = (int)BitConverter.ToUInt32($(PROPERTY_NAME_LOWER)ItemBytes);
			offset += 4;
		}
---
		/* ulong[] */
		byte[] $(PROPERTY_NAME_LOWER)SizeBytes = DiarkisPacket.Slice(bytes, offset, 2);
		Array.Reverse($(PROPERTY_NAME_LOWER)SizeBytes);
		int $(PROPERTY_NAME_LOWER)Size = (int)BitConverter.ToUInt16($(PROPERTY_NAME_LOWER)SizeBytes);
		if ($(PROPERTY_NAME_LOWER)Size + offset > bytes.Length)
    {
      return false;
    }
		offset += 2;
		$(PROPERTY_NAME_UPPER) = new ulong[$(PROPERTY_NAME_LOWER)Size];
		for (int i = 0; i < $(PROPERTY_NAME_LOWER)Size; i++)
		{
		  byte[] $(PROPERTY_NAME_LOWER)ItemBytes = DiarkisPacket.Slice(bytes, offset, 8);
			Array.Reverse($(PROPERTY_NAME_LOWER)ItemBytes);
			$(PROPERTY_NAME_UPPER)[i] = BitConverter.ToUInt64($(PROPERTY_NAME_LOWER)ItemBytes);
			offset += 8;
		}
---
		/* long[] */
		byte[] $(PROPERTY_NAME_LOWER)SizeBytes = DiarkisPacket.Slice(bytes, offset, 2);
		Array.Reverse($(PROPERTY_NAME_LOWER)SizeBytes);
		int $(PROPERTY_NAME_LOWER)Size = (int)BitConverter.ToUInt16($(PROPERTY_NAME_LOWER)SizeBytes);
		if ($(PROPERTY_NAME_LOWER)Size + offset > bytes.Length)
    {
      return false;
    }
		offset += 2;
		$(PROPERTY_NAME_UPPER) = new long[$(PROPERTY_NAME_LOWER)Size];
		for (int i = 0; i < $(PROPERTY_NAME_LOWER)Size; i++)
		{
		  byte[] $(PROPERTY_NAME_LOWER)ItemBytes = DiarkisPacket.Slice(bytes, offset, 8);
			Array.Reverse($(PROPERTY_NAME_LOWER)ItemBytes);
			$(PROPERTY_NAME_UPPER)[i] = (long)BitConverter.ToUInt64($(PROPERTY_NAME_LOWER)ItemBytes);
			offset += 8;
		}
---
		/* float[] */
		byte[] $(PROPERTY_NAME_LOWER)SizeBytes = DiarkisPacket.Slice(bytes, offset, 2);
		Array.Reverse($(PROPERTY_NAME_LOWER)SizeBytes);
		int $(PROPERTY_NAME_LOWER)Size = (int)BitConverter.ToUInt16($(PROPERTY_NAME_LOWER)SizeBytes);
		if ($(PROPERTY_NAME_LOWER)Size + offset > bytes.Length)
    {
      return false;
    }
		offset += 2;
		$(PROPERTY_NAME_UPPER) = new float[$(PROPERTY_NAME_LOWER)Size];
		for (int i = 0; i < $(PROPERTY_NAME_LOWER)Size; i++)
		{
		  byte[] $(PROPERTY_NAME_LOWER)ItemBytes = DiarkisPacket.Slice(bytes, offset, 4);
			Array.Reverse($(PROPERTY_NAME_LOWER)ItemBytes);
			$(PROPERTY_NAME_UPPER)[i] = BitConverter.ToSingle($(PROPERTY_NAME_LOWER)ItemBytes);
			offset += 4;
		}
---
		/* double[] */
		byte[] $(PROPERTY_NAME_LOWER)SizeBytes = DiarkisPacket.Slice(bytes, offset, 2);
		Array.Reverse($(PROPERTY_NAME_LOWER)SizeBytes);
		int $(PROPERTY_NAME_LOWER)Size = (int)BitConverter.ToUInt16($(PROPERTY_NAME_LOWER)SizeBytes);
		if ($(PROPERTY_NAME_LOWER)Size + offset > bytes.Length)
    {
      return false;
    }
		offset += 2;
		$(PROPERTY_NAME_UPPER) = new double[$(PROPERTY_NAME_LOWER)Size];
		for (int i = 0; i < $(PROPERTY_NAME_LOWER)Size; i++)
		{
		  byte[] $(PROPERTY_NAME_LOWER)ItemBytes = DiarkisPacket.Slice(bytes, offset, 8);
			Array.Reverse($(PROPERTY_NAME_LOWER)ItemBytes);
			$(PROPERTY_NAME_UPPER)[i] = BitConverter.ToDouble($(PROPERTY_NAME_LOWER)ItemBytes);
			offset += 8;
		}
---
		/* bool[] */
		byte[] $(PROPERTY_NAME_LOWER)SizeBytes = DiarkisPacket.Slice(bytes, offset, 2);
		Array.Reverse($(PROPERTY_NAME_LOWER)SizeBytes);
		int $(PROPERTY_NAME_LOWER)Size = (int)BitConverter.ToUInt16($(PROPERTY_NAME_LOWER)SizeBytes);
		if ($(PROPERTY_NAME_LOWER)Size + offset > bytes.Length)
    {
      return false;
    }
		offset += 2;
		$(PROPERTY_NAME_UPPER) = new bool[$(PROPERTY_NAME_LOWER)Size];
		for (int i = 0; i < $(PROPERTY_NAME_LOWER)Size; i++)
		{
			$(PROPERTY_NAME_UPPER)[i] = bytes[offset] == 0x01 ? true : false;
			offset++;
		}
---
		/* string[] */
		byte[] $(PROPERTY_NAME_LOWER)SizeBytes = DiarkisPacket.Slice(bytes, offset, 2);
		Array.Reverse($(PROPERTY_NAME_LOWER)SizeBytes);
		int $(PROPERTY_NAME_LOWER)Size = (int)BitConverter.ToUInt16($(PROPERTY_NAME_LOWER)SizeBytes);
		if ($(PROPERTY_NAME_LOWER)Size + offset > bytes.Length)
    {
      return false;
    }
		offset += 2;
		$(PROPERTY_NAME_UPPER) = new string[$(PROPERTY_NAME_LOWER)Size];
		for (int i = 0; i < $(PROPERTY_NAME_LOWER)Size; i++)
		{
		  byte[] $(PROPERTY_NAME_LOWER)ItemSizeBytes = DiarkisPacket.Slice(bytes, offset, 2);
			Array.Reverse($(PROPERTY_NAME_LOWER)ItemSizeBytes);
			int $(PROPERTY_NAME_LOWER)ItemSize = (int)BitConverter.ToUInt16($(PROPERTY_NAME_LOWER)ItemSizeBytes);
			offset += 2;
			$(PROPERTY_NAME_UPPER)[i] = Encoding.UTF8.GetString(DiarkisPacket.Slice(bytes, offset, $(PROPERTY_NAME_LOWER)ItemSize));
			offset += $(PROPERTY_NAME_LOWER)ItemSize;
		}
---
		/* byte[][] */
		byte[] $(PROPERTY_NAME_LOWER)SizeBytes = DiarkisPacket.Slice(bytes, offset, 2);
		Array.Reverse($(PROPERTY_NAME_LOWER)SizeBytes);
		int $(PROPERTY_NAME_LOWER)Size = (int)BitConverter.ToUInt16($(PROPERTY_NAME_LOWER)SizeBytes);
		if ($(PROPERTY_NAME_LOWER)Size + offset > bytes.Length)
    {
      return false;
    }
		offset += 2;
		$(PROPERTY_NAME_UPPER) = new byte[$(PROPERTY_NAME_LOWER)Size][];
		for (int i = 0; i < $(PROPERTY_NAME_LOWER)Size; i++)
		{
		  byte[] $(PROPERTY_NAME_LOWER)ItemSizeBytes = DiarkisPacket.Slice(bytes, offset, 2);
			Array.Reverse($(PROPERTY_NAME_LOWER)ItemSizeBytes);
			int $(PROPERTY_NAME_LOWER)ItemSize = (int)BitConverter.ToUInt16($(PROPERTY_NAME_LOWER)ItemSizeBytes);
			offset += 2;
			$(PROPERTY_NAME_UPPER)[i] = DiarkisPacket.Slice(bytes, offset, $(PROPERTY_NAME_LOWER)ItemSize);
			offset += $(PROPERTY_NAME_LOWER)ItemSize;
		}
---
		/* $(CUSTOM_DATA_TYPE)[] */
		byte[] $(PROPERTY_NAME_LOWER)SizeBytes = DiarkisPacket.Slice(bytes, offset, 2);
		Array.Reverse($(PROPERTY_NAME_LOWER)SizeBytes);
		int $(PROPERTY_NAME_LOWER)Size = (int)BitConverter.ToUInt16($(PROPERTY_NAME_LOWER)SizeBytes);
		if ($(PROPERTY_NAME_LOWER)Size + offset > bytes.Length)
    {
      return false;
    }
		offset += 2;
		$(PROPERTY_NAME_UPPER) = new Diarkis.$(CUSTOM_DATA_TYPE)[$(PROPERTY_NAME_LOWER)Size];
		for (int i = 0; i < $(PROPERTY_NAME_LOWER)Size; i++)
		{
			byte[] $(PROPERTY_NAME_LOWER)ItemSizeBytes = DiarkisPacket.Slice(bytes, offset, 2);
			Array.Reverse($(PROPERTY_NAME_LOWER)ItemSizeBytes);
			int $(PROPERTY_NAME_LOWER)ItemSize = (int)BitConverter.ToUInt16($(PROPERTY_NAME_LOWER)ItemSizeBytes);
			offset += 2;
      byte[] $(PROPERTY_NAME_LOWER)ItemBytes = DiarkisPacket.Slice(bytes, offset, $(PROPERTY_NAME_LOWER)ItemSize);
		  $(CUSTOM_DATA_TYPE) $(PROPERTY_NAME_LOWER)Item = new Diarkis.$(CUSTOM_DATA_TYPE)();
		  $(PROPERTY_NAME_LOWER)Item.Unpack($(PROPERTY_NAME_LOWER)ItemBytes);
			$(PROPERTY_NAME_UPPER)[i] = $(PROPERTY_NAME_LOWER)Item;
			offset += $(PROPERTY_NAME_LOWER)ItemSize;
		}
