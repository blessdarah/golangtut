package main

import (
	"blessdarah/tuts/internal/db/persistence"

	"gorm.io/gen"
)

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath:           "internal/db/query",
		Mode:              gen.WithDefaultQuery | gen.WithQueryInterface,
		FieldNullable:     true,
		FieldWithIndexTag: true,
		FieldWithTypeTag:  true,
	})

	g.ApplyBasic(persistence.User{}, persistence.Event{})
	g.Execute()
}
