package migrations

import (
	"os"
	"testing"

	"github.com/go-pg/pg/v9"
	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	tmp := os.TempDir()
	db := pg.Connect(&pg.Options{
		Addr:     "localhost:5432",
		User:     os.Getenv("TEST_DATABASE_USER"),
		Database: os.Getenv("TEST_DATABASE_NAME"),
	})

	err := Run(nil, tmp, []string{"cmd"})
	assert.Nil(t, err)

	err = Run(db, tmp, []string{"cmd", "migrate"})
	assert.Nil(t, err)

	err = Run(nil, tmp, []string{"cmd", "create"})
	assert.Equal(t, ErrCreateRequiresName, err)

	err = Run(nil, tmp, []string{"cmd", "create", "test_migration"})
	assert.Nil(t, err)

	err = Run(db, tmp, []string{"cmd", "rollback"})
	assert.Nil(t, err)
}
