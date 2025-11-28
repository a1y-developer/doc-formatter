package app

import (
	"github.com/a1y/doc-formatter/cmd/gateway/options"
	"github.com/a1y/doc-formatter/cmd/util"
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/i18n"
	"k8s.io/kubectl/pkg/util/templates"
)

func NewCmdGateway() *cobra.Command {
	o := options.NewOptions()

	cmd := &cobra.Command{
		Use:   "gateway",
		Short: i18n.T("Start gateway service"),
		Long: templates.LongDesc(i18n.T(`
			Start gateway service for doc-formatter.`)),
		Example: templates.Examples(i18n.T(`
			gateway --bind-address :8080 --auth-service localhost:8081`)),
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
