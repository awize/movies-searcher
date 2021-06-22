package service

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/awize/movies-searcher/config"
	"github.com/awize/movies-searcher/model"
	"github.com/go-resty/resty/v2"
)

var CHARACTERS_MAPPER = map[string]string{
	",": "</n>",
}

const GENRES_SEPARATOR = "-"
const GENRES_ATTRS_SEPARATOR = ":"
const SPOKEN_LANGUAGES_SEPARATOR = GENRES_SEPARATOR
const SPOKEN_LANGUAGES_ATTRS_SEPARATOR = GENRES_ATTRS_SEPARATOR

type MovieService struct {
	csvr   *csv.Reader
	csvw   *csv.Writer
	client *resty.Client
	config *config.Config
}

func NewMoviesService(fileR *csv.Reader, fileW *csv.Writer, client *resty.Client, config *config.Config) *MovieService {
	return &MovieService{
		fileR,
		fileW,
		client,
		config,
	}
}

func (r *MovieService) Get(id int) (*model.Movie, error) {
	movies := r.readMovies()

	for i := 0; i < len(movies); i++ {
		if movies[i].ID == id {
			return &movies[i], nil
		}
	}

	fmt.Println("Lets go to external API")
	movie, err := r.getExternalMovie(id)
	if err != nil {
		return nil, errors.New("something in external movie happened")
	}
	fmt.Println("Error writing in csv file", r.csvw.Write(r.getMovieValues(movie)))
	r.csvw.Flush()
	return &movie, nil
}

func (r *MovieService) GetAll() ([]model.Movie, error) {
	movies := r.readMovies()
	return movies, nil
}

func (r *MovieService) SearchMovies(query string, page int) ([]byte, error) {
	pageString := strconv.Itoa(page)
	fmt.Println(
		"page:", pageString,
		"query:", query)
	resp, err := r.client.R().
		SetQueryParams(map[string]string{
			"api_key":       r.config.MoviesAPI.ApiKey,
			"language":      r.config.MoviesAPI.Defaults["lang"],
			"include_adult": "false",
			"page":          pageString,
			"query":         query,
		}).
		Get(fmt.Sprint(r.config.MoviesAPI.BaseUrl, "search/movie"))

	if err != nil {
		return nil, err
	}
	fmt.Println("Response Info:")
	fmt.Println("  Error      :", err)
	fmt.Println("  Status Code:", resp.StatusCode())
	fmt.Println("  Status     :", resp.Status())
	fmt.Println("  Proto      :", resp.Proto())
	fmt.Println("  Time       :", resp.Time())
	fmt.Println("  Received At:", resp.ReceivedAt())
	fmt.Println("  Body       :\n", resp)
	fmt.Println()
	return resp.Body(), nil
}

func (r *MovieService) readMovies() []model.Movie {

	csvLines, err := r.csvr.ReadAll()
	movies := make([]model.Movie, 0)
	if err != nil {
		fmt.Println(err)
	}

	for _, movieLine := range csvLines {
		movies = append(movies, r.parseMovieLine(movieLine))
	}

	return movies
}

func (r *MovieService) getExternalMovie(id int) (model.Movie, error) {
	resp, err := r.client.R().
		SetQueryParams(map[string]string{
			"api_key":  r.config.MoviesAPI.ApiKey,
			"language": r.config.MoviesAPI.Defaults["lang"],
		}).
		Get(fmt.Sprint(r.config.MoviesAPI.BaseUrl, "/movie/", id))
	movie := model.Movie{}

	if err := json.Unmarshal(resp.Body(), &movie); err != nil {
		fmt.Println("err in unmarshall", err)
	}
	fmt.Println("err", err)
	return movie, nil
}

func (r *MovieService) getMovieValues(movie model.Movie) []string {
	movieValues := []string{}
	movieVal := reflect.ValueOf(movie)

	for i := 0; i < movieVal.NumField(); i++ {
		fieldName := movieVal.Type().Field(i).Name
		fieldInterface := movieVal.Field(i).Interface()
		value := ""
		switch {
		case fieldName == "Genres":
			formattedGenres := []string{}
			if genres, ok := fieldInterface.([]model.Genre); ok {
				for _, genre := range genres {
					formattedGenres = append(formattedGenres, fmt.Sprintf("%v:%v", genre.ID, genre.Name))
				}
			}

			value = strings.Join(formattedGenres, "-")

		case fieldName == "SpokenLanguages":
			formattedSpokenLanguages := []string{}
			if languages, ok := fieldInterface.([]model.Language); ok {
				for _, language := range languages {
					formattedSpokenLanguages = append(formattedSpokenLanguages, fmt.Sprintf("%v:%v", language.EnglishName, language.Name))
				}
			}
			value = strings.Join(formattedSpokenLanguages, SPOKEN_LANGUAGES_SEPARATOR)
		default:
			value = fmt.Sprintf("%v", fieldInterface)
		}
		movieValues = append(movieValues, transformContent(value, true))
	}

	return movieValues
}

func (r *MovieService) parseMovieLine(movieLine []string) model.Movie {
	movie := model.Movie{}
	movie.ID, _ = strconv.Atoi(movieLine[0])
	movie.Budget, _ = strconv.Atoi(movieLine[1])
	for _, genreLine := range strings.Split(movieLine[2], GENRES_SEPARATOR) {
		genreAttrs := strings.Split(genreLine, GENRES_ATTRS_SEPARATOR)
		ID, _ := strconv.Atoi(genreAttrs[0])
		Name := transformContent(genreAttrs[1], false)
		movie.Genres = append(movie.Genres, model.Genre{ID: ID, Name: Name})
	}
	movie.OriginalLanguage = transformContent(movieLine[3], false)
	movie.OriginalTitle = transformContent(movieLine[4], false)
	movie.Overview = transformContent(movieLine[5], false)
	movie.PosterPath = transformContent(movieLine[6], false)
	movie.ReleaseDate = transformContent(movieLine[7], false)
	movie.Revenue, _ = strconv.Atoi(movieLine[8])
	for _, spokenLanguageLine := range strings.Split(movieLine[9], SPOKEN_LANGUAGES_SEPARATOR) {
		spokenLanguageAttrs := strings.Split(spokenLanguageLine, SPOKEN_LANGUAGES_ATTRS_SEPARATOR)
		englishName := transformContent(spokenLanguageAttrs[0], false)
		name := transformContent(spokenLanguageAttrs[1], false)
		movie.SpokenLanguages = append(movie.SpokenLanguages, model.Language{Name: name, EnglishName: englishName})
	}
	movie.Status = transformContent(movieLine[10], false)
	movie.Title = transformContent(movieLine[11], false)
	movie.VoteAverage, _ = strconv.ParseFloat(movieLine[12], 64)
	movie.VoteCount, _ = strconv.Atoi(movieLine[13])

	return movie
}

func transformContent(content string, shouldScapeCharacters bool) string {
	scapedContent := content
	for character, scapedCharacter := range CHARACTERS_MAPPER {
		if shouldScapeCharacters {
			scapedContent = strings.ReplaceAll(scapedContent, character, scapedCharacter)
		} else {
			scapedContent = strings.ReplaceAll(scapedContent, scapedCharacter, character)
		}
	}
	return scapedContent
}
