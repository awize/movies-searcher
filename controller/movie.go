package controller

import (
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
			c.JSON(http.StatusBadRequest, gin.H{"message": "Id parameter not"})
			return
		}
		movie, err := mc.mu.GetMovie(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "something unexpected happened"})
			return
		}
		if movie == nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "not found"})
			return
		}
		c.JSON(http.StatusOK, movie)
	}
}

func (mc *MovieController) SearchMovie() gin.HandlerFunc {
	return func(c *gin.Context) {
		query := c.DefaultQuery("query", "")
		page := c.DefaultQuery("page", "0")

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
