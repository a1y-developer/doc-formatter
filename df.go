package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/a1y/doc-formatter/cmd"
)

// @title			AI Doc Formatter API
// @version		1.0
// @description	API for AI Doc Formatter
// @BasePath		/
func main() {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	command := cmd.NewDefaultDfctlCommand()

	if err := command.Execute(); err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
