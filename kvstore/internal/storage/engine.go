package storage

import "os"

type Engine struct {
	dataPath  string
	wal       *WAL
	memtable  *MemTable
	compactor *Compactor
	writeChan chan WriteRequest
	stopChan  chan struct{}
}

func NewEngine(path string, cfg Config) *Engine {
	os.MkdirAll(path, 0755)

	wal := NewWAL(path + "/wal.log")
	mem := NewMemTable()
	compactor := NewCompactor(path, cfg.CompactionSSTableThreshold)

	e := &Engine{
		dataPath:  path,
		wal:       wal,
		memtable:  mem,
		compactor: compactor,
		writeChan: make(chan WriteRequest, 10000),
		stopChan:  make(chan struct{}),
	}

	go e.walWriter(cfg)
	return e
}

func (e *Engine) Put(key, value string) error {
	done := make(chan error, 1)
	e.writeChan <- WriteRequest{Key: key, Value: value, Done: done}
	return <-done
}

func (e *Engine) Get(key string) (string, bool) {
	if v, ok := e.memtable.Get(key); ok {
		return v, true
	}
	return ReadFromSSTables(key)
}
