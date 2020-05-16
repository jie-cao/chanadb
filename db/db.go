package db

import "sync"

type DB struct {
	dbName string
	ownsCache bool
	ownsInfoLog bool
	logFileName uint64
	// State below is protected by mutex_
	mutex sync.Mutex
	shuttingDown bool
	mem *MemTable
	imm *MemTable
	hasImm bool
	logFileNumber uint64
	seed uint32 // For sampling.
	backgroundWorkFinishedSignal sync.Cond
	versions *VersionSet
	// Queue of writers.
	writers []*Writer
	tmpBatch *WriteBatch
	logWriter *LogWriter
	logfile *WritableFile
	bgError Status
}

func (db *DB) MakeRoomForWrite(force bool) Status  {
	allowDelay := !force
	var s Status
	for {
		if !db.bgError.OK() {
			s = db.bgError
			break
		} else if allowDelay && db.versions.

	}

}

func (db *DB) Write(options *WriteOption, updates *WriteBatch) Status {
	w := newWriter(db.mutex)
	w.batch = updates
	w.sync = options.sync
	w.done = false
	db.mutex.Lock()
	db.writers = append(db.writers, w)
	for !w.done && w != db.writers[0] {
		w.cv.Wait()
	}

	if w.done {
		return w.status
	}

	// May temporarily unlock and wait.
	status := db.MakeRoomForWrite(updates == nil)
	lastSequence := db.versions.LastSequence()
	lastWriter := w
	if status.OK() && updates != nil { // nullptr batch is for compactions
		writeBatch := db.BuildBatchGroup(&lastWriter)
		SetSequence(writeBatch, lastSequence+1)
		lastSequence += uint64(Count(writeBatch))

		// Add to log and apply to memtable.  We can release the lock
		// during this phase since &w is currently responsible for logging
		// and protects against concurrent loggers and concurrent writes
		// into mem_.
		{
			db.mutex.Unlock()
			db.logWriter.AddRecord(Contents(writeBatch))
			syncError := false
			if status.OK() && options.sync {
				status = db.logfile.Sync()
				if !status.OK() {
					syncError = true
				}
			}

			if status.OK() {
				status = InsertInto(writeBatch, db.mem)
			}

			db.mutex.Lock()
			if syncError {
				// The state of the log file is indeterminate: the log record we
				// just added may or may not show up when the DB is re-opened.
				// So we force the DB into a mode where all future writes fail.
				db.RecordBackgroundError(status)
			}
		}

		if writeBatch == db.tmpBatch {
			db.tmpBatch.Clear()
		}
		db.versions.SetLastSequence(lastSequence)

	}

	for {
		ready := db.writers[0]
		db.writers = db.writers[1:]
		if ready != w {
			ready.status = status
			ready.done = true
			ready.cv.Signal()
		}
		if ready == lastWriter {
			break
		}
	}

	// Notify new head of write queue
	if len(db.writers) != 0 {
		db.writers[0].cv.Signal()
	}

	return status
}

func (db *DB) RecordBackgroundError(status Status)  {
	if db.bgError.OK() {
		db.bgError = status
		db.backgroundWorkFinishedSignal.Broadcast()
	}
}

func (db *DB) BuildBatchGroup(lastWriter **Writer) *WriteBatch{
	first := db.writers[0]
	result := first.batch
	size := ByteSize(first.batch)

	maxSize := 1 << 20
	if size <= (128 << 10) {
		maxSize = size + (128 << 10)
	}

	*lastWriter = first
	writerIdx := 0
	for ; writerIdx != len(db.writers); writerIdx++ {
		w := db.writers[writerIdx]
		if w.sync && !first.sync {
			// do not include a sync writer into a batch handled by a non-sync writer
			break
		}

		if w.batch != nil {
			size += ByteSize(w.batch)
			if size > maxSize {

				// Do not make batch too big
				break
			}

			// Append to *result
			if result == first.batch {

				// Switch to temporary batch instead of disturbing callers's batch
				result = db.tmpBatch
				Append(result, first.batch)
			}
			Append(result, w.batch)
		}
		*lastWriter = w
	}

	return result
}