package main

import (
	"log"
	"net/http"
	"os"

	"github.com/kcsu/store/config"
	"github.com/kcsu/store/db"
	"github.com/kcsu/store/lookup"
)

var c *config.Config
var g db.GroupStore

func main() {
	// Initialise
	c = config.Init()
	var err error
	d, err := db.Init(c)
	if err != nil {
		log.Panic(err)
	}
	g = db.NewGroupStore(d)
	http.HandleFunc("/", handler)

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	// Start HTTP server.
	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	if err := lookup.Run(c, g); err != nil {
		w.WriteHeader(500)
	}
}
