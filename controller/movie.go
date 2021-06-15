package controller

import (
	"net/http"
	"strconv"

	"github.com/awize/movies-searcher/model"
	"github.com/gin-gonic/gin"
)

type MovieUseCase interface {
	GetMovie(id int) (*model.Movie, error)
	GetMovies() ([]model.Movie, error)
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
