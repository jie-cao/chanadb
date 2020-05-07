package db

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