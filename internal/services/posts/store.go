package posts

import (
	"context"
	"database/sql"
	"encoding/json"
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

func (s *Store) GetAllSitePostsBySlug(subdirectory, slug string) (*models.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var siteID string
	query := `
    SELECT id FROM sites
    WHERE subdirectory = $1
  `
	err := s.db.QueryRowContext(ctx, query, subdirectory).Scan(&siteID)
	if err != nil {
		return nil, err
	}

	query = `
    SELECT id, title, article_content, small_description, image, slug, created_at, updated_at, user_id, site_id
    FROM posts
    WHERE slug = $1 AND site_id = $2
  `
	post := new(models.Post)
	var marshalledArticleContent []byte
	err = s.db.QueryRowContext(ctx, query, slug, siteID).Scan(
		&post.ID,
		&post.Title,
		&marshalledArticleContent,
		&post.SmallDescription,
		&post.Image,
		&post.Slug,
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.UserID,
		&post.SiteID,
	)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(marshalledArticleContent, &post.ArticleContent); err != nil {
		return nil, err
	}

	return post, nil
}

func (s *Store) GetAllSitePostsBySubdirectory(subdirectory string) (*models.SitePosts, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	site := new(models.SiteName)
	query := `
    SELECT id, name FROM sites
    WHERE subdirectory = $1
  `
	err := s.db.QueryRowContext(ctx, query, subdirectory).Scan(&site.ID, &site.Name)
	if err != nil {
		return nil, err
	}

	query = `
    SELECT id, title, small_description, image, slug, created_at
    FROM posts
    WHERE site_id = $1
    ORDER BY created_at DESC;
  `
	rows, err := s.db.QueryContext(ctx, query, site.ID)
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		rows.Scan(
			&post.ID,
			&post.Title,
			&post.SmallDescription,
			&post.Image,
			&post.Slug,
			&post.CreatedAt,
		)
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &models.SitePosts{
		Site:  *site,
		Posts: posts,
	}, nil
}

func (s *Store) GetAllPostsByUserID(userID string, take int) ([]models.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var args []interface{}

	query := `
	    SELECT id, title, small_description, image, slug, created_at, user_id, site_id
	    FROM posts
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

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		rows.Scan(
			&post.ID,
			&post.Title,
			&post.SmallDescription,
			&post.Image,
			&post.Slug,
			&post.CreatedAt,
			&post.UserID,
			&post.SiteID,
		)
		posts = append(posts, post)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (s *Store) GetAllPostsByUserIDAndSiteID(userID, siteID string) (*models.PostSite, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
    SELECT p.id, p.title, p.image, p.created_at, s.id, s.subdirectory
    FROM posts p
    LEFT JOIN sites s
    ON p.site_id = s.id
    WHERE p.user_id = $1 AND p.site_id = $2
    ORDER BY created_at DESC;
  `

	rows, err := s.db.QueryContext(ctx, query, userID, siteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.PostSummary
	site := new(models.SiteSubdirectory)
	for rows.Next() {
		var post models.PostSummary
		rows.Scan(
			&post.ID,
			&post.Title,
			&post.Image,
			&post.CreatedAt,
			&site.ID,
			&site.Subdirectory,
		)
		posts = append(posts, post)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &models.PostSite{
		Posts: posts,
		Site:  *site,
	}, nil
}

func (s *Store) CreatePost(newPost models.CreatePostPayload, userID, siteID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stmt := `
		INSERT INTO posts
			(id, title, article_content, small_description, image, slug, created_at, updated_at, user_id, site_id)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	uuid, err := uuid.NewV7()
	if err != nil {
		return err
	}

	_, err = s.db.ExecContext(ctx, stmt,
		uuid,
		&newPost.Title,
		&newPost.ArticleContent,
		&newPost.SmallDescription,
		&newPost.Image,
		&newPost.Slug,
		time.Now(),
		time.Now(),
		userID,
		siteID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetPostBySlug(slug, userID, siteID string) (*models.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, title, article_content, small_description, image, slug, created_at, updated_at, user_id, site_id
		FROM posts
		WHERE slug = $1 AND user_id = $2 AND site_id = $3
		ORDER BY created_at DESC;
	`

	post := new(models.Post)
	var marshalledArticleContent []byte
	err := s.db.QueryRowContext(ctx, query, slug, userID, siteID).Scan(
		&post.ID,
		&post.Title,
		&marshalledArticleContent,
		&post.SmallDescription,
		&post.Image,
		&post.Slug,
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.UserID,
		&post.SiteID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if err = json.Unmarshal(marshalledArticleContent, &post.ArticleContent); err != nil {
		return nil, err
	}

	return post, nil
}

func (s *Store) GetPostByID(postID, siteID, userID string) (*models.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, title, article_content, small_description, image, slug, created_at, updated_at, user_id, site_id
		FROM posts
		WHERE id = $1 AND site_id = $2 AND user_id = $3
		ORDER BY created_at DESC;
	`

	post := new(models.Post)
	var marshalledArticleContent []byte
	err := s.db.QueryRowContext(ctx, query, postID, siteID, userID).Scan(
		&post.ID,
		&post.Title,
		&marshalledArticleContent,
		&post.SmallDescription,
		&post.Image,
		&post.Slug,
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.UserID,
		&post.SiteID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if err = json.Unmarshal(marshalledArticleContent, &post.ArticleContent); err != nil {
		return nil, err
	}

	return post, nil
}

func (s *Store) EditPost(post models.CreatePostPayload, postID, userID, siteID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	foundPost, err := s.GetPostByID(postID, siteID, userID)
	if err != nil {
		return err
	}
	if foundPost == nil {
		return errors.New("Post not found")
	}

	stmt := `
		UPDATE posts
		SET
			id = $1,
			title = $2,
			article_content = $3,
			small_description = $4,
			image = $5,
			slug =$6,
			created_at = $7,
			updated_at = $8,
			user_id = $9,
			site_id = $10
		WHERE
			id = $1 AND user_id = $9 AND site_id = $10
	`

	_, err = s.db.ExecContext(ctx, stmt,
		postID,
		&post.Title,
		&post.ArticleContent,
		&post.SmallDescription,
		&post.Image,
		&post.Slug,
		&foundPost.CreatedAt,
		time.Now(),
		userID,
		siteID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) DeletePost(postID, siteID, userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stmt := `
		DELETE FROM posts
		WHERE id = $1 AND site_id = $2 AND user_id = $3
	`

	_, err := s.db.ExecContext(ctx, stmt, postID, siteID, userID)
	if err != nil {
		return err
	}

	return nil
}
