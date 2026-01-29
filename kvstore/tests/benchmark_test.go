package tests

import (
	"fmt"
	"kvstore/internal/storage"
	"sort"
	"sync"
	"testing"
	"time"
)

func BenchmarkConcurrentWrites(b *testing.B) {
	engine := storage.NewEngine("./data", storage.LoadConfig("./config/config.json"))

	concurrency := 50
	latencies := make([]time.Duration, 0, b.N)
	var mu sync.Mutex
	var wg sync.WaitGroup

	opsPerWorker := b.N / concurrency
	if opsPerWorker == 0 {
		opsPerWorker = 1
	}

	startAll := time.Now()
	wg.Add(concurrency)

	for i := 0; i < concurrency; i++ {
		go func(worker int) {
			defer wg.Done()
			for j := 0; j < opsPerWorker; j++ {
				start := time.Now()
				key := fmt.Sprintf("key-%d-%d", worker, j)
				val := fmt.Sprintf("value-%d-%d", worker, j)
				engine.Put(key, val)
				elapsed := time.Since(start)

				mu.Lock()
				latencies = append(latencies, elapsed)
				mu.Unlock()
			}
		}(i)
	}

	wg.Wait()
	totalTime := time.Since(startAll)

	if len(latencies) == 0 {
		b.Fatal("no operations recorded")
	}

	sort.Slice(latencies, func(i, j int) bool {
		return latencies[i] < latencies[j]
	})

	p99Index := int(float64(len(latencies)) * 0.99)
	if p99Index >= len(latencies) {
		p99Index = len(latencies) - 1
	}

	p99 := latencies[p99Index]

	b.ReportMetric(float64(len(latencies))/totalTime.Seconds(), "writes/sec")
	b.ReportMetric(float64(p99.Microseconds()), "p99_latency_us")
}
