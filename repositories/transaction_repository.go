package repositories

import (
	"database/sql"
	"fmt"
	"kasir-api/models"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) CreateTransaction(items []models.CheckoutItem) (*models.Transaction, error) {
	var (
		res *models.Transaction
	)

	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()


	// initialize subtotal
	subtotal := 0
	totalAmount := 0
	// initialize transactionDetails model
	details := make([]models.TransactionDetail, 0)
	// loop every items
	for _, item := range items {
		var productName string
		var productID, price, stock int
		// get product to get price
		err := tx.QueryRow("SELECT id, name, price, stock FROM products WHERE id=$1", item.ProductID).Scan(&productID, &productName, &price, &stock )
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product id %d not found", item.ProductID)
		}

		if err != nil {
			return nil, err
		}

		// calculate current_total = quantity * price
		// sum it all into subtotal
		subtotal = item.Quantity * price
		totalAmount += subtotal

		// decrease stock number in product
		_, err = tx.Exec("UPDATE products SET stock = stock - $1 WHERE id = $2", item.Quantity, item.ProductID)
		if err != nil {
			return nil, err
		}

		// insert items into transactionDetails
		details = append(details, models.TransactionDetail{
			ProductID: productID,
			ProductName: productName,
			Quantity: item.Quantity,
			Subtotal: subtotal,
		})
	}
	
	// create transaction
	var transactionID int
	err = tx.QueryRow("INSERT INTO transactions (total_amount) VALUES ( $1 ) RETURNING ID", totalAmount).Scan(&transactionID)
	if err != nil {
		return nil, err
	}

	// insert transactionDetails and capture IDs
	for i, detail := range details {
		details[i].TransactionID = transactionID
		err := tx.QueryRow(
			"INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES ($1, $2, $3, $4) RETURNING id",
			transactionID, detail.ProductID, detail.Quantity, detail.Subtotal,
		).Scan(&details[i].ID)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	res = &models.Transaction{
		ID: transactionID,
		TotalAmount: totalAmount,
		Details: details, 
	}

	return res, nil 
}
