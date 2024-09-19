package sites

import (
	"context"
	"database/sql"
	"errors"
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

func (s *Store) CreateSite(newSite models.CreateSitePayload, userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stmt := `
    INSERT INTO sites
      (id, name, description, subdirectory, image_url, created_at, updated_at, user_id)
    VALUES
      ($1, $2, $3, $4, $5, $6, $7, $8)
  `

	uuid, err := uuid.NewV7()
	if err != nil {
		return err
	}

	_, err = s.db.ExecContext(
		ctx,
		stmt,
		uuid,
		newSite.Name,
		newSite.Description,
		newSite.Subdirectory,
		newSite.ImageUrl,
		time.Now(),
		time.Now(),
		userID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetSiteByID(siteID string) (*models.Site, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
    SELECT
      id, name, description, subdirectory, image_url, created_at, updated_at, user_id
    FROM sites
    WHERE id = $1
  `

	site := new(models.Site)
	err := s.db.QueryRowContext(ctx, query, siteID).Scan(
		&site.ID,
		&site.Name,
		&site.Description,
		&site.Subdirectory,
		&site.ImageUrl,
		&site.CreatedAt,
		&site.UpdatedAt,
		&site.UserID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return site, nil
}

func (s *Store) GetSiteBySubdirectory(subdirectory string) (*models.Site, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
    SELECT
      id, name, description, subdirectory, image_url, created_at, updated_at, user_id
    FROM sites
    WHERE subdirectory = $1
  `

	site := new(models.Site)
	err := s.db.QueryRowContext(ctx, query, subdirectory).Scan(
		&site.ID,
		&site.Name,
		&site.Description,
		&site.Subdirectory,
		&site.ImageUrl,
		&site.CreatedAt,
		&site.UpdatedAt,
		&site.UserID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return site, nil
}

func (s *Store) GetAllSitesByUserId(userID string, take int) ([]models.Site, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var args []interface{}

	query := `
	    SELECT id, name, description, subdirectory, image_url, created_at, updated_at, user_id
	    FROM sites
	    WHERE user_id = $1
	    ORDER BY created_at DESC
	`
	args = append(args, userID)

	if take != 0 {
		query += `LIMIT $2`
		args = append(args, take)
	}

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sites []models.Site

	for rows.Next() {
		var site models.Site
		err := rows.Scan(
			&site.ID,
			&site.Name,
			&site.Description,
			&site.Subdirectory,
			&site.ImageUrl,
			&site.CreatedAt,
			&site.UpdatedAt,
			&site.UserID,
		)
		if err != nil {
			return nil, err
		}
		sites = append(sites, site)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return sites, nil
}

func (s *Store) UpdateSiteImage(siteID, userID, imageUrl string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	site, err := s.GetSiteByID(siteID)
	if err != nil {
		return err
	}
	if site == nil {
		return errors.New("Site not found")
	}

	stmt := `
		UPDATE sites
		SET
			id = $1,
			name = $2,
			description = $3,
			subdirectory = $4,
			image_url = $5,
			created_at = $6,
			updated_at = $7,
			user_id = $8
		WHERE id = $1 AND user_id = $8
	`
	_, err = s.db.ExecContext(ctx, stmt,
		siteID,
		site.Name,
		site.Description,
		site.Subdirectory,
		imageUrl,
		site.CreatedAt,
		time.Now(),
		userID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) DeleteSite(siteID, userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stmt := `
		DELETE FROM sites
		WHERE id = $1 AND user_id = $2
	`

	_, err := s.db.ExecContext(ctx, stmt, siteID, userID)
	if err != nil {
		return err
	}

	return nil
}
