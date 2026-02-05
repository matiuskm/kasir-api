package handlers

import (
	"encoding/json"
	"kasir-api/models"
	"kasir-api/services"
	"net/http"
)

type TransactionHandler struct {
	service *services.TransactionService
}

func NewTransactionHandler(service *services.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

// receive items and qty
func (h *TransactionHandler) HandleCheckout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(contentTypeHeader, jsonContentType)
	switch r.Method {
		case http.MethodPost:
			h.Checkout(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *TransactionHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	var req models.CheckoutRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid product ID"})
		return
	}

	transaction, err := h.service.Checkout(req.Items)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid data"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.NewTransactionResponse(transaction))
}
