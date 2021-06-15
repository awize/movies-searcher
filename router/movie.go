package router

import (
	"github.com/awize/movies-searcher/controller"
	"github.com/gin-gonic/gin"
)

func MakeMovieHandlers(router *gin.Engine, movieController *controller.MovieController) {
	router.GET("/movies", movieController.GetMovies())
	router.GET("/movie/:id", movieController.GetMovie())
}
