package handlers

import "github.com/gin-gonic/gin"

type Handler interface {
	HandleIndex(c *gin.Context)
}
