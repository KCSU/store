package main

import (
	"fmt"
	"os"

	"github.com/kcsu/store/logging"
	"github.com/kcsu/store/route"
)

func main() {
	e := route.Init()
	port := os.Getenv("PORT")
	if port == "" {
		port = "1323"
	}
	e.HidePort = true
	logging.Infof("listening on port %s", port)
	logging.Fatalf("server stopped: %v", e.Start(fmt.Sprintf(":%s", port)))
}
