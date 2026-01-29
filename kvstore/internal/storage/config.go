package storage

import (
	"encoding/json"
	"os"
)

type Config struct {
	MemTableFlushSize          int `json:"memtable_flush_size"`
	WALBatchSize               int `json:"wal_batch_size"`
	WALFlushIntervalMs         int `json:"wal_flush_interval_ms"`
	CompactionSSTableThreshold int `json:"compaction_sstable_threshold"`
}

func LoadConfig(path string) Config {
	cfg := Config{
		MemTableFlushSize:          1000,
		WALBatchSize:               100,
		WALFlushIntervalMs:         10,
		CompactionSSTableThreshold: 4,
	}

	b, err := os.ReadFile(path)
	if err != nil {
		return cfg
	}

	_ = json.Unmarshal(b, &cfg)
	return cfg
}
