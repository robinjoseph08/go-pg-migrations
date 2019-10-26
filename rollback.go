package migrations

import (
	"fmt"
	"sort"

	"github.com/go-pg/pg/v9"
)

func rollback(db *pg.DB, directory string) error {
	// sort the registered migrations by name (which will sort by the
	// timestamp in their names)
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Name > migrations[j].Name
	})

	// look at the migrations table to see the already run migrations
	completed, err := getCompletedMigrations(db)
	if err != nil {
		return err
	}

	// acquire the migration lock from the migrations_lock table
	err = acquireLock(db)
	if err != nil {
		return err
	}
	defer releaseLock(db)

	batch, err := getLastBatchNumber(db)
	if err != nil {
		return err
	}
	// if no migrations have been run yet, exit early
	if batch == 0 {
		fmt.Println("No migrations have been run yet")
		return nil
	}

	rollback := getMigrationsForBatch(completed, batch)
	rollback = filterMigrations(migrations, rollback, true)

	fmt.Printf("Rolling back batch %d with %d migration(s)...\n", batch, len(rollback))

	for _, m := range rollback {
		var err error
		if m.DisableTransaction {
			err = m.Down(db)
		} else {
			err = db.RunInTransaction(func(tx *pg.Tx) error {
				return m.Down(tx)
			})
		}
		if err != nil {
			return fmt.Errorf("%s: %s", m.Name, err)
		}

		_, err = db.Model(&m).Where("name = ?", m.Name).Delete()
		if err != nil {
			return fmt.Errorf("%s: %s", m.Name, err)
		}
		fmt.Printf("Finished rolling back %q\n", m.Name)
	}

	return nil
}

func getMigrationsForBatch(migrations []migration, batch int32) []migration {
	var m []migration

	for _, migration := range migrations {
		if migration.Batch == batch {
			m = append(m, migration)
		}
	}

	return m
}
