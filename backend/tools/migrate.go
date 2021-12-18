package main

import (
	"log"

	"github.com/kcsu/store/config"
	"github.com/kcsu/store/db"
	"github.com/kcsu/store/model"
)

func main() {
	c := config.Init()
	d, err := db.Init(c)
	if err != nil {
		log.Fatal(err)
	}
	// TODO: manual migration
	d.AutoMigrate(
		&model.Formal{},
		&model.Ticket{},
	)
}
