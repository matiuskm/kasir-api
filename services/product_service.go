package services

import (
	"kasir-api/models"
	"kasir-api/repositories"
)

type ProductService struct {
	repo *repositories.ProductRepository
}

func NewProductService(repo *repositories.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) GetAllProducts(page, limit int, name string) ([]models.Product, error) {
	return s.repo.GetAllProducts(page, limit, name)
}

func (s *ProductService) CreateProduct(product *models.Product) error {
	return s.repo.CreateProduct(product)
}

func (s *ProductService) GetProductByID(id int) (*models.Product, error) {
	return s.repo.GetProductByID(id)
}

func (s *ProductService) UpdateProduct(product *models.Product) error {
	return s.repo.UpdateProduct(product)
}

func (s *ProductService) DeleteProduct(id int) error {
	return s.repo.DeleteProduct(id)
}

func (s *ProductService) CountProducts(name string) (int, error) {
	return s.repo.CountProducts(name)
}
