package movie

import "github.com/awize/movies-searcher/entity"

type Repository interface {
	Get(id int) (*entity.Movie, error)
	GetAll() ([]entity.Movie, error)
}

type UseCase interface {
	GetMovie(id int) (*entity.Movie, error)
	GetMovies() ([]entity.Movie, error)
}
