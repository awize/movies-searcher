package service

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

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

type filterFn func(model.Movie, map[string][]string) bool

type MovieService struct {
	file   *os.File
	csvr   *csv.Reader
	csvw   *csv.Writer
	client *resty.Client
	config *config.Config
}

func NewMoviesService(file *os.File, fileR *csv.Reader, fileW *csv.Writer, client *resty.Client, config *config.Config) *MovieService {
	return &MovieService{
		file,
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

	movie, err := r.getExternalMovie(id)
	if err != nil {
		return &model.Movie{}, fmt.Errorf("get %v: %v", id, err)
	}

	r.csvw.Write(r.getMovieValues(movie))
	r.csvw.Flush()
	return &movie, nil
}

func (r *MovieService) GetAll() ([]model.Movie, error) {
	movies := r.readMovies()
	return movies, nil
}

func (r *MovieService) GetFilteredMovies(params map[string][]string) ([]model.Movie, error) {
	movies := r.readMovies()
	filtersByName := map[string]filterFn{
		"genre": genreFilter,
		"lang":  langFilter,
	}

	filtersToApply := []filterFn{}

	for paramName, _ := range params {
		if filter := filtersByName[paramName]; filter != nil {
			filtersToApply = append(filtersToApply, filter)
		}
	}

	jobs := 2
	if jobParam := params["jobs"]; len(jobParam) > 0 {
		jobString := jobParam[0]
		jobsNumber, err := strconv.Atoi(jobString)
		if err == nil {
			jobs = jobsNumber
		}
	}

	// Lets divide the movies array in two by having 2 go routines
	c := make(chan []model.Movie)

	moviesToProcess := len(movies) / jobs
	startTime := time.Now()
	for i := 0; i < jobs; i++ {
		start := moviesToProcess * i
		end := moviesToProcess * (i + 1)
		isLastJob := i == jobs-1
		if isLastJob && end < len(movies) {
			end += 1
		}
		fmt.Println("i:", i, " moviesToProcess: ", moviesToProcess, " movies len", len(movies), "start: ", start, "end:", end)
		go func() {
			c <- filterMovies(movies[start:end], filtersToApply, params)
		}()
	}
	filteredMovies := []model.Movie{}
	for i := 0; i < jobs; i++ {
		moviesComing := <-c
		if len(moviesComing) > 0 {
			filteredMovies = append(filteredMovies, moviesComing...)
		}
		// select {
		// case movies := <-c:
		// 	if len(movies) > 0 {
		// 		filteredMovies = append(filteredMovies, movies...)
		// 	}
		// case <- time.After(800 * time.Millisecond)
		// 	return
		// }
	}
	elapsed := time.Since(startTime)
	fmt.Printf("Filtering movies took %s", elapsed)
	return filteredMovies, nil
}

func (r *MovieService) SearchMovies(query string, page int) ([]byte, error) {
	pageString := strconv.Itoa(page)
	resp, err := r.client.R().
		SetQueryParams(map[string]string{
			"api_key":       r.config.MoviesAPI.ApiKey,
			"language":      r.config.MoviesAPI.Defaults["lang"],
			"include_adult": "false",
			"page":          pageString,
			"query":         query,
		}).
		Get(fmt.Sprint(r.config.MoviesAPI.BaseUrl, "/search/movie"))

	if err != nil {
		fmt.Println("Response Info:")
		fmt.Println("  Error      :", err)
		fmt.Println("  Status Code:", resp.StatusCode())
		fmt.Println("  Status     :", resp.Status())
		fmt.Println("  Proto      :", resp.Proto())
		fmt.Println("  Time       :", resp.Time())
		fmt.Println("  Received At:", resp.ReceivedAt())
		return nil, err
	}

	return resp.Body(), nil
}

func (r *MovieService) readMovies() []model.Movie {
	r.file.Seek(0, io.SeekStart)
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

	statusCode := resp.StatusCode()
	if statusCode == http.StatusNotFound {
		return model.Movie{}, fmt.Errorf("external api: %v", model.ErrorNotFound)
	}

	if err != nil || statusCode != http.StatusOK {
		return model.Movie{}, fmt.Errorf("external api: %v", model.ErrorUnexpected)
	}

	movie := model.Movie{}

	if err := json.Unmarshal(resp.Body(), &movie); err != nil {
		return model.Movie{}, fmt.Errorf("getExternalMovie: %v", model.ErrorParsing)
	}

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
	if movieLine[2] != "" {
		for _, genreLine := range strings.Split(movieLine[2], GENRES_SEPARATOR) {
			genreAttrs := strings.Split(genreLine, GENRES_ATTRS_SEPARATOR)
			ID, _ := strconv.Atoi(genreAttrs[0])
			Name := transformContent(genreAttrs[1], false)
			movie.Genres = append(movie.Genres, model.Genre{ID: ID, Name: Name})
		}
	}

	movie.OriginalLanguage = transformContent(movieLine[3], false)
	movie.OriginalTitle = transformContent(movieLine[4], false)
	movie.Overview = transformContent(movieLine[5], false)
	movie.PosterPath = transformContent(movieLine[6], false)
	movie.ReleaseDate = transformContent(movieLine[7], false)
	movie.Revenue, _ = strconv.Atoi(movieLine[8])
	if movieLine[9] != "" {
		for _, spokenLanguageLine := range strings.Split(movieLine[9], SPOKEN_LANGUAGES_SEPARATOR) {
			spokenLanguageAttrs := strings.Split(spokenLanguageLine, SPOKEN_LANGUAGES_ATTRS_SEPARATOR)
			englishName := transformContent(spokenLanguageAttrs[0], false)
			name := transformContent(spokenLanguageAttrs[1], false)
			movie.SpokenLanguages = append(movie.SpokenLanguages, model.Language{Name: name, EnglishName: englishName})
		}
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

func filterMovies(movies []model.Movie, filters []filterFn, params map[string][]string) []model.Movie {

	filteredMovies := []model.Movie{}
	if len(filters) == 0 {
		return movies
	}
	for _, movie := range movies {
		fitsFilters := true
		for _, filter := range filters {
			if !filter(movie, params) {
				fitsFilters = false
				break
			}
		}
		if fitsFilters {
			filteredMovies = append(filteredMovies, movie)
		}
	}
	return filteredMovies
}

func genreFilter(movie model.Movie, params map[string][]string) bool {
	genreParameter, err := url.PathUnescape(params["genre"][0])
	fmt.Println(genreParameter)
	if err != nil {
		fmt.Println(err)
		return false
	}

	genreNames := strings.Split(genreParameter, ",")

	shouldFilter := len(genreNames) > 0
	if !shouldFilter {
		return true
	}

	if len(movie.Genres) == 0 && shouldFilter {
		fmt.Println("Empty genres and should filter")
		return false
	}

	var genreMap = make(map[string]bool)

	for _, genreName := range genreNames {
		genreMap[genreName] = true
	}

	genresMatched := 0
	for _, movieGenre := range movie.Genres {
		if genreMap[movieGenre.Name] {
			genresMatched += 1
		}
	}
	return len(genreNames) == genresMatched
}

func langFilter(movie model.Movie, params map[string][]string) bool {
	lang := params["lang"][0]
	for _, spokenLanguage := range movie.SpokenLanguages {
		if spokenLanguage.EnglishName == lang {
			return true
		}
	}
	return false
}
