package migrations

import (
	"fmt"

	"github.com/go-pg/pg/v10/orm"
)

func (m *migrator) ensureMigrationTables() error {
	exists, err := m.checkIfTableExists(m.opts.MigrationsTableName)
	if err != nil {
		return err
	}
	if !exists {
		_, err = m.db.Exec(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %q ("id" SERIAL PRIMARY KEY, "name" TEXT NOT NULL, "batch" INTEGER NOT NULL, "completed_at" TIMESTAMPTZ NOT NULL)`, escapeTableName(m.opts.MigrationsTableName)))
		if err != nil {
			return err
		}
	}

	exists, err = m.checkIfTableExists(m.opts.MigrationLockTableName)
	if err != nil {
		return err
	}
	if !exists {
		_, err = m.db.Exec(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %q ("id" TEXT PRIMARY KEY, "is_locked" BOOLEAN NOT NULL)`, escapeTableName(m.opts.MigrationLockTableName)))
		if err != nil {
			return err
		}
	}

	count, err := orm.NewQuery(m.db).
		Table(m.opts.MigrationLockTableName).
		Count()
	if err != nil {
		return err
	}
	if count == 0 {
		l := map[string]interface{}{
			"id":        lockID,
			"is_locked": false,
		}
		_, err = m.db.
			Model(&l).
			Table(m.opts.MigrationLockTableName).
			Insert()
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *migrator) checkIfTableExists(name string) (bool, error) {
	count, err := orm.NewQuery(m.db).
		Table("information_schema.tables").
		Where("table_name = ?", name).
		Where("table_schema = current_schema").
		Count()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
