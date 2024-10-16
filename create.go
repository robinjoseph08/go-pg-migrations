package migrations

import (
	"fmt"
	"os"
	"path"
	"time"
)

const timeFormat = "20060102150405"

var template = `package main

import (
	"github.com/go-pg/pg/v10/orm"
	migrations "github.com/robinjoseph08/go-pg-migrations/v3"
)

func init() {
	up := func(db orm.DB) error {
		_, err := db.Exec("")
		return err
	}

	down := func(db orm.DB) error {
		_, err := db.Exec("")
		return err
	}

	opts := migrations.MigrationOptions{}

	migrations.Register("%s", up, down, opts)
}
`

func create(directory, name string) error {
	version := time.Now().UTC().Format(timeFormat)
	fullname := fmt.Sprintf("%s_%s", version, name)
	filename := path.Join(directory, fullname+".go")

	fmt.Printf("Creating %s...\n", filename)

	return os.WriteFile(filename, []byte(fmt.Sprintf(template, fullname)), 0644)
}
