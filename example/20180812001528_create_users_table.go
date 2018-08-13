package main

import (
	"github.com/go-pg/pg/orm"
	"github.com/robinjoseph08/go-pg-migrations"
)

func init() {
	up := func(db orm.DB) error {
		_, err := db.Exec(`
			CREATE TABLE users (
				id SERIAL PRIMARY KEY,
				username TEXT NOT NULL UNIQUE,
				date_created TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
				date_modified TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
			)
		`)
		return err
	}

	down := func(db orm.DB) error {
		_, err := db.Exec("DROP TABLE users")
		return err
	}

	opts := migrations.MigrationOptions{}

	migrations.Register("20180812001528_create_users_table", up, down, opts)
}
