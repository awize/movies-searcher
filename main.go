package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/awize/movies-searcher/api/handler"
	"github.com/awize/movies-searcher/config"
	"github.com/awize/movies-searcher/repository"
	"github.com/awize/movies-searcher/usecase/movie"
)

func main() {
	moviesDB := repository.NewMoviesDB()

	movieService := movie.NewService(moviesDB)

	r := gin.Default()
	handler.MakeMovieHandlers(r, movieService)
	r.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Working"})
	})
	r.Run(":" + strconv.Itoa(config.API_PORT))
}
