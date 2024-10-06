package payments

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
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

func (s *Store) UpdatePayment(userID string, payload models.UpdatePaymentKhaltiPayload) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt := `
		UPDATE payments
		SET
			pidx = $1,
			transaction_id = $2,
			amount = $3,
			total_amount = $4,
			mobile = $5,
			status = $6,
			plan_id = $7
		WHERE pidx = $1
	`

	_, err = tx.Exec(stmt,
		payload.Pidx,
		payload.TransactionId,
		payload.Amount,
		payload.TotalAmount,
		payload.Mobile,
		payload.Status,
		payload.PlanId,
	)
	if err != nil {
		return err
	}

	uuid, err := uuid.NewV7()
	if err != nil {
		return err
	}

	var paymentID string
	query := `
		SELECT id from payments
		WHERE pidx = $1
	`
	if err = tx.QueryRow(query, payload.Pidx).Scan(&paymentID); err != nil {
		return err
	}

	stmt = `
		INSERT INTO subscriptions
			(id, start_date, end_date, user_id, plan_id, payment_id)
		VALUES
			($1, $2, $3, $4, $5, $6)
	`

	_, err = tx.Exec(stmt,
		uuid,
		time.Now(),
		time.Now().AddDate(0, 1, 0),
		userID,
		payload.PlanId,
		paymentID,
	)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *Store) GetUserByID(id string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, first_name, last_name, email, profile_image, created_at, updated_at
		FROM users
		WHERE id = $1;
	`

	user := new(models.User)
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.ProfileImage,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Store) GetPlanById(id int) (*models.Plan, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, plan_name, amount, interval, created_at, updated_at
		FROM plans
		WHERE id = $1
	`

	plan := new(models.Plan)
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&plan.ID,
		&plan.PlanName,
		&plan.Amount,
		&plan.Interval,
		&plan.CreatedAt,
		&plan.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return plan, nil
}

func (s *Store) InitiatePayment(data models.InitiatePaymentPayload) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stmt := `
		INSERT INTO payments
			(id, pidx, status)
		VALUES
			($1, $2, $3)
	`

	uuid, err := uuid.NewV7()
	if err != nil {
		return err
	}

	_, err = s.db.ExecContext(ctx, stmt, uuid, data.Pidx, data.Status)
	if err != nil {
		return err
	}

	return nil
}
