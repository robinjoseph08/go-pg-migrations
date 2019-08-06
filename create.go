package migrations

import (
	"fmt"
	"io/ioutil"
	"path"
	"time"
)

const timeFormat = "20060102150405"

var template = `package main

import (
	"github.com/go-pg/pg/orm"
	"github.com/robinjoseph08/go-pg-migrations"
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
	filename := path.Join(directory, fmt.Sprintf("%s.go", fullname))

	fmt.Printf("Creating %s...\n", filename)

	return ioutil.WriteFile(filename, []byte(fmt.Sprintf(template, fullname)), 0644)
}
