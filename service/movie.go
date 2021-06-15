package service

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/awize/movies-searcher/config"
	"github.com/awize/movies-searcher/model"
	"github.com/go-resty/resty/v2"
)

type MovieService struct {
	fileR  *os.File
	fileW  *os.File
	client *resty.Client
	config *config.Config
}

func NewMoviesService(fileR *os.File, fileW *os.File, client *resty.Client, config *config.Config) *MovieService {
	return &MovieService{
		fileR,
		fileW,
		client,
		config,
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
		movie.Budget, _ = strconv.Atoi(record[6])

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

func (r *MovieService) GetExternalMovie(id int) ([]model.Movie, error) {
	resp, err := r.client.R().
		SetQueryParams(map[string]string{
			"api_key":  r.config.MoviesAPI.ApiKey,
			"language": r.config.MoviesAPI.Defaults["lang"],
		}).
		Get(fmt.Sprint(r.config.MoviesAPI.BaseUrl, "/movie", id))

	fmt.Println(resp, err)

	// resp.budget
	// resp.genres
	// resp.id
	// resp.original_language
	// resp.original_title
	// resp.overview
	// resp.poster_path
	// resp.release_date
	// resp.revenue
	// resp.runtime
	// resp.spoken_languages
	// resp.status
	// resp.title
	// resp.vote_average
	// resp.vote_count
	return nil, nil
}
