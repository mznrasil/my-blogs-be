package users

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"

	"github.com/mznrasil/my-blogs-be/internal/helpers"
	"github.com/mznrasil/my-blogs-be/internal/models"
)

type Handler struct {
	store models.UserStore
}

func NewHandler(store models.UserStore) *Handler {
	return &Handler{
		store: store,
	}
}

func (h Handler) RegisterRoutes(router *mux.Router) {
	publicRouter := router.NewRoute().Subrouter()
	publicRouter.HandleFunc("/users", h.CreateUser).Methods(http.MethodPost)
	publicRouter.HandleFunc("/users/{id}", h.GetUserById).Methods(http.MethodGet)
	publicRouter.HandleFunc("/users/{id}", h.UpdateCustomerId).Methods(http.MethodPatch)
	publicRouter.HandleFunc("/customers/{id}", h.GetCustomerById).Methods(http.MethodGet)
}

func (h Handler) GetCustomerById(w http.ResponseWriter, r *http.Request) {
	customerID := mux.Vars(r)["id"]
	fmt.Println(customerID, "customerID")
	if customerID == "" {
		helpers.WriteJSONError(w, http.StatusBadRequest, "Customer ID not provided")
		return
	}

	customer, err := h.store.GetCustomerById(customerID)
	if err != nil {
		if err == sql.ErrNoRows {
			helpers.WriteJSONError(
				w,
				http.StatusNotFound,
				fmt.Sprintf("No customer found: %v", err.Error()),
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

	helpers.WriteJSONSuccess(w, http.StatusOK, "Customer fetched Successfully", customer)
}

func (h Handler) UpdateCustomerId(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["id"]
	if userID == "" {
		helpers.WriteJSONError(w, http.StatusBadRequest, "User ID not provided")
		return
	}

	var data struct {
		CustomerID string `json:"customer_id"`
	}
	helpers.DecodeJSONBody(w, r, &data)

	if err := helpers.Validate.Struct(data); err != nil {
		errors := err.(validator.ValidationErrors)
		helpers.WriteJSONError(
			w,
			http.StatusBadRequest,
			fmt.Sprintf("Invalid payload: %v", errors.Error()),
		)
		return
	}

	userCustomerID, err := h.store.UpdateCustomerId(userID, data.CustomerID)
	if err != nil {
		helpers.WriteJSONError(
			w,
			http.StatusInternalServerError,
			fmt.Sprintf("Server error: %v", err.Error()),
		)
		return
	}

	helpers.WriteJSONSuccess(
		w,
		http.StatusOK,
		"Updated Customer ID Successfully",
		userCustomerID,
	)
}

func (h Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	newUser := new(models.CreateUserPayload)
	helpers.DecodeJSONBody(w, r, newUser)

	if err := helpers.Validate.Struct(newUser); err != nil {
		errors := err.(validator.ValidationErrors)
		helpers.WriteJSONError(
			w,
			http.StatusBadRequest,
			fmt.Sprintf("Invalid payload: %v", errors.Error()),
		)
		return
	}

	// finally save the user in the database
	if err := h.store.CreateUser(*newUser); err != nil {
		helpers.WriteJSONError(
			w,
			http.StatusInternalServerError,
			fmt.Sprintf("Error saving user: %s", err.Error()),
		)
		return
	}

	helpers.WriteJSONSuccess(w, http.StatusCreated, "User Created Successfully", nil)
}

func (h Handler) GetUserById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	user, err := h.store.GetUserByID(id)
	if err != nil {
		if user == nil {
			helpers.WriteJSONError(
				w,
				http.StatusNotFound,
				fmt.Sprintf("User not found: %v", err.Error()),
			)
			return
		}
		helpers.WriteJSONError(
			w,
			http.StatusInternalServerError,
			fmt.Sprintf("Error: %v", err.Error()),
		)
		return
	}

	helpers.WriteJSONSuccess(w, http.StatusOK, "User Found", user)
}
