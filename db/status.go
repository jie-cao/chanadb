package db

const (
	kOK              byte = 0
	kNotFount        byte = 1
	kCorruption      byte = 2
	kNotSupported    byte = 3
	kInvalidArgument byte = 4
	kIOError         byte = 5
)

type Status struct {
	state byte
}

// Kind gets the kind of the datum.
func (d *Datum) Kind() byte {
	return d.k
}
