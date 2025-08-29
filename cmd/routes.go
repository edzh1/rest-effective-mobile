package main

import (
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	// mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	w.WriteHeader(200)
	// 	w.Write([]byte(""))
	// })

	mux.HandleFunc("POST /subscriptions", app.subscriptionCreate)
	mux.HandleFunc("GET /subscriptions/{id}", app.subscriptionView)

	return mux
}
