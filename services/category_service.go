package services

import (
	"kasir-api/models"
	"kasir-api/repositories"
)

type CategoryService struct {
	repo *repositories.CategoryRepository
}

func NewCategoryService(repo *repositories.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) GetAllCategories(page, limit int) ([]models.Category, error) {
	return s.repo.GetAllCategories(page, limit)
}

func (s *CategoryService) CreateCategory(category *models.Category) error {
	return s.repo.CreateCategory(category)
}

func (s *CategoryService) GetCategoryByID(id int) (*models.Category, error) {
	return s.repo.GetCategoryByID(id)
}

func (s *CategoryService) GetCategoryByIDWithProducts(id int) (*models.Category, error) {
	return s.repo.GetCategoryByIDWithProducts(id)
}

func (s *CategoryService) UpdateCategory(category *models.Category) error {
	return s.repo.UpdateCategory(category)
}

func (s *CategoryService) DeleteCategory(id int) error {
	return s.repo.DeleteCategory(id)
}

func (s *CategoryService) CountCategories() (int, error) {
	return s.repo.CountCategories()
}
