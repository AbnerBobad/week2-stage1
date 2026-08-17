package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Consumer struct {
	CN_ID        string    `json:"id"`
	CN_Name      string    `json:"name"`
	CN_Email     string    `json:"email"`
	CN_Status    string    `json:"status"`
	CN_Version   int       `json:"version"`
	CN_CreatedAt time.Time `json:"created_at"`
	CN_UpdatedAt time.Time `json:"updated_at"`
}

type ConsumerModel struct {
	DB *sql.DB
}

// Consumer Struct
// Creates a new consumer in the database and returns the newly created consumer with its ID, status, version, created_at, and updated_at fields populated.
func (m ConsumerModel) Insert(consumer *Consumer) error {
	query := `
		INSERT INTO consumers (name, email)
		VALUES ($1, $2)
		RETURNING id, status, version, created_at, updated_at`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{consumer.CN_Name, consumer.CN_Email}

	return m.DB.QueryRowContext(ctx, query, args...).Scan(
		&consumer.CN_ID,
		&consumer.CN_Status,
		&consumer.CN_Version,
		&consumer.CN_CreatedAt,
		&consumer.CN_UpdatedAt,
	)
}

// Fetches a consumer from the database based on the provided ID. If the consumer is found,
// it returns the consumer struct; otherwise, it returns an error indicating that the record was not found.
func (m ConsumerModel) Get(id string) (*Consumer, error) {
	query := `
		SELECT id, name, email, status, version, created_at, updated_at
		FROM consumers
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var consumer Consumer

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&consumer.CN_ID,
		&consumer.CN_Name,
		&consumer.CN_Email,
		&consumer.CN_Status,
		&consumer.CN_Version,
		&consumer.CN_CreatedAt,
		&consumer.CN_UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	return &consumer, nil
}
