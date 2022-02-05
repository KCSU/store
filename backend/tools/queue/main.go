package main

import (
	"log"

	"github.com/kcsu/store/config"
	"github.com/kcsu/store/db"
	"github.com/kcsu/store/queue"
)

func main() {
	// Initialise data
	c := config.Init()
	d, err := db.Init(c)
	if err != nil {
		log.Panic(err)
	}
	f := db.NewFormalStore(d)
	if err := queue.Run(c, d, f); err != nil {
		log.Panic(err)
	}
}
