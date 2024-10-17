package migrations

import (
	"fmt"
	"sort"

	"github.com/go-pg/pg/v10"
)

func (m *migrator) rollback() error {
	// sort the registered migrations by name (which will sort by the
	// timestamp in their names)
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Name > migrations[j].Name
	})

	// look at the migrations table to see the already run migrations
	completed, err := m.getCompletedMigrations()
	if err != nil {
		return err
	}

	// acquire the migration lock from the migrations_lock table
	err = m.acquireLock()
	if err != nil {
		return err
	}
	defer m.releaseLock()

	batch, err := m.getLastBatchNumber()
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

	for _, mig := range rollback {
		var err error
		if mig.DisableTransaction {
			err = mig.Down(m.db)
		} else {
			err = m.db.RunInTransaction(m.db.Context(), func(tx *pg.Tx) error {
				return mig.Down(tx)
			})
		}
		if err != nil {
			return fmt.Errorf("%s: %s", mig.Name, err)
		}

		_, err = m.db.
			Exec(fmt.Sprintf("DELETE FROM %q WHERE name = ?", escapeTableName(m.opts.MigrationsTableName)), mig.Name)
		if err != nil {
			return fmt.Errorf("%s: %s", mig.Name, err)
		}
		fmt.Printf("Finished rolling back %q\n", mig.Name)
	}

	return nil
}

func getMigrationsForBatch(migrations []*migration, batch int32) []*migration {
	var m []*migration
	for _, mig := range migrations {
		if mig.Batch == batch {
			m = append(m, mig)
		}
	}

	return m
}
