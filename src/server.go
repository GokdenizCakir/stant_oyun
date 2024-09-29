package src

import (
	"log"

	"github.com/GokdenizCakir/stant_oyun/src/routes"
	"github.com/GokdenizCakir/stant_oyun/src/utils"
)

func Init(config utils.Config) {
	r := routes.Router

	err := r.Run(":" + config.PORT)
	if err != nil {
		log.Fatalf("failed to run server: %v\n", err)
	}
}
