package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/awize/movies-searcher/model"
	"github.com/gin-gonic/gin"
)

type MovieUseCase interface {
	GetMovie(id int) (*model.Movie, error)
	GetMovies() ([]model.Movie, error)
	SearchMovies(query string, page int) ([]byte, error)
	FilterMovies(params map[string][]string) ([]model.Movie, error)
}

type MovieController struct {
	mu MovieUseCase
}

func NewMovieController(u MovieUseCase) *MovieController {
	return &MovieController{
		mu: u,
	}
}

func (mc *MovieController) GetMovies() gin.HandlerFunc {
	return func(c *gin.Context) {
		movies, _ := mc.mu.GetMovies()
		c.JSON(http.StatusOK, movies)
	}
}

func (mc *MovieController) GetMovie() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			generateHttpResponse(c, http.StatusBadRequest, fmt.Sprintf("id: %v", model.ErrorMistmatchType.Description))
			return
		}
		movie, err := mc.mu.GetMovie(id)

		if errors.Is(err, model.ErrorUnexpected.Err) {
			generateHttpResponse(c, http.StatusInternalServerError, model.ErrorUnexpected.Description)
			return
		}

		if errors.Is(err, model.ErrorNotFound.Err) {
			generateHttpResponse(c, http.StatusNotFound, model.ErrorNotFound.Description)
			return
		}

		c.JSON(http.StatusOK, movie)
	}
}

func (mc *MovieController) SearchMovie() gin.HandlerFunc {
	return func(c *gin.Context) {
		query := c.DefaultQuery("query", "")
		page := c.DefaultQuery("page", "1")

		pageNumber, err := strconv.Atoi(page)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "page should be a number"})
		}
		result, err := mc.mu.SearchMovies(query, pageNumber)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "something unexpected happened"})
			return
		}
		fmt.Println(result)

		c.Data(http.StatusOK, "application/json", result)
	}
}

func (mc *MovieController) FilterMovies() gin.HandlerFunc {
	return func(c *gin.Context) {
		queries := c.Request.URL.Query()

		movies, err := mc.mu.FilterMovies(queries)

		if errors.Is(err, model.ErrorUnexpected.Err) {
			generateHttpResponse(c, http.StatusInternalServerError, model.ErrorUnexpected.Description)
			return
		}

		c.JSON(http.StatusOK, movies)
	}
}

func generateHttpResponse(c *gin.Context, httpCode int, message string) {
	c.JSON(httpCode, gin.H{"message": message})
}
