package db

import (
	"os"
	"time"
	"strings"
)

func SleepForMicroseconds(micros uint)  {
	time.Sleep(time.Duration(micros) * time.Microsecond)
}

func NewWritableFile(fileName string, result **WritableFile) *Status {
	var file, err = os.OpenFile(fileName, os.O_RDWR, 0644)
	if err != nil {

	}

	writableFilePtr := new(WritableFile)
	writableFilePtr.file = file
	writableFilePtr.pos = 0
	writableFilePtr.isManifest = isManifest(fileName)
	writableFilePtr.fileName = fileName
	writableFilePtr.dirName = dirName(fileName)
	*result = writableFilePtr
	return StatusOK()
}

func dirName(fileName string) string  {
	separatorIdx := strings.LastIndex(fileName, "/")
	if separatorIdx {
		return "."
	}
	 return fileName[0:separatorIdx]
}

func isManifest(fileName string) bool {
	return strings.HasPrefix(Basement(fileName), "MANIFEST")
}

func Basement(fileName string) string  {

}


