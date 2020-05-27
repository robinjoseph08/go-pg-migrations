package migrations

import (
	"os"
	"testing"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/stretchr/testify/assert"
)

func TestEnsureMigrationTables(t *testing.T) {
	db := pg.Connect(&pg.Options{
		Addr:     "localhost:5432",
		User:     os.Getenv("TEST_DATABASE_USER"),
		Database: os.Getenv("TEST_DATABASE_NAME"),
	})

	// drop tables to start from a clean database
	dropMigrationTables(t, db)

	err := ensureMigrationTables(db)
	assert.Nil(t, err)

	tables := []string{"migrations", "migration_lock"}

	for _, table := range tables {
		assertTable(t, db, table, true)
	}

	assertOneLock(t, db)

	// with existing tables, ensureMigrationTables should do anything
	err = ensureMigrationTables(db)
	assert.Nil(t, err)

	for _, table := range tables {
		assertTable(t, db, table, true)
	}

	assertOneLock(t, db)
}

func dropMigrationTables(t *testing.T, db *pg.DB) {
	t.Helper()

	_, err := db.Exec("DROP TABLE migrations")
	assert.Nil(t, err)
	_, err = db.Exec("DROP TABLE migration_lock")
	assert.Nil(t, err)
}

func assertTable(t *testing.T, db *pg.DB, table string, exists bool) {
	t.Helper()

	want := 0
	msg := "expected %q table to not exist"
	if exists {
		want = 1
		msg = "expected %q table to exist"
	}

	count, err := orm.NewQuery(db).
		Table("information_schema.tables").
		Where("table_name = ?", table).
		Where("table_schema = current_schema").
		Count()
	assert.Nil(t, err)
	assert.Equalf(t, want, count, msg, table)
}

func assertOneLock(t *testing.T, db *pg.DB) {
	t.Helper()

	count, err := orm.NewQuery(db).
		Table("migration_lock").
		Count()
	assert.Nil(t, err)
	assert.Equal(t, 1, count, "expected migraions_lock to have a row")
}
