package member

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List() ([]Member, error) {
	return s.repo.FindAll()
}

func (s *Service) GetByID(id int64) (*Member, error) {
	return s.repo.FindByID(id)
}

func (s *Service) Create(req CreateMemberRequest) (*Member, error) {
	return s.repo.Create(req)
}

func (s *Service) Update(id int64, req UpdateMemberRequest) (*Member, error) {
	return s.repo.Update(id, req)
}

func (s *Service) Delete(id int64) error {
	return s.repo.Delete(id)
}