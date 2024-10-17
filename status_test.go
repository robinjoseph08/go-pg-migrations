package migrations

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/go-pg/pg/v10"
	"github.com/stretchr/testify/require"
)

func TestStatus(t *testing.T) {
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

	completed := []migration{
		{Name: "2021_02_26_151503_dump", Up: noopMigration, Down: noopMigration, Batch: 1},
		{Name: "2021_02_26_151504_create_a_dump_table_for_test", Up: noopMigration, Down: noopMigration, Batch: 2},
	}
	uncompleted := []migration{
		{Name: "2021_02_26_151502_create_2nd_dump_table", Up: noopMigration, Down: noopMigration},
		{Name: "2021_02_26_151505_create_3rd_dump_table", Up: noopMigration, Down: noopMigration},
	}
	expected := strings.TrimSpace(`
+---------+------------------------------------------------+-------+
| Applied | Migration                                      | Batch |
+---------+------------------------------------------------+-------+
|    √    | 2021_02_26_151503_dump                         |     1 |
|    √    | 2021_02_26_151504_create_a_dump_table_for_test |     2 |
|         | 2021_02_26_151502_create_2nd_dump_table        |       |
|         | 2021_02_26_151505_create_3rd_dump_table        |       |
+---------+------------------------------------------------+-------+
`)

	t.Run("status_command", func(tt *testing.T) {
		clearMigrations(tt, db)
		resetMigrations(tt)

		migrations = completed[:1]
		err := migrate(db)
		require.Nil(tt, err, "migrate: %v", err)
		migrations = completed[:2]
		err = migrate(db)
		require.Nil(tt, err, "migrate: %v", err)

		migrations = append(migrations, uncompleted...)
		bf := bytes.NewBuffer(nil)
		err = status(db, bf)
		require.Nil(tt, err, "status: %v", err)

		got := strings.TrimSpace(bf.String())
		if got != expected {
			tt.Errorf("status table not match:\nEXPECTED:\n%s\nACTUAL:\n%s", expected, got)
		}
	})

	t.Run("write_status_table", func(tt *testing.T) {
		bf := bytes.NewBuffer(nil)
		err := writeStatusTable(bf, completed, uncompleted)
		require.Nil(tt, err, "write_status_table: %v", err)

		got := strings.TrimSpace(bf.String())
		if got != expected {
			tt.Errorf("status table not match:\nEXPECTED:\n%s\nACTUAL:\n%s", expected, got)
		}
	})

	t.Run("no_migrations_found", func(tt *testing.T) {
		expected := strings.TrimSpace(`No migrations found`)

		bf := bytes.NewBuffer(nil)
		err := writeStatusTable(bf, nil, nil)
		require.Nil(tt, err, "write_status_table: %v", err)

		got := strings.TrimSpace(bf.String())
		if got != expected {
			tt.Errorf("status table not match:\nEXPECTED:\n%s\nACTUAL:\n%s", expected, got)
		}
	})
}
