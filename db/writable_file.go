package db

import (
	"log"
	"math"
	"os"
)

const (
	kWritableFileBufferSize  = 65536
)

type WritableFile struct {
	fileName string
	dirName string
	isManifest bool
	pos int
	file *os.File
	buf []byte
}

func (writableFile *WritableFile) Sync() Status {
	s := writableFile.SyncDirIfManifest()
	if !s.OK() {
		return s
	}
	s = writableFile.FlushBuffer()
	if !s.OK() {
		return s
	}

	return writableFile.SyncFD()
}

func (writableFile *WritableFile) Append(data *Slice) Status {
	writeSize := data.Size()
	writeData := data.Data()

	// Fit as much as possible into buffer
	copySize := int(math.Max(float64(writeSize), float64(kWritableFileBufferSize - writableFile.pos)))


	writeData = writeData[copySize:]
	writeSize -= uint(copySize)

	writableFile.pos += copySize

	if writeSize  == 0 {
		return *(StatusOK())
	}

	// Can't fit in buffer, so need to do at least one write
	s := writableFile.FlushBuffer()
	if !s.OK() {
		return s
	}

	if writeSize < kWritableFileBufferSize {
		copy(writableFile.buf, writeData[:writeSize])
		writableFile.pos = int(writeSize)
		return *(StatusOK())
	}

	return writableFile.WriteUnbuffered(writeData, int(writeSize))
}

func (writableFile *WritableFile) Flush() Status {
	return writableFile.FlushBuffer()
}

func (writableFile *WritableFile) SyncDirIfManifest() Status{
	s := Status{}
	if !writableFile.isManifest {
		return s
	}

	file, err := os.Open(writableFile.dirName) // For read access.
	if err != nil {
		log.Fatal(err)
	}

	err = file.Sync()
	if err != nil {
		log.Fatal(err)
	}

	err = file.Close()
	if err != nil {
		log.Fatal(err)
	}
	return s
}

func (writableFile *WritableFile) FlushBuffer() Status{
	s := writableFile.WriteUnbuffered(writableFile.buf, writableFile.pos)
	return s
}

func (writableFile *WritableFile) SyncFD() Status{
	err := writableFile.file.Sync()
	if err != nil {
		log.Fatal(err)
	}
	return Status{}
}

func (writableFile *WritableFile) WriteUnbuffered(data []byte, size int) Status{
	for size > 0 {
		writeResult, err := writableFile.file.Write(data)
		if err != nil {

		}
		data = data[writeResult:]
		size -= writeResult
	}

	return Status{}
}