package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/a1y/doc-formatter/cmd/storage/app"
	"github.com/sirupsen/logrus"
)

func main() {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	cmd := app.NewCmdStorage()

	if err := cmd.Execute(); err != nil {
		logrus.Fatal(err)
		os.Exit(1)
	}
	os.Exit(0)
}
