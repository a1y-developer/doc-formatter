package options

import (
	"github.com/a1y/doc-formatter/internal/gateway"
	"github.com/a1y/doc-formatter/internal/gateway/route"
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/i18n"
)

type Options struct {
	Address     string
	AuthService string
}

func NewOptions() *Options {
	return &Options{}
}

func (o *Options) Config() (*gateway.Config, error) {
	cfg := gateway.NewConfig()
	cfg.Address = o.Address
	cfg.AuthService = o.AuthService
	return cfg, nil
}

func (o *Options) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.Address, "bind-address", ":8080", i18n.T("the address to bind the gateway to"))
	cmd.Flags().StringVar(&o.AuthService, "auth-service", ":8081", i18n.T("the address of the authentication service"))
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
