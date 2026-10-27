package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/kcsu/store/alerts"
	"github.com/kcsu/store/config"
	"github.com/kcsu/store/db"
	"github.com/kcsu/store/logging"
	"github.com/kcsu/store/lookup"
	lookupv1 "github.com/kcsu/store/proto/kcsu/lookup/v1"
	"google.golang.org/protobuf/proto"
)

var c *config.Config
var g db.GroupStore

func main() {
	// Initialise
	c = config.Init()
	var err error
	d, err := db.Init(c)
	if err != nil {
		alerts.Fatal("connecting to the database failed", err)
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
		alerts.Fatal("server stopped", err)
	}
}

// Handle a Pub/Sub push from kcsu-admin's hourly lookup sync
func handler(w http.ResponseWriter, r *http.Request) {
	var push struct {
		Message struct {
			Data []byte `json:"data"`
		} `json:"message"`
	}
	snapshot := &lookupv1.GroupSnapshot{}
	err := json.NewDecoder(r.Body).Decode(&push)
	if err == nil {
		err = proto.Unmarshal(push.Message.Data, snapshot)
	}
	if err != nil || snapshot.Name == "" {
		// Redelivering can't fix a malformed message
		alerts.Send("lookup-sync-failed", "Dropped a malformed Lookup snapshot", err, nil)
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if err := lookup.ApplySnapshot(g, snapshot); err != nil {
		// Probably transient; reply with 500 for Pub/Sub to retry delivery
		alerts.Send("lookup-sync-failed", "Lookup sync failed", err, map[string]string{"lookup_group": snapshot.Name})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
