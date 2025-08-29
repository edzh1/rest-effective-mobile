package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID          uuid.UUID  `json:"id"`
	UserID      uuid.UUID  `json:"user_id"`
	ServiceName string     `json:"service_name"`
	Price       int        `json:"price"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date,omitempty"`
}

type SubscriptionModel struct {
	DB *sql.DB
}

type SubscriptionFilter struct {
	UserID      *uuid.UUID
	ServiceName *string
	StartDate   *time.Time
	EndDate     *time.Time
	Page        *int
}

func (m *SubscriptionModel) Get(id uuid.UUID) (Subscription, error) {
	var s Subscription
	stmt := `
		SELECT id, user_id, service_name, price, start_date, end_date
		FROM subscriptions
		WHERE id = $1
	`
	row := m.DB.QueryRow(stmt, id)
	err := row.Scan(&s.ID, &s.UserID, &s.ServiceName, &s.Price, &s.StartDate, &s.EndDate)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Subscription{}, ErrNoRecord
		} else {
			return Subscription{}, err
		}
	}

	return s, nil
}

func (m *SubscriptionModel) List(filter SubscriptionFilter) ([]Subscription, error) {
	var subscriptions []Subscription
	limit := 20
	argIndex := 1
	args := []interface{}{}
	stmt := `
		SELECT id, user_id, service_name, price, start_date, end_date
		FROM subscriptions
		WHERE 1 = 1
	`

	if filter.UserID != nil {
		stmt += fmt.Sprintf(" AND user_id = $%d", argIndex)
		args = append(args, *filter.UserID)
		argIndex++
	}

	if filter.ServiceName != nil {
		stmt += fmt.Sprintf(" AND service_name = $%d", argIndex)
		args = append(args, *filter.ServiceName)
		argIndex++
	}

	if filter.StartDate != nil {
		stmt += fmt.Sprintf(" AND start_date >= $%d", argIndex)
		args = append(args, *filter.StartDate)
		argIndex++
	}

	if filter.EndDate != nil {
		stmt += fmt.Sprintf(" AND end_date <= $%d", argIndex)
		args = append(args, *filter.EndDate)
		argIndex++
	}

	if filter.Page != nil {
		stmt += fmt.Sprintf(" AND offset = $%d AND limit = $%d", argIndex, argIndex+1)
		args = append(args, *filter.Page*limit, limit)
		argIndex += 2
	}

	rows, err := m.DB.Query(stmt, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var s Subscription
		err = rows.Scan(&s.ID, &s.UserID, &s.ServiceName, &s.Price, &s.StartDate, &s.EndDate)
		if err != nil {
			return nil, err
		}
		subscriptions = append(subscriptions, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return subscriptions, nil
}

func (m *SubscriptionModel) Insert(userID, serviceName string, price int, startDate time.Time, endDate *time.Time) (uuid.UUID, error) {
	var id uuid.UUID
	stmt := `
		INSERT INTO subscriptions
		(user_id, service_name, price, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	err := m.DB.QueryRow(stmt, userID, serviceName, price, startDate, endDate).Scan(&id)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (m *SubscriptionModel) Update(id uuid.UUID, userID, serviceName string, price int, startDate time.Time, endDate *time.Time) (uuid.UUID, error) {
	stmt := `
		INSERT INTO subscriptions
		SET user_id = $2, service_name = $3, price = $4, start_date = $5, end_date = $6 
		WHERE id = $1
	`
	result, err := m.DB.Exec(stmt, id, userID, serviceName, price, startDate, endDate)
	if err != nil {
		return uuid.Nil, err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return uuid.Nil, err
	}
	if rows == 0 {
		return uuid.Nil, ErrNoRecord
	}
	return id, nil
}

func (m *SubscriptionModel) Delete(id uuid.UUID) error {
	stmt := `
		DELETE FROM subscriptions
		WHERE id = $1
	`
	result, err := m.DB.Exec(stmt, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNoRecord
	}
	return nil
}
