package main

import (
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /subscriptions", app.subscriptionCreate)
	mux.HandleFunc("GET /subscriptions", app.subscriptionViewList)
	mux.HandleFunc("GET /subscriptions/total", app.subscriptionTotal)
	mux.HandleFunc("GET /subscriptions/{id}", app.subscriptionView)
	mux.HandleFunc("PUT /subscriptions/{id}", app.subscriptionUpdate)
	mux.HandleFunc("DELETE /subscriptions/{id}", app.subscriptionDelete)

	return app.recoverPanic(app.logRequest(commonHeaders(mux)))
}
