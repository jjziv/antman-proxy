package handlers

import (
	"fmt"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	imageManager "antman-proxy/managers/image"
)

type Config struct {
	ImageManager imageManager.Manager
}

type ImageHandler struct {
	imageManager imageManager.Manager
}

func NewHandler(cfg *Config) (Handler, error) {
	if cfg == nil {
		return nil, fmt.Errorf("ImageHandler Config is nil!")
	}

	if cfg.ImageManager == nil {
		return nil, fmt.Errorf("cfg.ImageManager is nil!")
	}

	return &ImageHandler{
		imageManager: cfg.ImageManager,
	}, nil
}

func validFormats() []string {
	return strings.Split(os.Getenv("VALID_FORMATS"), ",")
}

func (h *ImageHandler) HandleResize(c *gin.Context) {
	url := c.Query("url")
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing URL parameter"})
		return
	}

	if !h.imageManager.IsURLAllowed(url) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL domain"})
		return
	}

	width, _ := strconv.Atoi(c.Query("width"))
	height, _ := strconv.Atoi(c.Query("height"))
	format := c.DefaultQuery("format", "jpeg")
	quality, _ := strconv.Atoi(c.DefaultQuery("quality", fmt.Sprintf("%d", imageManager.DefaultQualityPercent)))

	if width <= 0 && height <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "At least one of width or height must be specified"})
		return
	}

	if (width > 0 && width > 2000) || (height > 0 && height > 2000) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dimensions must be in the range 1-2000"})
		return
	}

	if !slices.Contains(validFormats(), format) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format must be one of jpeg, png, webp"})
		return
	}

	if quality < 1 || quality > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Quality must be between 1 and 100"})
		return
	}

	path, err := h.imageManager.ProcessImage(url, width, height, format, quality)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mimeType := fmt.Sprintf("image/%s", format)
	if format == "jpeg" {
		mimeType = "image/jpeg"
	}
	c.Header("Content-Type", mimeType)

	fileInfo, err := os.Stat(path)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	etag := fmt.Sprintf(`W/"%d"`, fileInfo.ModTime().Unix())
	if c.Request.Header.Get("If-None-Match") == etag {
		c.JSON(http.StatusNotModified, gin.H{})
		return
	}

	c.Header("ETag", etag)
	c.File(path)
}
