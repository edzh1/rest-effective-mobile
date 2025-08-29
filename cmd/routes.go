package main

import (
	"net/http"
	"os"

	_ "github.com/edzh1/rest-effective-mobile/docs"
	"github.com/justinas/alice"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	mux.Handle("POST /subscriptions", standard.ThenFunc(app.subscriptionCreate))
	mux.Handle("GET /subscriptions", standard.ThenFunc(app.subscriptionViewList))
	mux.Handle("GET /subscriptions/total", standard.ThenFunc(app.subscriptionTotal))
	mux.Handle("GET /subscriptions/{id}", standard.ThenFunc(app.subscriptionView))
	mux.Handle("PUT /subscriptions/{id}", standard.ThenFunc(app.subscriptionUpdate))
	mux.Handle("DELETE /subscriptions/{id}", standard.ThenFunc(app.subscriptionDelete))

	if os.Getenv("ENV") == "dev" {
		mux.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)
	}

	return mux
}
