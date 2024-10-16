package migrations

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	tmp := os.TempDir()
	name := fmt.Sprintf("create_test_migration_%d", r.Int())

	err := create(tmp, name)
	assert.Nil(t, err)

	files, err := os.ReadDir(tmp)
	assert.Nil(t, err)

	found := false
	for _, f := range files {
		if strings.Contains(f.Name(), name) {
			found = true
		}
	}

	assert.Truef(t, found, "expected to find %q migration in %q", name, tmp)
}
