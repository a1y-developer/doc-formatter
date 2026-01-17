package middleware

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"sync/atomic"

	"github.com/gin-gonic/gin"
)

type (
	contextKeyTraceID int
)

const (
	TraceIDKey contextKeyTraceID = 0
)

var TraceIDHeader = "x-df-trace"

var (
	prefix string
	reqid  uint64
)

func init() {
	hostname, err := os.Hostname()
	if hostname == "" || err != nil {
		hostname = "localhost"
	}
	var buffer [12]byte
	var b64 string
	for len(b64) < 10 {
		_, err := rand.Read(buffer[:])
		if err != nil {
			return
		}
		b64 = base64.StdEncoding.EncodeToString(buffer[:])
		b64 = strings.NewReplacer("+", "", "/", "").Replace(b64)
	}
	prefix = fmt.Sprintf("%s/%s", hostname, b64[0:10])
}

func TraceID(c *gin.Context) {
	ctx := c.Request.Context()
	traceID := c.GetHeader(TraceIDHeader)
	if traceID == "" {
		id := atomic.AddUint64(&reqid, 1)
		traceID = fmt.Sprintf("%s-%06d", prefix, id)
	}
	ctx = context.WithValue(ctx, TraceIDKey, traceID)
	c.Request = c.Request.WithContext(ctx)
	c.Next()
}

func GetTraceID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if traceID, ok := ctx.Value(TraceIDKey).(string); ok {
		return traceID
	}
	return ""
}
