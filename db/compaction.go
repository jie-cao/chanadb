package db

type Compaction struct {
	level int
	maxOutFileSize uint64
	inputVersion Ver
}