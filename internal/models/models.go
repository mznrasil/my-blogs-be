package models

import (
	"time"
)

type Payment struct {
	Id            string    `json:"id"`
	Pidx          string    `json:"pidx"`
	Status        string    `json:"status"`
	TransactionId string    `json:"transaction_id"`
	Amount        float64   `json:"amount"`
	Mobile        string    `json:"mobile"`
	TotalAmount   float64   `json:"total_amount"`
	PlanId        int       `json:"plan_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type UpdatePaymentKhaltiPayload struct {
	Pidx          string  `json:"pidx"           validate:"required"`
	TransactionId string  `json:"transaction_id" validate:"required"`
	Amount        float64 `json:"amount"         validate:"required"`
	TotalAmount   float64 `json:"total_amount"   validate:"required"`
	Mobile        string  `json:"mobile"         validate:"required"`
	Status        string  `json:"status"         validate:"required"`
	PlanId        int     `json:"plan_id"        validate:"required"`
}

type CustomerInfo struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type InitiatePaymentKhaltiPayload struct {
	ReturnUrl         string       `json:"return_url"`
	WebsiteUrl        string       `json:"website_url"`
	Amount            string       `json:"amount"`
	PurchaseOrderID   string       `json:"purchase_order_id"`
	PurchaseOrderName string       `json:"purchase_order_name"`
	CustomerInfo      CustomerInfo `json:"customer_info"`
}

type InitiatePaymentPayload struct {
	Pidx   string `json:"pidx"`
	Status string `json:"status"`
}

type KhaltiPaymentResponse struct {
	Pidx       string `json:"pidx"`
	PaymentUrl string `json:"payment_url"`
	ExpiresAt  string `json:"expires_at"`
	ExpiresIn  int    `json:"expires_in"`
}

type PaymentStore interface {
	InitiatePayment(data InitiatePaymentPayload) error
	GetPlanById(id int) (*Plan, error)
	GetUserByID(id string) (*User, error)
	UpdatePayment(userID string, payload UpdatePaymentKhaltiPayload) error
}

type Plan struct {
	ID        int       `json:"id"`
	PlanName  string    `json:"plan_name"`
	Amount    float64   `json:"amount"`
	Interval  string    `json:"interval"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Subscription struct {
	Id        string    `json:"id"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	UserId    string    `json:"user_id"`
	PlanId    int       `json:"plan_id"`
	PaymentId string    `json:"payment_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SubscriptionStore interface {
	CheckSubscriptionStatus(userID string) (bool, error)
	GetSubscriptionDetails(userID string) (*Subscription, error)
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
