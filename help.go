package migrations

import "fmt"

const helpText = `Usage:
  go run %s/*.go [command]

Commands:
  create   - create a new migration in %s with the provided name
  migrate  - run any migrations that haven't been run yet
  rollback - roll back the previous run batch of migrations
  help     - print this help text

Examples:
  go run %s/*.go create create_users_table
  go run %s/*.go migrate
  go run %s/*.go rollback
  go run %s/*.go help
`

func help(directory string) {
	fmt.Printf(helpText, directory, directory, directory, directory, directory, directory)
}
