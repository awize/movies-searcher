package service

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/awize/movies-searcher/model"
)

type MovieService struct {
	fileR *os.File
	fileW *os.File
}

func NewMoviesService(fileR *os.File, fileW *os.File) *MovieService {
	return &MovieService{
		fileR,
		fileW,
	}
}

func (r *MovieService) readRecords() []model.Movie {
	fileInfo, err := r.fileR.Stat()

	fmt.Println(fileInfo.ModTime(), err)
	r.fileR.Seek(0, io.SeekStart)
	csvLines, err := csv.NewReader(r.fileR).ReadAll()
	records := make([]model.Movie, 0)
	if err != nil {
		fmt.Println(err)
	}

	for index, record := range csvLines {
		if index == 0 {
			continue
		}

		movie := model.Movie{}
		movie.ID, _ = strconv.Atoi(record[1][2:])
		movie.Title = record[2]
		movie.Year, _ = strconv.Atoi(record[0])
		movie.Budget, _ = strconv.Atoi(record[6])
		movie.Domgross, _ = strconv.Atoi(record[7])

		records = append(records, movie)
	}

	return records
}

func (r *MovieService) Get(id int) (*model.Movie, error) {
	movies := r.readRecords()
	fmt.Println(movies)
	for i := 0; i < len(movies); i++ {
		if movies[i].ID == id {
			return &movies[i], nil
		}
	}
	return &model.Movie{}, errors.New("movie not found")
}

func (r *MovieService) GetAll() ([]model.Movie, error) {
	movies := r.readRecords()
	return movies, nil
}
