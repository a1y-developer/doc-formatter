package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/a1y/doc-formatter/pkg/version"
)

func TestMakeUpdateVersionGoFile_ContainsExpectedFields(t *testing.T) {
	v := &version.Info{
		ReleaseVersion: "v1.2.3",
		GitInfo: &version.GitInfo{
			LatestTag: "v1.2.3",
			Commit:    "abcdef1234567890",
			TreeState: "clean",
		},
		BuildInfo: &version.BuildInfo{
			GoVersion: "go1.21.0",
			GOOS:      "linux",
			GOARCH:    "amd64",
			NumCPU:    8,
			Compiler:  "gc",
			BuildTime: "2026-01-01 00:00:00",
		},
	}

	got := makeUpdateVersionGoFile(v)

	assert.True(t, strings.Contains(got, "package version"))
	assert.True(t, strings.Contains(got, `ReleaseVersion: "v1.2.3"`))
	assert.True(t, strings.Contains(got, `LatestTag:   "v1.2.3"`))
	assert.True(t, strings.Contains(got, `Commit:      "abcdef1234567890"`))
	assert.True(t, strings.Contains(got, `TreeState:   "clean"`))
	assert.True(t, strings.Contains(got, `GoVersion: "go1.21.0"`))
	assert.True(t, strings.Contains(got, `GOOS:      "linux"`))
	assert.True(t, strings.Contains(got, `GOARCH:    "amd64"`))
}
