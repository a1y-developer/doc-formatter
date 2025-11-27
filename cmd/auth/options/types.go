package options

import (
	"os"

	"github.com/spf13/pflag"
)

type Options interface {
	// Validate checks Options and return a slice of found error(s)
	Validate() error
	// AddFlags adds flags for a specific Option to the specified FlagSet
	AddFlags(fs *pflag.FlagSet)
}

const (
	DefaultDBPort = 5432
	DefaultPort   = 8081
)

var (
	DBHostEnv      = os.Getenv("AUTH_DB_HOST")
	DBPortEnv      = os.Getenv("AUTH_DB_PORT")
	DBUserEnv      = os.Getenv("AUTH_DB_USER")
	DBPassEnv      = os.Getenv("AUTH_DB_PASS")
	DBNameEnv      = os.Getenv("AUTH_DB_NAME")
	PortEnv        = os.Getenv("AUTH_PORT")
	AutoMigrateEnv = os.Getenv("AUTH_AUTO_MIGRATE")
)
