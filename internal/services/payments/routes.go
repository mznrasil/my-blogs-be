package payments

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"

	"github.com/mznrasil/my-blogs-be/internal/helpers"
	"github.com/mznrasil/my-blogs-be/internal/models"
)

type Handler struct {
	store models.PaymentStore
}

func NewHandler(store models.PaymentStore) *Handler {
	return &Handler{
		store: store,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/payment/initiate", h.InitiatePayment).Methods(http.MethodPost)
	router.HandleFunc("/payment", h.UpdatePayment).Methods(http.MethodPatch)
}

func (h *Handler) UpdatePayment(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	payload := new(models.UpdatePaymentKhaltiPayload)
	helpers.DecodeJSONBody(w, r, payload)

	if err := helpers.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		helpers.WriteJSONError(
			w,
			http.StatusBadRequest,
			fmt.Sprintf("Invalid payload: %v", errors.Error()),
		)
		return
	}

	err := h.store.UpdatePayment(userID, *payload)
	if err != nil {
		if err == sql.ErrNoRows {
			helpers.WriteJSONError(
				w,
				http.StatusNotFound,
				fmt.Sprintf("Payment not found: %v", err.Error()),
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

	helpers.WriteJSONSuccess(w, http.StatusOK, "Payment updated succesfully", nil)
}

func (h *Handler) InitiatePayment(w http.ResponseWriter, r *http.Request) {
	var data struct {
		PlanID int `json:"plan_id"`
	}
	helpers.DecodeJSONBody(w, r, &data)

	if err := helpers.Validate.Struct(&data); err != nil {
		errors := err.(validator.ValidationErrors)
		helpers.WriteJSONError(
			w,
			http.StatusBadRequest,
			fmt.Sprintf("Invalid payload: %v", errors.Error()),
		)
		return
	}

	plan, err := h.store.GetPlanById(data.PlanID)
	if err != nil {
		if err == sql.ErrNoRows {
			helpers.WriteJSONError(w, http.StatusNotFound, "Plan not found")
			return
		}

		helpers.WriteJSONError(
			w,
			http.StatusInternalServerError,
			fmt.Sprintf("Server error: %v", err.Error()),
		)
		return
	}

	// get user info from user_id
	userID := r.Context().Value("userID").(string)
	user, err := h.store.GetUserByID(userID)
	if err != nil {
		if err == sql.ErrNoRows {
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
			fmt.Sprintf("Internal server error: %v", err.Error()),
		)
		return
	}

	// construct payload to send to khalti payment gateway
	initiatePaymentPayload := models.InitiatePaymentKhaltiPayload{
		ReturnUrl:         os.Getenv("KHALTI_RETURN_URL"),
		WebsiteUrl:        os.Getenv("KHALTI_WEBSITE_URL"),
		Amount:            strconv.FormatFloat(plan.Amount*100, 'f', 1, 64),
		PurchaseOrderID:   strconv.Itoa(plan.ID),
		PurchaseOrderName: plan.PlanName,
		CustomerInfo: models.CustomerInfo{
			Name:  fmt.Sprintf("%v %v", user.FirstName, user.LastName),
			Email: user.Email,
		},
	}

	jsonData, err := json.Marshal(initiatePaymentPayload)
	if err != nil {
		helpers.WriteJSONError(
			w,
			http.StatusInternalServerError,
			fmt.Sprintf("Failed to marshal JSON: %v", err.Error()),
		)
		return
	}

	req, err := http.NewRequest(
		http.MethodPost,
		os.Getenv("KHALTI_PAYMENT_INITIATE_API"),
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		helpers.WriteJSONError(
			w,
			http.StatusInternalServerError,
			fmt.Sprintf("Failed to construct request: %v", err.Error()),
		)
		return
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Key %v", os.Getenv("KHALTI_SECRET")))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		helpers.WriteJSONError(
			w,
			http.StatusInternalServerError,
			fmt.Sprintf("Failed to send request: %v", err.Error()),
		)
		return
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		helpers.WriteJSONError(
			w,
			http.StatusInternalServerError,
			fmt.Sprintf("Failed to read response body: %v", err.Error()),
		)
		return
	}

	khaltiPaymentResponse := new(models.KhaltiPaymentResponse)
	if err := json.Unmarshal(body, khaltiPaymentResponse); err != nil {
		helpers.WriteJSONError(
			w,
			http.StatusInternalServerError,
			fmt.Sprintf("Error unmarshalling JSON: %v", err.Error()),
		)
		return
	}

	// after successful response, save the pending state of payment in the payments table
	err = h.store.InitiatePayment(models.InitiatePaymentPayload{
		Pidx:   khaltiPaymentResponse.Pidx,
		Status: "Initiated",
	})
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
		"Payment Initiated successfully",
		khaltiPaymentResponse,
	)
}
