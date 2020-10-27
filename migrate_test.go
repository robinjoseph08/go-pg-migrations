package migrations

import (
	"context"
	"errors"
	"log"
	"os"
	"testing"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegister(t *testing.T) {
	resetMigrations(t)

	name := "foo"
	trxOpts := MigrationOptions{DisableTransaction: false}
	noTrxOpts := MigrationOptions{DisableTransaction: true}

	cases := []struct {
		name     string
		up, down func(orm.DB) error
		opts     MigrationOptions
	}{
		{name, noopMigration, noopMigration, trxOpts},
		{name, noopMigration, noopMigration, noTrxOpts},
	}

	for i, tt := range cases {
		Register(tt.name, tt.up, tt.down, tt.opts)

		require.Len(t, migrations, i+1)
		m := migrations[i]
		assert.Equal(t, tt.name, m.Name)
		assert.NotNil(t, m.Up)
		assert.NotNil(t, m.Down)
		assert.Equal(t, tt.opts.DisableTransaction, m.DisableTransaction)
	}
}

func TestMigrate(t *testing.T) {
	db := pg.Connect(&pg.Options{
		Addr:     "localhost:5432",
		User:     os.Getenv("TEST_DATABASE_USER"),
		Database: os.Getenv("TEST_DATABASE_NAME"),
	})

	db.AddQueryHook(logQueryHook{})

	err := ensureMigrationTables(db)
	require.Nil(t, err)

	defer clearMigrations(t, db)
	defer resetMigrations(t)

	t.Run("sorts migrations", func(tt *testing.T) {
		clearMigrations(tt, db)
		resetMigrations(tt)
		migrations = []migration{
			{Name: "456", Up: noopMigration, Down: noopMigration},
			{Name: "123", Up: noopMigration, Down: noopMigration},
		}

		err := migrate(db)
		assert.Nil(tt, err)

		assert.Equal(tt, "123", migrations[0].Name)
		assert.Equal(tt, "456", migrations[1].Name)
	})

	t.Run("only runs uncompleted migrations", func(tt *testing.T) {
		clearMigrations(tt, db)
		resetMigrations(tt)
		migrations = []migration{
			{Name: "123", Up: noopMigration, Down: noopMigration, Batch: 1, CompletedAt: time.Now()},
			{Name: "456", Up: noopMigration, Down: noopMigration},
		}

		_, err := db.Model(&migrations[0]).Insert()
		assert.Nil(tt, err)

		err = migrate(db)
		assert.Nil(tt, err)

		var m []migration
		err = db.Model(&m).Order("name").Select()
		assert.Nil(tt, err)
		require.Len(tt, m, 2)
		assert.Equal(tt, m[0].Batch, int32(1))
		assert.Equal(tt, m[1].Batch, int32(2))
	})

	t.Run("exits early if there aren't any migrations to run", func(tt *testing.T) {
		clearMigrations(tt, db)
		resetMigrations(tt)
		migrations = []migration{
			{Name: "123", Up: noopMigration, Down: noopMigration, Batch: 1, CompletedAt: time.Now()},
			{Name: "456", Up: noopMigration, Down: noopMigration, Batch: 1, CompletedAt: time.Now()},
		}

		_, err := db.Model(&migrations).Insert()
		assert.Nil(tt, err)

		err = migrate(db)
		assert.Nil(tt, err)

		count, err := db.Model(&migration{}).Where("batch = 2").Count()
		assert.Nil(tt, err)
		assert.Equal(tt, 0, count)
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

		err = migrate(db)
		assert.Equal(tt, ErrAlreadyLocked, err)
	})

	t.Run("increments batch number for each run and associates all migrations with it", func(tt *testing.T) {
		clearMigrations(tt, db)
		resetMigrations(tt)
		migrations = []migration{
			{Name: "123", Up: noopMigration, Down: noopMigration, Batch: 5, CompletedAt: time.Now()},
			{Name: "456", Up: noopMigration, Down: noopMigration},
			{Name: "789", Up: noopMigration, Down: noopMigration},
		}

		_, err := db.Model(&migrations[0]).Insert()
		assert.Nil(tt, err)

		err = migrate(db)
		assert.Nil(tt, err)

		batch, err := getLastBatchNumber(db)
		assert.Nil(tt, err)
		assert.Equal(tt, batch, int32(6))

		count, err := db.Model(&migration{}).Where("batch = ?", batch).Count()
		assert.Nil(tt, err)
		assert.Equal(tt, 2, count)
	})

	t.Run(`runs "up" within a transaction if specified`, func(tt *testing.T) {
		clearMigrations(tt, db)
		resetMigrations(tt)
		migrations = []migration{
			{Name: "123", Up: erringMigration, Down: noopMigration, DisableTransaction: false},
		}

		err := migrate(db)
		assert.EqualError(tt, err, "123: error")

		assertTable(tt, db, "test_table", false)
	})

	t.Run(`doesn't run "up" within a transaction if specified`, func(tt *testing.T) {
		clearMigrations(tt, db)
		resetMigrations(tt)
		migrations = []migration{
			{Name: "123", Up: erringMigration, Down: noopMigration, DisableTransaction: true},
		}

		err := migrate(db)
		assert.EqualError(tt, err, "123: error")

		assertTable(tt, db, "test_table", true)
	})
}

type logQueryHook struct{}

func (logQueryHook) BeforeQuery(ctx context.Context, event *pg.QueryEvent) (context.Context, error) {
	return ctx, nil
}

func (qh logQueryHook) AfterQuery(ctx context.Context, event *pg.QueryEvent) error {
	query, err := event.FormattedQuery()
	if err != nil {
		return err
	}

	log.Println(string(query))

	return nil
}

func resetMigrations(t *testing.T) {
	t.Helper()
	migrations = []migration{}
}

func clearMigrations(t *testing.T, db *pg.DB) {
	t.Helper()

	_, err := db.Exec("DELETE FROM migrations")
	assert.Nil(t, err)
	_, err = db.Exec("DROP TABLE IF EXISTS test_table")
	assert.Nil(t, err)
}

func noopMigration(db orm.DB) error {
	return nil
}

func erringMigration(db orm.DB) error {
	_, err := db.Exec("CREATE TABLE test_table (id integer)")
	if err != nil {
		return err
	}
	return errors.New("error")
}
