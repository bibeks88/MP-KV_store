package tests

import (
	"kvstore/internal/storage"
	"testing"
	"time"
)

func TestPutGet(t *testing.T) {
	e := storage.NewEngine("./data", storage.LoadConfig("./config/config.json"))
	e.Put("x", "1")
	time.Sleep(20 * time.Millisecond)
	v, ok := e.Get("x")
	if !ok || v != "1" {
		t.Fail()
	}
}
