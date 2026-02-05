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

func (r *TransactionRepository) CreateTransaction(items []models.CheckoutItem, reportTZ string) (*models.Transaction, error) {
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
		err := tx.QueryRow("SELECT id, name, price, stock FROM products WHERE id=$1", item.ProductID).Scan(&productID, &productName, &price, &stock)
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
			ProductID:   productID,
			ProductName: productName,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	// create transaction
	var (
		transactionID   int
		transactionDate string
	)
	err = tx.QueryRow(
		"INSERT INTO transactions (total_amount) VALUES ( $1 ) RETURNING id, TO_CHAR(created_at AT TIME ZONE $2, 'YYYY-MM-DD')",
		totalAmount, reportTZ,
	).Scan(&transactionID, &transactionDate)
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
		ID:              transactionID,
		TotalAmount:     totalAmount,
		TransactionDate: transactionDate,
		Details:         details,
	}

	return res, nil
}

func (r *TransactionRepository) GetTransactions(page, limit int, reportTZ string) ([]models.TransactionListItem, error) {
	offset := (page - 1) * limit

	rows, err := r.db.Query(`
		SELECT
			t.id,
			COALESCE(SUM(td.quantity), 0) AS total_products,
			t.total_amount,
			TO_CHAR(t.created_at AT TIME ZONE $3, 'YYYY-MM-DD') AS transaction_date
		FROM transactions t
		LEFT JOIN transaction_details td ON td.transaction_id = t.id
		GROUP BY t.id, t.total_amount, t.created_at
		ORDER BY t.id ASC
		LIMIT $1 OFFSET $2
	`, limit, offset, reportTZ)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transactions := make([]models.TransactionListItem, 0)
	for rows.Next() {
		var item models.TransactionListItem
		if err := rows.Scan(&item.ID, &item.TotalProducts, &item.TotalAmount, &item.TransactionDate); err != nil {
			return nil, err
		}
		transactions = append(transactions, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}

func (r *TransactionRepository) CountTransactions() (int, error) {
	var total int
	err := r.db.QueryRow("SELECT COUNT(*) FROM transactions").Scan(&total)
	return total, err
}

func (r *TransactionRepository) SumTransactions() (int, error) {
	var total int
	err := r.db.QueryRow("SELECT COALESCE(SUM(total_amount), 0) FROM transactions").Scan(&total)
	return total, err
}

func (r *TransactionRepository) GetTransactionsByDateRange(page, limit int, startDate, endDate, reportTZ string) ([]models.TransactionListItem, error) {
	offset := (page - 1) * limit

	rows, err := r.db.Query(`
		WITH bounds AS (
			SELECT
				(to_date($1, 'YYYY-MM-DD')::timestamp AT TIME ZONE $3) AS start_ts,
				((to_date($2, 'YYYY-MM-DD')::timestamp + INTERVAL '1 day') AT TIME ZONE $3) AS end_ts
		)
		SELECT
			t.id,
			COALESCE(SUM(td.quantity), 0) AS total_products,
			t.total_amount,
			TO_CHAR(t.created_at AT TIME ZONE $3, 'YYYY-MM-DD') AS transaction_date
		FROM transactions t
		JOIN bounds b ON true
		LEFT JOIN transaction_details td ON td.transaction_id = t.id
		WHERE t.created_at >= b.start_ts
			AND t.created_at < b.end_ts
		GROUP BY t.id, t.total_amount, t.created_at
		ORDER BY t.id ASC
		LIMIT $4 OFFSET $5
	`, startDate, endDate, reportTZ, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transactions := make([]models.TransactionListItem, 0)
	for rows.Next() {
		var item models.TransactionListItem
		if err := rows.Scan(&item.ID, &item.TotalProducts, &item.TotalAmount, &item.TransactionDate); err != nil {
			return nil, err
		}
		transactions = append(transactions, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}

func (r *TransactionRepository) CountTransactionsByDateRange(startDate, endDate, reportTZ string) (int, error) {
	var total int
	err := r.db.QueryRow(`
		WITH bounds AS (
			SELECT
				(to_date($1, 'YYYY-MM-DD')::timestamp AT TIME ZONE $3) AS start_ts,
				((to_date($2, 'YYYY-MM-DD')::timestamp + INTERVAL '1 day') AT TIME ZONE $3) AS end_ts
		)
		SELECT COUNT(*)
		FROM transactions t
		JOIN bounds b ON true
		WHERE t.created_at >= b.start_ts
			AND t.created_at < b.end_ts
	`, startDate, endDate, reportTZ).Scan(&total)
	return total, err
}

func (r *TransactionRepository) SumTransactionsByDateRange(startDate, endDate, reportTZ string) (int, error) {
	var total int
	err := r.db.QueryRow(`
		WITH bounds AS (
			SELECT
				(to_date($1, 'YYYY-MM-DD')::timestamp AT TIME ZONE $3) AS start_ts,
				((to_date($2, 'YYYY-MM-DD')::timestamp + INTERVAL '1 day') AT TIME ZONE $3) AS end_ts
		)
		SELECT COALESCE(SUM(total_amount), 0)
		FROM transactions t
		JOIN bounds b ON true
		WHERE t.created_at >= b.start_ts
			AND t.created_at < b.end_ts
	`, startDate, endDate, reportTZ).Scan(&total)
	return total, err
}

func (r *TransactionRepository) GetTransactionByID(id int, reportTZ string) (*models.Transaction, error) {
	var transaction models.Transaction
	err := r.db.QueryRow(`
		SELECT id, total_amount, TO_CHAR(created_at AT TIME ZONE $2, 'YYYY-MM-DD') AS transaction_date
		FROM transactions
		WHERE id = $1
	`, id, reportTZ).Scan(&transaction.ID, &transaction.TotalAmount, &transaction.TransactionDate)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(`
		SELECT p.name, td.quantity, td.subtotal
		FROM transaction_details td
		JOIN products p ON p.id = td.product_id
		WHERE td.transaction_id = $1
		ORDER BY td.id ASC
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	details := make([]models.TransactionDetail, 0)
	for rows.Next() {
		var detail models.TransactionDetail
		if err := rows.Scan(&detail.ProductName, &detail.Quantity, &detail.Subtotal); err != nil {
			return nil, err
		}
		details = append(details, detail)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	transaction.Details = details
	return &transaction, nil
}

func (r *TransactionRepository) GetTodayReport(reportTZ string) (*models.DailyReport, error) {
	report := &models.DailyReport{
		ProdukTerlaris: make([]models.BestProduct, 0),
	}

	err := r.db.QueryRow(`
		WITH bounds AS (
			SELECT
				(date_trunc('day', now() AT TIME ZONE $1) AT TIME ZONE $1) AS start_ts,
				((date_trunc('day', now() AT TIME ZONE $1) + INTERVAL '1 day') AT TIME ZONE $1) AS end_ts
		)
		SELECT
			COALESCE(SUM(total_amount), 0) AS total_revenue,
			COUNT(*) AS total_transaksi
		FROM transactions t
		JOIN bounds b ON true
		WHERE t.created_at >= b.start_ts
			AND t.created_at < b.end_ts
	`, reportTZ).Scan(&report.TotalRevenue, &report.TotalTransaksi)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(`
		WITH bounds AS (
			SELECT
				(date_trunc('day', now() AT TIME ZONE $1) AT TIME ZONE $1) AS start_ts,
				((date_trunc('day', now() AT TIME ZONE $1) + INTERVAL '1 day') AT TIME ZONE $1) AS end_ts
		),
		product_totals AS (
			SELECT p.id, p.name AS nama, SUM(td.quantity) AS qty_terjual
			FROM transaction_details td
			JOIN transactions t ON t.id = td.transaction_id
			JOIN products p ON p.id = td.product_id
			JOIN bounds b ON true
			WHERE t.created_at >= b.start_ts
				AND t.created_at < b.end_ts
			GROUP BY p.id, p.name
		),
		max_qty AS (
			SELECT MAX(qty_terjual) AS max_qty FROM product_totals
		)
		SELECT pt.nama, pt.qty_terjual
		FROM product_totals pt
		JOIN max_qty mq ON pt.qty_terjual = mq.max_qty
		WHERE mq.max_qty IS NOT NULL
		ORDER BY pt.nama
	`, reportTZ)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var product models.BestProduct
		if err := rows.Scan(&product.Nama, &product.QtyTerjual); err != nil {
			return nil, err
		}
		report.ProdukTerlaris = append(report.ProdukTerlaris, product)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return report, nil
}

func (r *TransactionRepository) GetReportRange(startDate, endDate, reportTZ string) (*models.DailyReport, error) {
	report := &models.DailyReport{
		ProdukTerlaris: make([]models.BestProduct, 0),
	}

	err := r.db.QueryRow(`
		WITH bounds AS (
			SELECT
				(to_date($1, 'YYYY-MM-DD')::timestamp AT TIME ZONE $3) AS start_ts,
				((to_date($2, 'YYYY-MM-DD')::timestamp + INTERVAL '1 day') AT TIME ZONE $3) AS end_ts
		)
		SELECT
			COALESCE(SUM(total_amount), 0) AS total_revenue,
			COUNT(*) AS total_transaksi
		FROM transactions t
		JOIN bounds b ON true
		WHERE t.created_at >= b.start_ts
			AND t.created_at < b.end_ts
	`, startDate, endDate, reportTZ).Scan(&report.TotalRevenue, &report.TotalTransaksi)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(`
		WITH bounds AS (
			SELECT
				(to_date($1, 'YYYY-MM-DD')::timestamp AT TIME ZONE $3) AS start_ts,
				((to_date($2, 'YYYY-MM-DD')::timestamp + INTERVAL '1 day') AT TIME ZONE $3) AS end_ts
		),
		product_totals AS (
			SELECT p.id, p.name AS nama, SUM(td.quantity) AS qty_terjual
			FROM transaction_details td
			JOIN transactions t ON t.id = td.transaction_id
			JOIN products p ON p.id = td.product_id
			JOIN bounds b ON true
			WHERE t.created_at >= b.start_ts
				AND t.created_at < b.end_ts
			GROUP BY p.id, p.name
		),
		max_qty AS (
			SELECT MAX(qty_terjual) AS max_qty FROM product_totals
		)
		SELECT pt.nama, pt.qty_terjual
		FROM product_totals pt
		JOIN max_qty mq ON pt.qty_terjual = mq.max_qty
		WHERE mq.max_qty IS NOT NULL
		ORDER BY pt.nama
	`, startDate, endDate, reportTZ)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var product models.BestProduct
		if err := rows.Scan(&product.Nama, &product.QtyTerjual); err != nil {
			return nil, err
		}
		report.ProdukTerlaris = append(report.ProdukTerlaris, product)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return report, nil
}
