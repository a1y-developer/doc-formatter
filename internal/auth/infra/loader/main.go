package main

import (
	"io"
	"os"

	"ariga.io/atlas-provider-gorm/gormschema"
	"github.com/sirupsen/logrus"

	"github.com/a1y/doc-formatter/internal/auth/infra/persistence"
)

func main() {
	stmts, err := gormschema.New("postgres").Load(&persistence.UserModel{})
	if err != nil {
		logrus.Errorf("failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}
	_, err = io.WriteString(os.Stdout, stmts)
	if err != nil {
		os.Exit(1)
		return
	}
}
