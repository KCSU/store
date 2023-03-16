package main

import (
	"log"

	"github.com/kcsu/store/config"
	"github.com/kcsu/store/db"
	"github.com/kcsu/store/model"
)

// Migrate the database model
func main() {
	c := config.Init()
	d, err := db.Init(c)
	if err != nil {
		log.Panic(err)
	}
	// TODO: manual migration
	err = d.AutoMigrate(
		&model.Bill{},
		&model.Formal{},
		&model.Ticket{},
		&model.ManualTicket{},
		&model.User{},
		&model.Group{},
		&model.GroupUser{},
		&model.Role{},
		&model.Permission{},
		&model.AccessLog{},
		&model.Pigeonhole{},
		&model.ExtraCharge{},
	)
	if err != nil {
		log.Panic(err)
	}
}
