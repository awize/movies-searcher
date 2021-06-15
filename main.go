package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/awize/movies-searcher/config"
	"github.com/awize/movies-searcher/controller"
	"github.com/awize/movies-searcher/router"
	"github.com/awize/movies-searcher/service"
	usecase "github.com/awize/movies-searcher/usecase/movie"
	"github.com/gin-gonic/gin"
)

func main() {

	fileR, err := os.Open("/Users/alexis.jimenez/Wizening/movies-searcher/movies.csv")
	if err != nil {
		fmt.Println(err)
	}
	fileW, err := os.OpenFile("/Users/alexis.jimenez/Wizening/movies-searcher/movies.csv", os.O_RDWR, 0644)
	if err != nil {
		fmt.Println(err)
	}
	defer fileR.Close()
	defer fileW.Close()
	movieService := service.NewMoviesService(fileR, fileW)
	movieUseCase := usecase.NewMovieUsecase(movieService)
	movieController := controller.NewMovieController(movieUseCase)

	r := gin.Default()
	router.MakeMovieHandlers(r, movieController)
	r.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Working"})
	})
	r.Run(":" + strconv.Itoa(config.API_PORT))
}

// Viper -- set config file - file nameyaml
// validator
