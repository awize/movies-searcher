package router

import (
	"github.com/gin-gonic/gin"
)

type controller interface {
	GetMovies() gin.HandlerFunc
	GetMovie() gin.HandlerFunc
	SearchMovie() gin.HandlerFunc
	FilterMovies() gin.HandlerFunc
}

func MakeMovieHandlers(router *gin.Engine, movieController controller) {
	router.GET("/movies", movieController.GetMovies())
	router.GET("/movie/:id", movieController.GetMovie())
	router.GET("/search/movie", movieController.SearchMovie())
	router.GET("/filter", movieController.FilterMovies())
}
