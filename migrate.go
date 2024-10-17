package migrations

import (
	"fmt"
	"sort"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

var migrations []*migration

// Register accepts a name, up, down, and options and adds the migration to the
// global migrations slice.
func Register(name string, up, down func(orm.DB) error, opts MigrationOptions) {
	migrations = append(migrations, &migration{
		Name:               name,
		Up:                 up,
		Down:               down,
		DisableTransaction: opts.DisableTransaction,
	})
}

func (m *migrator) migrate() (err error) {
	// sort the registered migrations by name (which will sort by the
	// timestamp in their names)
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Name < migrations[j].Name
	})

	// look at the migrations table to see the already run migrations
	completed, err := m.getCompletedMigrations()
	if err != nil {
		return err
	}

	// diff the completed migrations from the registered migrations to find
	// the migrations we still need to run
	uncompleted := filterMigrations(migrations, completed, false)

	// if there are no migrations that need to be run, exit early
	if len(uncompleted) == 0 {
		fmt.Println("Migrations already up to date")
		return nil
	}

	// acquire the migration lock from the migrations_lock table
	err = m.acquireLock()
	if err != nil {
		return err
	}
	defer func() {
		e := m.releaseLock()
		if e != nil && err == nil {
			err = e
		}
	}()

	// find the last batch number
	batch, err := m.getLastBatchNumber()
	if err != nil {
		return err
	}
	batch++

	fmt.Printf("Running batch %d with %d migration(s)...\n", batch, len(uncompleted))

	for _, mig := range uncompleted {
		var err error
		if mig.DisableTransaction {
			err = mig.Up(m.db)
		} else {
			err = m.db.RunInTransaction(m.db.Context(), func(tx *pg.Tx) error {
				return mig.Up(tx)
			})
		}
		if err != nil {
			return fmt.Errorf("%s: %s", mig.Name, err)
		}

		migrationMap := map[string]interface{}{
			"name":         mig.Name,
			"batch":        batch,
			"completed_at": time.Now(),
		}
		_, err = m.db.
			Model(&migrationMap).
			Table(m.opts.MigrationsTableName).
			Insert()
		if err != nil {
			return fmt.Errorf("%s: %s", mig.Name, err)
		}
		fmt.Printf("Finished running %q\n", mig.Name)
	}

	return nil
}

func (m *migrator) getCompletedMigrations() ([]*migration, error) {
	var completed []*migration

	err := orm.NewQuery(m.db).
		Table(m.opts.MigrationsTableName).
		Order("id").
		Select(&completed)
	if err != nil {
		return nil, err
	}

	return completed, nil
}

func (m *migrator) acquireLock() error {
	l := map[string]interface{}{"is_locked": true}
	result, err := m.db.
		Model(&l).
		Table(m.opts.MigrationLockTableName).
		Column("is_locked").
		Where("id = ?", lockID).
		Where("is_locked = ?", false).
		Update()
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrAlreadyLocked
	}

	return nil
}

func (m *migrator) releaseLock() error {
	l := map[string]interface{}{"is_locked": true}
	_, err := m.db.
		Model(&l).
		Table(m.opts.MigrationLockTableName).
		Column("is_locked").
		Where("id = ?", lockID).
		Update()
	return err
}

func (m *migrator) getLastBatchNumber() (int32, error) {
	var res struct{ Batch int32 }
	err := orm.NewQuery(m.db).
		Table(m.opts.MigrationsTableName).
		ColumnExpr("COALESCE(MAX(batch), 0) AS batch").
		Select(&res)
	if err != nil {
		return 0, err
	}
	return res.Batch, nil
}

func filterMigrations(all, subset []*migration, wantCompleted bool) []*migration {
	subsetMap := map[string]bool{}
	for _, c := range subset {
		subsetMap[c.Name] = true
	}

	var d []*migration
	for _, a := range all {
		if subsetMap[a.Name] == wantCompleted {
			d = append(d, a)
		}
	}

	return d
}
