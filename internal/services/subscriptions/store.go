package subscriptions

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/mznrasil/my-blogs-be/internal/models"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) GetSubscriptionDetails(userID string) (*models.Subscription, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, start_date, end_date, user_id, plan_id, payment_id, created_at, updated_at
		FROM subscriptions
		WHERE user_id = $1
	`
	subscription := new(models.Subscription)
	err := s.db.QueryRowContext(ctx, query, userID).Scan(
		&subscription.Id,
		&subscription.StartDate,
		&subscription.EndDate,
		&subscription.UserId,
		&subscription.PlanId,
		&subscription.PaymentId,
		&subscription.CreatedAt,
		&subscription.UpdatedAt,
	)
	fmt.Println(subscription)
	if err != nil {
		return nil, err
	}

	return subscription, nil
}

func (s *Store) CheckSubscriptionStatus(userID string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT end_date
		FROM subscriptions
		WHERE user_id = $1
	`

	var endDate time.Time
	err := s.db.QueryRowContext(ctx, query, userID).Scan(&endDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	if time.Now().After(endDate) {
		return false, nil
	}

	return true, nil
}
