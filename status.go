package migrations

import (
	"fmt"

	"github.com/go-pg/pg/v10"
)

func migrationStatus(db *pg.DB) error {
	uncompleted, err := getUncompletedMigrations(db)
	if err != nil {
		return nil
	}

	if len(uncompleted) == 0 {
		fmt.Println("Migrations already up to date")
		return nil
	}

	return ErrPendingMigrations{len(uncompleted)}
}

// ErrPendingMigrations is returned by the 'status' command when there is at
// least one pending migration
type ErrPendingMigrations struct {
	N int
}

func (p ErrPendingMigrations) Error() string {
	if p.N == 1 {
		return "1 migration is pending"
	}
	return fmt.Sprintf("%d migrations are pending", p.N)
}
