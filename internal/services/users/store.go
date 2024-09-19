package users

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

func (s *Store) GetCustomerById(customerID string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
    SELECT
      id, first_name, last_name, email, profile_image, created_at, updated_at, customer_id
    FROM users
    WHERE customer_id = $1;
  `

	customer := new(models.User)
	err := s.db.QueryRowContext(ctx, query, customerID).Scan(
		&customer.ID,
		&customer.FirstName,
		&customer.LastName,
		&customer.Email,
		&customer.ProfileImage,
		&customer.CreatedAt,
		&customer.UpdatedAt,
		&customer.CustomerID,
	)
	if err != nil {
		return nil, err
	}

	return customer, nil
}

func (s *Store) UpdateCustomerId(userID, customerID string) (*models.UserCustomerID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user, err := s.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	stmt := `
    UPDATE users
    SET 
      id = $1,
      first_name = $2,
      last_name = $3,
      email = $4,
      profile_image = $5,
      created_at = $6,
      updated_at = $7,
      customer_id = $8
    WHERE id = $1
    RETURNING customer_id;
  `

	user_customer_id := new(models.UserCustomerID)
	err = s.db.QueryRowContext(
		ctx,
		stmt,
		userID,
		user.FirstName,
		user.LastName,
		user.Email,
		user.ProfileImage,
		user.CreatedAt,
		user.UpdatedAt,
		customerID,
	).Scan(&user_customer_id.CustomerID)
	if err != nil {
		return nil, err
	}

	return user_customer_id, nil
}

func (s *Store) CreateUser(newUser models.CreateUserPayload) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stmt := `
    INSERT INTO users
      (id, first_name, last_name, email, profile_image, created_at, updated_at)
    VALUES
      ($1, $2, $3, $4, $5, $6, $7)
  `

	_, err := s.db.ExecContext(
		ctx,
		stmt,
		newUser.ID,
		newUser.FirstName,
		newUser.LastName,
		newUser.Email,
		newUser.ProfileImage,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetUserByID(id string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user := new(models.User)
	query := `
    SELECT id, first_name, last_name, email, profile_image, created_at, updated_at
    FROM users
    WHERE id=$1
  `

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
