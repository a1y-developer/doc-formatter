package main

import (
	"io"
	"os"

	"ariga.io/atlas-provider-gorm/gormschema"
	"github.com/a1y/doc-formatter/pkg/storage/infra/persistence"
	"github.com/sirupsen/logrus"
)

func main() {
	stmts, err := gormschema.New("postgres").Load(&persistence.DocumentModel{})
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
