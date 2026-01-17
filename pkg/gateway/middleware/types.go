package middleware

type contextKey struct {
	name string
}

type LoggerOptions struct {
	ServiceName string
	Environment string
	Level       string
	Format      string
	FilePath    string
	MaxSize     int
	MaxBackups  int
	MaxAge      int
	Compress    bool
	Sample      bool
}

var (
	APILoggerKey       = &contextKey{"logger"}
	RunLoggerKey       = &contextKey{"runLogger"}
	RunLoggerBufferKey = &contextKey{"runLoggerBuffer"}
)
