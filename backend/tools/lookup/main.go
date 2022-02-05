package main

import (
	"log"

	"github.com/kcsu/store/config"
	"github.com/kcsu/store/db"
	"github.com/kcsu/store/lookup"
)

// Run the lookup task once.
// This can be useful for testing with local/development databases.
func main() {
	// Initialise
	c := config.Init()
	d, err := db.Init(c)
	if err != nil {
		log.Panic(err)
	}

	if err := lookup.Run(c, d); err != nil {
		log.Panic(err)
	}
}
