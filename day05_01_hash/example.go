package main

import (
	_ "github.com/icattlecoder/godaemon"
	"net/http"
	"log"
)

func main() {
	mux:=http.NewServeMux()
	mux.HandleFunc("/index", func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte("hello,golang!\n"))
	})
	log.Fatalln(http.ListenAndServe(":7070",mux))
}

