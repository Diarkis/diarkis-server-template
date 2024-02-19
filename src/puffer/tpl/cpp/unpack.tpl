---
		/* uint8_t */
		$(PROPERTY_NAME_LOWER) = (uint8_t)_data[0 + _offset];
		_offset += sizeof(uint8_t);
---
		/* uint16_t */
		$(PROPERTY_NAME_LOWER) = GetUint16(_data + _offset);
		_offset += sizeof(uint16_t);
---
		/* uint32_t */
		$(PROPERTY_NAME_LOWER) = GetUint32(_data + _offset);
		_offset += sizeof(uint32_t);
---
		/* uint64_t */
		$(PROPERTY_NAME_LOWER) = GetUint64(_data + _offset);
		_offset += sizeof(uint64_t);
---
		/* int8_t */
		$(PROPERTY_NAME_LOWER) = (int8_t)_data[0 + _offset];
		_offset += sizeof(int8_t);
---
		/* int16_t */
		$(PROPERTY_NAME_LOWER) = (int16_t)GetUint16(_data + _offset);
		_offset += sizeof(int16_t);
---
		/* int32_t */
		$(PROPERTY_NAME_LOWER) = (int32_t)GetUint32(_data + _offset);
		_offset += sizeof(int32_t);
---
		/* int64_t */
		$(PROPERTY_NAME_LOWER) = (int64_t)GetUint64(_data + _offset);
		_offset += sizeof(int64_t);
---
		/* float */
		$(PROPERTY_NAME_LOWER) = GetFloat(_data + _offset);
		_offset += sizeof(float);
---
		/* double */
		$(PROPERTY_NAME_LOWER) = GetDouble(_data + _offset);
		_offset += sizeof(double);
---
		/* $(CUSTOM_DATA_TYPE) */
		uint16_t _$(PROPERTY_NAME_LOWER)Size = GetUint16(_data + _offset);
		_offset += sizeof(uint16_t);
		size_t _$(PROPERTY_NAME_LOWER)UnpackedBytes = 0;
		bool _$(PROPERTY_NAME_LOWER)Packed = $(PROPERTY_NAME_LOWER).Unpack(_data + _offset, _$(PROPERTY_NAME_LOWER)Size, _$(PROPERTY_NAME_LOWER)UnpackedBytes);
		if (!_$(PROPERTY_NAME_LOWER)Packed)
		{
			return false;
		}
		_offset += _$(PROPERTY_NAME_LOWER)UnpackedBytes;
---
		/* bool */
		$(PROPERTY_NAME_LOWER) = (bool)_data[0 + _offset];
		_offset += sizeof(bool);
---
		/* Diarkis::StdString */
		uint16_t _$(PROPERTY_NAME_LOWER)Size = GetUint16(_data + _offset);
		_offset += sizeof(uint16_t);
		$(PROPERTY_NAME_LOWER) = Diarkis::StdString(_data + _offset, _data + _offset + _$(PROPERTY_NAME_LOWER)Size);
		_offset += _$(PROPERTY_NAME_LOWER)Size;
---
		/* Diarkis::StdVector<uint8_t> */
		uint16_t _$(PROPERTY_NAME_LOWER)Size = GetUint16(_data + _offset);
		_offset += sizeof(uint16_t);
		$(PROPERTY_NAME_LOWER).resize(_$(PROPERTY_NAME_LOWER)Size);
		std::copy(_data + _offset, _data + _offset + _$(PROPERTY_NAME_LOWER)Size, $(PROPERTY_NAME_LOWER).data());
		_offset += _$(PROPERTY_NAME_LOWER)Size;
---
		/* Diarkis::StdVector<uint16_t> */
		$(PROPERTY_NAME_LOWER).clear();
		uint16_t _$(PROPERTY_NAME_LOWER)Size = GetUint16(_data + _offset);
		_offset += sizeof(uint16_t);
		$(PROPERTY_NAME_LOWER).reserve(_$(PROPERTY_NAME_LOWER)Size);
		for (int i = 0; i < _$(PROPERTY_NAME_LOWER)Size; i++)
		{
			$(PROPERTY_NAME_LOWER).push_back(GetUint16(_data + _offset));
			_offset += sizeof(uint16_t);
		}
---
		/* Diarkis::StdVector<uint32_t> */
		$(PROPERTY_NAME_LOWER).clear();
		uint16_t _$(PROPERTY_NAME_LOWER)Size = GetUint16(_data + _offset);
		_offset += sizeof(uint16_t);
		$(PROPERTY_NAME_LOWER).reserve(_$(PROPERTY_NAME_LOWER)Size);
		for (int i = 0; i < _$(PROPERTY_NAME_LOWER)Size; i++)
		{
			$(PROPERTY_NAME_LOWER).push_back(GetUint32(_data + _offset));
			_offset += sizeof(uint32_t);
		}
---
		/* Diarkis::StdVector<uint64_t> */
		$(PROPERTY_NAME_LOWER).clear();
		uint16_t _$(PROPERTY_NAME_LOWER)Size = GetUint16(_data + _offset);
		_offset += sizeof(uint16_t);
		$(PROPERTY_NAME_LOWER).reserve(_$(PROPERTY_NAME_LOWER)Size);
		for (int i = 0; i < _$(PROPERTY_NAME_LOWER)Size; i++)
		{
			$(PROPERTY_NAME_LOWER).push_back(GetUint64(_data + _offset));
			_offset += sizeof(uint64_t);
		}
---
		/* Diarkis::StdVector<int8_t> */
		uint16_t _$(PROPERTY_NAME_LOWER)Size = GetUint16(_data + _offset);
		_offset += sizeof(uint16_t);
		$(PROPERTY_NAME_LOWER).resize(_$(PROPERTY_NAME_LOWER)Size);
		std::copy(_data + _offset, _data + _offset + _$(PROPERTY_NAME_LOWER)Size, $(PROPERTY_NAME_LOWER).data());
		_offset += _$(PROPERTY_NAME_LOWER)Size;
---
		/* Diarkis::StdVector<int16_t> */
		$(PROPERTY_NAME_LOWER).clear();
		uint16_t _$(PROPERTY_NAME_LOWER)Size = GetUint16(_data + _offset);
		_offset += sizeof(uint16_t);
		$(PROPERTY_NAME_LOWER).reserve(_$(PROPERTY_NAME_LOWER)Size);
		for (int i = 0; i < _$(PROPERTY_NAME_LOWER)Size; i++)
		{
			$(PROPERTY_NAME_LOWER).push_back((int16_t)GetUint16(_data + _offset));
			_offset += sizeof(uint16_t);
		}
---
		/* Diarkis::StdVector<int32_t> */
		$(PROPERTY_NAME_LOWER).clear();
		uint16_t _$(PROPERTY_NAME_LOWER)Size = GetUint16(_data + _offset);
		_offset += sizeof(uint16_t);
		$(PROPERTY_NAME_LOWER).reserve(_$(PROPERTY_NAME_LOWER)Size);
		for (int i = 0; i < _$(PROPERTY_NAME_LOWER)Size; i++)
		{
			$(PROPERTY_NAME_LOWER).push_back((int32_t)GetUint32(_data + _offset));
			_offset += sizeof(int32_t);
		}
---
		/* Diarkis::StdVector<int64_t> */
		$(PROPERTY_NAME_LOWER).clear();
		uint16_t _$(PROPERTY_NAME_LOWER)Size = GetUint16(_data + _offset);
		_offset += sizeof(uint16_t);
		$(PROPERTY_NAME_LOWER).reserve(_$(PROPERTY_NAME_LOWER)Size);
		for (int i = 0; i < _$(PROPERTY_NAME_LOWER)Size; i++)
		{
			$(PROPERTY_NAME_LOWER).push_back((int64_t)GetUint64(_data + _offset));
			_offset += sizeof(int64_t);
		}
---
		/* Diarkis::StdVector<float> */
		$(PROPERTY_NAME_LOWER).clear();
		uint16_t _$(PROPERTY_NAME_LOWER)Size = GetUint16(_data + _offset);
		_offset += sizeof(uint16_t);
		$(PROPERTY_NAME_LOWER).reserve(_$(PROPERTY_NAME_LOWER)Size);
		for (int i = 0; i < _$(PROPERTY_NAME_LOWER)Size; i++)
		{
			$(PROPERTY_NAME_LOWER).push_back(GetFloat(_data + _offset));
			_offset += sizeof(float);
		}
---
		/* Diarkis::StdVector<double> */
		$(PROPERTY_NAME_LOWER).clear();
		uint16_t _$(PROPERTY_NAME_LOWER)Size = GetUint16(_data + _offset);
		_offset += sizeof(uint16_t);
		$(PROPERTY_NAME_LOWER).reserve(_$(PROPERTY_NAME_LOWER)Size);
		for (int i = 0; i < _$(PROPERTY_NAME_LOWER)Size; i++)
		{
			$(PROPERTY_NAME_LOWER).push_back(GetDouble(_data + _offset));
			_offset += sizeof(double);
		}
---
		/* Diarkis::StdVector<bool> */
		$(PROPERTY_NAME_LOWER).clear();
		uint16_t _$(PROPERTY_NAME_LOWER)Size = GetUint16(_data + _offset);
		_offset += sizeof(uint16_t);
		$(PROPERTY_NAME_LOWER).reserve(_$(PROPERTY_NAME_LOWER)Size);
		for (int i = 0; i < _$(PROPERTY_NAME_LOWER)Size; i++)
		{
			uint8_t v = _data[0 + _offset];
			if (v == 0x01)
			{
				$(PROPERTY_NAME_LOWER).push_back(true);
			}
			else if (v == 0x02)
			{
				$(PROPERTY_NAME_LOWER).push_back(false);
			}
			_offset += sizeof(uint8_t);
		}
---
		/* Diarkis::StdVector<Diarkis::StdString> */
		$(PROPERTY_NAME_LOWER).clear();
		uint16_t _$(PROPERTY_NAME_LOWER)Size = GetUint16(_data + _offset);
		_offset += sizeof(uint16_t);
		$(PROPERTY_NAME_LOWER).reserve(_$(PROPERTY_NAME_LOWER)Size);
		for (int i = 0; i < _$(PROPERTY_NAME_LOWER)Size; i++)
		{
			uint16_t _strSize = GetUint16(_data + _offset);
			_offset += sizeof(uint16_t);
			Diarkis::StdString str = Diarkis::StdString(_data + _offset, _data + _offset + _strSize);
			_offset += _strSize;
			$(PROPERTY_NAME_LOWER).push_back(str);
		}
---
		/* Diarkis::StdVector<Diarkis::StdVector<uint8_t>> */
		$(PROPERTY_NAME_LOWER).clear();
		uint16_t _$(PROPERTY_NAME_LOWER)Size = GetUint16(_data + _offset);
		_offset += sizeof(uint16_t);
		$(PROPERTY_NAME_LOWER).reserve(_$(PROPERTY_NAME_LOWER)Size);
		for (int i = 0; i < _$(PROPERTY_NAME_LOWER)Size; i++)
		{
			uint16_t _bytesSize = GetUint16(_data + _offset);
			_offset += sizeof(uint16_t);
			Diarkis::StdVector<uint8_t> _bytes(_data + _offset, _data + _offset + _bytesSize);
			_offset += _bytesSize;
			$(PROPERTY_NAME_LOWER).push_back(_bytes);
		}
---
		/* Diarkis::StdVector<$(CUSTOM_DATA_TYPE)> */
		$(PROPERTY_NAME_LOWER).clear();
		uint16_t _$(PROPERTY_NAME_LOWER)Size = GetUint16(_data + _offset);
		_offset += sizeof(uint16_t);
		$(PROPERTY_NAME_LOWER).reserve(_$(PROPERTY_NAME_LOWER)Size);
		for (int i = 0; i < _$(PROPERTY_NAME_LOWER)Size; i++)
		{
			uint16_t _elementSize = GetUint16(_data + _offset);
			_offset += sizeof(uint16_t);
			$(CUSTOM_DATA_TYPE) e;
			size_t _unpackedBytes = 0;
			bool _unpacked = e.Unpack(_data + _offset, _size - _offset, _unpackedBytes);
			if (!_unpacked)
			{
				return false;
			}
			$(PROPERTY_NAME_LOWER).push_back(e);
			_offset += _unpackedBytes;
		}
