package main

import (
	"log"

	"github.com/GokdenizCakir/stant_oyun/src/db"
	"github.com/GokdenizCakir/stant_oyun/src/routes"
	"github.com/GokdenizCakir/stant_oyun/src/utils"

	server "github.com/GokdenizCakir/stant_oyun/src"
)

func main() {

	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal(err)
	}

	db.Init(config)
	routes.Init()
	server.Init(config)
}
