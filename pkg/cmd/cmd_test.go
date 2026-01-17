package cmd

import (
	"bytes"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
	"k8s.io/cli-runtime/pkg/genericiooptions"
)

func TestNewDefaultDfctlCommand_WithNoSubcommand(t *testing.T) {
	t.Parallel()

	origArgs := os.Args
	t.Cleanup(func() { os.Args = origArgs })

	os.Args = []string{"df"}

	cmd := NewDefaultDfctlCommand()
	require.NotNil(t, cmd)
	require.Equal(t, "df", cmd.Use)
}

func TestNewDefaultDfctlCommand_WithExistingSubcommand(t *testing.T) {
	t.Parallel()

	origArgs := os.Args
	t.Cleanup(func() { os.Args = origArgs })

	os.Args = []string{"df", "gateway"}

	cmd := NewDefaultDfctlCommand()
	require.NotNil(t, cmd)

	sub, _, err := cmd.Find([]string{"gateway"})
	require.NoError(t, err)
	require.NotNil(t, sub)
	require.Equal(t, "gateway", sub.Name())
}

func TestNewDefaultDfctlCommandWithArgs_FallsBackToRootOnUnknownSubcommand(t *testing.T) {
	t.Parallel()

	args := []string{"df", "unknown-subcommand"}
	o := DfctlOptions{
		Arguments: args,
		IOStreams: genericiooptions.IOStreams{
			In:     bytes.NewBuffer(nil),
			Out:    &bytes.Buffer{},
			ErrOut: &bytes.Buffer{},
		},
	}

	cmd := NewDefaultDfctlCommandWithArgs(o)
	require.NotNil(t, cmd)

	// The root command should not have a direct subcommand named "unknown-subcommand"
	_, _, err := cmd.Find([]string{"unknown-subcommand"})
	require.Error(t, err)
}

func TestNewDfctlCmd_StructureAndHooks(t *testing.T) {
	t.Parallel()

	o := DfctlOptions{
		Arguments: []string{"df"},
		IOStreams: genericiooptions.IOStreams{
			In:     bytes.NewBuffer(nil),
			Out:    &bytes.Buffer{},
			ErrOut: &bytes.Buffer{},
		},
	}

	cmd := NewDfctlCmd(o)
	require.NotNil(t, cmd)
	require.Equal(t, "df", cmd.Use)
	require.NotNil(t, cmd.PersistentPreRunE)
	require.NotNil(t, cmd.PersistentPostRunE)

	// Exercise the help runner directly
	runHelp(cmd, nil)
}

func TestRunHelp_CallsCobraHelp(t *testing.T) {
	t.Parallel()

	var called bool
	testCmd := &cobra.Command{
		Use: "df",
	}
	testCmd.SetHelpFunc(func(*cobra.Command, []string) {
		called = true
	})

	runHelp(testCmd, nil)
	require.True(t, called, "expected custom Help function to be called")
}
