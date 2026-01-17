package request

import (
	"github.com/gin-gonic/gin"
)

func decode(c *gin.Context, payload interface{}) error {
	return c.ShouldBindJSON(payload)
}
