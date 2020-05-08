package db

const (
	kOK              byte = 0
	kNotFound        byte = 1
	kCorruption      byte = 2
	kNotSupported    byte = 3
	kInvalidArgument byte = 4
	kIOError         byte = 5
)

type Status struct {
	state byte
}

// Return error status of an appropriate type.
func  NotFound(msg *Slice, msg2 *Slice) *Status {
	status := new(Status)
	status.state = kNotFound
	return status
}
