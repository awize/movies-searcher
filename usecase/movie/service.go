package movie

import (
	"errors"

	"github.com/awize/movies-searcher/entity"
)

type Service struct {
	repo Repository
}

func NewService(r Repository) *Service {
	return &Service{
		repo: r,
	}
}

func (s *Service) GetMovie(id int) (*entity.Movie, error) {
	movie, err := s.repo.Get(id)
	if err != nil {
		return nil, err
	}

	if movie == nil {
		return nil, errors.New("not-found")
	}

	return movie, nil
}

func (s *Service) GetMovies() ([]entity.Movie, error) {
	movies, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	return movies, nil
}
