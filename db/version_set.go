package db

type VersionSet struct {
	dbName string
	options *Options
	tableCache *TableCache
	icmp *InternalKeyComparator
	nextFileNumber uint64
	manifestFileNumber uint64
	lastSequence uint64
	logNumber uint64
	prevLogNumber uint64

	descriptorFile *WritableFile
	descriptorLog *LogWriter
	dummyVersion Version
	current *Version

	compactPointer []string


}

func (versionSet *VersionSet) LastSequence() uint64{

}

func (versionSet *VersionSet) SetLastSequence(s uint64) {

}