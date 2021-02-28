// Package migrations provides a robust mechanism for registering, creating, and
// running migrations using go-pg-pg.
package migrations

import (
	"errors"
	"os"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

// Errors that can be returned from Run.
var (
	ErrAlreadyLocked      = errors.New("migration table is already locked")
	ErrCreateRequiresName = errors.New("migration name is required for create")
)

// MigrationOptions allows settings to be configured on a per-migration basis.
type MigrationOptions struct {
	DisableTransaction bool
}

type migration struct {
	tableName struct{} `pg:"migrations,alias:migrations"`

	ID          int32
	Name        string
	Batch       int32
	CompletedAt time.Time
	Up          func(orm.DB) error `pg:"-"`
	Down        func(orm.DB) error `pg:"-"`

	DisableTransaction bool `pg:"-"`
}

type lock struct {
	tableName struct{} `pg:"migration_lock,alias:migration_lock"`

	ID       string
	IsLocked bool `pg:",use_zero,notnull"`
}

const lockID = "lock"

// Run takes in a directory and an argument slice and runs the appropriate command.
func Run(db *pg.DB, directory string, args []string) error {
	cmd := ""

	if len(args) > 1 {
		cmd = args[1]
	}

	switch cmd {
	case "migrate":
		err := ensureMigrationTables(db)
		if err != nil {
			return err
		}

		return migrate(db)
	case "create":
		if len(args) < 3 {
			return ErrCreateRequiresName
		}
		name := args[2]
		return create(directory, name)
	case "rollback":
		err := ensureMigrationTables(db)
		if err != nil {
			return err
		}

		return rollback(db)
	case "status":
		err := ensureMigrationTables(db)
		if err != nil {
			return err
		}

		return status(db, os.Stdout)
	default:
		help(directory)
		return nil
	}
}
