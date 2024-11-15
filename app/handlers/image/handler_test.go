package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	imageManager "antman-proxy/managers/image/mock_manager"
)

const (
	testURL    = "http://imgur.com/image.jpg"
	testWidth  = 100
	testHeight = 100
)

func setupTest(t *testing.T) (*gomock.Controller, *imageManager.MockManager, *gin.Engine) {
	ctrl := gomock.NewController(t)
	mockManager := imageManager.NewMockManager(ctrl)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	t.Setenv("ALLOWED_DOMAINS", "imgur.com,i.imgur.com,flickr.com,staticflickr.com,images.unsplash.com")
	t.Setenv("VALID_FORMATS", "jpeg,png,webp")
	return ctrl, mockManager, router
}

func TestImageHandler_NewHandler(t *testing.T) {
	t.Run("returns error when config is nil", func(t *testing.T) {
		handler, err := NewHandler(nil)
		assert.Error(t, err)
		assert.Nil(t, handler)
		assert.Contains(t, err.Error(), "ImageHandler Config is nil!")
	})

	t.Run("returns error when image manager is nil", func(t *testing.T) {
		handler, err := NewHandler(&Config{
			ImageManager: nil,
		})
		assert.Error(t, err)
		assert.Nil(t, handler)
		assert.Contains(t, err.Error(), "cfg.ImageManager is nil!")
	})

	t.Run("successfully creates handler", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockManager := imageManager.NewMockManager(ctrl)
		handler, err := NewHandler(&Config{
			ImageManager: mockManager,
		})

		assert.NoError(t, err)
		assert.NotNil(t, handler)
	})
}

func TestImageHandler_HandleResize(t *testing.T) {
	ctrl, mockManager, router := setupTest(t)
	defer ctrl.Finish()

	handler, err := NewHandler(&Config{ImageManager: mockManager})
	require.NoError(t, err)

	router.GET("/resize", handler.HandleResize)

	t.Run("missing URL parameter", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/resize?width=100&height=100", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Missing URL parameter")
	})

	t.Run("invalid format parameter", func(t *testing.T) {
		w := httptest.NewRecorder()
		mockManager.EXPECT().IsURLAllowed(testURL).Return(true)
		req := httptest.NewRequest("GET", "/resize?url=http://imgur.com/image.jpg&width=100&height=100&format=invalid", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Format must be one of jpeg, png, webp")
	})

	t.Run("disallowed domain", func(t *testing.T) {
		mockManager.EXPECT().IsURLAllowed("http://unsafe.com/image.jpg").Return(false)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/resize?url=http://unsafe.com/image.jpg&width=100&height=100", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid URL domain")
	})
}

func TestImageHandler_HandleResize_Success(t *testing.T) {
	ctrl, mockManager, router := setupTest(t)
	defer ctrl.Finish()

	tempDir := t.TempDir()
	err := os.WriteFile(filepath.Join(tempDir, "image.jpg"), []byte("image.jpg"), 0777)
	if err != nil {
		return
	}

	handler, err := NewHandler(&Config{ImageManager: mockManager})
	require.NoError(t, err)

	router.GET("/resize", handler.HandleResize)

	t.Run("successful image processing", func(t *testing.T) {
		mockManager.EXPECT().IsURLAllowed(testURL).Return(true)
		mockManager.EXPECT().ProcessImage(
			testURL,
			testWidth,
			testHeight,
			"jpeg",
			80,
		).Return(filepath.Join(tempDir, "image.jpg"), nil)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", fmt.Sprintf("/resize?url=%s&width=%d&height=%d&format=jpeg&quality=80", testURL, testWidth, testHeight), nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("successful processing with default parameters", func(t *testing.T) {
		mockManager.EXPECT().IsURLAllowed(testURL).Return(true)
		mockManager.EXPECT().ProcessImage(
			testURL,
			testWidth,
			testHeight,
			"jpeg", // default format
			85,     // default quality
		).Return(filepath.Join(tempDir, "image.jpg"), nil)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", fmt.Sprintf("/resize?url=%s&width=%d&height=%d", testURL, testWidth, testHeight), nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestImageHandler_HandleResize_ProcessingError(t *testing.T) {
	ctrl, mockManager, router := setupTest(t)
	defer ctrl.Finish()

	handler, err := NewHandler(&Config{ImageManager: mockManager})
	require.NoError(t, err)

	router.GET("/resize", handler.HandleResize)

	t.Run("processing error", func(t *testing.T) {
		mockManager.EXPECT().IsURLAllowed(testURL).Return(true)
		mockManager.EXPECT().ProcessImage(
			testURL,
			testWidth,
			testHeight,
			"webp",
			85,
		).Return("", errors.New("processing failed"))

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", fmt.Sprintf("/resize?url=%s&width=%d&height=%d&format=%s", testURL, testWidth, testHeight, "webp"), nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "processing failed")
	})
}
