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
	rows, err := r.db.Query(`
		SELECT
			p.id,
			p.name,
			p.price,
			p.stock,
			p.category_id,
			c.id,
			c.name,
			c.description
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		ORDER BY p.id ASC`,
	)
	if err != nil {
		log.Println("Error executing query:", err)
		return nil, err
	}
	defer rows.Close()

	products := make([]models.Product, 0)
	for rows.Next() {
		var p models.Product
		var productCategoryID sql.NullInt64
		var categoryID sql.NullInt64
		var categoryName sql.NullString
		var categoryDescription sql.NullString
		if err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Price,
			&p.Stock,
			&productCategoryID,
			&categoryID,
			&categoryName,
			&categoryDescription,
		); err != nil {
			return nil, err
		}
		if productCategoryID.Valid {
			id := int(productCategoryID.Int64)
			p.CategoryID = &id
		}
		if categoryID.Valid {
			p.Category = &models.Category{
				ID:          int(categoryID.Int64),
				Name:        categoryName.String,
				Description: categoryDescription.String,
			}
		}
		products = append(products, p)
	}
	return products, nil
}

func (r *ProductRepository) CreateProduct(product *models.Product) error {
	err := r.db.QueryRow(
		"INSERT INTO products (name, price, stock, category_id) VALUES ($1, $2, $3, $4) RETURNING id",
		product.Name, product.Price, product.Stock, product.CategoryID,
	).Scan(&product.ID)
	return err
}

func (r *ProductRepository) GetProductByID(id int) (*models.Product, error) {
	var p models.Product
	var productCategoryID sql.NullInt64
	var categoryID sql.NullInt64
	var categoryName sql.NullString
	var categoryDescription sql.NullString
	err := r.db.QueryRow(`
		SELECT
			p.id,
			p.name,
			p.price,
			p.stock,
			p.category_id,
			c.id,
			c.name,
			c.description
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE p.id = $1`,
		id,
	).Scan(
		&p.ID,
		&p.Name,
		&p.Price,
		&p.Stock,
		&productCategoryID,
		&categoryID,
		&categoryName,
		&categoryDescription,
	)
	if err != nil {
		return nil, err
	}
	if productCategoryID.Valid {
		cid := int(productCategoryID.Int64)
		p.CategoryID = &cid
	}
	if categoryID.Valid {
		p.Category = &models.Category{
			ID:          int(categoryID.Int64),
			Name:        categoryName.String,
			Description: categoryDescription.String,
		}
	}
	return &p, nil
}

func (r *ProductRepository) UpdateProduct(product *models.Product) error {
	_, err := r.db.Exec(
		"UPDATE products SET name = $1, price = $2, stock = $3, category_id = $4 WHERE id = $5",
		product.Name, product.Price, product.Stock, product.CategoryID, product.ID,
	)
	return err
}

func (r *ProductRepository) DeleteProduct(id int) error {
	_, err := r.db.Exec("DELETE FROM products WHERE id = $1", id)
	return err
}
