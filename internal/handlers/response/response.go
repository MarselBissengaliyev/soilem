package response

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type errorResponse struct {
	Error string `json:"error"`
}

func NewErrorResponse(c *gin.Context, statusCode int, message string) {
	log.Error(message)
	c.JSON(statusCode, errorResponse{message})
}

type DataResponse gin.H

func NewDataResponse(c *gin.Context, statusCode int, data DataResponse) {
	c.JSON(statusCode, data)
}
