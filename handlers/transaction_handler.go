package handlers

import (
	"database/sql"
	"encoding/json"
	"kasir-api/models"
	"kasir-api/services"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type TransactionHandler struct {
	service *services.TransactionService
}

func NewTransactionHandler(service *services.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

// HandleCheckout godoc
// @Summary Create transaction (checkout)
// @Description Creates a transaction and returns the simplified details.
// @Tags transactions
// @Accept json
// @Produce json
// @Param checkout body models.CheckoutRequest true "Checkout payload"
// @Success 201 {object} models.TransactionResponse
// @Failure 400 {object} map[string]string
// @Router /api/checkout [post]
func (h *TransactionHandler) HandleCheckout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(contentTypeHeader, jsonContentType)
	switch r.Method {
	case http.MethodPost:
		h.Checkout(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// HandleReportHariIni godoc
// @Summary Today's sales report
// @Description Returns total revenue, total transactions, and top-selling products for today based on REPORT_TZ (default Asia/Jakarta).
// @Tags reports
// @Produce json
// @Success 200 {object} models.DailyReport
// @Failure 500 {object} map[string]string
// @Router /api/report/hari-ini [get]
func (h *TransactionHandler) HandleReportHariIni(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(contentTypeHeader, jsonContentType)
	switch r.Method {
	case http.MethodGet:
		h.ReportHariIni(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// HandleReportRange godoc
// @Summary Report by date range
// @Description Returns report for a date range (inclusive) based on REPORT_TZ (default Asia/Jakarta).
// @Tags reports
// @Produce json
// @Param start_date query string true "Start date (YYYY-MM-DD)"
// @Param end_date query string true "End date (YYYY-MM-DD)"
// @Success 200 {object} models.DailyReport
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/report [get]
func (h *TransactionHandler) HandleReportRange(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(contentTypeHeader, jsonContentType)
	switch r.Method {
	case http.MethodGet:
		h.ReportRange(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// HandleTransactions godoc
// @Summary List transactions
// @Description Returns transactions with total products and total amount.
// @Tags transactions
// @Produce json
// @Param page query int false "Page number (starts at 1)"
// @Param limit query int false "Items per page (max 100)"
// @Param month query string false "Filter by month (YYYY-MM)"
// @Param start_date query string false "Filter start date (YYYY-MM-DD)"
// @Param end_date query string false "Filter end date (YYYY-MM-DD)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/transactions [get]
func (h *TransactionHandler) HandleTransactions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(contentTypeHeader, jsonContentType)
	switch r.Method {
	case http.MethodGet:
		h.ListTransactions(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// HandleTransactionByID godoc
// @Summary Get transaction detail
// @Description Returns transaction detail with simplified item fields.
// @Tags transactions
// @Produce json
// @Param id path int true "Transaction ID"
// @Success 200 {object} models.TransactionResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/transactions/{id} [get]
func (h *TransactionHandler) HandleTransactionByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(contentTypeHeader, jsonContentType)
	switch r.Method {
	case http.MethodGet:
		h.GetTransactionByID(w, r)
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

func (h *TransactionHandler) ReportHariIni(w http.ResponseWriter, r *http.Request) {
	report, err := h.service.GetTodayReport()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to retrieve report"})
		return
	}

	json.NewEncoder(w).Encode(report)
}

func (h *TransactionHandler) ReportRange(w http.ResponseWriter, r *http.Request) {
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")
	if startDate == "" || endDate == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "start_date and end_date are required"})
		return
	}

	startParsed, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid start_date format"})
		return
	}
	endParsed, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid end_date format"})
		return
	}
	if endParsed.Before(startParsed) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "end_date must be after or equal to start_date"})
		return
	}

	report, err := h.service.GetReportRange(startDate, endDate)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to retrieve report"})
		return
	}

	json.NewEncoder(w).Encode(report)
}

func (h *TransactionHandler) ListTransactions(w http.ResponseWriter, r *http.Request) {
	page, limit, ok := parsePagination(r)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid pagination params"})
		return
	}

	q := r.URL.Query()
	month := q.Get("month")
	startDate := q.Get("start_date")
	endDate := q.Get("end_date")

	if month != "" && (startDate != "" || endDate != "") {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Use either month or start_date/end_date"})
		return
	}

	if month != "" {
		parsedMonth, err := time.Parse("2006-01", month)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid month format"})
			return
		}
		startDate = time.Date(parsedMonth.Year(), parsedMonth.Month(), 1, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
		endDate = time.Date(parsedMonth.Year(), parsedMonth.Month()+1, 0, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
	}

	if startDate != "" || endDate != "" {
		if startDate == "" || endDate == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "start_date and end_date are required"})
			return
		}
		startParsed, err := time.Parse("2006-01-02", startDate)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid start_date format"})
			return
		}
		endParsed, err := time.Parse("2006-01-02", endDate)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid end_date format"})
			return
		}
		if endParsed.Before(startParsed) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "end_date must be after or equal to start_date"})
			return
		}
	}

	var (
		total        int
		totalRevenue int
		transactions []models.TransactionListItem
		err          error
	)

	if startDate != "" && endDate != "" {
		total, err = h.service.CountTransactionsByDateRange(startDate, endDate)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to count transactions"})
			return
		}

		totalRevenue, err = h.service.SumTransactionsByDateRange(startDate, endDate)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to sum transactions"})
			return
		}

		transactions, err = h.service.GetTransactionsByDateRange(page, limit, startDate, endDate)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to retrieve transactions"})
			return
		}
	} else {
		total, err = h.service.CountTransactions()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to count transactions"})
			return
		}

		totalRevenue, err = h.service.SumTransactions()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to sum transactions"})
			return
		}

		transactions, err = h.service.GetTransactions(page, limit)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to retrieve transactions"})
			return
		}
	}

	meta := PaginationMeta{
		Total:    total,
		Page:     page,
		PageSize: limit,
		HasMore:  page*limit < total,
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":          transactions,
		"meta":          meta,
		"total_revenue": totalRevenue,
	})
}

func (h *TransactionHandler) GetTransactionByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/transactions/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid transaction ID"})
		return
	}

	transaction, err := h.service.GetTransactionByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Transaction not found"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to retrieve transaction"})
		return
	}

	json.NewEncoder(w).Encode(models.NewTransactionResponse(transaction))
}
