package models

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID          uuid.UUID  `json:"id"`
	UserId      uuid.UUID  `json:"user_id"`
	ServiceName string     `json:"service_name"`
	Price       int        `json:"price"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date,omitempty"`
}

type SubscriptionModel struct {
	DB *sql.DB
}

func (m *SubscriptionModel) Get(id uuid.UUID) (Subscription, error) {
	var s Subscription
	stmt := "SELECT id, service_name FROM subscriptions WHERE id = $1"
	row := m.DB.QueryRow(stmt, id)
	err := row.Scan(&s.ID, &s.ServiceName)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Subscription{}, ErrNoRecord
		} else {
			return Subscription{}, err
		}
	}

	return s, nil
}

func (m *SubscriptionModel) Insert(userID, serviceName string, price int, startDate time.Time, endDate *time.Time) (uuid.UUID, error) {
	var id uuid.UUID
	stmt := "INSERT INTO subscriptions (user_id, service_name, price, start_date, end_date) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	err := m.DB.QueryRow(stmt, userID, serviceName, price, startDate, endDate).Scan(&id)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}
