package main

import (
	"log"
	"os"

	"github.com/go-pg/pg/v10"
	migrations "github.com/robinjoseph08/go-pg-migrations/v3"
)

const directory = "example"

func main() {
	db := pg.Connect(&pg.Options{
		Addr:     "localhost:5432",
		User:     "test_user",
		Database: "test",
		Password: "",
	})

	err := migrations.Run(db, directory, os.Args)
	if err != nil {
		log.Fatalln(err)
	}
}
