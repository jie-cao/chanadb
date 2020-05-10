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

func (s *SkipList) Get(key *LookupKey, value *string, status *Status) (bool, *Status) {
	memKey := key.memtableKey()
	iter := newIterator(s)
	iter.Seek(string(memKey.Data()))
	if iter.Valid() {
		// entry format is:
		//    klength  varint32
		//    userkey  char[klength]
		//    tag      uint64
		//    vlength  varint32
		//    value    char[vlength]
		// Check that it belongs to same user key.  We do not check the
		// sequence number since the Seek() call above should have skipped
		// all entries with overly large sequence numbers.
		entry := iter.Key()
		var keyLength uint32
		keyPtr := GetVarint32Ptr([]byte(entry), 5, &keyLength)
		if NewSlice(keyPtr, int(keyLength - 8)).Compare(key.userKey()) > 0 {
			// Correct user key
			tag := DecodeFixed64(keyPtr[keyLength - 8:])
			switch byte(tag & 0xff) {
			case kTypeValue:
				{
					v := GetLengthPrefixedSlice(keyPtr[keyLength:])
					stringValue := string(v.Data()[0:v.Size()])
					value = &stringValue
					return true, nil
				}
			case kTypeDeletion:
				status := NotFound(nil, nil)
				return true, status
			}
		}
	}
	return false, nil
}
