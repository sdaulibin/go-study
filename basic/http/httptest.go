package main

import (
	"log"
	"net/http"
)

type server string

func (http *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)
	w.Write([]byte("Hello World"))
}

func main()  {
	var server server
	http.ListenAndServe("127.0.0.1:9999",&server)
}