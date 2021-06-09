package repository

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/awize/movies-searcher/entity"
)

type MoviesDB struct {
	movies []entity.Movie
}

func NewMoviesDB() *MoviesDB {
	instance := &MoviesDB{
		movies: make([]entity.Movie, 0),
	}
	instance.LoadMovies()
	return instance
}

func (r *MoviesDB) LoadMovies() {
	csvFile, err := os.Open("/Users/alexis.jimenez/Wizening/movies-searcher/movies.csv")
	if err != nil {
		fmt.Println(err)
	}
	defer csvFile.Close()

	csvLines, err := csv.NewReader(csvFile).ReadAll()

	if err != nil {
		fmt.Println(err)
	}

	for index, record := range csvLines {
		if index == 0 {
			continue
		}

		movie := entity.Movie{}
		movie.ID, _ = strconv.Atoi(record[1][2:])
		movie.Title = record[2]
		movie.Year, _ = strconv.Atoi(record[0])
		movie.Budget, _ = strconv.Atoi(record[6])
		movie.Domgross, _ = strconv.Atoi(record[7])

		r.movies = append(r.movies, movie)
	}
}

func (r *MoviesDB) Get(id int) (*entity.Movie, error) {
	for i := 0; i < len(r.movies); i++ {
		if r.movies[i].ID == id {
			return &r.movies[i], nil
		}
	}
	return &entity.Movie{}, errors.New("movie not found")
}

func (r *MoviesDB) GetAll() ([]entity.Movie, error) {
	return r.movies, nil
}
