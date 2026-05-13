package main

import (
	"fmt"

	models "encomers/license/internal/domain/valueObjects"
	"encomers/license/internal/endpoint/app"
	"encomers/license/internal/endpoint/logger"
)

var vocab []rune = []rune("АЕТОРНУКХСВМ")
var region, _ = models.NewRegionPart("116", "RUS")

func main() {

	log := logger.New(
		logger.WithLevel("release"),
		logger.WithAppName("plate-generator"),
		logger.WithVersion("1.0.0"),
	)

	defer log.Sync()

	app := app.New(log)
	if err := app.Run(":8080"); err != nil {
		fmt.Printf("failed to run the server: %v\n", err)
	}
}
