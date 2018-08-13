package main

import (
	"log"
	"os"

	"github.com/go-pg/pg"
	"github.com/robinjoseph08/go-pg-migrations"
)

const directory = "example"

func main() {
	db := pg.Connect(&pg.Options{
		Addr:     "localhost:5432",
		User:     "test_user",
		Database: "test",
	})

	err := migrations.Run(db, directory, os.Args)
	if err != nil {
		log.Fatalln(err)
	}
}
