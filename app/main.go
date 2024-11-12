package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	html "antman-proxy/handlers/html"
	image "antman-proxy/handlers/image"
	cache "antman-proxy/managers/cache"
	"antman-proxy/server"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	maxAge, err := strconv.ParseInt(os.Getenv("CACHE_MAX_AGE"), 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	cacheManager, err := cache.NewManager(&cache.Config{
		CacheDir: os.Getenv("CACHE_DIR"),
		MaxAge:   maxAge,
	})
	if err != nil {
		log.Fatal(err)
	}

	htmlHandler, err := html.NewHandler()
	if err != nil {
		log.Fatal(err)
	}

	imageHandler, err := image.NewHandler()
	if err != nil {
		log.Fatal(err)
	}

	s := server.NewServer(&server.Config{
		CacheManager: cacheManager,
		HtmlHandler:  htmlHandler,
		ImageHandler: imageHandler,
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
