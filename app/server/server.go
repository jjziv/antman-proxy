package server

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	htmlHandler "antman-proxy/handlers/html"
	imageHandler "antman-proxy/handlers/image"
)

type Config struct {
	Port         string
	HtmlHandler  htmlHandler.Handler
	ImageHandler imageHandler.Handler
}

type Server struct {
	router       *gin.Engine
	htmlHandler  htmlHandler.Handler
	imageHandler imageHandler.Handler
}

func NewServer(cfg *Config) *http.Server {
	// Gin router config with custom HTTP configuration and graceful shutdown of the server built-in
	// Creates a router without any middleware by default
	router := gin.Default()

	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default, gin.DefaultWriter = os.Stdout
	router.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	router.Use(gin.Recovery())

	router.LoadHTMLGlob("static/templates/*")

	s := &http.Server{
		Addr:           ":" + cfg.Port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	router.GET("/", cfg.HtmlHandler.HandleIndex)
	router.GET("/resize", cfg.ImageHandler.HandleResize)

	return s
}
