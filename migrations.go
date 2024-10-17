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

// RunOptions allows settings to be configured for the environment that the migrations are run.
type RunOptions struct {
	// Set this to configure the table name of the migrations table. The default is `migrations`. Changing this after
	// you already have a migrations table does NOT rename it; it assumes you're starting fresh.
	MigrationsTableName string
	// Set this to configure the table name of the lock table. The default is `migration_lock`. Changing this after you
	// already have a migration lock table does NOT rename it; it just creates a new one and leaves any existing ones
	// alone.
	MigrationLockTableName string
}

// migration doesn't map to the table that we create to keep track of migrations. To see details of that table, see
// setup.go. This struct has tableName and pg tags because it makes tests easier, but it's not used in non-test code.
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

const lockID = "lock"

// Run takes in a directory and an argument slice and runs the appropriate command with default options.
func Run(db *pg.DB, directory string, args []string) error {
	return RunWithOptions(db, directory, args, RunOptions{})
}

// RunWithOptions takes in a directory, an argument slice, and run options and runs the appropriate command.
func RunWithOptions(db *pg.DB, directory string, args []string, opts RunOptions) error {
	cmd := ""
	if len(args) > 1 {
		cmd = args[1]
	}

	m := newMigrator(db, opts)

	switch cmd {
	case "migrate":
		err := m.ensureMigrationTables()
		if err != nil {
			return err
		}

		return m.migrate()
	case "create":
		if len(args) < 3 {
			return ErrCreateRequiresName
		}
		name := args[2]
		return m.create(directory, name)
	case "rollback":
		err := m.ensureMigrationTables()
		if err != nil {
			return err
		}
		return m.rollback()
	case "status":
		err := m.ensureMigrationTables()
		if err != nil {
			return err
		}

		return status(db, os.Stdout)
	default:
		help(directory)
		return nil
	}
}
