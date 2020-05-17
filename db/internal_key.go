package db

import (
	"bytes"
)

type InternalKey struct {
	rep string
}

func newInternalKey(userKey *Slice, sequenceName uint64, valueType byte)  *InternalKey{

}

func AddInternalKey(result string, ) *InternalKey {

}

type ParsedInternalKey struct {
	userKey Slice
	sequence uint64
	valueType byte
}

func (parsedInternalKey *ParsedInternalKey) DebugString() string {
	var debugStringBuffer bytes.Buffer
	debugStringBuffer.WriteString("'")
	debugStringBuffer.WriteString(parsedInternalKey.EscapeString())
	debugStringBuffer.WriteString("' @ ")
	debugStringBuffer.WriteString(string(parsedInternalKey.sequence))
	debugStringBuffer.WriteString(" : ")
	debugStringBuffer.WriteString(string(parsedInternalKey.valueType))

	return debugStringBuffer.String()
}

func (parsedInternalKey *ParsedInternalKey) EscapeString() string {

}