package borrowing

import (
	"go-library-rest-api/internal/response"
	"math"
	"strings"
	"time"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List(page, limit int, filter BorrowingFilter) ([]Borrowing, response.PaginationMeta, error) {
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

	borrowings, err := s.repo.FindAll(limit, offset, filter)
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

	return borrowings, meta, nil
}

func (s *Service) GetByID(id int64) (*Borrowing, error) {
	return s.repo.FindByID(id)
}

func (s *Service) Create(req CreateBorrowingRequest) (*Borrowing, error) {
	var dueDate *time.Time

	if req.DueDate != nil && strings.TrimSpace(*req.DueDate) != "" {
		parsed, err := time.Parse("2006-01-02", strings.TrimSpace(*req.DueDate))
		if err != nil {
			return nil, ErrInvalidDueDate
		}
		dueDate = &parsed
	}

	params := CreateBorrowingParams{
		MemberID: req.MemberID,
		BookID:   req.BookID,
		DueDate:  dueDate,
	}

	return s.repo.Create(params)
}

func (s *Service) Return(id int64) (*Borrowing, error) {
	return s.repo.Return(id)
}

func (s *Service) Delete(id int64) error {
	return s.repo.Delete(id)
}
