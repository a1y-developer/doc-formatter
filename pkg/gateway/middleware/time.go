package middleware

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

var StartTimeKey = &contextKey{"startTime"}

func Time(c *gin.Context) {
	ctx := c.Request.Context()
	if GetStartTime(ctx).IsZero() {
		ctx = context.WithValue(ctx, StartTimeKey, time.Now())
	}
	c.Request = c.Request.WithContext(ctx)
	c.Next()
}

func GetStartTime(ctx context.Context) time.Time {
	if ctx == nil {
		return time.Time{}
	}
	if startTime, ok := ctx.Value(StartTimeKey).(time.Time); ok {
		return startTime
	}
	return time.Time{}
}
