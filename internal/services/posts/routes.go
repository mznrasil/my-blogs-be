package posts

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"

	"github.com/mznrasil/my-blogs-be/internal/helpers"
	"github.com/mznrasil/my-blogs-be/internal/middleware"
	"github.com/mznrasil/my-blogs-be/internal/models"
)

type Handler struct {
	store models.PostStore
}

func NewHandler(store models.PostStore) *Handler {
	return &Handler{
		store: store,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	authRouter := router.NewRoute().Subrouter()
	authRouter.Use(middleware.WithAuth)
	authRouter.HandleFunc("/posts", h.GetAllPosts).Methods(http.MethodGet)
	authRouter.HandleFunc("/{siteID}/posts", h.GetAllPostsBySiteID).Methods(http.MethodGet)
	authRouter.HandleFunc("/{siteID}/posts/{postID}", h.GetPostByID).Methods(http.MethodGet)
	authRouter.HandleFunc("/{siteID}/posts", h.CreatePost).Methods(http.MethodPost)
	authRouter.HandleFunc("/{siteID}/posts/{postID}", h.EditPost).Methods(http.MethodPatch)
	authRouter.HandleFunc("/{siteID}/posts/{postID}", h.DeletePost).Methods(http.MethodDelete)

	publicRouter := router.NewRoute().Subrouter()
	publicRouter.HandleFunc("/posts/{subdirectory}", h.GetAllSitePostsBySubdirectory).
		Methods(http.MethodGet)
	publicRouter.HandleFunc("/posts/{subdirectory}/{slug}", h.GetAllSitePostsBySlug).
		Methods(http.MethodGet)
}

func (h *Handler) GetAllSitePostsBySlug(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subdirectory := vars["subdirectory"]
	slug := vars["slug"]
	if subdirectory == "" {
		helpers.WriteJSONError(w, http.StatusBadRequest, "Subdirectory not provided")
		return
	}
	if slug == "" {
		helpers.WriteJSONError(w, http.StatusBadRequest, "Slug not provided")
		return
	}

	post, err := h.store.GetAllSitePostsBySlug(subdirectory, slug)
	if err != nil {
		if err == sql.ErrNoRows {
			helpers.WriteJSONError(
				w,
				http.StatusNotFound,
				fmt.Sprintf("Post not found: %v", err.Error()),
			)
			return
		}

		helpers.WriteJSONError(
			w,
			http.StatusInternalServerError,
			fmt.Sprintf("Server error: %v", err.Error()),
		)
		return
	}

	helpers.WriteJSONSuccess(w, http.StatusOK, "Posts fetched successfully", post)
}

func (h *Handler) GetAllSitePostsBySubdirectory(w http.ResponseWriter, r *http.Request) {
	subdirectory := mux.Vars(r)["subdirectory"]
	if subdirectory == "" {
		helpers.WriteJSONError(w, http.StatusBadRequest, "Subdirectory not provided")
		return
	}

	sitePosts, err := h.store.GetAllSitePostsBySubdirectory(subdirectory)
	if err != nil {
		helpers.WriteJSONError(
			w,
			http.StatusInternalServerError,
			fmt.Sprintf("Server error: %v", err.Error()),
		)
		return
	}

	helpers.WriteJSONSuccess(w, http.StatusOK, "Site posts fetched successfully", sitePosts)
}

func (h *Handler) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	if userID == "" {
		helpers.WriteJSONError(w, http.StatusBadRequest, "User ID not found.")
		return
	}

	var take int
	takeParam := r.URL.Query().Get("take")
	if takeParam != "" {
		takeInt, err := strconv.Atoi(takeParam)
		if err != nil {
			helpers.WriteJSONError(
				w,
				http.StatusInternalServerError,
				fmt.Sprintf("Failed to convert take to integer: %v", err.Error()),
			)
			return
		}
		take = takeInt
	}

	posts, err := h.store.GetAllPostsByUserID(userID, take)
	if err != nil {
		helpers.WriteJSONError(
			w,
			http.StatusInternalServerError,
			fmt.Sprintf("Server error: %v", err.Error()),
		)
		return
	}

	helpers.WriteJSONSuccess(w, http.StatusOK, "Posts fetched successfully", posts)
}

func (h *Handler) GetAllPostsBySiteID(w http.ResponseWriter, r *http.Request) {
	userID, siteID, err := getUserIDAndSiteID(w, r)
	if err != nil {
		helpers.WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	posts, err := h.store.GetAllPostsByUserIDAndSiteID(userID, siteID)
	if err != nil {
		helpers.WriteJSONError(
			w,
			http.StatusInternalServerError,
			fmt.Sprintf("Error:%v", err.Error()),
		)
		return
	}

	helpers.WriteJSONSuccess(w, http.StatusOK, "Posts fetched successfully", posts)
}

func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	userID, siteID, err := getUserIDAndSiteID(w, r)
	if err != nil {
		helpers.WriteJSONError(w, http.StatusNotFound, err.Error())
		return
	}

	newPost := new(models.CreatePostPayload)
	helpers.DecodeJSONBody(w, r, newPost)

	if err := helpers.Validate.Struct(newPost); err != nil {
		errors := err.(validator.ValidationErrors)
		helpers.WriteJSONError(
			w,
			http.StatusBadRequest,
			fmt.Sprintf("Invalid Payload: %v", errors.Error()),
		)
		return
	}

	post, err := h.store.GetPostBySlug(newPost.Slug, userID, siteID)
	if err != nil {
		helpers.WriteJSONError(
			w,
			http.StatusInternalServerError,
			fmt.Sprintf("Server error: %v", err.Error()),
		)
		return
	}

	if post != nil {
		helpers.WriteJSONError(w, http.StatusConflict, "Post with this slug already exists")
		return
	}

	if err = h.store.CreatePost(*newPost, userID, siteID); err != nil {
		helpers.WriteJSONError(
			w,
			http.StatusInternalServerError,
			fmt.Sprintf("Server error: %v", err.Error()),
		)
		return
	}

	helpers.WriteJSONSuccess(w, http.StatusCreated, "Post Created Successfully", nil)
}

func (h *Handler) GetPostByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["postID"]
	if postID == "" {
		helpers.WriteJSONError(w, http.StatusBadRequest, "Post ID not specified")
		return
	}
	userID, siteID, err := getUserIDAndSiteID(w, r)
	if err != nil {
		helpers.WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	post, err := h.store.GetPostByID(postID, siteID, userID)
	if err != nil {
		helpers.WriteJSONError(
			w,
			http.StatusInternalServerError,
			fmt.Sprintf("Server Error: %v", err.Error()),
		)
		return
	}

	if post == nil {
		helpers.WriteJSONError(w, http.StatusNotFound, "Post not found")
		return
	}

	helpers.WriteJSONSuccess(w, http.StatusOK, "Post fetched successfully", post)
}

func (h *Handler) EditPost(w http.ResponseWriter, r *http.Request) {
	userID, siteID, err := getUserIDAndSiteID(w, r)
	if err != nil {
		helpers.WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	postID := mux.Vars(r)["postID"]
	if postID == "" {
		helpers.WriteJSONError(w, http.StatusBadRequest, "Post ID not found")
		return
	}

	postPayload := new(models.CreatePostPayload)
	helpers.DecodeJSONBody(w, r, postPayload)

	if err := helpers.Validate.Struct(postPayload); err != nil {
		errors := err.(validator.ValidationErrors)
		helpers.WriteJSONError(
			w,
			http.StatusBadRequest,
			fmt.Sprintf("Invalid payload: %v", errors.Error()),
		)
		return
	}

	post, err := h.store.GetPostBySlug(postPayload.Slug, userID, siteID)
	if err != nil {
		helpers.WriteJSONError(
			w,
			http.StatusInternalServerError,
			fmt.Sprintf("Server error: %v", err.Error()),
		)
		return
	}

	if post != nil {
		helpers.WriteJSONError(w, http.StatusNotFound, "Post with this slug already exists.")
		return
	}

	if err = h.store.EditPost(*postPayload, postID, userID, siteID); err != nil {
		helpers.WriteJSONError(
			w,
			http.StatusInternalServerError,
			fmt.Sprintf("Server error: %v", err.Error()),
		)
		return
	}

	helpers.WriteJSONSuccess(w, http.StatusCreated, "Post Created Successfully", nil)
}

func (h *Handler) DeletePost(w http.ResponseWriter, r *http.Request) {
	userID, siteID, err := getUserIDAndSiteID(w, r)
	if err != nil {
		helpers.WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	postID := mux.Vars(r)["postID"]
	if postID == "" {
		helpers.WriteJSONError(w, http.StatusBadRequest, "Post ID not found")
		return
	}

	if err := h.store.DeletePost(postID, siteID, userID); err != nil {
		helpers.WriteJSONError(
			w,
			http.StatusInternalServerError,
			fmt.Sprintf("Server error: %v", err.Error()),
		)
		return
	}

	helpers.WriteJSONSuccess(w, http.StatusOK, "Post Deleted Successfully", nil)
}

func getUserIDAndSiteID(_ http.ResponseWriter, r *http.Request) (string, string, error) {
	userID := r.Context().Value("userID").(string)
	vars := mux.Vars(r)
	siteID := vars["siteID"]

	if userID == "" {
		return userID, siteID, errors.New("User not found")
	}

	if siteID == "" {
		return userID, siteID, errors.New("Site not found")
	}

	return userID, siteID, nil
}
