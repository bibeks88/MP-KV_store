package storage

import (
	"os"
	"sort"
	"strings"
)

func WriteSSTable(basePath string, data map[string]string) string {
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	path := basePath + "/sst_" + RandString() + ".db"
	f, _ := os.Create(path)

	for _, k := range keys {
		f.WriteString(k + ":" + data[k] + "\n")
	}
	f.Close()

	return path
}

func ReadFromSSTables(key string) (string, bool) {
	files, _ := os.ReadDir("data")
	for i := len(files) - 1; i >= 0; i-- {
		b, err := os.ReadFile("data/" + files[i].Name())
		if err != nil {
			continue
		}
		for _, line := range strings.Split(string(b), "\n") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 && parts[0] == key {
				return parts[1], true
			}
		}
	}
	return "", false
}
