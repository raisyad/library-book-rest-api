package book

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List() ([]Book, error) {
	return s.repo.FindAll()
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