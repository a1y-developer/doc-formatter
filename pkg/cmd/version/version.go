package version

import (
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/templates"

	"github.com/a1y/doc-formatter/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/i18n"
)

func NewCmdVersion() *cobra.Command {
	var (
		versionShort   = i18n.T(`Print the Doc Formatter version information for the current context`)
		versionExample = i18n.T(`
		# Print the Doc Formatter version
		df version`)
	)

	o := NewVersionOptions()

	cmd := &cobra.Command{
		Use:     "version",
		Short:   versionShort,
		Long:    templates.LongDesc(versionShort),
		Example: templates.Examples(versionExample),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			defer util.RecoverErr(&err)
			util.CheckErr(o.Validate())
			o.Run()
			return
		},
	}

	cmd.Flags().StringVarP(&o.Output, "output", "o", "",
		i18n.T("Output format. Only json format is supported for now"))

	return cmd
}
