package main

import (
	"fmt"
	"flag"
	"os"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	model "github.com/mazeForGit/WordlistStorage/model"
	routers "github.com/mazeForGit/WordlistStorage/routers"
)

func main() {

	//
	// define and handle flags
	var flagServerName = flag.String("name", "server", "server name")
	var flagServerPort = flag.String("port", "6001", "server port")
	var fileConfig = flag.String("frConfig", "./data/config.json", "file containing config")
	var fileWordList = flag.String("frWL", "./data/wordList.json", "file containing wordList")
	var fileWordListStorage = flag.String("frWLS", "./data/wordListStorage.json", "file containing wordListStorage")
	flag.Parse()
	
	//
	// handle flags
	if *fileConfig != "" {
		err := model.ReadConfigurationFromFile(*fileConfig)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
	}
	if *fileWordList != "" {
		err := model.ReadWordListFromFile(*fileWordList)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
	}
	if *fileWordListStorage != "" {
		err := model.ReadWordListStorageFromFile(*fileWordListStorage)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
	}
	
	//
	// init middleware
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)

	router := gin.Default()
	router.RedirectTrailingSlash = false

	router.LoadHTMLGlob("public/*.html")
	router.Use(static.Serve("/", static.LocalFile("./public", false)))
	router.GET("/index", routers.Index)
	router.NoRoute(routers.NotFoundError)
	router.GET("/500", routers.InternalServerError)
	router.GET("/reporter", routers.Reporter)
	router.GET("/player", routers.Player)
	router.GET("/health", routers.HealthGET)

	// global config
	router.GET("/config", routers.ConfigGET)
	router.POST("/config", routers.ConfigPOST)
	router.GET("/wordliststorage", routers.WordListStorageGET)
	router.PUT("/wordliststorage", routers.WordListStoragePUT)

	// session based
	router.GET("/domain", routers.DomainGET)
	router.GET("/word", routers.WordGET)
	router.POST("/word", routers.WordPOST)
	router.GET("/wordlist", routers.WordListGET)
	router.PUT("/wordlist", routers.WordListPUT)
	router.GET("/result", routers.ResultGET)

	log.Info("run server name = " + *flagServerName + " on port = " + port(*flagServerPort))
	router.Run(port(*flagServerPort))
}
//
// get port from environment or cli
//
func port(flagServerPort string) string {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = flagServerPort
	}
	return ":" + port
}