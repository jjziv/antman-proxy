package server

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	htmlHandler "antman-proxy/handlers/html"
	imageHandler "antman-proxy/handlers/image"
	cacheManager "antman-proxy/managers/cache"
	imageManager "antman-proxy/managers/image"
	"antman-proxy/middlewares"
)

type Config struct {
	HtmlHandler  htmlHandler.Handler
	ImageHandler imageHandler.Handler
	CacheManager cacheManager.Manager
	ImageManager imageManager.Manager
	Port         string
}

type Server struct {
	htmlHandler  htmlHandler.Handler
	imageHandler imageHandler.Handler
	cacheManager cacheManager.Manager
	imageManager imageManager.Manager
	router       *gin.Engine
}

func NewServer(cfg *Config) *http.Server {
	if cfg == nil {
		log.Fatal("Server Config is nil!")
	}

	if cfg.HtmlHandler == nil {
		log.Fatal("cfg.HtmlHandler is nil!")
	}

	if cfg.ImageHandler == nil {
		log.Fatal("cfg.ImageHandler is nil!")
	}

	if cfg.CacheManager == nil {
		log.Fatal("cfg.CacheManager is nil!")
	}

	if cfg.ImageManager == nil {
		log.Fatal("cfg.ImageManager is nil!")
	}

	// Gin router config with custom HTTP configuration and graceful shutdown of the server built-in
	// Creates a router without any middleware by default
	router := gin.Default()

	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default, gin.DefaultWriter = os.Stdout
	router.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	router.Use(gin.Recovery())

	router.Use(middlewares.Headers())

	capacity, err := strconv.ParseFloat(os.Getenv("REQUEST_CAPACITY"), 64)
	if err != nil {
		capacity = float64(60)
	}

	refillRate, err := strconv.ParseFloat(os.Getenv("REQUEST_REFILL_RATE"), 64)
	if err != nil {
		refillRate = float64(1)
	}

	router.Use(middlewares.RateLimiter(&middlewares.RateLimiterConfig{Capacity: capacity, RefillRate: refillRate}))

	router.LoadHTMLGlob("static/templates/*")

	router.GET("/", cfg.HtmlHandler.HandleIndex)
	router.GET("/resize", cfg.ImageHandler.HandleResize)

	s := &http.Server{
		Addr:           ":" + cfg.Port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return s
}
