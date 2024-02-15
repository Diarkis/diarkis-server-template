---
		/* uint8_t */
		_data[0 + _offset] = $(PROPERTY_NAME_LOWER);
		_offset += sizeof(uint8_t);
---
		/* uint16_t */
		PutUint16(_data + _offset, $(PROPERTY_NAME_LOWER));
		_offset += sizeof(uint16_t);
---
		/* uint32_t */
		PutUint32(_data + _offset, $(PROPERTY_NAME_LOWER));
		_offset += sizeof(uint32_t);
---
		/* uint64_t */
		PutUint64(_data + _offset, $(PROPERTY_NAME_LOWER));
		_offset += sizeof(uint64_t);
---
		/* int8_t */
		_data[0 + _offset] = $(PROPERTY_NAME_LOWER);
		_offset += sizeof(int8_t);
---
		/* int16_t */
		PutUint16(_data + _offset, (uint16_t)$(PROPERTY_NAME_LOWER));
		_offset += sizeof(int16_t);
---
		/* int32_t */
		PutUint32(_data + _offset, (uint32_t)$(PROPERTY_NAME_LOWER));
		_offset += sizeof(int32_t);
---
		/* int64_t */
		PutUint64(_data + _offset, (uint64_t)$(PROPERTY_NAME_LOWER));
		_offset += sizeof(int64_t);
---
		/* float */
		PutFloat(_data + _offset, $(PROPERTY_NAME_LOWER));
		_offset += sizeof(float);
---
		/* double */
		PutDouble(_data + _offset, $(PROPERTY_NAME_LOWER));
		_offset += sizeof(double);
---
		/* $(CUSTOM_DATA_TYPE) */
		size_t _$(PROPERTY_NAME_LOWER)Size = $(PROPERTY_NAME_LOWER).Length();
		PutUint16(_data + _offset, (uint16_t)_$(PROPERTY_NAME_LOWER)Size);
		_offset += sizeof(uint16_t);
		size_t _$(PROPERTY_NAME_LOWER)PackedBytes = 0;
		bool _$(PROPERTY_NAME_LOWER)Packed = $(PROPERTY_NAME_LOWER).Pack(_data + _offset, _size - _offset, _$(PROPERTY_NAME_LOWER)PackedBytes);
		if (!_$(PROPERTY_NAME_LOWER)Packed)
		{
			return false;
		}
		_offset += _$(PROPERTY_NAME_LOWER)PackedBytes;
---
		/* bool */
		_data[0 + _offset] = (uint8_t)$(PROPERTY_NAME_LOWER);
		_offset += sizeof(uint8_t);
---
		/* Diarkis::StdString */
		size_t _$(PROPERTY_NAME_LOWER)Size = $(PROPERTY_NAME_LOWER).length();
		PutUint16(_data + _offset, (uint16_t)$(PROPERTY_NAME_LOWER).length());
		_offset += sizeof(uint16_t);
		std::copy($(PROPERTY_NAME_LOWER).data(), $(PROPERTY_NAME_LOWER).data() + $(PROPERTY_NAME_LOWER).length(), _data + _offset);
		_offset += $(PROPERTY_NAME_LOWER).length();
---
		/* Diarkis::StdVector<uint8_t> */
		size_t _$(PROPERTY_NAME_LOWER)Size = $(PROPERTY_NAME_LOWER).size();
		PutUint16(_data + _offset, (uint16_t)_$(PROPERTY_NAME_LOWER)Size);
		_offset += sizeof(uint16_t);
		std::copy($(PROPERTY_NAME_LOWER).data(), $(PROPERTY_NAME_LOWER).data() + _$(PROPERTY_NAME_LOWER)Size, _data + _offset);
		_offset += _$(PROPERTY_NAME_LOWER)Size;
---
		/* Diarkis::StdVector<uint16_t> */
		size_t _$(PROPERTY_NAME_LOWER)Size = $(PROPERTY_NAME_LOWER).size();
		PutUint16(_data + _offset, (uint16_t)_$(PROPERTY_NAME_LOWER)Size);
		_offset += sizeof(uint16_t);
		for (const auto& e : $(PROPERTY_NAME_LOWER))
		{
			PutUint16(_data + _offset, e);
			_offset += sizeof(uint16_t);
		}
---
		/* Diarkis::StdVector<uint32_t> */
		size_t _$(PROPERTY_NAME_LOWER)Size = $(PROPERTY_NAME_LOWER).size();
		PutUint16(_data + _offset, (uint16_t)_$(PROPERTY_NAME_LOWER)Size);
		_offset += sizeof(uint16_t);
		for (const auto& e : $(PROPERTY_NAME_LOWER))
		{
			PutUint32(_data + _offset, e);
			_offset += sizeof(uint32_t);
		}
---
		/* Diarkis::StdVector<uint64_t> */
		size_t _$(PROPERTY_NAME_LOWER)Size = $(PROPERTY_NAME_LOWER).size();
		PutUint16(_data + _offset, (uint16_t)_$(PROPERTY_NAME_LOWER)Size);
		_offset += sizeof(uint16_t);
		for (const auto& e : $(PROPERTY_NAME_LOWER))
		{
			PutUint64(_data + _offset, e);
			_offset += sizeof(uint64_t);
		}
---
		/* Diarkis::StdVector<int8_t> */
		size_t _$(PROPERTY_NAME_LOWER)Size = $(PROPERTY_NAME_LOWER).size();
		PutUint16(_data + _offset, (uint16_t)_$(PROPERTY_NAME_LOWER)Size);
		_offset += sizeof(uint16_t);
		std::copy($(PROPERTY_NAME_LOWER).data(), $(PROPERTY_NAME_LOWER).data() + _$(PROPERTY_NAME_LOWER)Size, _data + _offset);
		_offset += _$(PROPERTY_NAME_LOWER)Size;
---
		/* Diarkis::StdVector<int16_t> */
		size_t _$(PROPERTY_NAME_LOWER)Size = $(PROPERTY_NAME_LOWER).size();
		PutUint16(_data + _offset, (uint16_t)_$(PROPERTY_NAME_LOWER)Size);
		_offset += sizeof(uint16_t);
		for (const auto& e : $(PROPERTY_NAME_LOWER))
		{
			PutUint16(_data + _offset, (uint16_t)e);
			_offset += sizeof(int16_t);
		}
---
		/* Diarkis::StdVector<int32_t> */
		size_t _$(PROPERTY_NAME_LOWER)Size = $(PROPERTY_NAME_LOWER).size();
		PutUint16(_data + _offset, (uint16_t)_$(PROPERTY_NAME_LOWER)Size);
		_offset += sizeof(uint16_t);
		for (const auto& e : $(PROPERTY_NAME_LOWER))
		{
			PutUint32(_data + _offset, (uint32_t)e);
			_offset += sizeof(int32_t);
		}
---
		/* Diarkis::StdVector<int64_t> */
		size_t _$(PROPERTY_NAME_LOWER)Size = $(PROPERTY_NAME_LOWER).size();
		PutUint16(_data + _offset, (uint16_t)_$(PROPERTY_NAME_LOWER)Size);
		_offset += sizeof(uint16_t);
		for (const auto& e : $(PROPERTY_NAME_LOWER))
		{
			PutUint64(_data + _offset, (uint64_t)e);
			_offset += sizeof(int64_t);
		}
---
		/* Diarkis::StdVector<float> */
		size_t _$(PROPERTY_NAME_LOWER)Size = $(PROPERTY_NAME_LOWER).size();
		PutUint16(_data + _offset, (uint16_t)_$(PROPERTY_NAME_LOWER)Size);
		_offset += sizeof(uint16_t);
		for (const auto& e : $(PROPERTY_NAME_LOWER))
		{
			PutFloat(_data + _offset, e);
			_offset += sizeof(float);
		}
---
		/* Diarkis::StdVector<double> */
		size_t _$(PROPERTY_NAME_LOWER)Size = $(PROPERTY_NAME_LOWER).size();
		PutUint16(_data + _offset, (uint16_t)_$(PROPERTY_NAME_LOWER)Size);
		_offset += sizeof(uint16_t);
		for (const auto& e : $(PROPERTY_NAME_LOWER))
		{
			PutDouble(_data + _offset, e);
			_offset += sizeof(double);
		}
---
		/* Diarkis::StdVector<bool> */
		size_t _$(PROPERTY_NAME_LOWER)Size = $(PROPERTY_NAME_LOWER).size();
		PutUint16(_data + _offset, (uint16_t)_$(PROPERTY_NAME_LOWER)Size);
		_offset += sizeof(uint16_t);
		for (bool b : $(PROPERTY_NAME_LOWER))
		{
			if (b) // true
			{
				_data[0 + _offset] = (uint8_t)0x01;
			}
			else // false
			{
				_data[0 + _offset] = (uint8_t)0x02;
			}
			_offset += sizeof(uint8_t);
		}
---
		/* Diarkis::StdVector<Diarkis::StdString> */
		size_t _$(PROPERTY_NAME_LOWER)Size = $(PROPERTY_NAME_LOWER).size();
		PutUint16(_data + _offset, (uint16_t)_$(PROPERTY_NAME_LOWER)Size);
		_offset += sizeof(uint16_t);
		for (const auto& e : $(PROPERTY_NAME_LOWER))
		{
			size_t len = e.length();
			PutUint16(_data + _offset, (uint16_t)len);
			_offset += sizeof(uint16_t);
			std::copy(e.data(), e.data() + len, _data + _offset);
			_offset += len;
		}
---
		/* Diarkis::StdVector<Diarkis::StdVector<uint8_t>> */
		size_t _$(PROPERTY_NAME_LOWER)Size = $(PROPERTY_NAME_LOWER).size();
		PutUint16(_data + _offset, (uint16_t)_$(PROPERTY_NAME_LOWER)Size);
		_offset += sizeof(uint16_t);
		for (const auto& v : $(PROPERTY_NAME_LOWER))
		{
			size_t len = v.size();
			PutUint16(_data + _offset, (uint16_t)len);
			_offset += sizeof(uint16_t);
			std::copy(v.data(), v.data() + len, _data + _offset);
			_offset += len;
		}
---
		/* Diarkis::StdVector<$(CUSTOM_DATA_TYPE)> */
		size_t _$(PROPERTY_NAME_LOWER)Size = $(PROPERTY_NAME_LOWER).size();
		PutUint16(_data + _offset, (uint16_t)_$(PROPERTY_NAME_LOWER)Size);
		_offset += sizeof(uint16_t);
		for (const auto& e : $(PROPERTY_NAME_LOWER))
		{
			size_t _elementSize = e.Length();
			PutUint16(_data + _offset, (uint16_t)_elementSize);
			_offset += sizeof(uint16_t);
			size_t _packedBytes = 0;
			bool _packed = e.Pack(_data + _offset, _size - _offset, _packedBytes);
			if (!_packed)
			{
				return false;
			}
			_offset += _packedBytes;
		}