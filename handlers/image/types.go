package handlers

import (
	"github.com/gin-gonic/gin"
)

type Handler interface {
	HandleResize(c *gin.Context)
}
