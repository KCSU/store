package main

import (
	"fmt"
	"os"

	"github.com/kcsu/store/route"
)

func main() {
	e := route.Init()
	port := os.Getenv("PORT")
	if port == "" {
		port = "1323"
	}
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))
}
