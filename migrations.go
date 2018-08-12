package migrations

import "fmt"

// Run takes in an argument slice and runs the appropriate command.
func Run(args []string) error {
	cmd := ""

	if len(args) > 1 {
		cmd = args[1]
	}

	switch cmd {
	case "help", "":
		fmt.Println("help")
	case "migrate":
		fmt.Println("migrate")
	case "create":
		fmt.Println("create")
	case "rollback":
		fmt.Println("rollback")
	}

	return nil
}
