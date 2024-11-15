package managers

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"net/url"
	"strings"

	"github.com/chai2010/webp"
	"github.com/nfnt/resize"

	cacheManager "antman-proxy/managers/cache"
)

const (
	DefaultQualityPercent = 85
)

type Config struct {
	AllowedDomains []string
	CacheManager   cacheManager.Manager
}

type ImageManager struct {
	allowedDomains []string
	cacheManager   cacheManager.Manager
}

func NewManager(cfg *Config) (*ImageManager, error) {
	if cfg == nil {
		return nil, fmt.Errorf("ImageManager Config is nil!")
	}

	if cfg.AllowedDomains == nil {
		return nil, fmt.Errorf("cfg.AllowedDomains is nil!")
	}

	if cfg.CacheManager == nil {
		return nil, fmt.Errorf("cfg.CacheManager is nil!")
	}

	return &ImageManager{
		allowedDomains: cfg.AllowedDomains,
		cacheManager:   cfg.CacheManager,
	}, nil
}

func (m *ImageManager) IsURLAllowed(imageURL string) bool {
	parsedURL, err := url.Parse(imageURL)
	if err != nil {
		return false
	}

	domain := strings.ToLower(parsedURL.Host)
	for _, allowed := range m.allowedDomains {
		if strings.Contains(domain, allowed) {
			return true
		}
	}
	return false
}

// Determine what the final dimensions should be accounting for original aspect ratio.
func calculateDimensions(image image.Image, width, height int) (int, int) {
	bounds := image.Bounds()
	originalWidth := bounds.Max.X - bounds.Min.X
	originalHeight := bounds.Max.Y - bounds.Min.Y

	aspectRatio := float64(originalWidth) / float64(originalHeight)
	if width == 0 && height > 0 {
		return int(float64(height) / aspectRatio), height
	} else if height == 0 && width > 0 {
		return width, int(float64(width) * aspectRatio)
	} else {
		return width, height
	}
}

func (m *ImageManager) ProcessImage(imageURL string, width, height int, format string, quality int) (string, error) {
	cacheKey := m.generateCacheKey(imageURL, width, height, format)

	cached := m.cacheManager.Get(cacheKey, format)
	if cached != "" {
		return cached, nil
	}

	resp, err := http.Get(imageURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return "", err
	}

	// @TODO: handle resizing by percentage, check cache key generation
	finalWidth, finalHeight := calculateDimensions(img, height, width)
	resizedImage := resize.Resize(uint(finalWidth), uint(finalHeight), img, resize.Lanczos3)

	output := new(bytes.Buffer)

	switch format {
	case "jpeg":
		err = jpeg.Encode(output, resizedImage, &jpeg.Options{Quality: quality})
		if err != nil {
			return "", fmt.Errorf("jpeg.Encode: %s", err)
		}
	case "png":
		err = png.Encode(output, resizedImage)
		if err != nil {
			return "", fmt.Errorf("png.Encode: %s", err)
		}
	case "webp":
		options := &webp.Options{
			Lossless: false,
			Quality:  float32(quality),
		}

		err = webp.Encode(output, resizedImage, options)
		if err != nil {
			return "", fmt.Errorf("webp.Encode: %v", err)
		}
	default:
		return "", fmt.Errorf("unsupported output format: %s", format)
	}

	return m.cacheManager.Set(cacheKey, output.Bytes(), format)
}

func (m *ImageManager) generateCacheKey(url string, width, height int, format string) string {
	data := fmt.Sprintf("%s_%d_%d_%s", url, width, height, format)
	return fmt.Sprintf("%x", md5.Sum([]byte(data)))
}
