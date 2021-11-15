package main

import "github.com/kcsu/store/route"

func main() {
	e := route.Init()
	e.Logger.Fatal(e.Start(":1323"))
}
