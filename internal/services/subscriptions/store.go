package subscriptions

import (
	"context"
	"database/sql"
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

func (s *Store) UpdateSubscription(id string, subscription models.UpdateSubscriptionPayload) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	gotSubscription, err := s.GetSubscriptionByID(id)
	if err != nil {
		return err
	}

	stmt := `
    UPDATE subscriptions
    SET
      stripe_subscription_id = $1, 
      interval = $2, 
      status = $3, 
      plan_id = $4, 
      current_period_start = $5, 
      current_period_end = $6, 
      created_at = $7, 
      updated_at = $8,
      user_id = $9
    WHERE stripe_subscription_id = $1;
  `

	_, err = s.db.ExecContext(
		ctx,
		stmt,
		gotSubscription.StripeSubscriptionId,
		gotSubscription.Interval,
		subscription.Status,
		subscription.PlanId,
		subscription.CurrentPeriodStart,
		subscription.CurrentPeriodEnd,
		gotSubscription.CreatedAt,
		time.Now(),
		gotSubscription.UserId,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetSubscriptionByID(id string) (*models.Subscription, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
    SELECT 
      stripe_subscription_id, interval, status, plan_id, 
      current_period_start, current_period_end, created_at, updated_at, user_id
    FROM subscriptions
    WHERE stripe_subscription_id = $1 
  `

	subscription := new(models.Subscription)
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&subscription.StripeSubscriptionId,
		&subscription.Interval,
		&subscription.Status,
		&subscription.PlanId,
		&subscription.CurrentPeriodStart,
		&subscription.CurrentPeriodEnd,
		&subscription.CreatedAt,
		&subscription.UpdatedAt,
		&subscription.UserId,
	)
	if err != nil {
		return nil, err
	}

	return subscription, nil
}

func (s *Store) CreateSubscription(
	userID string,
	newSubscription models.CreateSubscriptionPayload,
) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stmt := `
    INSERT INTO subscriptions
      (stripe_subscription_id, 
      interval, 
      status, 
      plan_id, 
      current_period_start, 
      current_period_end, 
      created_at, 
      updated_at,
      user_id) 
    VALUES
      (?, ?, ?, ?, ?, ?, ?, ?, ?)
  `

	_, err := s.db.ExecContext(ctx, stmt,
		newSubscription.StripeSubscriptionId,
		newSubscription.Interval,
		newSubscription.Status,
		newSubscription.PlanId,
		newSubscription.CurrentPeriodStart,
		newSubscription.CurrentPeriodEnd,
		time.Now(),
		time.Now(),
		userID,
	)
	if err != nil {
		return err
	}

	return nil
}
