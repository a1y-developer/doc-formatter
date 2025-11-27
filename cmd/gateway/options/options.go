package options

import (
	"github.com/a1y/doc-formatter/internal/gateway"
	"github.com/a1y/doc-formatter/internal/gateway/route"
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/i18n"
)

type Options struct {
	HTTPAddr     string
	AuthGRPCAddr string
}

func NewOptions() *Options {
	return &Options{}
}

func (o *Options) Config() (*gateway.Config, error) {
	cfg := gateway.NewConfig()
	cfg.HTTPAddr = o.HTTPAddr
	cfg.AuthGRPCAddr = o.AuthGRPCAddr
	return cfg, nil
}

func (o *Options) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.HTTPAddr, "http-addr", ":8080", i18n.T("specify the address to listen on"))
	cmd.Flags().StringVar(&o.AuthGRPCAddr, "auth-grpc-addr", ":8081",
		i18n.T("authentication service address to listen on"))
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

	_, err = route.NewRouter(config)
	if err != nil {
		return err
	}
	return nil
}
