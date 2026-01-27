package repositories

import (
	"database/sql"
	"kasir-api/models"
	"log"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) GetAllProducts() ([]models.Product, error) {
	rows, err := r.db.Query("SELECT id, name, price, stock FROM products ORDER BY id ASC")
	if err != nil {
		log.Println("Error executing query:", err)
		return nil, err
	}
	defer rows.Close()

	products := make([]models.Product, 0)
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (r *ProductRepository) CreateProduct(product *models.Product) error {
	err := r.db.QueryRow("INSERT INTO products (name, price, stock) VALUES ($1, $2, $3) RETURNING id",
		product.Name, product.Price, product.Stock).Scan(&product.ID)
	return err
}

func (r *ProductRepository) GetProductByID(id int) (*models.Product, error) {
	var p models.Product
	err := r.db.QueryRow("SELECT id, name, price, stock FROM products WHERE id = $1", id).
		Scan(&p.ID, &p.Name, &p.Price, &p.Stock)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProductRepository) UpdateProduct(product *models.Product) error {
	_, err := r.db.Exec("UPDATE products SET name = $1, price = $2, stock = $3 WHERE id = $4",
		product.Name, product.Price, product.Stock, product.ID)
	return err
}

func (r *ProductRepository) DeleteProduct(id int) error {
	_, err := r.db.Exec("DELETE FROM products WHERE id = $1", id)
	return err
}