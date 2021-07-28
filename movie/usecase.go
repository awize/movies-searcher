package usecase

import (
	"fmt"

	"github.com/awize/movies-searcher/model"
)

type MovieService interface {
	Get(id int) (*model.Movie, error)
	GetAll() ([]model.Movie, error)
	SearchMovies(query string, page int) ([]byte, error)
	GetFilteredMovies(params map[string][]string) ([]model.Movie, error)
}

type MovieUsecase struct {
	service MovieService
}

func NewMovieUsecase(service MovieService) *MovieUsecase {
	return &MovieUsecase{
		service: service,
	}
}

func (u *MovieUsecase) GetMovie(id int) (*model.Movie, error) {
	movie, err := u.service.Get(id)

	if err != nil {
		return nil, err
	}

	if movie == nil {
		return nil, fmt.Errorf("get movie %v: %v", id, model.ErrorNotFound)
	}

	return movie, nil
}

func (u *MovieUsecase) GetMovies() ([]model.Movie, error) {
	movies, err := u.service.GetAll()
	if err != nil {
		return nil, err
	}
	return movies, nil
}

func (u *MovieUsecase) SearchMovies(query string, page int) ([]byte, error) {
	result, err := u.service.SearchMovies(query, page)
	return result, err
}

func (u *MovieUsecase) FilterMovies(params map[string][]string) ([]model.Movie, error) {
	result, err := u.service.GetFilteredMovies(params)
	fmt.Println("hello")
	return result, err
}
