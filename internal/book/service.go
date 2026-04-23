package book

import (
	"go-library-rest-api/internal/response"
	"math"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List(page, limit int, filter BookFilter) ([]Book, response.PaginationMeta, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	offset := (page - 1) * limit

	totalItems, err := s.repo.Count(filter)
	if err != nil {
		return nil, response.PaginationMeta{}, err
	}

	books, err := s.repo.FindAll(limit, offset, filter)
	if err != nil {
		return nil, response.PaginationMeta{}, err
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(limit)))

	meta := response.PaginationMeta{
		CurrentPage: page,
		PageSize:    limit,
		TotalItems:  totalItems,
		TotalPages:  totalPages,
	}

	return books, meta, nil
}

func (s *Service) GetByID(id int64) (*Book, error) {
	return s.repo.FindByID(id)
}

func (s *Service) Create(req CreateBookRequest) (*Book, error) {
	return s.repo.Create(req)
}

func (s *Service) Update(id int64, req UpdateBookRequest) (*Book, error) {
	return s.repo.Update(id, req)
}

func (s *Service) Delete(id int64) error {
	return s.repo.Delete(id)
}
