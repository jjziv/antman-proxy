package handlers

import (
	"fmt"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"

	imageManager "antman-proxy/managers/image"
)

type Config struct {
	ImageManager imageManager.Manager
	WorkerPool   *WorkerPool
}

type ImageHandler struct {
	imageManager imageManager.Manager
	mu           sync.RWMutex
	workerPool   *WorkerPool
}

func NewHandler(cfg *Config) (Handler, error) {
	if cfg == nil {
		return nil, fmt.Errorf("ImageHandler Config is nil!")
	}

	if cfg.ImageManager == nil {
		return nil, fmt.Errorf("cfg.ImageManager is nil!")
	}

	if cfg.WorkerPool == nil {
		return nil, fmt.Errorf("cfg.WorkerPool is nil!")
	}

	return &ImageHandler{
		imageManager: cfg.ImageManager,
		mu:           sync.RWMutex{},
		workerPool:   cfg.WorkerPool,
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

	var path string
	var err error
	var resultMu sync.Mutex

	width, _ := strconv.Atoi(c.Query("width"))
	height, _ := strconv.Atoi(c.Query("height"))
	format := c.DefaultQuery("format", "jpeg")
	quality, _ := strconv.Atoi(c.DefaultQuery("quality", fmt.Sprintf("%d", imageManager.DefaultQualityPercent)))

	h.workerPool.Submit(func() {
		h.mu.RLock()
		defer h.mu.RUnlock()

		if !h.imageManager.IsURLAllowed(url) {
			resultMu.Lock()
			err = fmt.Errorf("Invalid URL domain")
			resultMu.Unlock()
			return
		}

		if width <= 0 && height <= 0 {
			err = fmt.Errorf("At least one of width or height must be specified")
			return
		}

		if (width > 0 && width > 2000) || (height > 0 && height > 2000) {
			err = fmt.Errorf("Dimensions must be in the range 1-2000")
			return
		}

		if !slices.Contains(validFormats(), format) {
			err = fmt.Errorf("Format must be one of jpeg, png, webp")
			return
		}

		if quality < 1 || quality > 100 {
			err = fmt.Errorf("Quality must be between 1 and 100")
			return
		}

		resultMu.Lock()
		path, err = h.imageManager.ProcessImage(url, width, height, format, quality)
		if err != nil {
			err = fmt.Errorf(err.Error())
			return
		}
		resultMu.Unlock()

		mimeType := fmt.Sprintf("image/%s", format)
		if format == "jpeg" {
			mimeType = "image/jpeg"
		}
		c.Header("Content-Type", mimeType)

		fileInfo, err := os.Stat(path)
		if err != nil {
			err = fmt.Errorf(err.Error())
			return
		}

		etag := fmt.Sprintf(`W/"%d"`, fileInfo.ModTime().Unix())
		if c.Request.Header.Get("If-None-Match") == etag {
			c.JSON(http.StatusNotModified, gin.H{})
			return
		}

		c.Header("ETag", etag)
	})

	h.workerPool.Wait()

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.File(path)
}
