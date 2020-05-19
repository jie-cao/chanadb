package db

type Compaction struct {
	level int
	maxOutFileSize uint64
	inputVersion *Version
	edit VersionEdit
	inputs [2]FileMetaData
	grandParents FileMetaData
	grandParentIndex uint
	seenKey bool
	overLappedBytes int64
	levels []int
}