// Package migrations provides a robust mechanism for registering, creating, and
// running migrations using go-pg-pg.
package migrations

import (
	"errors"
	"time"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
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
	tableName struct{} `sql:"migrations,alias:migrations"`

	ID          int32
	Name        string
	Batch       int32
	CompletedAt time.Time
	Up          func(orm.DB) error `sql:"-"`
	Down        func(orm.DB) error `sql:"-"`

	DisableTransaction bool `sql:"-"`
}

type lock struct {
	tableName struct{} `sql:"migration_lock,alias:migration_lock"`

	ID       string
	IsLocked bool `sql:",notnull"`
}

const lockID = "lock"

// Run takes in a directory and an argument slice and runs the appropriate command.
func Run(db *pg.DB, directory string, args []string) error {
	err := ensureMigrationTables(db)
	if err != nil {
		return err
	}

	cmd := ""

	if len(args) > 1 {
		cmd = args[1]
	}

	switch cmd {
	case "migrate":
		return migrate(db, directory)
	case "create":
		if len(args) < 3 {
			return ErrCreateRequiresName
		}
		name := args[2]
		return create(directory, name)
	case "rollback":
		return rollback(db, directory)
	default:
		help(directory)
		return nil
	}
}
