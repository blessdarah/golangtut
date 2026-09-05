package main

import (
	"blessdarah/tuts/internal/db/persistence"
	"fmt"
	"io"
	"os"

	"ariga.io/atlas-provider-gorm/gormschema"
)

func main() {
	stmts, err := gormschema.New("postgres").Load(&persistence.User{}, &persistence.Event{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}

	_, _ = io.WriteString(os.Stdout, stmts)
}
