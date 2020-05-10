package db

import "encoding/binary"

func VarintLength(v uint64) int {
	len := 1
	for v >= 128 {
		v >>= 7
		len++
	}
	return len
}

func EncodeFixed64(buffer []byte, value uint64){
// Recent clang and gcc optimize this to a single mov / str instruction.
buffer[0] = uint8(value)
buffer[1] = uint8(value >> 8)
buffer[2] = uint8(value >> 16)
buffer[3] = uint8(value >> 24)
buffer[4] = uint8(value >> 32)
buffer[5] = uint8(value >> 40)
buffer[6] = uint8(value >> 48)
buffer[7] = uint8(value >> 56)
}

func EncodeVarint32(dst []byte, v uint32) []byte{
	ptr := 0
	const bitsFlag = 128
	if v < (1 << 7) {
		dst[ptr] = uint8(v)
		ptr++
	} else if v < (1 << 14) {
		dst[ptr] = uint8(v | bitsFlag)
		ptr++
		dst[ptr] = uint8(v >> 7)
		ptr++
	} else if v < (1 << 21) {
		dst[ptr] = uint8(v | bitsFlag)
		ptr++
		dst[ptr] = uint8(v >> 7) | bitsFlag
		ptr++
		dst[ptr] = uint8(v >> 14)
		ptr++
	} else if v < (1 << 28) {
		dst[ptr] = uint8(v | bitsFlag)
		ptr++
		dst[ptr] = uint8(v >> 7) | bitsFlag
		ptr++
		dst[ptr] = uint8(v >> 14)| bitsFlag
		ptr++
		dst[ptr] = uint8(v >> 21)
		ptr++
	} else {
		dst[ptr] = uint8(v | bitsFlag)
		ptr++
		dst[ptr] = uint8(v>>7) | bitsFlag
		ptr++
		dst[ptr] = uint8(v>>14) | bitsFlag
		ptr++
		dst[ptr] = uint8(v>>21) | bitsFlag
		ptr++
		dst[ptr] = uint8(v >> 28)
		ptr++
	}
	return dst[ptr:]
}

func GetVarint32Ptr(p []byte, limit int, value *uint32) []byte {
	var ptrIndex = 0
	if ptrIndex < limit {
		result := binary.LittleEndian.Uint32(p)
		if (result & 128) == 0 {
			*value = result
			return p[1:]
		}
	}
	return GetVarint32PtrFallback(p, limit, value)
}

func GetVarint32PtrFallback(p []byte, limit int, value *uint32) []byte {
	var result uint32 = 0
	pIndex := 0
	var shift uint
	for shift = 0; shift <= 28 && pIndex < limit; shift += 7 {
		byteValue := binary.LittleEndian.Uint32(p)
		pIndex++
		if byteValue&uint32(128) != 0 {
			// More bytes are present
			result |= (byteValue & 127) << shift
		} else {
			result |= byteValue << shift
			*value = result
			return p[pIndex:]
		}
	}
	return nil
}

func DecodeFixed64(buffer []byte) uint64 {
	// Recent clang and gcc optimize this to a single mov / ldr instruction.
	return binary.LittleEndian.Uint64(buffer)
}

func DecodeFixed32(buffer []byte) uint32 {
// Recent clang and gcc optimize this to a single mov / ldr instruction.
	return binary.LittleEndian.Uint32(buffer)
}

func EncodeFixed32(buffer []byte, value uint32) {
	// Recent clang and gcc optimize this to a single mov / str instruction.
	buffer[0] = uint8(value)
	buffer[1] = uint8(value >> 8)
	buffer[2] = uint8(value >> 16)
	buffer[3] = uint8(value >> 24)
}


func GetLengthPrefixedSlice(data []byte) *Slice{
	var len uint32
	p := GetVarint32Ptr(data, 5, &len)  // +5: we assume "p" is not corrupted
	return NewSlice(p, int(len))
}

func PutLengthPrefixedSlice(dst []byte, value *Slice) {
	PutVarint32(dst, uint32(value.Size()))
	dst = append(dst, uint32(value.Size()))
}


func PutVarint32(dst []byte, v uint32) {
	buf := make([]byte, 5)
	ptr := EncodeVarint32(buf, v)
	dst = append(dst, ptr[5:]...)
}