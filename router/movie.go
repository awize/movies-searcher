package router

import (
	"github.com/gin-gonic/gin"
)

type controller interface {
	GetMovies() gin.HandlerFunc
	GetMovie() gin.HandlerFunc
}

func MakeMovieHandlers(router *gin.Engine, movieController controller) {
	router.GET("/movies", movieController.GetMovies())
	router.GET("/movie/:id", movieController.GetMovie())
}
