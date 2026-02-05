package models

type Transaction struct {
	ID              int                 `json:"id"`
	TotalAmount     int                 `json:"total_amount"`
	TransactionDate string              `json:"transaction_date"`
	Details         []TransactionDetail `json:"details"`
}

type TransactionDetail struct {
	ID            int    `json:"id"`
	TransactionID int    `json:"transaction_id"`
	ProductID     int    `json:"product_id"`
	ProductName   string `json:"product_name"`
	Quantity      int    `json:"quantity"`
	Subtotal      int    `json:"subtotal"`
}

type TransactionListItem struct {
	ID              int    `json:"id"`
	TotalProducts   int    `json:"total_products"`
	TotalAmount     int    `json:"total_amount"`
	TransactionDate string `json:"transaction_date"`
}

type TransactionResponse struct {
	ID              int                         `json:"id"`
	TotalAmount     int                         `json:"total_amount"`
	TransactionDate string                      `json:"transaction_date"`
	Details         []TransactionDetailResponse `json:"details"`
}

type TransactionDetailResponse struct {
	ProductName string `json:"product_name"`
	Quantity    int    `json:"quantity"`
	Subtotal    int    `json:"subtotal"`
}

func NewTransactionResponse(t *Transaction) TransactionResponse {
	details := make([]TransactionDetailResponse, 0, len(t.Details))
	for _, detail := range t.Details {
		details = append(details, TransactionDetailResponse{
			ProductName: detail.ProductName,
			Quantity:    detail.Quantity,
			Subtotal:    detail.Subtotal,
		})
	}

	return TransactionResponse{
		ID:              t.ID,
		TotalAmount:     t.TotalAmount,
		TransactionDate: t.TransactionDate,
		Details:         details,
	}
}

type CheckoutRequest struct {
	Items []CheckoutItem `json:"items"`
}

type CheckoutItem struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}
