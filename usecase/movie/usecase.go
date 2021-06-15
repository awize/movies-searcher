package usecase

import (
	"errors"
	"fmt"

	"github.com/awize/movies-searcher/model"
)

type MovieService interface {
	Get(id int) (*model.Movie, error)
	GetAll() ([]model.Movie, error)
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
	fmt.Println(movie)
	if err != nil {
		return nil, err
	}

	if movie == nil {
		return nil, errors.New("not-found")
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
