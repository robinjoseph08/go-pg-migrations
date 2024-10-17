package migrations

import (
	"os"
	"testing"

	"github.com/go-pg/pg/v10"
	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	tmp := os.TempDir()
	db := pg.Connect(&pg.Options{
		Addr:     "localhost:5432",
		User:     os.Getenv("TEST_DATABASE_USER"),
		Database: os.Getenv("TEST_DATABASE_NAME"),
	})

	err := Run(nil, tmp, []string{"cmd"})
	assert.Nil(t, err)

	err = Run(db, tmp, []string{"cmd", "migrate"})
	assert.Nil(t, err)

	err = Run(nil, tmp, []string{"cmd", "create"})
	assert.Equal(t, ErrCreateRequiresName, err)

	err = Run(nil, tmp, []string{"cmd", "create", "test_migration"})
	assert.Nil(t, err)

	err = Run(db, tmp, []string{"cmd", "rollback"})
	assert.Nil(t, err)
}

func TestRunWithOptions(t *testing.T) {
	tmp := os.TempDir()
	db := pg.Connect(&pg.Options{
		Addr:     "localhost:5432",
		User:     os.Getenv("TEST_DATABASE_USER"),
		Database: os.Getenv("TEST_DATABASE_NAME"),
	})
	db.AddQueryHook(logQueryHook{})

	t.Run("default", func(tt *testing.T) {
		dropMigrationTables(tt, db)

		err := RunWithOptions(db, tmp, []string{"cmd", "migrate"}, RunOptions{})
		assert.Nil(tt, err)
		assertTable(tt, db, "migrations", true)
		assertTable(tt, db, "migration_lock", true)
		assertTable(tt, db, "custom_migrations", false)
		assertTable(tt, db, "custom_migration_lock", false)
	})

	t.Run("custom tables - migrate", func(tt *testing.T) {
		dropMigrationTables(tt, db)

		err := RunWithOptions(db, tmp, []string{"cmd", "migrate"}, RunOptions{
			MigrationsTableName:    "custom_migrations",
			MigrationLockTableName: "custom_migration_lock",
		})
		assert.Nil(tt, err)
		assertTable(tt, db, "custom_migrations", true)
		assertTable(tt, db, "custom_migration_lock", true)
		assertTable(tt, db, "migrations", false)
		assertTable(tt, db, "migration_lock", false)
	})

	t.Run("custom tables - rollback", func(tt *testing.T) {
		dropMigrationTables(tt, db)

		err := RunWithOptions(db, tmp, []string{"cmd", "rollback"}, RunOptions{
			MigrationsTableName:    "custom_migrations",
			MigrationLockTableName: "custom_migration_lock",
		})
		assert.Nil(tt, err)
		assertTable(tt, db, "custom_migrations", true)
		assertTable(tt, db, "custom_migration_lock", true)
		assertTable(tt, db, "migrations", false)
		assertTable(tt, db, "migration_lock", false)
	})
}
