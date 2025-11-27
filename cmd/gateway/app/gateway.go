package app

import (
	"github.com/a1y/doc-formatter/cmd/gateway/options"
	"github.com/a1y/doc-formatter/cmd/util"
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/i18n"
	"k8s.io/kubectl/pkg/util/templates"
)

func NewCmdGateway() *cobra.Command {
	var (
		gatewayShort = i18n.T(`Start api gateway.`)

		gatewayLong = i18n.T(`Start api gateway.`)

		gatewayExample = i18n.T(`
		# Start gateway
		gateway --http-addr :8080 --auth-grpc-addr localhost:8081`)
	)

	o := options.NewOptions()
	cmd := &cobra.Command{
		Use:     "gateway",
		Short:   gatewayShort,
		Long:    templates.LongDesc(gatewayLong),
		Example: templates.Examples(gatewayExample),
		RunE: func(_ *cobra.Command, args []string) (err error) {
			defer util.RecoverErr(&err)
			o.Complete(args)
			util.CheckErr(o.Validate())
			util.CheckErr(o.Run())
			return
		},
	}

	o.AddFlags(cmd)

	return cmd
}
