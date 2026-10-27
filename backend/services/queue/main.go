package main

import (
	"net/http"
	"os"

	"github.com/kcsu/store/alerts"
	"github.com/kcsu/store/config"
	"github.com/kcsu/store/db"
	"github.com/kcsu/store/logging"
	"github.com/kcsu/store/queue"
	"gorm.io/gorm"
)

var c *config.Config
var d *gorm.DB
var f db.FormalStore

func main() {
	// Initialise data
	c = config.Init()
	var err error
	d, err = db.Init(c)
	if err != nil {
		alerts.Fatal("connecting to the database failed", err)
	}
	f = db.NewFormalStore(d)
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
		alerts.Fatal("server stopped", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	if err := queue.Run(c, d, f); err != nil {
		alerts.Send("queue-failed", "Ticket queue run failed", err, nil)
		w.WriteHeader(500)
	}
}
