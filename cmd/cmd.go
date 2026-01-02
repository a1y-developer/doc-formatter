package cmd

import (
	"os"

	authapp "github.com/a1y/doc-formatter/cmd/auth"
	gatewayapp "github.com/a1y/doc-formatter/cmd/gateway"
	storageapp "github.com/a1y/doc-formatter/cmd/storage"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericiooptions"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/kubectl/pkg/util/templates"
)

type DfctlOptions struct {
	Arguments []string
	genericiooptions.IOStreams
}

func NewDefaultDfctlCommand() *cobra.Command {
	return NewDefaultDfctlCommandWithArgs(DfctlOptions{
		Arguments: os.Args,
		IOStreams: genericiooptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr},
	})
}

func NewDefaultDfctlCommandWithArgs(o DfctlOptions) *cobra.Command {
	cmd := NewDfctlCmd(o)

	if len(o.Arguments) > 1 {
		cmdPathPieces := o.Arguments[1:]
		if _, _, err := cmd.Find(cmdPathPieces); err == nil {
			// sub command exist
			return cmd
		}
	}

	return cmd
}

func NewDfctlCmd(o DfctlOptions) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "df",
		Short: "dfctl is a command line tool for the Doc Formatter project",
		Long: `dfctl is a command line tool for the Doc Formatter project.
It is used to manage the Doc Formatter project.`,
		SilenceErrors: true,
		Run:           runHelp,
		PersistentPreRunE: func(*cobra.Command, []string) error {
			return initProfiling()
		},
		PersistentPostRunE: func(*cobra.Command, []string) error {
			if err := flushProfiling(); err != nil {
				return err
			}
			return nil
		},
	}
	rootCmd.SetGlobalNormalizationFunc(cliflag.WarnWordSepNormalizeFunc)

	flags := rootCmd.PersistentFlags()

	addProfilingFlags(flags)
	groups := templates.CommandGroups{
		{
			Message: "Gateway Commands",
			Commands: []*cobra.Command{
				gatewayapp.NewCmdGateway(),
			},
		},
		{
			Message: "Auth Commands",
			Commands: []*cobra.Command{
				authapp.NewCmdAuth(),
			},
		},
		{
			Message: "Storage Commands",
			Commands: []*cobra.Command{
				storageapp.NewCmdStorage(),
			},
		},
	}
	groups.Add(rootCmd)

	filters := []string{"options"}

	templates.ActsAsRootCommand(rootCmd, filters, groups...)
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	return rootCmd
}

func runHelp(cmd *cobra.Command, args []string) {
	_ = cmd.Help()
}
