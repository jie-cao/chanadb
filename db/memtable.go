package db

type MemTable struct {
	table *SkipList
	ref   int
	arena Arena
}

const (
	kTypeDeletion byte = 0
	kTypeValue    byte = 1
)

func (memtable *MemTable) Add(sequenceNumber uint64, valueType byte, key *Slice, value *Slice) {

	// Format of an entry is concatenation of:
	//  key_size     : varint32 of internal_key.size()
	//  key bytes    : char[internal_key.size()]
	//  value_size   : varint32 of value.size()
	//  value bytes  : char[value.size()]
	keySize := key.size
	valSize := value.size
	internalKeySize := int(keySize) + 8
	encodedLen := VarintLength(uint64(internalKeySize)) + int(internalKeySize) + VarintLength(uint64(valSize)) + int(valSize)
	buffer := memtable.arena.Allocate(encodedLen)
	originBuf := buffer
	buffer = EncodeVarint32(buffer, uint32(internalKeySize))
	copy(buffer, key.data)
	buffer = buffer[keySize:]
	EncodeFixed64(buffer, uint64((sequenceNumber << 8) | uint64(valueType)))
	buffer = buffer[8:]
	EncodeVarint32(buffer, uint32(valSize))
	copy(buffer, value.data)
	memtable.table.Insert(string(originBuf))
}
