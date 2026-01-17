package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleResult(c *gin.Context, err error, data any) {
	if err != nil {
		status := http.StatusInternalServerError
		c.AbortWithStatusJSON(status, FailureResponse(c, err))
		return
	}

	c.JSON(http.StatusOK, SuccessResponse(c, data))
}

func HandleResultWithStatus(c *gin.Context, err error, data any, status int) {
	if err != nil {
		c.AbortWithStatusJSON(status, FailureResponse(c, err))
		return
	}

	c.JSON(status, SuccessResponse(c, data))
}
