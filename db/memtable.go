package db

type Area struct {
	allocPtr *byte
	allocBytesRemaining int64
	blocks []*byte
	memoryUsage int64
}

func (are*Area) Allocate(bytes int64) *byte {

}