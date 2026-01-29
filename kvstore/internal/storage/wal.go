package storage

import (
	"os"
	"sync"
)

type WAL struct {
	file *os.File
	mu   sync.Mutex
}

func NewWAL(path string) *WAL {
	f, _ := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	return &WAL{file: f}
}

func (w *WAL) Append(key, value string) {
	w.mu.Lock()
	w.file.WriteString(key + ":" + value + "\n")
	w.file.Sync()
	w.mu.Unlock()
}
