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
	WritableFile* logfile_;
	uint64_t logfile_number_ GUARDED_BY(mutex_);
	log::Writer* log_;

	// Queue of writers.
	std::deque<Writer*> writers_ GUARDED_BY(mutex_);
	WriteBatch* tmp_batch_ GUARDED_BY(mutex_);

	SnapshotList snapshots_ GUARDED_BY(mutex_);

	// Set of table files to protect from deletion because they are
	// part of ongoing compactions.
	std::set<uint64_t> pending_outputs_ GUARDED_BY(mutex_);

	// Has a background compaction been scheduled or is running?
	bool background_compaction_scheduled_ GUARDED_BY(mutex_);

	ManualCompaction* manual_compaction_ GUARDED_BY(mutex_);

	VersionSet* const versions_ GUARDED_BY(mutex_);

	// Have we encountered a background error in paranoid mode?
	Status bg_error_ GUARDED_BY(mutex_);

	CompactionStats stats_[config::kNumLevels] GUARDED_BY(mutex_);
}

func (db *DB) Write(options *WriteOption, updates *WriteBatch) Status {
	w := newWriter(&db.mutex)
	w.batch = updates
	w.sync = options.sync
	w.done = false

MutexLock l(&mutex_);
writers_.push_back(&w);
while (!w.done && &w != writers_.front()) {
w.cv.Wait();
}
if (w.done) {
return w.status;
}

// May temporarily unlock and wait.
Status status = MakeRoomForWrite(updates == nullptr);
uint64_t last_sequence = versions_->LastSequence();
Writer* last_writer = &w;
if (status.ok() && updates != nullptr) {  // nullptr batch is for compactions
WriteBatch* write_batch = BuildBatchGroup(&last_writer);
WriteBatchInternal::SetSequence(write_batch, last_sequence + 1);
last_sequence += WriteBatchInternal::Count(write_batch);

// Add to log and apply to memtable.  We can release the lock
// during this phase since &w is currently responsible for logging
// and protects against concurrent loggers and concurrent writes
// into mem_.
{
mutex_.Unlock();
status = log_->AddRecord(WriteBatchInternal::Contents(write_batch));
bool sync_error = false;
if (status.ok() && options.sync) {
status = logfile_->Sync();
if (!status.ok()) {
sync_error = true;
}
}
if (status.ok()) {
status = WriteBatchInternal::InsertInto(write_batch, mem_);
}
mutex_.Lock();
if (sync_error) {
// The state of the log file is indeterminate: the log record we
// just added may or may not show up when the DB is re-opened.
// So we force the DB into a mode where all future writes fail.
RecordBackgroundError(status);
}
}
if (write_batch == tmp_batch_) tmp_batch_->Clear();

versions_->SetLastSequence(last_sequence);
}

while (true) {
Writer* ready = writers_.front();
writers_.pop_front();
if (ready != &w) {
ready->status = status;
ready->done = true;
ready->cv.Signal();
}
if (ready == last_writer) break;
}

// Notify new head of write queue
if (!writers_.empty()) {
writers_.front()->cv.Signal();
}

return status;
}