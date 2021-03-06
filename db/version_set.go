package db

type FileMetaData struct {
	refs int
	allowedSeeks int
	number uint64
	fileSize uint64
	smallest InternalKey
	largest InternalKey
}

type Version struct {
	vSet *VersionSet //所属的VersionSet
	next *Version // 在所处的VersionSet中，后一个Version指针
	prev *Version // 在所处的VersionSet中，前一个Version指针
	refs int //当前版本引用数，只有为了0才会被释放
	files [][]FileMetaData //当前版本中每一个level所包含的文件信息
}


type VersionEdit struct {
	comparator *Slice //比较器的名称
	logNumber uint64 //日志编号
	prevLogNumber uint64 //上一个日志编号
	nextFileNumber uint64 //下一个文件编号
	lastSequence uint64 //最后的序列号
	hasComparator bool //是否有比较器
	hasLogNumber bool
	hasPrevLogNumber bool
	hasNextFileNumber bool
	hasLastSequence bool

	//压缩点 <层数， InternalKey键>
	compactPointer []CompactPointer
	deletedFileSet map[DeletedFile]bool //删除文件标志
	newFiles []NewFile //新增的文件，记录了level和FileMetaData
}

func (versionEdit *VersionEdit) AddFile(level int, file uint64, fileSize uint64, smallest InternalKey, largest InternalKey)  {
	f := new(FileMetaData)
	f.number = file
	f.fileSize = fileSize
	f.smallest = smallest
	f.largest = largest
	newFile := NewFile{level , *f}
	versionEdit.newFiles = append(versionEdit.newFiles, newFile)
}

func (versionEdit *VersionEdit) RemoveFile(level int, file uint64) {
	deletedFile := DeletedFile{level, file}
	versionEdit.deletedFileSet[deletedFile] = true
}

func (versionEdit *VersionEdit) EncodeTo(dest string) {
	if versionEdit.hasComparator {
		PutVarint32([]byte(dest), kComparator)
		PutLengthPrefixedSlice([]byte(dest), versionEdit.comparator)
	}
}


type CompactPointer struct{
	index int
	internalKey InternalKey
}

type NewFile struct {
	level int //层数
	fileMetaDat FileMetaData //文件描述信息
}

//需要删除的文件
type DeletedFile struct {
	level int //层数
	fileNumber uint64 //file number
}

type VersionSet struct {
	dummyVersions Version // head of circular doubly-linked list of versions
	current *Version // == dummyVersion.prev
	dbName string
	nextFileNumber uint64
	manifestFileNumber uint64
	lastSequence uint64
	logNumber uint64
	prevLogNumber uint64
	descriptorFile *WritableFile
	descriptorLog *LogWriter
	compactPointer []string
}

func (versionSet *VersionSet) LastSequence() uint64{
	return versionSet.lastSequence
}

func (versionSet *VersionSet) SetLastSequence(s uint64) {
	versionSet.lastSequence = s
}

func (versionSet *VersionSet) NumLevelFiles(level int) int {
	return len(versionSet.current.files[level])
}

func (versionSet *VersionSet) NewFileNumber() uint64 {
	versionSet.nextFileNumber++
	return versionSet.nextFileNumber
}