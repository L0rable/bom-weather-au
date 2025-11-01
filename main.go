package main

import (
	"log"

	"github.com/L0rable/bom-weather-au/internal"
)

func main() {
	station := internal.GetObservationSummary("3000")
	log.Println(station)
}
