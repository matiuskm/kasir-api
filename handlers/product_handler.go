package handlers

import (
	"encoding/json"
	"kasir-api/models"
	"kasir-api/services"
	"net/http"
	"strconv"
	"strings"
)

type ProductHandler struct {
	service *services.ProductService
}

const (
	contentTypeHeader = "Content-Type"
	jsonContentType   = "application/json"
)

func NewProductHandler(service *services.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) getAllProducts(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set(contentTypeHeader, jsonContentType)
	products, err := h.service.GetAllProducts()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Failed to retrieve products"})
		return
	}
	_ = json.NewEncoder(w).Encode(products)
}

func (h *ProductHandler) createProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(contentTypeHeader, jsonContentType)
	var product models.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ =json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	err = h.service.CreateProduct(&product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create product"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(contentTypeHeader, jsonContentType)
	// fetch id from url
	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")
	// convert id to integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Invalid product ID"})
		return
	}
	// fetch product by id
	product, err := h.service.GetProductByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Failed to retrieve product"})
		return
	}
	_ = json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	// fetch id from url
	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")
	// convert id to integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid product ID"})
		return
	}
	// read data from request body
	var productUpdate models.Product
	err = json.NewDecoder(r.Body).Decode(&productUpdate)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}
	productUpdate.ID = id
	// update product
	err = h.service.UpdateProduct(&productUpdate)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to update product"})
		return
	}
	json.NewEncoder(w).Encode(productUpdate)
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	// fetch id from url
	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")
	// convert id to integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid product ID"})
		return
	}
	// delete product
	err = h.service.DeleteProduct(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to delete product"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// HandleProducts godoc
// @Summary List or create products
// @Description GET returns all products. POST creates a new product.
// @Tags products
// @Accept json
// @Produce json
// @Param product body models.Product false "Product payload (POST only)"
// @Success 200 {array} models.Product
// @Success 201 {object} models.Product
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/products [get]
// @Router /api/products [post]
func (h *ProductHandler) HandleProducts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getAllProducts(w, r)
	case http.MethodPost:
		h.createProduct(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// HandleProductByID godoc
// @Summary Get, update, or delete product by ID
// @Description GET retrieves a product, PUT updates it, DELETE removes it.
// @Tags products
// @Produce json
// @Param id path int true "Product ID"
// @Param product body models.Product false "Product payload (PUT only)"
// @Success 200 {object} models.Product
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/products/{id} [get]
// @Router /api/products/{id} [put]
// @Router /api/products/{id} [delete]
func (h *ProductHandler) HandleProductByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(contentTypeHeader, jsonContentType)
	
	switch r.Method {
		case http.MethodGet:
			h.GetProductByID(w, r)
		case http.MethodPut:
			h.UpdateProduct(w, r)
		case http.MethodDelete:
			h.DeleteProduct(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
