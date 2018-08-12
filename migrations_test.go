package migrations

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	tmp := os.TempDir()

	err := Run(tmp, []string{"cmd"})
	assert.Nil(t, err)

	err = Run(tmp, []string{"cmd", "migrate"})
	assert.Nil(t, err)

	err = Run(tmp, []string{"cmd", "create"})
	assert.Equal(t, ErrCreateRequiresName, err)

	err = Run(tmp, []string{"cmd", "create", "test_migration"})
	assert.Nil(t, err)

	err = Run(tmp, []string{"cmd", "rollback"})
	assert.Nil(t, err)
}
