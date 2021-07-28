package main

import (
	"crypto/tls"
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/awize/movies-searcher/config"
	"github.com/awize/movies-searcher/controller"
	"github.com/awize/movies-searcher/router"
	"github.com/awize/movies-searcher/service"
	usecase "github.com/awize/movies-searcher/usecase/movie"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
)

func main() {
	public := viper.New()
	public.SetConfigFile(config.ConfigFile)
	if err := public.ReadInConfig(); err != nil {
		fmt.Println(err)
	}
	config := &config.Config{}
	fmt.Println(public.AllKeys())
	err := public.Unmarshal(config)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("config:", config)
	workspaceDir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	fileNameAbs := filepath.Join(workspaceDir, config.MoviesFileName)
	fileR, err := os.Open(fileNameAbs)
	if err != nil {
		fmt.Println(err)
	}
	fileW, err := os.OpenFile(fileNameAbs, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println(err)
	}
	defer fileR.Close()
	defer fileW.Close()
	client := resty.New()
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	csvr := csv.NewReader(fileR)
	csvw := csv.NewWriter(fileW)
	movieService := service.NewMoviesService(fileR, csvr, csvw, client, config)
	movieUseCase := usecase.NewMovieUsecase(movieService)
	movieController := controller.NewMovieController(movieUseCase)

	r := gin.Default()
	router.MakeMovieHandlers(r, movieController)
	// TODO: Change to api stats
	r.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Working"})
	})
	r.Run(":" + strconv.Itoa(config.Port))
}

// Viper -- set config file - file nameyaml
// validator
