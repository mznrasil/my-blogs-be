package subscriptions

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"

	"github.com/mznrasil/my-blogs-be/internal/helpers"
	"github.com/mznrasil/my-blogs-be/internal/middleware"
	"github.com/mznrasil/my-blogs-be/internal/models"
)

type Handler struct {
	store models.SubscriptionStore
}

func NewHandler(store models.SubscriptionStore) *Handler {
	return &Handler{
		store: store,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.Use(middleware.WithAuth)
	router.HandleFunc("/subscriptions", h.CreateSubscription).Methods(http.MethodPost)
	router.HandleFunc("/subscriptions/{id}", h.GetSubscriptionByID).Methods(http.MethodGet)
	router.HandleFunc("/subscriptions/{id}", h.UpdateSubscription).Methods(http.MethodPatch)
}

func (h *Handler) GetSubscriptionByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if id == "" {
		helpers.WriteJSONError(w, http.StatusBadRequest, "Subscription ID not provided")
		return
	}

	subscription, err := h.store.GetSubscriptionByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			helpers.WriteJSONError(
				w,
				http.StatusNotFound,
				fmt.Sprintf("Subscription not found: %v", err.Error()),
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

	helpers.WriteJSONSuccess(w, http.StatusOK, "Subscription fetched successfully", subscription)
}

func (h *Handler) UpdateSubscription(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if id == "" {
		helpers.WriteJSONError(w, http.StatusBadRequest, "Subscription Id not provided")
		return
	}

	subscription := new(models.UpdateSubscriptionPayload)
	helpers.DecodeJSONBody(w, r, subscription)

	if err := helpers.Validate.Struct(subscription); err != nil {
		errors := err.(validator.ValidationErrors)
		helpers.WriteJSONError(
			w,
			http.StatusBadRequest,
			fmt.Sprintf("Invalid payload: %v", errors.Error()),
		)
		return
	}

	if err := h.store.UpdateSubscription(id, *subscription); err != nil {
		if err == sql.ErrNoRows {
			helpers.WriteJSONError(
				w,
				http.StatusNotFound,
				fmt.Sprintf("Subscription not found: %v", err.Error()),
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

	helpers.WriteJSONSuccess(w, http.StatusOK, "Updated subscription successfully", nil)
}

func (h *Handler) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	newSubscription := new(models.CreateSubscriptionPayload)
	helpers.DecodeJSONBody(w, r, newSubscription)

	if err := helpers.Validate.Struct(newSubscription); err != nil {
		errors := err.(validator.ValidationErrors)
		helpers.WriteJSONError(
			w,
			http.StatusBadRequest,
			fmt.Sprintf("Invalid payload: %v", errors.Error()),
		)
		return
	}

	if err := h.store.CreateSubscription(userID, *newSubscription); err != nil {
		helpers.WriteJSONError(
			w,
			http.StatusInternalServerError,
			fmt.Sprintf("Server Error: %v", err.Error()),
		)
		return
	}

	helpers.WriteJSONSuccess(w, http.StatusOK, "Subscription created successfully", nil)
}
