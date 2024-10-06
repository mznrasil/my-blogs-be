package subscriptions

import (
	"database/sql"
	"fmt"
	"net/http"

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
	router.HandleFunc("/subscriptions/status", h.CheckSubscriptionStatus).Methods(http.MethodGet)
	router.HandleFunc("/subscriptions", h.GetSubscriptionDetails).Methods(http.MethodGet)
}

func (h *Handler) GetSubscriptionDetails(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	subscription, err := h.store.GetSubscriptionDetails(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			helpers.WriteJSONError(w, http.StatusNotFound, "Subscription not found")
			return
		}
		helpers.WriteJSONError(w, http.StatusInternalServerError, fmt.Sprintf("Server error: %v", err.Error()))
		return
	}

	helpers.WriteJSONSuccess(w, http.StatusOK, "Subscription fetched successfully", subscription)
}

func (h *Handler) CheckSubscriptionStatus(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	isActive, err := h.store.CheckSubscriptionStatus(userID)
	if err != nil {
		helpers.WriteJSONError(w, http.StatusInternalServerError, fmt.Sprintf("Server error: %v", err.Error()))
		return
	}

	helpers.WriteJSONSuccess(w, http.StatusOK, "Fetch status successfully", isActive)
}
