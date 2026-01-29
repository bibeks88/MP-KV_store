package storage

import (
	"fmt"
	"time"
)

type WriteRequest struct {
	Key   string
	Value string
	Done  chan error
}

func (e *Engine) walWriter(cfg Config) {

	batch := make([]WriteRequest, 0, cfg.WALBatchSize)
	ticker := time.NewTicker(time.Duration(cfg.WALFlushIntervalMs) * time.Millisecond)
	defer ticker.Stop()

	flush := func() {
		if len(batch) == 0 {
			return
		}

		for _, r := range batch {
			e.wal.Append(r.Key, r.Value)
		}

		for _, r := range batch {
			e.memtable.Put(r.Key, r.Value)
			r.Done <- nil
		}

		if e.memtable.Size() >= cfg.MemTableFlushSize {
			file := WriteSSTable(e.dataPath, e.memtable.Snapshot())
			fmt.Println("SSTABLE CREATED:", file)
			e.memtable.Clear()
			e.compactor.Register(file)
		}

		batch = batch[:0]
	}

	for {
		select {
		case r := <-e.writeChan:
			batch = append(batch, r)
			if len(batch) >= cfg.WALBatchSize {
				flush()
			}

		case <-ticker.C:
			flush()

		case <-e.stopChan:
			flush()
			return
		}
	}
}
