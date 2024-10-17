package migrations

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"strings"
	"unicode/utf8"
)

func (m *migrator) status(w io.Writer) error {
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

	return writeStatusTable(w, completed, uncompleted)
}

func writeStatusTable(w io.Writer, completed []*migration, uncompleted []*migration) error {
	if len(completed)+len(uncompleted) == 0 {
		_, err := fmt.Fprintln(w, "No migrations found")
		return err
	}

	maxNameLength := 20
	for _, m := range completed {
		maxNameLength = maxInt(maxNameLength, utf8.RuneCountInString(m.Name))
	}
	for _, m := range uncompleted {
		maxNameLength = maxInt(maxNameLength, utf8.RuneCountInString(m.Name))
	}

	bf := bytes.NewBuffer(nil)

	// write header
	bf.WriteString("+---------+" + strings.Repeat("-", maxNameLength+2) + "+-------+\n")
	bf.WriteString("| Applied | Migration" + strings.Repeat(" ", maxNameLength-8) + "| Batch |\n")
	bf.WriteString("+---------+" + strings.Repeat("-", maxNameLength+2) + "+-------+\n")

	// write completed migrations
	for _, m := range completed {
		bf.WriteString("|    √    | " + m.Name + strings.Repeat(" ", maxNameLength-len(m.Name)) + " | " + fmt.Sprintf("%5d", m.Batch) + " |\n")
	}

	// write uncompleted migrations
	for _, m := range uncompleted {
		bf.WriteString("|         | " + m.Name + strings.Repeat(" ", maxNameLength-len(m.Name)) + " |       |\n")
	}

	// write footer
	bf.WriteString("+---------+" + strings.Repeat("-", maxNameLength+2) + "+-------+\n")

	_, err := bf.WriteTo(w)
	return err
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
