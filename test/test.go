package main

import (
	"fmt"
	"github.com/nickham-su/go_memory_cache"
)

func main() {
	cache := go_memory_cache.NewMemoryCache(true)
	cache.SetMaxMemory("100MB")
	cache.Set("int", 1)
	cache.Set("bool", false)
	cache.Set("data", map[string]interface{}{"a": 1})
	fmt.Println(cache.Get("int"))      // 1 true
	fmt.Println(cache.Get("bool"))     // false true
	fmt.Println(cache.Get("data"))     // map[a:1] true
	fmt.Println("keys:", cache.Keys()) // keys: 3
	cache.Del("int")
	fmt.Println(cache.Get("int"))      // <nil> false
	fmt.Println("keys:", cache.Keys()) // keys: 2
	cache.Flush()
	fmt.Println("keys:", cache.Keys()) // keys: 0
}
