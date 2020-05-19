package db

const (
	// Zero is reserved for preallocated filtes
	kComparator = 1
	kLogNumber = 2

	kNextFileNumber = 3
	kLastSequence = 4
	kCompactPointer = 5
	kDeletedFIle = 6
	kNewFile = 7
	//8 was used for large value refs
	kPrevLogNumber = 9
)
