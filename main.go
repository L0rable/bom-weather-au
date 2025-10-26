package main

import (
	"io"
	"log"

	"github.com/L0rable/bom-weather-au/internal"
)

func main() {
	loc := internal.SearchLocation("3000")

	if *loc == (internal.Location{}) {
		log.Fatalln("Not a valid postcode (main.go)")
		return
	}

	conn := internal.OpenFtpServer()
	obsURL := internal.AusObservationState[loc.State]
	resp, err := conn.Retr(obsURL)
	if err != nil {
		log.Fatal(err)
	}

	data, err := io.ReadAll(resp)
	if err != nil {
		log.Println(data)
	}

	stnData := internal.UnmarshalXML(data)
	stn := internal.GetClosetStation(loc, stnData)
	log.Println("closets station:", stn)

	internal.CloseFtpServer(conn)
}
