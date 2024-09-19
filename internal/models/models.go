package models

import (
	"time"
)

type Subscription struct {
	StripeSubscriptionId string    `json:"stripe_subscription_id"`
	Interval             string    `json:"interval"`
	Status               string    `json:"status"`
	PlanId               string    `json:"plan_id"`
	CurrentPeriodStart   int       `json:"current_period_start"`
	CurrentPeriodEnd     int       `json:"current_period_end"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
	UserId               string    `json:"user_id"`
}

type CreateSubscriptionPayload struct {
	StripeSubscriptionId string `json:"stripe_subscription_id" validate:"required"`
	Interval             string `json:"interval"               validate:"required"`
	Status               string `json:"status"                 validate:"required"`
	PlanId               string `json:"plan_id"                validate:"required"`
	CurrentPeriodStart   int    `json:"current_period_start"   validate:"required"`
	CurrentPeriodEnd     int    `json:"current_period_end"     validate:"required"`
}

type UpdateSubscriptionPayload struct {
	PlanId             string `json:"plan_id"              validate:"required"`
	CurrentPeriodStart int    `json:"current_period_start" validate:"required"`
	CurrentPeriodEnd   int    `json:"current_period_end"   validate:"required"`
	Status             string `json:"status"               validate:"required"`
}

type SubscriptionStore interface {
	CreateSubscription(userID string, newSubscription CreateSubscriptionPayload) error
	GetSubscriptionByID(id string) (*Subscription, error)
	UpdateSubscription(id string, subscription UpdateSubscriptionPayload) error
}

type User struct {
	ID           string    `json:"id"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Email        string    `json:"email"`
	ProfileImage string    `json:"profile_image"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	CustomerID   string    `json:"customer_id"`
}

type UserCustomerID struct {
	CustomerID string `json:"customer_id"`
}

type CreateUserPayload struct {
	ID           string `json:"id"            validate:"required"`
	FirstName    string `json:"first_name"    validate:"required"`
	LastName     string `json:"last_name"     validate:"required"`
	Email        string `json:"email"         validate:"email"`
	ProfileImage string `json:"profile_image"`
}

type UserStore interface {
	CreateUser(newUser CreateUserPayload) error
	GetUserByID(id string) (*User, error)
	UpdateCustomerId(userID, customerID string) (*UserCustomerID, error)
	GetCustomerById(customerID string) (*User, error)
}

type Site struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Subdirectory string    `json:"subdirectory"`
	ImageUrl     string    `json:"image_url"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	UserID       string    `json:"user_id"`
}

type SiteName struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type SiteSubdirectory struct {
	ID           string `json:"id"`
	Subdirectory string `json:"subdirectory"`
}

type SitePosts struct {
	Site  SiteName `json:"site"`
	Posts []Post   `json:"posts"`
}

type CreateSitePayload struct {
	Name         string `json:"name"         validate:"required"`
	Description  string `json:"description"`
	Subdirectory string `json:"subdirectory" validate:"required"`
	ImageUrl     string `json:"image_url"`
}

type SiteStore interface {
	CreateSite(newSite CreateSitePayload, userID string) error
	GetSiteByID(siteID string) (*Site, error)
	GetSiteBySubdirectory(subdirectory string) (*Site, error)
	GetAllSitesByUserId(userID string, take int) ([]Site, error)
	UpdateSiteImage(siteID, userID, imageUrl string) error
	DeleteSite(siteID, userID string) error
}

type CreatePostPayload struct {
	Title            string `json:"title"`
	ArticleContent   any    `json:"article_content"`
	SmallDescription string `json:"small_description"`
	Image            string `json:"image"`
	Slug             string `json:"slug"`
}

type Post struct {
	ID               string    `json:"id"`
	Title            string    `json:"title"`
	ArticleContent   any       `json:"article_content"`
	SmallDescription string    `json:"small_description"`
	Image            string    `json:"image"`
	Slug             string    `json:"slug"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	UserID           string    `json:"user_id"`
	SiteID           string    `json:"site_id"`
}

type PostSummary struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Image     string    `json:"image"`
	CreatedAt time.Time `json:"created_at"`
}

type PostSite struct {
	Posts []PostSummary    `json:"posts"`
	Site  SiteSubdirectory `json:"site"`
}

type PostStore interface {
	GetAllPostsByUserID(userID string, take int) ([]Post, error)
	GetAllPostsByUserIDAndSiteID(userID, siteID string) (*PostSite, error)
	CreatePost(newPost CreatePostPayload, userID, siteID string) error
	GetPostBySlug(slug, userID, siteID string) (*Post, error)
	GetPostByID(postID, siteID, userID string) (*Post, error)
	EditPost(post CreatePostPayload, postID, userID, siteID string) error
	DeletePost(postID, siteID, userID string) error
	GetAllSitePostsBySubdirectory(subdirectory string) (*SitePosts, error)
	GetAllSitePostsBySlug(subdirectory, slug string) (*Post, error)
}
