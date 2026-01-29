package main

import (
	"fmt"
	"kvstore/internal/storage"
	"time"
)

func main() {
	cfg := storage.LoadConfig("./config/config.json")
	engine := storage.NewEngine("./data", cfg)

	fmt.Printf("CONFIG: %+v\n", cfg)

	engine.Put("a", "1")
	engine.Put("b", "2")
	engine.Put("a", "3")

	time.Sleep(50 * time.Millisecond)

	v, _ := engine.Get("a")
	fmt.Println(v)

	select {}
}
