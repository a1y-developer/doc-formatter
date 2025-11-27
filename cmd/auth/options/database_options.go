package options

import (
	"fmt"
	"strconv"

	"github.com/a1y/doc-formatter/cmd/auth/util"
	"github.com/a1y/doc-formatter/internal/auth/infra"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/spf13/pflag"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	ErrDBHostNotSpecified = errors.New("--db-host must be specified")
	ErrDBNameNotSpecified = errors.New("--db-name must be specified")
	ErrDBUserNotSpecified = errors.New("--db-user must be specified")
	ErrDBPortNotSpecified = errors.New("--db-port must be specified")
)

// DatabaseOptions holds the database access layer configurations.
type DatabaseOptions struct {
	DBName      string `json:"dbName,omitempty" yaml:"dbName,omitempty"`
	DBUser      string `json:"dbUser,omitempty" yaml:"dbUser,omitempty"`
	DBPassword  string `json:"dbPassword,omitempty" yaml:"dbPassword,omitempty"`
	DBHost      string `json:"dbHost,omitempty" yaml:"dbHost,omitempty"`
	DBPort      int    `json:"dbPort,omitempty" yaml:"dbPort,omitempty"`
	AutoMigrate bool   `json:"autoMigrate,omitempty" yaml:"autoMigrate,omitempty"`
}

// InstallDB uses the run options to generate and open a db session.
func (o *DatabaseOptions) InstallDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		o.DBHost, o.DBUser, o.DBPassword, o.DBName,
		o.DBPort, "disable", "Asia/Ho_Chi_Minh",
	)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

// ApplyTo uses the run options to generate and open a db session.
func (o *DatabaseOptions) ApplyTo(db **gorm.DB) error {
	d, err := o.InstallDB()
	if err != nil {
		return err
	}
	// TODO(hoangndst): add metric push addr
	*db = d

	if *db != nil && o.AutoMigrate {
		if err := infra.AutoMigrate(*db); err != nil {
			logrus.Fatalf("Failed to auto migrate: %+v", err)
		}
	}
	return nil
}

// Validate checks validation of DatabaseOptions
func (o *DatabaseOptions) Validate() error {
	var errs []error
	if len(o.DBHost) == 0 {
		errs = append(errs, ErrDBHostNotSpecified)
	}
	if len(o.DBName) == 0 {
		errs = append(errs, ErrDBNameNotSpecified)
	}
	if len(o.DBUser) == 0 {
		errs = append(errs, ErrDBUserNotSpecified)
	}
	if o.DBPort == 0 {
		errs = append(errs, ErrDBPortNotSpecified)
	}
	if errs != nil {
		err := util.AggregateError(errs)
		return errors.Wrap(err, "invalid db options")
	}
	return nil
}

// AddFlags adds flags related to DB to a specified FlagSet
func (o *DatabaseOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.DBName, "db-name", DBNameEnv, "the database name")
	fs.StringVar(&o.DBUser, "db-user", DBUserEnv, "the user name used to access database")
	fs.StringVar(&o.DBPassword, "db-pass", DBPassEnv, "the user password used to access database")
	fs.StringVar(&o.DBHost, "db-host", DBHostEnv, "database host")
	dbPort, err := strconv.Atoi(DBPortEnv)
	if err != nil {
		dbPort = DefaultDBPort
	}
	fs.IntVar(&o.DBPort, "db-port", dbPort, "database port")
	autoMigrate, err := strconv.ParseBool(AutoMigrateEnv)
	if err != nil {
		autoMigrate = false
	}
	fs.BoolVar(&o.AutoMigrate, "auto-migrate", autoMigrate, "whether to enable automatic migration")
}
