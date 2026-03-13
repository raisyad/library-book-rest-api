package borrowing

import (
	"strings"
	"time"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List() ([]Borrowing, error) {
	return s.repo.FindAll()
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