package main

import (
	"math/rand"
	"os"
	"time"

	app2 "github.com/a1y/doc-formatter/cmd/gateway/app"
	"github.com/sirupsen/logrus"
)

// @title AI Doc Formatter API Gateway
// @version 1.0
// @description API Gateway for AI Doc Formatter
// @BasePath /
func main() {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	cmd := app2.NewCmdGateway()

	if err := cmd.Execute(); err != nil {
		logrus.Fatal(err)
		os.Exit(1)
	}
	os.Exit(0)
}
