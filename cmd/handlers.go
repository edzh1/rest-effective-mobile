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
	UserID      string `json:"user_id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	ServiceName string `json:"service_name" example:"Yandex Plus"`
	Price       int    `json:"price" example:"400"`
	StartDate   string `json:"start_date" example:"2025-07-01"`
	EndDate     string `json:"end_date,omitempty" example:""`
}

type subscriptionUpdateBody struct {
	subscriptionCreateBody
}

type IDResponse struct {
	ID uuid.UUID `json:"id" example:"a3509860-d66f-4be4-8984-0b7a15b8f10c"`
}

type TotalResponse struct {
	Total int `json:"total" example:"100500"`
}

// subscriptionView godoc
// @Summary Get subscription by ID
// @Description Get a single subscription by its ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "Subscription ID" Format(uuid)
// @Success 200 {object} models.Subscription
// @Failure 400 {string} string "Invalid UUID format"
// @Failure 404 {string} string "404 page not found"
// @Router /subscriptions/{id} [get]
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
		Subscription models.Subscription `json:"subscription"`
	}{
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

// subscriptionViewList godoc
// @Summary List subscriptions with filters
// @Description Get list of subscriptions with optional filters
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param user_id query string false "User ID filter" Format(uuid)
// @Param service_name query string false "Service name filter"
// @Param start_date query string false "Start date filter (YYYY-MM-DD)" Format(date)
// @Param end_date query string false "End date filter (YYYY-MM-DD)" Format(date)
// @Param page query int false "Page number" minimum(1) maximum(100)
// @Success 200 {object} []models.Subscription
// @Failure 400 {string} string "Invalid parameter format"
// @Router /subscriptions [get]
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
		Subscriptions []models.Subscription `json:"subscriptions"`
	}{
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

// subscriptionTotal godoc
// @Summary Calculate total subscription cost
// @Description Calculate total cost of subscriptions for a period with filters
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param user_id query string false "User ID filter" Format(uuid)
// @Param service_name query string false "Service name filter"
// @Param start_date query string false "Period start (YYYY-MM-DD)" Format(date)
// @Param end_date query string false "Period end (YYYY-MM-DD)" Format(date)
// @Success 200 {object} TotalResponse
// @Failure 400 {string} string "Invalid parameter format"
// @Router /subscriptions/total [get]
func (app *application) subscriptionTotal(w http.ResponseWriter, r *http.Request) {
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

	total, err := app.subscriptions.CountTotal(filter)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	data := struct {
		Total int `json:"total"`
	}{
		Total: total,
	}

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
}

// subscriptionCreate godoc
// @Summary Create new subscription
// @Description Create a new subscription record
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body subscriptionCreateBody true "Subscription data"
// @Success 200 {object} IDResponse "{"id": "a3509860-d66f-4be4-8984-0b7a15b8f10c"}"
// @Failure 400 {string} string "Bad Request"
// @Router /subscriptions [post]
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

	data := IDResponse{
		ID: id,
	}

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
}

// subscriptionUpdate godoc
// @Summary Update subscription
// @Description Update an existing subscription
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "Subscription ID" Format(uuid)
// @Param subscription body subscriptionUpdateBody true "Updated subscription data"
// @Success 200 {object} IDResponse
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Not Found"
// @Router /subscriptions/{id} [put]
func (app *application) subscriptionUpdate(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid UUID format", http.StatusBadRequest)
		return
	}

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

	_, err = app.subscriptions.Update(id, reqBody.UserID, reqBody.ServiceName, reqBody.Price, startDate, endDate)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	data := IDResponse{
		ID: id,
	}

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
}

// subscriptionDelete godoc
// @Summary Delete subscription
// @Description Delete a subscription by ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "Subscription ID" Format(uuid)
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Invalid UUID format"
// @Failure 404 {string} string "Not Found"
// @Router /subscriptions/{id} [delete]
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

	w.WriteHeader(http.StatusOK)
}
