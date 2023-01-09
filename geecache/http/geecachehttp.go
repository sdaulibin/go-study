package http

import (
	"fmt"
	"geecache/cache"
	"log"
	"net/http"
)

var db = map[string]string{
	"tom":  "687",
	"jack": "612",
	"sam":  "543",
}

func main() {
	cache.NewGroup("scores", 2<<10, cache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key ", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))
	addr := "127.0.0.1:9999"
	peers := cache.NewHttpPool(addr)
	log.Println("geecache is running at ", addr)
	log.Fatal(http.ListenAndServe(addr, peers))
}
