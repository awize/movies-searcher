package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/awize/movies-searcher/usecase/movie"
	"github.com/gin-gonic/gin"
)

func getMovies(service movie.UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		movies, _ := service.GetMovies()
		c.JSON(http.StatusOK, movies)
	}
}

func getMovie(service movie.UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			fmt.Println(err)
		}
		movie, err := service.GetMovie(id)
		if err != nil {
			fmt.Println(err)
		}
		c.JSON(http.StatusOK, movie)
	}
}

func MakeMovieHandlers(router *gin.Engine, service movie.UseCase) {
	router.GET("/movies", getMovies(service))
	router.GET("/movie/:id", getMovie(service))
}
