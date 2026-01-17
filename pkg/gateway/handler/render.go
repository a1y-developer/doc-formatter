package handler

import (
	"time"

	"github.com/a1y/doc-formatter/pkg/gateway/middleware"
	"github.com/gin-gonic/gin"
)

const SuccessMessage = "OK"

func GenerateResponse(c *gin.Context, data any, err error) *Response {
	response := &Response{}
	if err == nil {
		response.Success = true
		response.Message = SuccessMessage
		response.Data = data
	} else {
		response.Success = false
		response.Message = err.Error()
	}

	ctx := c.Request.Context()

	if traceID := middleware.GetTraceID(ctx); len(traceID) > 0 {
		response.TraceID = traceID
	}

	if startTime := middleware.GetStartTime(ctx); !startTime.IsZero() {
		endTime := time.Now()
		response.StartTime = &startTime
		response.EndTime = &endTime
		response.CostTime = Duration(endTime.Sub(startTime))
	}
	return response
}

func FailureResponse(c *gin.Context, err error) *Response {
	return GenerateResponse(c, nil, err)
}

func SuccessResponse(c *gin.Context, data any) *Response {
	return GenerateResponse(c, data, nil)
}
