package db

import "sync"

type Writer struct {
	status Status
	batch *WriteBatch
	sync bool
	done bool
	cv sync.Cond
}

func newWriter(mu sync.Mutex) *Writer {
	writer := new(Writer)
	writer.batch = nil
	writer.sync = false
	writer.cv = sync.NewCond(mu)

	return writer
}

