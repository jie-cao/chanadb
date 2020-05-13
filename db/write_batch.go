package db

type WriteBatch struct {
	rep []byte
}

func (writeBatch *WriteBatch) Put(key *Slice, value *Slice) {
	SetCount(writeBatch, Count(writeBatch)+1)
	writeBatch.rep = append(writeBatch.rep, kTypeValue)
	PutLengthPrefixedSlice(writeBatch.rep, key)
	PutLengthPrefixedSlice(writeBatch.rep, value)
}

func (writeBatch *WriteBatch) Delete(key *Slice) {
	SetCount(writeBatch, Count(writeBatch)+1)
	writeBatch.rep = append(writeBatch.rep, kTypeDeletion)
	PutLengthPrefixedSlice(writeBatch.rep, key)
}

func (writeBatch *WriteBatch) Append(source *WriteBatch) {
	Append(writeBatch, source)
}

func (writeBatch *WriteBatch) Clear()  {
	writeBatch.rep = make([]byte, kHeader)
}
