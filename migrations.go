package migrations

import (
	"errors"
	"fmt"
)

// Errors that can be returned from Run.
var (
	ErrCreateRequiresName = errors.New("migration name is required for create")
)

// Run takes in a directory and an argument slice and runs the appropriate command.
func Run(directory string, args []string) error {
	cmd := ""

	if len(args) > 1 {
		cmd = args[1]
	}

	switch cmd {
	case "migrate":
		fmt.Println("migrate")
		return nil
	case "create":
		if len(args) < 3 {
			return ErrCreateRequiresName
		}
		name := args[2]
		return create(directory, name)
	case "rollback":
		fmt.Println("rollback")
		return nil
	default:
		fmt.Println("help")
		return nil
	}
}
