package main

import (
	"os"
	
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/static"
	log "github.com/sirupsen/logrus"
	
	routers "github.com/mazeForGit/WordlistStorage/routers"
)

func port() string {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "5000"
	}
	return ":" + port
}

func main() {

	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)

	router := gin.Default()
	router.RedirectTrailingSlash = false

	router.LoadHTMLGlob("public/*.html")
	router.Use(static.Serve("/", static.LocalFile("./public", false)))
	router.GET("/index", routers.Index)
	router.NoRoute(routers.NotFoundError)
	router.GET("/500", routers.InternalServerError)
	router.GET("/health", routers.HealthGET)

	// global config
	router.GET("/config", routers.ConfigGET)
	router.POST("/config", routers.ConfigPOST)
	router.GET("/wordliststorage", routers.WordListStorageGET)
	router.PUT("/wordliststorage", routers.WordListStoragePUT)
	
	// session based
	router.GET("/wordlist", routers.WordListGET)
	router.PUT("/wordlist", routers.WordListPUT)
	router.GET("/result", routers.ResultGET)

	log.Info("Starting server on port " + port())
	router.Run(port())
}