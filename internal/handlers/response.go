package handlers

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type errorResponse struct {
	Error string `json:"error"`
}

func newErrorResponse(c *gin.Context, statusCode int, message string) {
	log.Error(message)
	c.JSON(statusCode, errorResponse{message})
}

type dataResponse gin.H

func newDataResponse(c *gin.Context, statusCode int, data dataResponse) {
	c.JSON(statusCode, data)
}
