package db

type LookupKey struct {
	// We construct a char array of the form:
	//    klength  varint32               <-- start_
	//    userkey  char[klength]          <-- kstart_
	//    tag      uint64
	//                                    <-- end_
	// The array is a suitable MemTable key.
	// The suffix starting with "userkey" can be used as an InternalKey.
	start  []byte
	kstart []byte
	end    []byte
	space  []byte // Avoid allocation for short keys
}

// Initialize *this for looking up user_key at a snapshot with
// the specified sequence number.
func newLookupKey(userKey *Slice, sequence uint64) *LookupKey{
	lookupKey := new(LookupKey)

	userKeySize := userKey.Size()
	var needed int = int(userKeySize + 13) // A conservative estimate
	var dst []byte
	if needed <= len(lookupKey.space) {
		dst = lookupKey.space
	} else {
		dst = make([]byte, needed)
	}
	lookupKey.start = dst
	dst = EncodeVarint32(dst, uint32(userKeySize+ 8))
	lookupKey.kstart = dst
	copy(dst, userKey.Data())
	dst = dst[userKeySize:]
	EncodeFixed64(dst, PackSequenceAndType(sequence, kTypeValue))
	dst = dst[8:]
	lookupKey.end = dst

	return lookupKey
}

func PackSequenceAndType(seq uint64, valueType byte) uint64 {
	return (seq << 8) | uint64(valueType)
}

// Return a key suitable for lookup in a MemTable.
func (lookupKey *LookupKey) memtableKey() *Slice {
	return NewSlice(lookupKey.start, len(lookupKey.start) - len(lookupKey.end))
}

// Return an internal key (suitable for passing to an internal iterator)
func (lookupKey *LookupKey)  internalKey() *Slice {
	return NewSlice(lookupKey.kstart, len(lookupKey.start) - len(lookupKey.end))
}

// Return the user key
func (lookupKey *LookupKey) userKey() *Slice{
	return NewSlice(lookupKey.kstart, len(lookupKey.start) - len(lookupKey.kstart) - 8)
}