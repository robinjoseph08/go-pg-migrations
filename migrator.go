package migrations

import (
	"strings"

	"github.com/go-pg/pg/v10"
)

type migrator struct {
	db   *pg.DB
	opts RunOptions
}

func newMigrator(db *pg.DB, opts RunOptions) *migrator {
	if opts.MigrationsTableName == "" {
		opts.MigrationsTableName = "migrations"
	}
	if opts.MigrationLockTableName == "" {
		opts.MigrationLockTableName = "migration_lock"
	}

	return &migrator{
		db:   db,
		opts: opts,
	}
}

func escapeTableName(name string) string {
	return strings.ReplaceAll(name, `"`, `""`)
}
