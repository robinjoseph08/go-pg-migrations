package migrations

import (
	"os"
	"testing"
	"time"

	"github.com/go-pg/pg/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRollback(t *testing.T) {
	tmp := os.TempDir()
	db := pg.Connect(&pg.Options{
		Addr:     "localhost:5432",
		User:     os.Getenv("TEST_DATABASE_USER"),
		Database: os.Getenv("TEST_DATABASE_NAME"),
	})

	err := ensureMigrationTables(db)
	require.Nil(t, err)

	defer clearMigrations(t, db)
	defer resetMigrations(t)

	t.Run("sorts migrations in reverse order", func(tt *testing.T) {
		clearMigrations(tt, db)
		resetMigrations(tt)
		migrations = []migration{
			{Name: "123", Up: noopMigration, Down: noopMigration},
			{Name: "456", Up: noopMigration, Down: noopMigration},
		}

		err := rollback(db, tmp)
		assert.Nil(tt, err)

		assert.Equal(tt, "456", migrations[0].Name)
		assert.Equal(tt, "123", migrations[1].Name)
	})

	t.Run("returns an error if the migration lock is already held", func(tt *testing.T) {
		clearMigrations(tt, db)
		resetMigrations(tt)
		migrations = []migration{
			{Name: "123", Up: noopMigration, Down: noopMigration},
			{Name: "456", Up: noopMigration, Down: noopMigration},
		}

		err := acquireLock(db)
		assert.Nil(tt, err)
		defer releaseLock(db)

		err = rollback(db, tmp)
		assert.Equal(tt, ErrAlreadyLocked, err)
	})

	t.Run("exits early if there aren't any migrations to rollback", func(tt *testing.T) {
		clearMigrations(tt, db)
		resetMigrations(tt)
		migrations = []migration{
			{Name: "123", Up: noopMigration, Down: noopMigration},
			{Name: "456", Up: noopMigration, Down: noopMigration},
		}

		err := rollback(db, tmp)
		assert.Nil(tt, err)

		count, err := db.Model(&migration{}).Count()
		assert.Nil(tt, err)
		assert.Equal(tt, 0, count)
	})

	t.Run("only rolls back the last batch", func(tt *testing.T) {
		clearMigrations(tt, db)
		resetMigrations(tt)
		migrations = []migration{
			{Name: "123", Up: noopMigration, Down: noopMigration, Batch: 4, CompletedAt: time.Now()},
			{Name: "456", Up: noopMigration, Down: noopMigration, Batch: 5, CompletedAt: time.Now()},
			{Name: "789", Up: noopMigration, Down: noopMigration, Batch: 5, CompletedAt: time.Now()},
			{Name: "010", Up: noopMigration, Down: noopMigration},
		}

		m := migrations[:2]
		err := db.Insert(&m)
		assert.Nil(tt, err)

		err = rollback(db, tmp)
		assert.Nil(tt, err)

		batch, err := getLastBatchNumber(db)
		assert.Nil(tt, err)
		assert.Equal(tt, batch, int32(4))

		count, err := db.Model(&migration{}).Count()
		assert.Nil(tt, err)
		assert.Equal(tt, 1, count)
	})

	t.Run(`runs "down" within a transaction if specified`, func(tt *testing.T) {
		clearMigrations(tt, db)
		resetMigrations(tt)
		migrations = []migration{
			{Name: "123", Up: noopMigration, Down: erringMigration, DisableTransaction: false, Batch: 1, CompletedAt: time.Now()},
		}

		err := db.Insert(&migrations)
		assert.Nil(tt, err)

		err = rollback(db, tmp)
		assert.EqualError(tt, err, "123: error")

		assertTable(tt, db, "test_table", false)
	})

	t.Run(`doesn't run "down" within a transaction if specified`, func(tt *testing.T) {
		clearMigrations(tt, db)
		resetMigrations(tt)
		migrations = []migration{
			{Name: "123", Up: noopMigration, Down: erringMigration, DisableTransaction: true, Batch: 1, CompletedAt: time.Now()},
		}

		err := db.Insert(&migrations)
		assert.Nil(tt, err)

		err = rollback(db, tmp)
		assert.EqualError(tt, err, "123: error")

		assertTable(tt, db, "test_table", true)
	})
}
