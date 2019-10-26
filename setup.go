package migrations

import (
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
)

func ensureMigrationTables(db *pg.DB) error {
	exists, err := checkIfTableExists("migrations", db)
	if err != nil {
		return err
	}
	if !exists {
		err = createTable(&migration{}, db)
		if err != nil {
			return err
		}
	}

	exists, err = checkIfTableExists("migration_lock", db)
	if err != nil {
		return err
	}
	if !exists {
		err = createTable(&lock{}, db)
		if err != nil {
			return err
		}
	}

	count, err := db.Model(&lock{}).Count()
	if err != nil {
		return err
	}
	if count == 0 {
		l := lock{ID: lockID, IsLocked: false}
		err = db.Insert(&l)
		if err != nil {
			return err
		}
	}

	return nil
}

func checkIfTableExists(name string, db orm.DB) (bool, error) {
	count, err := orm.NewQuery(db).
		Table("information_schema.tables").
		Where("table_name = ?", name).
		Where("table_schema = current_schema").
		Count()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func createTable(model interface{}, db *pg.DB) error {
	opts := orm.CreateTableOptions{IfNotExists: true}
	return db.CreateTable(model, &opts)
}
