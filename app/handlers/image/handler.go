package handlers

import (
	"github.com/gin-gonic/gin"
)

type ImageHandler struct {
}

func NewHandler() (Handler, error) {
	return &ImageHandler{}, nil
}

func (h *ImageHandler) HandleResize(c *gin.Context) {

}
