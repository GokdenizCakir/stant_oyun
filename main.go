package main

import (
	"github.com/GokdenizCakir/stant_oyun/src/db"
	"github.com/GokdenizCakir/stant_oyun/src/routes"

	server "github.com/GokdenizCakir/stant_oyun/src"
)

func main() {
	db.Init()
	routes.Init()
	server.Init()
}
