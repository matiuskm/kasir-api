package repositories

import (
	"database/sql"
	"kasir-api/models"
)

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) GetAllCategories(page, limit int) ([]models.Category, error) {
	offset := (page - 1) * limit
	rows, err := r.db.Query(
		"SELECT id, name, description FROM categories ORDER BY id ASC LIMIT $1 OFFSET $2",
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]models.Category, 0)
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Description); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}

func (r *CategoryRepository) CreateCategory(category *models.Category) error {
	err := r.db.QueryRow("INSERT INTO categories (name, description) VALUES ($1, $2) RETURNING id",
		category.Name, category.Description).Scan(&category.ID)
	return err
}

func (r *CategoryRepository) GetCategoryByID(id int) (*models.Category, error) {
	var c models.Category
	err := r.db.QueryRow("SELECT id, name, description FROM categories WHERE id = $1", id).
		Scan(&c.ID, &c.Name, &c.Description)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CategoryRepository) GetCategoryByIDWithProducts(id int) (*models.Category, error) {
	category, err := r.GetCategoryByID(id)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(`
		SELECT
			p.id,
			p.name,
			p.price,
			p.stock,
			p.category_id
		FROM products p
		WHERE p.category_id = $1
		ORDER BY p.id ASC`,
		id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]models.Product, 0)
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CategoryID); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	category.Products = products
	return category, nil
}

func (r *CategoryRepository) UpdateCategory(category *models.Category) error {
	_, err := r.db.Exec("UPDATE categories SET name = $1, description = $2 WHERE id = $3",
		category.Name, category.Description, category.ID)
	return err
}

func (r *CategoryRepository) DeleteCategory(id int) error {
	_, err := r.db.Exec("DELETE FROM categories WHERE id = $1", id)
	return err
}

func (r *CategoryRepository) CountCategories() (int, error) {
	var total int
	err := r.db.QueryRow("SELECT COUNT(*) FROM categories").Scan(&total)
	return total, err
}
