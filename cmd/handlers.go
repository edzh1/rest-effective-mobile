package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/edzh1/rest-effective-mobile/internal/models"
	"github.com/google/uuid"
)

type subscriptionCreateBody struct {
	UserID      string `json:"user_id"`
	ServiceName string `json:"service_name"`
	Price       int    `json:"price"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date,omitempty"`
}

type subscriptionUpdateBody struct {
	ID uuid.UUID `json:"id"`
	subscriptionCreateBody
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
			app.serverError(w, r, err)
		}
		return
	}

	data := struct {
		Status       string              `json:"status"`
		Subscription models.Subscription `json:"subscription"`
	}{
		Status:       "success",
		Subscription: subscription,
	}

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
}

func (app *application) subscriptionViewList(w http.ResponseWriter, r *http.Request) {
	var filter models.SubscriptionFilter
	query := r.URL.Query()

	if userIDStr := query.Get("user_id"); userIDStr != "" {
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			http.Error(w, "Invalid user_id format", http.StatusBadRequest)
			return
		}
		filter.UserID = &userID
	}

	if serviceName := query.Get("service_name"); serviceName != "" {
		filter.ServiceName = &serviceName
	}

	if startDateStr := query.Get("start_date"); startDateStr != "" {
		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			http.Error(w, "Invalid start_date format", http.StatusBadRequest)
			return
		}
		filter.StartDate = &startDate
	}

	if endDateStr := query.Get("end_date"); endDateStr != "" {
		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			http.Error(w, "Invalid end_date format", http.StatusBadRequest)
			return
		}
		filter.EndDate = &endDate
	}

	if pageStr := query.Get("page"); pageStr != "" {
		page, err := strconv.Atoi(pageStr)
		if err != nil || page == 0 {
			http.Error(w, "Invalid page format", http.StatusBadRequest)
			return
		}
		filter.Page = &page
	}

	subscriptions, err := app.subscriptions.List(filter)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	data := struct {
		Status        string                `json:"status"`
		Subscriptions []models.Subscription `json:"subscriptions"`
	}{
		Status:        "success",
		Subscriptions: subscriptions,
	}

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
}

func (app *application) subscriptionCreate(w http.ResponseWriter, r *http.Request) {
	reader := r.Body
	defer reader.Close()

	body, err := io.ReadAll(reader)
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

	id, err := app.subscriptions.Insert(reqBody.UserID, reqBody.ServiceName, reqBody.Price, startDate, endDate)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	data := struct {
		Status string    `json:"status"`
		Id     uuid.UUID `json:"id"`
	}{
		Status: "success",
		Id:     id,
	}

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
}

func (app *application) subscriptionUpdate(w http.ResponseWriter, r *http.Request) {
	reader := r.Body
	defer reader.Close()

	body, err := io.ReadAll(reader)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var reqBody subscriptionUpdateBody

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

	id, err := app.subscriptions.Update(reqBody.ID, reqBody.UserID, reqBody.ServiceName, reqBody.Price, startDate, endDate)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	data := struct {
		Status string    `json:"status"`
		Id     uuid.UUID `json:"id"`
	}{
		Status: "success",
		Id:     id,
	}

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
}

func (app *application) subscriptionDelete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid UUID format", http.StatusBadRequest)
		return
	}

	err = app.subscriptions.Delete(id)

	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	data := struct {
		Status string `json:"status"`
	}{
		Status: "success",
	}

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
}
