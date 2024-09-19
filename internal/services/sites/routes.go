package sites

import (
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
	store models.SiteStore
}

func NewHandler(store models.SiteStore) *Handler {
	return &Handler{
		store: store,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.Use(middleware.WithAuth)
	router.HandleFunc("/sites", h.CreateSite).Methods(http.MethodPost)
	router.HandleFunc("/sites", h.GetAllSites).Methods(http.MethodGet)
	router.HandleFunc("/sites/{siteID}", h.UpdateSiteImage).Methods(http.MethodPatch)
	router.HandleFunc("/sites/{siteID}", h.DeleteSite).Methods(http.MethodDelete)
}

func (h *Handler) CreateSite(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	newSite := new(models.CreateSitePayload)
	helpers.DecodeJSONBody(w, r, newSite)

	if err := helpers.Validate.Struct(newSite); err != nil {
		errors := err.(validator.ValidationErrors)
		helpers.WriteJSONError(
			w,
			http.StatusBadRequest,
			fmt.Sprintf("Invalid Payload: %v", errors.Error()),
		)
		return
	}

	site, err := h.store.GetSiteBySubdirectory(newSite.Subdirectory)
	if err != nil {
		helpers.WriteJSONError(
			w,
			http.StatusInternalServerError,
			fmt.Sprintf("Internal Server Error: %v", err.Error()),
		)
		return
	}

	if site != nil {
		helpers.WriteJSONError(
			w,
			http.StatusBadRequest,
			"Site with this subdirectory already exists",
		)
		return
	}

	if err = h.store.CreateSite(*newSite, userID); err != nil {
		helpers.WriteJSONError(
			w,
			http.StatusInternalServerError,
			fmt.Sprintf("Error: %v", err.Error()),
		)
		return
	}

	helpers.WriteJSONSuccess(w, http.StatusCreated, "Site Created", nil)
}

func (h *Handler) GetAllSites(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	if userID == "" {
		helpers.WriteJSONError(w, http.StatusNotFound, "User not found")
		return
	}

	takeParam := r.URL.Query().Get("take")
	var take int
	if takeParam != "" {
		takeInt, err := strconv.Atoi(takeParam)
		if err != nil {
			helpers.WriteJSONError(
				w,
				http.StatusInternalServerError,
				fmt.Sprintf("Failed to convert take parameter to integer: %v", err.Error()),
			)
			return
		}
		take = takeInt
	}

	sites, err := h.store.GetAllSitesByUserId(userID, take)
	if err != nil {
		helpers.WriteJSONError(
			w,
			http.StatusInternalServerError,
			fmt.Sprintf("Error: %v", err.Error()),
		)
		return
	}

	helpers.WriteJSONSuccess(w, http.StatusOK, "Sites fetched successfully", sites)
}

func (h *Handler) UpdateSiteImage(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	siteID := mux.Vars(r)["siteID"]
	if siteID == "" {
		helpers.WriteJSONError(w, http.StatusBadRequest, "Site Id not provided")
		return
	}

	var data struct {
		ImageUrl string `json:"image_url"`
	}
	helpers.DecodeJSONBody(w, r, &data)
	if data.ImageUrl == "" {
		helpers.WriteJSONError(w, http.StatusBadRequest, "Image url not provided")
		return
	}

	err := h.store.UpdateSiteImage(siteID, userID, data.ImageUrl)
	if err != nil {
		helpers.WriteJSONError(
			w,
			http.StatusInternalServerError,
			fmt.Sprintf("Server Error: %v", err.Error()),
		)
		return
	}

	helpers.WriteJSONSuccess(w, http.StatusOK, "Updated site image", nil)
}

func (h *Handler) DeleteSite(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	siteID := mux.Vars(r)["siteID"]
	if userID == "" {
		helpers.WriteJSONError(w, http.StatusBadRequest, "User ID not found")
		return
	}
	if siteID == "" {
		helpers.WriteJSONError(w, http.StatusBadRequest, "Site ID not found")
		return
	}

	err := h.store.DeleteSite(siteID, userID)
	if err != nil {
		helpers.WriteJSONError(
			w,
			http.StatusInternalServerError,
			fmt.Sprintf("Server error: %v", err.Error()),
		)
		return
	}

	helpers.WriteJSONSuccess(w, http.StatusOK, "Deleted site successfully", nil)
}
