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
	MaskString    = "******"
	DefaultDBPort = 5432
	DefaultPort   = 3000
)

var (
	DBHostEnv      = os.Getenv("DB_HOST")
	DBPortEnv      = os.Getenv("DB_PORT")
	DBUserEnv      = os.Getenv("DB_USER")
	DBPassEnv      = os.Getenv("DB_PASS")
	DBNameEnv      = os.Getenv("DB_NAME")
	PortEnv        = os.Getenv("PORT")
	AutoMigrateEnv = os.Getenv("AUTO_MIGRATE")
)
