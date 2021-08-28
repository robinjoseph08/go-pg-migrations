package migrations

import (
	"os"
	"testing"

	"github.com/go-pg/pg/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMigrationStatus(t *testing.T) {
	db := pg.Connect(&pg.Options{
		Addr:     "localhost:5432",
		User:     os.Getenv("TEST_DATABASE_USER"),
		Database: os.Getenv("TEST_DATABASE_NAME"),
	})

	err := ensureMigrationTables(db)
	require.Nil(t, err)

	defer clearMigrations(t, db)
	defer resetMigrations(t)

	t.Run("Returns nil if migrations are up to date", func(tt *testing.T) {
		clearMigrations(tt, db)
		resetMigrations(tt)

		migrations = []migration{
			{Name: "123", Up: noopMigration, Down: noopMigration},
		}

		_, err := db.Model(&migrations[0]).Insert()
		assert.Nil(tt, err)

		pendingErr := migrationStatus(db)
		assert.Nil(tt, pendingErr)
	})

	t.Run("Returns an error if migrations are not up to date", func(tt *testing.T) {
		clearMigrations(tt, db)
		resetMigrations(tt)

		migrations = []migration{
			{Name: "123", Up: noopMigration, Down: noopMigration},
		}

		pendingErr := migrationStatus(db)
		assert.EqualError(tt, pendingErr, "1 migration is pending")

		migrations = append(migrations, migration{
			Name: "123", Up: noopMigration, Down: noopMigration,
		})

		pendingErr = migrationStatus(db)
		assert.EqualError(tt, pendingErr, "2 migrations are pending")
	})
}
