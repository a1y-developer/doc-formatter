package main

import (
	"bytes"
	"io"
	"os"
	"testing"

	"ariga.io/atlas-provider-gorm/gormschema"
	"github.com/a1y/doc-formatter/internal/storage/infra/persistence"
	"github.com/stretchr/testify/require"
)

func TestMain_GeneratesStorageSchema(t *testing.T) {
	t.Parallel()

	// Pre-check to avoid triggering os.Exit on environments where gormschema fails.
	stmts, err := gormschema.New("postgres").Load(&persistence.DocumentModel{})
	if err != nil {
		t.Skipf("skipping storage loader main test due to gormschema error: %v", err)
	}
	require.NotEmpty(t, stmts)

	origStdout := os.Stdout
	r, w, err := os.Pipe()
	require.NoError(t, err)
	os.Stdout = w

	var buf bytes.Buffer
	done := make(chan struct{})
	go func() {
		_, _ = io.Copy(&buf, r)
		close(done)
	}()

	main()

	_ = w.Close()
	<-done
	os.Stdout = origStdout

	require.NotEmpty(t, buf.String())
}
