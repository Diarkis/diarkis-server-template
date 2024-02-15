---
		/* byte */
		byte[] $(PROPERTY_NAME_LOWER)Bytes = new byte[1];
		$(PROPERTY_NAME_LOWER)Bytes[0] = $(PROPERTY_NAME_UPPER);
		bytes = DiarkisPacket.Combine(bytes, $(PROPERTY_NAME_LOWER)Bytes);
---
		/* ushort */
		byte[] $(PROPERTY_NAME_LOWER)Bytes = BitConverter.GetBytes($(PROPERTY_NAME_UPPER));
		Array.Reverse($(PROPERTY_NAME_LOWER)Bytes);
		bytes = DiarkisPacket.Combine(bytes, $(PROPERTY_NAME_LOWER)Bytes);
---
		/* uint */
		byte[] $(PROPERTY_NAME_LOWER)Bytes = BitConverter.GetBytes($(PROPERTY_NAME_UPPER));
		Array.Reverse($(PROPERTY_NAME_LOWER)Bytes);
		bytes = DiarkisPacket.Combine(bytes, $(PROPERTY_NAME_LOWER)Bytes);
---
		/* ulong */
		byte[] $(PROPERTY_NAME_LOWER)Bytes = BitConverter.GetBytes($(PROPERTY_NAME_UPPER));
		Array.Reverse($(PROPERTY_NAME_LOWER)Bytes);
		bytes = DiarkisPacket.Combine(bytes, $(PROPERTY_NAME_LOWER)Bytes);
---
		/* sbyte */
		byte[] $(PROPERTY_NAME_LOWER)Bytes = new byte[1];
		$(PROPERTY_NAME_LOWER)Bytes[0] = (byte)$(PROPERTY_NAME_UPPER);
		bytes = DiarkisPacket.Combine(bytes, $(PROPERTY_NAME_LOWER)Bytes);
---
		/* short */
		byte[] $(PROPERTY_NAME_LOWER)Bytes = BitConverter.GetBytes((ushort)$(PROPERTY_NAME_UPPER));
		Array.Reverse($(PROPERTY_NAME_LOWER)Bytes);
		bytes = DiarkisPacket.Combine(bytes, $(PROPERTY_NAME_LOWER)Bytes);
---
		/* int */
		byte[] $(PROPERTY_NAME_LOWER)Bytes = BitConverter.GetBytes((uint)$(PROPERTY_NAME_UPPER));
		Array.Reverse($(PROPERTY_NAME_LOWER)Bytes);
		bytes = DiarkisPacket.Combine(bytes, $(PROPERTY_NAME_LOWER)Bytes);
---
		/* long */
		byte[] $(PROPERTY_NAME_LOWER)Bytes = BitConverter.GetBytes((ulong)$(PROPERTY_NAME_UPPER));
		Array.Reverse($(PROPERTY_NAME_LOWER)Bytes);
		bytes = DiarkisPacket.Combine(bytes, $(PROPERTY_NAME_LOWER)Bytes);
---
		/* float */
		byte[] $(PROPERTY_NAME_LOWER)Bytes = BitConverter.GetBytes((float)$(PROPERTY_NAME_UPPER));
		Array.Reverse($(PROPERTY_NAME_LOWER)Bytes);
		bytes = DiarkisPacket.Combine(bytes, $(PROPERTY_NAME_LOWER)Bytes);
---
		/* double */
		byte[] $(PROPERTY_NAME_LOWER)Bytes = BitConverter.GetBytes((double)$(PROPERTY_NAME_UPPER));
		Array.Reverse($(PROPERTY_NAME_LOWER)Bytes);
		bytes = DiarkisPacket.Combine(bytes, $(PROPERTY_NAME_LOWER)Bytes);
---
		/* $(CUSTOM_DATA_TYPE) */
		byte[] $(PROPERTY_NAME_LOWER)Packed = $(PROPERTY_NAME_UPPER).Pack();
		byte[] $(PROPERTY_NAME_LOWER)SizeBytes = BitConverter.GetBytes((ushort)$(PROPERTY_NAME_LOWER)Packed.Length);
		Array.Reverse($(PROPERTY_NAME_LOWER)SizeBytes);
		bytes = DiarkisPacket.Combine(bytes, $(PROPERTY_NAME_LOWER)SizeBytes);
		bytes = DiarkisPacket.Combine(bytes, $(PROPERTY_NAME_LOWER)Packed);
---
		/* bool */
		byte[] $(PROPERTY_NAME_LOWER)Bytes = new byte[1];
		if ($(PROPERTY_NAME_UPPER))
		{
			$(PROPERTY_NAME_LOWER)Bytes[0] = (byte)1; /* true */
		}
		else
		{
			$(PROPERTY_NAME_LOWER)Bytes[0] = (byte)2; /* false */
		}
		bytes = DiarkisPacket.Combine(bytes, $(PROPERTY_NAME_LOWER)Bytes);
---
		/* string */
		byte[] $(PROPERTY_NAME_LOWER)Bytes = Encoding.UTF8.GetBytes($(PROPERTY_NAME_UPPER));
		byte[] $(PROPERTY_NAME_LOWER)SizeBytes = BitConverter.GetBytes((ushort)$(PROPERTY_NAME_LOWER)Bytes.Length);
		Array.Reverse($(PROPERTY_NAME_LOWER)SizeBytes);
		bytes = DiarkisPacket.Combine(bytes, $(PROPERTY_NAME_LOWER)SizeBytes);
		bytes = DiarkisPacket.Combine(bytes, $(PROPERTY_NAME_LOWER)Bytes);
---
		/* byte[] */
		byte[] $(PROPERTY_NAME_LOWER)SizeBytes = BitConverter.GetBytes((ushort)$(PROPERTY_NAME_UPPER).Length);
		Array.Reverse($(PROPERTY_NAME_LOWER)SizeBytes);
		bytes = DiarkisPacket.Combine(bytes, $(PROPERTY_NAME_LOWER)SizeBytes);
		bytes = DiarkisPacket.Combine(bytes, $(PROPERTY_NAME_UPPER));
---
		/* ushort[] */
		byte[] $(PROPERTY_NAME_LOWER)Bytes = BitConverter.GetBytes((ushort)$(PROPERTY_NAME_UPPER).Length);
		Array.Reverse($(PROPERTY_NAME_LOWER)Bytes);
		bytes = DiarkisPacket.Combine(bytes, $(PROPERTY_NAME_LOWER)Bytes);
		for (int i = 0; i < $(PROPERTY_NAME_UPPER).Length; i++)
		{
			byte[] $(PROPERTY_NAME_LOWER)ItemBytes = BitConverter.GetBytes((ushort)$(PROPERTY_NAME_UPPER)[i]);
			Array.Reverse($(PROPERTY_NAME_LOWER)ItemBytes);
			bytes = DiarkisPacket.Combine(bytes, $(PROPERTY_NAME_LOWER)ItemBytes);
		}
---
		/* uint[] */
		byte[] $(PROPERTY_NAME_LOWER)Bytes = BitConverter.GetBytes((ushort)$(PROPERTY_NAME_UPPER).Length);
		Array.Reverse($(PROPERTY_NAME_LOWER)Bytes);
		bytes = DiarkisPacket.Combine(bytes, $(PROPERTY_NAME_LOWER)Bytes);
		for (int i = 0; i < $(PROPERTY_NAME_UPPER).Length; i++)
		{
			byte[] $(PROPERTY_NAME_LOWER)ItemBytes = BitConverter.GetBytes((uint)$(PROPERTY_NAME_UPPER)[i]);
			Array.Reverse($(PROPERTY_NAME_LOWER)ItemBytes);
			bytes = DiarkisPacket.Combine(bytes, $(PROPERTY_NAME_LOWER)ItemBytes);
		}
---
		/* ulong[] */
		byte[] $(PROPERTY_NAME_LOWER)Bytes = BitConverter.GetBytes((ushort)$(PROPERTY_NAME_UPPER).Length);
		Array.Reverse($(PROPERTY_NAME_LOWER)Bytes);
		bytes = DiarkisPacket.Combine(bytes, $(PROPERTY_NAME_LOWER)Bytes);
		for (int i = 0; i < $(PROPERTY_NAME_UPPER).Length; i++)
		{
			byte[] $(PROPERTY_NAME_LOWER)ItemBytes = BitConverter.GetBytes((ulong)$(PROPERTY_NAME_UPPER)[i]);
			Array.Reverse($(PROPERTY_NAME_LOWER)ItemBytes);
			bytes = DiarkisPacket.Combine(bytes, $(PROPERTY_NAME_LOWER)ItemBytes);
		}
---
		/* sbyte[] */
		byte[] $(PROPERTY_NAME_LOWER)Bytes = BitConverter.GetBytes((ushort)$(PROPERTY_NAME_UPPER).Length);
		Array.Reverse($(PROPERTY_NAME_LOWER)Bytes);
		bytes = DiarkisPacket.Combine(bytes, $(PROPERTY_NAME_LOWER)Bytes);
		for (int i = 0; i < $(PROPERTY_NAME_UPPER).Length; i++)
		{
			byte[] $(PROPERTY_NAME_LOWER)ItemBytes = new byte[1];
			$(PROPERTY_NAME_LOWER)ItemBytes[0] = (byte)$(PROPERTY_NAME_UPPER)[i];
			bytes = DiarkisPacket.Combine(bytes, $(PROPERTY_NAME_LOWER)ItemBytes);
		}
---
		/* short[] */
		byte[] $(PROPERTY_NAME_LOWER)Bytes = BitConverter.GetBytes((ushort)$(PROPERTY_NAME_UPPER).Length);
		Array.Reverse($(PROPERTY_NAME_LOWER)Bytes);
		bytes = DiarkisPacket.Combine(bytes, $(PROPERTY_NAME_LOWER)Bytes);
		for (int i = 0; i < $(PROPERTY_NAME_UPPER).Length; i++)
		{
			byte[] $(PROPERTY_NAME_LOWER)ItemBytes = BitConverter.GetBytes((ushort)$(PROPERTY_NAME_UPPER)[i]);
			Array.Reverse($(PROPERTY_NAME_LOWER)ItemBytes);
			bytes = DiarkisPacket.Combine(bytes, $(PROPERTY_NAME_LOWER)ItemBytes);
		}
---
		/* int[] */
		byte[] $(PROPERTY_NAME_LOWER)Bytes = BitConverter.GetBytes((ushort)$(PROPERTY_NAME_UPPER).Length);
		Array.Reverse($(PROPERTY_NAME_LOWER)Bytes);
		bytes = DiarkisPacket.Combine(bytes, $(PROPERTY_NAME_LOWER)Bytes);
		for (int i = 0; i < $(PROPERTY_NAME_UPPER).Length; i++)
		{
			byte[] $(PROPERTY_NAME_LOWER)ItemBytes = BitConverter.GetBytes((int)$(PROPERTY_NAME_UPPER)[i]);
			Array.Reverse($(PROPERTY_NAME_LOWER)ItemBytes);
			bytes = DiarkisPacket.Combine(bytes, $(PROPERTY_NAME_LOWER)ItemBytes);
		}
---
		/* long[] */
		byte[] $(PROPERTY_NAME_LOWER)Bytes = BitConverter.GetBytes((ushort)$(PROPERTY_NAME_UPPER).Length);
		Array.Reverse($(PROPERTY_NAME_LOWER)Bytes);
		bytes = DiarkisPacket.Combine(bytes, $(PROPERTY_NAME_LOWER)Bytes);
		for (int i = 0; i < $(PROPERTY_NAME_UPPER).Length; i++)
		{
			byte[] $(PROPERTY_NAME_LOWER)ItemBytes = BitConverter.GetBytes((long)$(PROPERTY_NAME_UPPER)[i]);
			Array.Reverse($(PROPERTY_NAME_LOWER)ItemBytes);
			bytes = DiarkisPacket.Combine(bytes, $(PROPERTY_NAME_LOWER)ItemBytes);
		}
---
		/* float[] */
		byte[] $(PROPERTY_NAME_LOWER)Bytes = BitConverter.GetBytes((ushort)$(PROPERTY_NAME_UPPER).Length);
		Array.Reverse($(PROPERTY_NAME_LOWER)Bytes);
		bytes = DiarkisPacket.Combine(bytes, $(PROPERTY_NAME_LOWER)Bytes);
		for (int i = 0; i < $(PROPERTY_NAME_UPPER).Length; i++)
		{
			byte[] $(PROPERTY_NAME_LOWER)ItemBytes = BitConverter.GetBytes((float)$(PROPERTY_NAME_UPPER)[i]);
			Array.Reverse($(PROPERTY_NAME_LOWER)ItemBytes);
			bytes = DiarkisPacket.Combine(bytes, $(PROPERTY_NAME_LOWER)ItemBytes);
		}
---
		/* double[] */
		byte[] $(PROPERTY_NAME_LOWER)Bytes = BitConverter.GetBytes((ushort)$(PROPERTY_NAME_UPPER).Length);
		Array.Reverse($(PROPERTY_NAME_LOWER)Bytes);
		bytes = DiarkisPacket.Combine(bytes, $(PROPERTY_NAME_LOWER)Bytes);
		for (int i = 0; i < $(PROPERTY_NAME_UPPER).Length; i++)
		{
			byte[] $(PROPERTY_NAME_LOWER)ItemBytes = BitConverter.GetBytes((double)$(PROPERTY_NAME_UPPER)[i]);
			Array.Reverse($(PROPERTY_NAME_LOWER)ItemBytes);
			bytes = DiarkisPacket.Combine(bytes, $(PROPERTY_NAME_LOWER)ItemBytes);
		}
---
		/* bool[] */
		byte[] $(PROPERTY_NAME_LOWER)Bytes = BitConverter.GetBytes((ushort)$(PROPERTY_NAME_UPPER).Length);
		Array.Reverse($(PROPERTY_NAME_LOWER)Bytes);
		bytes = DiarkisPacket.Combine(bytes, $(PROPERTY_NAME_LOWER)Bytes);
		for (int i = 0; i < $(PROPERTY_NAME_UPPER).Length; i++)
		{
			byte[] $(PROPERTY_NAME_LOWER)ItemBytes = new byte[1];
			if ($(PROPERTY_NAME_UPPER)[i])
			{
				$(PROPERTY_NAME_LOWER)ItemBytes[0] = (byte)1; /* true */
			}
			else
			{
				$(PROPERTY_NAME_LOWER)ItemBytes[0] = (byte)2; /* false */
			}
			bytes = DiarkisPacket.Combine(bytes, $(PROPERTY_NAME_LOWER)ItemBytes);
		}
---
		/* string[] */
		byte[] $(PROPERTY_NAME_LOWER)Bytes = BitConverter.GetBytes((ushort)$(PROPERTY_NAME_UPPER).Length);
		Array.Reverse($(PROPERTY_NAME_LOWER)Bytes);
		bytes = DiarkisPacket.Combine(bytes, $(PROPERTY_NAME_LOWER)Bytes);
		for (int i = 0; i < $(PROPERTY_NAME_UPPER).Length; i++)
		{
			byte[] $(PROPERTY_NAME_LOWER)ItemBytes = Encoding.UTF8.GetBytes($(PROPERTY_NAME_UPPER)[i]);
			byte[] $(PROPERTY_NAME_LOWER)ItemSizeBytes = BitConverter.GetBytes((ushort)$(PROPERTY_NAME_LOWER)ItemBytes.Length);
			Array.Reverse($(PROPERTY_NAME_LOWER)ItemSizeBytes);
			bytes = DiarkisPacket.Combine(bytes, $(PROPERTY_NAME_LOWER)ItemSizeBytes);
			bytes = DiarkisPacket.Combine(bytes, $(PROPERTY_NAME_LOWER)ItemBytes);
		}
---
		/* byte[][] */
		byte[] $(PROPERTY_NAME_LOWER)Bytes = BitConverter.GetBytes((ushort)$(PROPERTY_NAME_UPPER).Length);
		Array.Reverse($(PROPERTY_NAME_LOWER)Bytes);
		bytes = DiarkisPacket.Combine(bytes, $(PROPERTY_NAME_LOWER)Bytes);
		for (int i = 0; i < $(PROPERTY_NAME_UPPER).Length; i++)
		{
			byte[] $(PROPERTY_NAME_LOWER)ItemSizeBytes = BitConverter.GetBytes((ushort)$(PROPERTY_NAME_UPPER)[i].Length);
			Array.Reverse($(PROPERTY_NAME_LOWER)ItemSizeBytes);
			bytes = DiarkisPacket.Combine(bytes, $(PROPERTY_NAME_LOWER)ItemSizeBytes);
			bytes = DiarkisPacket.Combine(bytes, $(PROPERTY_NAME_UPPER)[i]);
		}
---
		/* $(CUSTOM_DATA_TYPE)[] */
		byte[] $(PROPERTY_NAME_LOWER)Bytes = BitConverter.GetBytes((ushort)$(PROPERTY_NAME_UPPER).Length);
		Array.Reverse($(PROPERTY_NAME_LOWER)Bytes);
		bytes = DiarkisPacket.Combine(bytes, $(PROPERTY_NAME_LOWER)Bytes);
		for (int i = 0; i < $(PROPERTY_NAME_UPPER).Length; i++)
		{
			byte[] $(PROPERTY_NAME_LOWER)ItemBytes = $(PROPERTY_NAME_UPPER)[i].Pack();
			byte[] $(PROPERTY_NAME_LOWER)ItemSizeBytes = BitConverter.GetBytes((ushort)$(PROPERTY_NAME_LOWER)ItemBytes.Length);
			Array.Reverse($(PROPERTY_NAME_LOWER)ItemSizeBytes);
			bytes = DiarkisPacket.Combine(bytes, $(PROPERTY_NAME_LOWER)ItemSizeBytes);
			bytes = DiarkisPacket.Combine(bytes, $(PROPERTY_NAME_LOWER)ItemBytes);
		}
