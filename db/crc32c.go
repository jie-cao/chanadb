package db

const kMaskDelta = 0xa282ead8

func Mask(crc uint32) uint32 {
	return ((crc >>15)|(crc <<17)) + kMaskDelta
}

func UnMask(maskedCrc uint32) uint32 {
	rot := maskedCrc - kMaskDelta
	return (rot >> 17) | (rot << 15)
}
