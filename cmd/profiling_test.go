package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/require"
)

func TestAddProfilingFlags_SetsDefaults(t *testing.T) {
	t.Parallel()

	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)

	addProfilingFlags(flags)

	profileFlag := flags.Lookup("profile")
	require.NotNil(t, profileFlag)
	require.Equal(t, "none", profileFlag.DefValue)

	profileOutputFlag := flags.Lookup("profile-output")
	require.NotNil(t, profileOutputFlag)
	require.Equal(t, "profile.pprof", profileOutputFlag.DefValue)
}

func TestInitProfiling_NoneProfile(t *testing.T) {
	t.Parallel()

	profileName = "none"
	profileOutput = "ignored.pprof"

	err := initProfiling()
	require.NoError(t, err)

	err = flushProfiling()
	require.NoError(t, err)
}

func TestInitProfiling_InvalidProfileName(t *testing.T) {
	t.Parallel()

	profileName = "invalid-profile-name"
	profileOutput = "ignored.pprof"

	err := initProfiling()
	require.Error(t, err)
}

func TestCPUProfiling_Lifecycle(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	outFile := filepath.Join(dir, "cpu.pprof")

	profileName = "cpu"
	profileOutput = outFile

	err := initProfiling()
	require.NoError(t, err)

	// Stopping CPU profiling should succeed even if no significant work was done.
	err = flushProfiling()
	require.NoError(t, err)
}

func TestHeapProfiling_WritesProfile(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	outFile := filepath.Join(dir, "heap.pprof")

	profileName = "heap"
	profileOutput = outFile

	// initProfiling should set up signal handling and validate the profile
	err := initProfiling()
	require.NoError(t, err)

	err = flushProfiling()
	require.NoError(t, err)

	info, err := os.Stat(outFile)
	require.NoError(t, err)
	require.Greater(t, info.Size(), int64(0))
}
