package migrations

import (
	"fmt"
	"sort"
	"time"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

var migrations []migration

// Register accepts a name, up, down, and options and adds the migration to the
// global migrations slice.
func Register(name string, up, down func(orm.DB) error, opts MigrationOptions) {
	migrations = append(migrations, migration{
		Name:               name,
		Up:                 up,
		Down:               down,
		DisableTransaction: opts.DisableTransaction,
	})
}

func migrate(db *pg.DB, directory string) error {
	// sort the registered migrations by name (which will sort by the
	// timestamp in their names)
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Name < migrations[j].Name
	})

	// look at the migrations table to see the already run migrations
	completed, err := getCompletedMigrations(db)
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
	err = acquireLock(db)
	if err != nil {
		return err
	}
	defer releaseLock(db)

	// find the last batch number
	batch, err := getLastBatchNumber(db)
	if err != nil {
		return err
	}
	batch = batch + 1

	fmt.Printf("Running batch %d with %d migration(s)...\n", batch, len(uncompleted))

	for _, m := range uncompleted {
		m.Batch = batch
		var err error
		if m.DisableTransaction {
			err = m.Up(db)
		} else {
			err = db.RunInTransaction(func(tx *pg.Tx) error {
				return m.Up(tx)
			})
		}
		if err != nil {
			return fmt.Errorf("%s: %s", m.Name, err)
		}

		m.CompletedAt = time.Now()
		err = db.Insert(&m)
		if err != nil {
			return fmt.Errorf("%s: %s", m.Name, err)
		}
		fmt.Printf("Finished running %q\n", m.Name)
	}

	return nil
}

func getCompletedMigrations(db orm.DB) ([]migration, error) {
	var completed []migration

	err := db.
		Model(&completed).
		Order("id").
		Select()
	if err != nil {
		return nil, err
	}

	return completed, nil
}

func filterMigrations(all, subset []migration, wantCompleted bool) []migration {
	subsetMap := map[string]bool{}

	for _, c := range subset {
		subsetMap[c.Name] = true
	}

	var d []migration

	for _, a := range all {
		if subsetMap[a.Name] == wantCompleted {
			d = append(d, a)
		}
	}

	return d
}

func acquireLock(db *pg.DB) error {
	l := lock{ID: lockID, IsLocked: true}

	result, err := db.Model(&l).
		Column("is_locked").
		WherePK().
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

func releaseLock(db orm.DB) error {
	l := lock{ID: lockID, IsLocked: false}
	return db.Update(&l)
}

func getLastBatchNumber(db orm.DB) (int32, error) {
	var res struct{ Batch int32 }
	err := db.Model(&migration{}).
		ColumnExpr("COALESCE(MAX(batch), 0) AS batch").
		Select(&res)
	if err != nil {
		return 0, err
	}
	return res.Batch, nil
}
