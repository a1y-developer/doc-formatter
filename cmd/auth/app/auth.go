package app

import (
	"github.com/a1y/doc-formatter/cmd/auth/options"
	"github.com/a1y/doc-formatter/cmd/util"
	"github.com/spf13/cobra"

	"k8s.io/kubectl/pkg/util/i18n"
	"k8s.io/kubectl/pkg/util/templates"
)

func NewCmdAuth() *cobra.Command {
	var (
		serverShort = i18n.T(`Start authentication service.`)

		serverLong = i18n.T(`Start authentication service.`)

		serverExample = i18n.T(`
		# Start authentication service
		auth --db-host localhost --db-port 5432 --db-name auth --db-user root --db-pass 123456`)
	)

	o := options.NewAuthOptions()
	cmd := &cobra.Command{
		Use:     "auth",
		Short:   serverShort,
		Long:    templates.LongDesc(serverLong),
		Example: templates.Examples(serverExample),
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
