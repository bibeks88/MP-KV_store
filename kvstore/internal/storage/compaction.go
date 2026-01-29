package storage

import (
	"os"
	"sort"
	"strings"
	"sync"
)

type Compactor struct {
	path      string
	threshold int
	files     []string
	mu        sync.Mutex
}

func NewCompactor(path string, threshold int) *Compactor {
	return &Compactor{
		path:      path,
		threshold: threshold,
	}
}

func (c *Compactor) Register(file string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.files = append(c.files, file)

	if len(c.files) >= c.threshold {
		c.compact()
	}
}

func (c *Compactor) compact() {
	merged := make(map[string]string)
	for i := len(c.files) - 1; i >= 0; i-- {
		b, _ := os.ReadFile(c.files[i])
		for _, line := range strings.Split(string(b), "\n") {
			p := strings.SplitN(line, ":", 2)
			if len(p) == 2 {
				if _, ok := merged[p[0]]; !ok {
					merged[p[0]] = p[1]
				}
			}
		}
		os.Remove(c.files[i])
	}

	keys := make([]string, 0, len(merged))
	for k := range merged {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	f, _ := os.Create(c.path + "/compacted.db")
	for _, k := range keys {
		f.WriteString(k + ":" + merged[k] + "\n")
	}
	f.Close()

	c.files = nil
}
