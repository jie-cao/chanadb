package db

const kHeader = 12

type WriteBatchInternal struct {

}

func Count(b *WriteBatch) int {
	data := []byte(b.rep)
	return int(DecodeFixed32(data[8:]))
}

func SetCount(b *WriteBatch, n int) {
	EncodeFixed32([]byte(b.rep[8:]), uint32(n))
}

func Append(dst *WriteBatch, src *WriteBatch) {
	SetCount(dst, Count(dst) + Count(src))
	dst.rep = append(dst.rep, src.rep[kHeader:kHeader + len(src.rep) - kHeader]...)
}

func ByteSize(batch *WriteBatch) int {
	return len(batch.rep)
}

func SetSequence(batch *WriteBatch, seq uint64) {
	EncodeFixed64(batch.rep, seq)
}

func Contents(batch *WriteBatch) *Slice{
	return NewSliceFromBytes(batch.rep)
}

func InsertInto(batch *WriteBatch, memtTable *MemTable) Status {

}