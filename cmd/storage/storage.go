package storage

import (
	"github.com/a1y/doc-formatter/cmd/storage/options"
	"github.com/a1y/doc-formatter/cmd/util"
	"github.com/spf13/cobra"

	"k8s.io/kubectl/pkg/util/i18n"
	"k8s.io/kubectl/pkg/util/templates"
)

func NewCmdStorage() *cobra.Command {
	var (
		serverShort = i18n.T(`Start storage service.`)

		serverLong = i18n.T(`Start storage service.`)

		serverExample = i18n.T(`
		# Start storage service
		storage --db-host localhost --db-port 5432 --db-name storage --db-user root --db-pass 123456 --s3-endpoint http://localhost:9000 --s3-bucket my-bucket`)
	)

	o := options.NewStorageOptions()
	cmd := &cobra.Command{
		Use:     "storage",
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
