package util

import (
	"bytes"
	"errors"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

func TestRecoverErr_CapturesPanicError(t *testing.T) {
	t.Parallel()

	var err error
	func() {
		defer RecoverErr(&err)
		panic(errors.New("boom"))
	}()

	require.EqualError(t, err, "boom")
}

func TestRecoverErr_CapturesPanicString(t *testing.T) {
	t.Parallel()

	var err error
	func() {
		defer RecoverErr(&err)
		panic("string-panic")
	}()

	require.EqualError(t, err, "string-panic")
}

func TestRecoverErr_UnknownType(t *testing.T) {
	t.Parallel()

	var err error
	func() {
		defer RecoverErr(&err)
		panic(123)
	}()

	require.EqualError(t, err, "unknown panic")
}

func TestCheckErr_PanicsOnError(t *testing.T) {
	t.Parallel()

	require.Panics(t, func() {
		CheckErr(errors.New("should-panic"))
	})
}

func TestUsageErrorf_FormatsMessageAndCommandPath(t *testing.T) {
	t.Parallel()

	rootCmd := &cobra.Command{
		Use: "root",
	}
	childCmd := &cobra.Command{
		Use: "child",
	}
	rootCmd.AddCommand(childCmd)

	err := UsageErrorf(childCmd, "invalid %s", "arg")
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid arg")
	require.Contains(t, err.Error(), "See 'root child -h' for help and examples")
}

func TestRequireNoArguments_AllowsNoArgs(t *testing.T) {
	t.Parallel()

	cmd := &cobra.Command{Use: "test"}
	require.NotPanics(t, func() {
		RequireNoArguments(cmd, nil)
	})
}

func TestRequireNoArguments_PanicsOnExtraArgs(t *testing.T) {
	t.Parallel()

	cmd := &cobra.Command{Use: "test"}

	require.Panics(t, func() {
		RequireNoArguments(cmd, []string{"extra"})
	})
}

func TestDefaultSubCommandRun_PrintsHelpOnNoArgs(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	cmd := &cobra.Command{
		Use:   "test",
		Short: "short description",
	}

	run := DefaultSubCommandRun(&buf)

	require.NotPanics(t, func() {
		run(cmd, nil)
	})

	out := buf.String()
	require.NotEmpty(t, out)
	require.Contains(t, out, "short description")
}
