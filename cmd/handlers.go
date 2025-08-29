package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/edzh1/rest-effective-mobile/internal/models"
	"github.com/google/uuid"
)

type subscriptionCreateBody struct {
	UserId      string `json:"user_id"`
	ServiceName string `json:"service_name"`
	Price       int    `json:"price"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date,omitempty"`
}

func (app *application) subscriptionView(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid UUID format", http.StatusBadRequest)
		return
	}

	subscription, err := app.subscriptions.Get(id)

	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")

	data := struct {
		Status       string              `json:"status"`
		Subscription models.Subscription `json:"subscription"`
	}{
		Status:       "success",
		Subscription: subscription,
	}

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Write(jsonBytes)
}

func (app *application) subscriptionCreate(w http.ResponseWriter, r *http.Request) {
	reader := r.Body
	defer reader.Close()

	body, err := io.ReadAll(reader)
	// TODO check for server errors
	// TODO add validator
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var reqBody subscriptionCreateBody

	err = json.Unmarshal(body, &reqBody)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	startDate, err := time.Parse("2006-01-02", reqBody.StartDate)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var endDate *time.Time
	if reqBody.EndDate != "" {
		t, err := time.Parse("2006-01-02", reqBody.EndDate)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		endDate = &t
	}

	id, err := app.subscriptions.Insert(reqBody.UserId, reqBody.ServiceName, reqBody.Price, startDate, endDate)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	data := struct {
		Status string    `json:"status"`
		Id     uuid.UUID `json:"id"`
	}{
		Status: "success",
		Id:     id,
	}

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Write(jsonBytes)
}
