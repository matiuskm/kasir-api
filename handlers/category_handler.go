package handlers

import (
	"encoding/json"
	"kasir-api/models"
	"kasir-api/services"
	"net/http"
	"strconv"
	"strings"
)

type CategoryHandler struct {
	service *services.CategoryService
}

func NewCategoryHandler(service *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

func (h *CategoryHandler) getAllCategories(w http.ResponseWriter, _ *http.Request) {
	categories, err := h.service.GetAllCategories()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string {
			"error": "failed to retrieve categories",
		})
		return
	}
	_ = json.NewEncoder(w).Encode(categories)
}

func (h *CategoryHandler) createCategory(w http.ResponseWriter, r *http.Request) {
	var category models.Category
	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string {
			"error": "invalid request body",
		})
		return
	}

	err = h.service.CreateCategory(&category)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ =json.NewEncoder(w).Encode(map[string]string {
			"error": "failed to create category",
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(category)
}

func (h *CategoryHandler) getCategoryByID(w http.ResponseWriter, r *http.Request) {
	// fetch id from url
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string {
			"error": "invalid category id",
		})
		return
	}

	include := r.URL.Query().Get("include")
	var category *models.Category
	if include == "products" {
		category, err = h.service.GetCategoryByIDWithProducts(id)
	} else {
		category, err = h.service.GetCategoryByID(id)
	}
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string {
			"error": "category not found",
		})
		return
	}

	_ = json.NewEncoder(w).Encode(category)
}

func (h *CategoryHandler) updateCategoryByID(w http.ResponseWriter, r *http.Request) {
	// fetch id from url
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string {
			"error": "invalid category id",
		})
		return
	}

	var categoryUpdate models.Category
	err = json.NewDecoder(r.Body).Decode(&categoryUpdate)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string {
			"error": "invalid request body",
		})
		return
	}

	categoryUpdate.ID = id
	err = h.service.UpdateCategory(&categoryUpdate)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string {
			"error": "failed to update category",
		})
		return
	}

	_ = json.NewEncoder(w).Encode(categoryUpdate)
}

func (h *CategoryHandler) deleteCategoryByID(w http.ResponseWriter, r *http.Request) {
	// fetch id from url
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string {
			"error": "invalid category id",
		})
		return
	}

	err = h.service.DeleteCategory(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string {
			"error": "failed to delete category",
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// HandleCategories godoc
// @Summary List or create categories
// @Description GET returns all categories. POST creates a new category.
// @Tags categories
// @Accept json
// @Produce json
// @Param category body models.Category false "Category payload (POST only)"
// @Success 200 {array} models.Category
// @Success 201 {object} models.Category
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/categories [get]
// @Router /api/categories [post]
func (h *CategoryHandler) HandleCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(contentTypeHeader, jsonContentType)

	switch r.Method {
		case http.MethodGet:
			h.getAllCategories(w, r)
		case http.MethodPost:
			h.createCategory(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// HandleCategoryByID godoc
// @Summary Get, update, or delete category by ID
// @Description GET retrieves a category; use include=products to include products. PUT updates it, DELETE removes it.
// @Tags categories
// @Produce json
// @Param id path int true "Category ID"
// @Param include query string false "Optional includes" Enums(products)
// @Param category body models.Category false "Category payload (PUT only)"
// @Success 200 {object} models.Category
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/categories/{id} [get]
// @Router /api/categories/{id} [put]
// @Router /api/categories/{id} [delete]
func (h *CategoryHandler) HandleCategoryByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(contentTypeHeader, jsonContentType)

	switch r.Method {
		case http.MethodGet:
			h.getCategoryByID(w, r)
		case http.MethodPut:
			h.updateCategoryByID(w, r)
		case http.MethodDelete:
			h.deleteCategoryByID(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
