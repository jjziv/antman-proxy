package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HtmlHandler struct {
}

func NewHandler() (Handler, error) {
	return &HtmlHandler{}, nil
}

func (h *HtmlHandler) HandleIndex(c *gin.Context) {
	c.Header("Cache-Control", "public, max-age=300")
	c.HTML(http.StatusOK, "index.html", nil)
}
