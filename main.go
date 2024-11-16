package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	htmlHandler "antman-proxy/handlers/html"
	imageHandler "antman-proxy/handlers/image"
	cacheManager "antman-proxy/managers/cache"
	imageManager "antman-proxy/managers/image"
	"antman-proxy/server"
)

func main() {
	//err := godotenv.Load(".env")
	//if err != nil {
	//	log.Fatal("Error loading .env file")
	//}

	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	html, err := htmlHandler.NewHandler()
	if err != nil {
		log.Fatal(err)
	}

	maxAge, err := strconv.ParseInt(os.Getenv("CACHE_MAX_AGE"), 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	cache, err := cacheManager.NewManager(&cacheManager.Config{
		CacheDir: os.Getenv("CACHE_DIR"),
		MaxAge:   maxAge,
	})
	if err != nil {
		log.Fatal(err)
	}

	allowedDomains := strings.Split(os.Getenv("ALLOWED_DOMAINS"), ",")
	imgManager, err := imageManager.NewManager(&imageManager.Config{
		AllowedDomains: allowedDomains,
		CacheManager:   cache,
	})
	if err != nil {
		log.Fatal(err)
	}

	image, err := imageHandler.NewHandler(&imageHandler.Config{
		ImageManager: imgManager,
	})
	if err != nil {
		log.Fatal(err)
	}

	s := server.NewServer(&server.Config{
		HtmlHandler:  html,
		ImageHandler: image,
		CacheManager: cache,
		ImageManager: imgManager,
		Port:         port,
	})

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	log.Println("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}
