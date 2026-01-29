# KV Store (Go)

A persistent, disk-backed key/value storage engine implemented in Go. It is
optimized for write-heavy workloads and designed for predictable latency,
crash safety, and scalability.

---

## Features

- Durable writes via Write-Ahead Log (WAL)
- Low-latency in-memory writes with MemTable
- Immutable SSTables for efficient on-disk storage
- Background compaction to reduce read amplification
- Channel-based batching for high throughput
- Single-writer WAL model for correctness
- Concurrent reads and safe background operations

---

## Requirements

- Go 1.20+

---

## Quick Start

Clone and prepare the project:

```bash
git clone <repository-url>
cd kvstore
go mod tidy
```

Run the server (development/demo):

```bash
go run cmd/server/main.go
```

Alternatively, download the repository ZIP, extract it, and run the same
commands from the extracted directory.

---

## Configuration

The engine uses a small JSON configuration file at `config/config.json`.
Adjust these values to tune flushing, batching, and compaction behavior.

Example `config/config.json`:

```json
{
  "memtable_flush_size": 1000,
  "wal_batch_size": 100,
  "wal_flush_interval_ms": 10,
  "compaction_sstable_threshold": 4
}
```

Notes:

- Lower `memtable_flush_size` causes more frequent SSTable creation (useful for demos).
- `compaction_sstable_threshold` controls when background compaction runs.

---

## Usage & Behavior Examples

### WAL-only Behavior (Small Dataset)

When the dataset fits in memory and the MemTable does not reach the flush
threshold, writes are durable thanks to the WAL and no SSTables are created.

Steps:

1.  Run: `go run cmd/server/main.go`
2.  Verify that `data/wal.log` exists and no `sst_*.db` files are present.
    Expected: Low-latency writes and only `wal.log` in the data directory.

### Disk Spill & SSTable Creation (Large Dataset)

When MemTable reaches the configured flush size it is flushed to an immutable
SSTable on disk. After multiple SSTables, the background compaction merges
SSTables into a compacted file.

Demo steps:

1. Reduce `memtable_flush_size` and `compaction_sstable_threshold` in `config/config.json` to low values (e.g. 10 and 2).
2. Run loads (concurrent writes) to trigger flushes and compaction.
3. Inspect `tests/data/` for artifacts: `wal.log`, `sst_*.db`, `compacted.db`.

Example benchmark output (indicative):

```
BenchmarkConcurrentWrites-20    5569    211507 ns/op    12024 p99_latency_us    4719 writes/sec
```

---

## Testing & Benchmarks

Run unit tests:

```bash
go test ./...
```

Run benchmark tests:

```bash
go test -bench=. -run=^$
```

---

## Project Layout

- `cmd/server/` — server entry point
- `internal/storage/` — core engine (WAL, MemTable, SSTable, compaction)
- `config/config.json` — runtime configuration
- `data/` — default runtime data directory
- `tests/` — tests and benchmark inputs

---
