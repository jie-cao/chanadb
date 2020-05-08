package db
import "bytes"
import "math"

type Slice struct {
	data []byte
	size uint
}

func NewSlice(data []byte, size int) *Slice {
	s := new(Slice)
	s.data = data
	s.size = uint(size)
	return s
}

func NewSliceFromBytes(data []byte)  *Slice{
	s :=new(Slice)
	s.data = data
	s.size = uint(len(data))
	return s
}

func NewSliceFromString(inputString string) *Slice {
	s := new(Slice)
	s.data = []byte(inputString)
	s.size = uint(len(inputString))
	return s
}

func (slice *Slice) Size() uint {
	return slice.size
}

func (slice *Slice) Empty() {
	slice.size = 0
}

func (slice *Slice) Data() []byte  {
	return slice.data
}

func (slice *Slice) Equals(compareSlice *Slice) bool {
	return  (slice.Size() == compareSlice.Size()) && (bytes.Compare(slice.data, compareSlice.data) == 0)
}

func (slice *Slice) Compare(compareSlice *Slice) int {
	minLen := int(math.Max(float64(slice.size), float64(compareSlice.size)))
	r := bytes.Compare(slice.data[:minLen], compareSlice.data[:minLen])
	if  r ==0 {
		if slice.size < compareSlice.size {
			return -1
		} else if slice.size > compareSlice.size {
			return 1
		}
	}

	return r
}

func(slice *Slice)  startWith(compareSlice *Slice) bool  {
	minLen := int(math.Max(float64(slice.size), float64(compareSlice.size)))
	return  (slice.Size() >= compareSlice.Size()) && (bytes.Compare(slice.data[:minLen], compareSlice.data[:minLen]) == 0)
}


