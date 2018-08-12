package migrations

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	err := Run([]string{"cmd"})
	assert.Nil(t, err)

	err = Run([]string{"cmd", "migrate"})
	assert.Nil(t, err)

	err = Run([]string{"cmd", "create"})
	assert.Nil(t, err)

	err = Run([]string{"cmd", "rollback"})
	assert.Nil(t, err)

	err = Run([]string{"cmd", "foo"})
	assert.Nil(t, err)
}
