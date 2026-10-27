package main

import (
	"net/http"
	"os"

	"github.com/kcsu/store/config"
	"github.com/kcsu/store/db"
	"github.com/kcsu/store/logging"
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
		logging.Fatalf("connecting to the database failed: %v", err)
	}
	g = db.NewGroupStore(d)
	http.HandleFunc("/", handler)

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		logging.Infof("defaulting to port %s", port)
	}

	// Start HTTP server.
	logging.Infof("listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		logging.Fatalf("server stopped: %v", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	if err := lookup.Run(c, g); err != nil {
		logging.Errorf("syncing with Lookup failed: %v", err)
		w.WriteHeader(500)
	}
}
