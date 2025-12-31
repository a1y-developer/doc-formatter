package options

import (
	"github.com/a1y/doc-formatter/internal/gateway"
	"github.com/a1y/doc-formatter/internal/gateway/route"
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/i18n"
)

type Options struct {
	Address        string
	AuthService    string
	StorageService string
	LogLevel       string
	LogFormat      string
	LogFilePath    string
	LogMaxSize     int
	LogMaxBackups  int
	LogMaxAge      int
	LogCompress    bool
	LogEnvironment string
	LogSample      bool
}

func NewOptions() *Options {
	return &Options{}
}

func (o *Options) Config() (*gateway.Config, error) {
	cfg := gateway.NewConfig()
	cfg.Address = o.Address
	cfg.AuthService = o.AuthService
	cfg.StorageService = o.StorageService
	cfg.Logging.Level = o.LogLevel
	cfg.Logging.Format = o.LogFormat
	cfg.Logging.FilePath = o.LogFilePath
	cfg.Logging.MaxSize = o.LogMaxSize
	cfg.Logging.MaxBackups = o.LogMaxBackups
	cfg.Logging.MaxAge = o.LogMaxAge
	cfg.Logging.Compress = o.LogCompress
	cfg.Logging.Environment = o.LogEnvironment
	cfg.Logging.Sample = o.LogSample
	return cfg, nil
}

func (o *Options) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.Address, "bind-address", ":8080", i18n.T("the address to bind the gateway to"))
	cmd.Flags().StringVar(&o.AuthService, "auth-service", ":8081", i18n.T("the address of the authentication service"))
	cmd.Flags().StringVar(&o.StorageService, "storage-service", ":8082", i18n.T("the address of the storage service"))

	cmd.Flags().StringVar(&o.LogLevel, "log-level", "info", i18n.T("log level: debug, info, warn, error"))
	cmd.Flags().StringVar(&o.LogFormat, "log-format", "json", i18n.T("log format: json or console"))
	cmd.Flags().StringVar(&o.LogFilePath, "log-file", "", i18n.T("log file path; empty means stdout/stderr only"))
	cmd.Flags().IntVar(&o.LogMaxSize, "log-max-size", 100, i18n.T("maximum size in MB of the log file before rotation"))
	cmd.Flags().IntVar(&o.LogMaxBackups, "log-max-backups", 3, i18n.T("maximum number of rotated log files to retain"))
	cmd.Flags().IntVar(&o.LogMaxAge, "log-max-age", 7, i18n.T("maximum number of days to retain old log files"))
	cmd.Flags().BoolVar(&o.LogCompress, "log-compress", true, i18n.T("compress rotated log files"))
	cmd.Flags().StringVar(&o.LogEnvironment, "log-env", "dev", i18n.T("logging environment: dev or prod"))
	cmd.Flags().BoolVar(&o.LogSample, "log-sample", false, i18n.T("enable log sampling (recommended in prod)"))
}

func (o *Options) Complete(args []string) {}

func (o *Options) Validate() error {
	return nil
}

func (o *Options) Run() error {
	config, err := o.Config()
	if err != nil {
		return err
	}

	r, err := route.NewRouter(config)
	if err != nil {
		return err
	}

	return r.Run(config.Address)
}
