package db

type WritableFile struct {

}

func (writableFile *WritableFile) Sync() Status {

}

func (writableFile *WritableFile) Append(slice *Slice) Status {

}

func (writableFile *WritableFile) Flush() Status {

}