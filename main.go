package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("hello"))
	})

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal("can't start server " + err.Error())
	}
}
