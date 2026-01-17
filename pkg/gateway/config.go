package gateway

type Config struct {
	Address        string
	AuthService    string
	StorageService string
	Logging        LoggingConfig
}

// LoggingConfig holds structured logging configuration for the gateway.
type LoggingConfig struct {
	// Level is the minimum log level: debug, info, warn, error.
	Level string
	// Format is either "json" (recommended) or "console" (for local dev).
	Format string
	// FilePath is the optional log file path. If empty, logs go only to stdout/stderr.
	FilePath string
	// MaxSize is the maximum size in megabytes of the log file before it gets rotated.
	MaxSize int
	// MaxBackups is the maximum number of old log files to retain.
	MaxBackups int
	// MaxAge is the maximum number of days to retain old log files.
	MaxAge int
	// Compress determines if the rotated log files should be compressed.
	Compress bool
	// Environment indicates the runtime environment, e.g. "dev" or "prod".
	Environment string
	// Sample enables log sampling in production to reduce volume.
	Sample bool
}

func NewConfig() *Config {
	return &Config{
		Logging: LoggingConfig{
			Level:       "info",
			Format:      "json",
			FilePath:    "",
			MaxSize:     100,
			MaxBackups:  3,
			MaxAge:      7,
			Compress:    true,
			Environment: "dev",
			Sample:      false,
		},
	}
}
