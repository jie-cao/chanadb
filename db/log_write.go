package db

const (
	// Zero is reserved for preallocated filtes
	kZeroType byte = 0
	kFullType byte = 1

	//For fragments
	kFirstType byte = 2
	kMiddleType byte = 3
	kLastType byte = 4
)


const kHeaderSize = 4 + 2 + 1

type LogWriter struct {
	blockOffset int
	typeCrC [kLastType + 1]uint32
	dest *WritableFile
}

func (logWriter *LogWriter) AddRecord(slice *Slice) {
	data := slice.Data()
	left := slice.Size()

	var s Status
	var begin = true
	for {
		leftOver := kBlockSize - logWriter.blockOffset
		if leftOver < kHeaderSize {
			if leftOver > 0 {
				logWriter.dest.Append(NewSlice("\x00\x00\x00\x00\x00\x00", leftOver))
			}
			logWriter.blockOffset = 0
		}
		available := kBlockSize - logWriter.blockOffset - kHeaderSize
		var fragmentLength uint
		if leftOver < available {
			fragmentLength = left
		} else {
			fragmentLength = available
		}

		var recordType byte
		end := left == fragmentLength

		if begin && end {
			recordType = kFullType
		} else if begin {
			recordType = kFirstType
		} else if end {
			recordType = kLastType
		} else {
			recordType = kMiddleType
		}

		s = logWriter.EmitPhysicalRecord(recordType, data, fragmentLength)


		if s.OK() && left > 0 {
			break
		}
	}
}

func (logWriter *LogWriter) EmitPhysicalRecord(t byte, data []byte, length uint) Status{


}