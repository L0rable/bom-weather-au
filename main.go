package main

import (
	"encoding/xml"
	"io"
	"log"
	"strconv"
	"time"

	"github.com/jlaffaye/ftp"
)

const FTP_URL = "ftp.bom.gov.au:21"

type ObservationSummary struct {
	XMLName     xml.Name `xml:"product"`
	Observation struct {
		Station []struct {
			WmoId              string `xml:"wmo-id,attr"`
			BomId              string `xml:"bom-id,attr"`
			Timezone           string `xml:"tz,attr"`
			StandardName       string `xml:"stn-name,attr"`
			StandardHeight     string `xml:"stn-height,attr"`
			Type               string `xml:"type,attr"`
			Latitude           string `xml:"lat,attr"`
			Longitude          string `xml:"lon,attr"`
			ForecastDistrictId string `xml:"forecast-district-id,attr"`
			Description        string `xml:"description,attr"`
		} `xml:"station"`
	} `xml:"observations"`
}

type Station struct {
	name        string
	description string
	timezone    string
	latitude    float64
	longitude   float64
}

func openFtpServer() *ftp.ServerConn {
	conn, err := ftp.Dial(FTP_URL, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		log.Fatal(err)
	}

	err = conn.Login("anonymous", "")
	if err != nil {
		log.Fatal(err)
	}
	return conn
}

func closeFtpServer(conn *ftp.ServerConn) {
	err := conn.Quit()
	if err != nil {
		log.Fatal(err)
	}
}

func UnmarshalXML(data []byte) []*Station {
	var obsSummary *ObservationSummary
	err := xml.Unmarshal([]byte(data), &obsSummary)
	if err != nil {
		log.Fatal("Error unmarshalling XML:", err)
	}

	var stations []*Station
	summaryStations := obsSummary.Observation.Station
	for _, summaryStation := range summaryStations {
		lat, err1 := strconv.ParseFloat(summaryStation.Latitude, 64)
		long, err2 := strconv.ParseFloat(summaryStation.Longitude, 64)
		if err1 != nil || err2 != nil {
			log.Println("Err1 lat: ", err1, " UnmarshalXML() (main.go)")
			log.Println("Err2 long: ", err2, " UnmarshalXML() (main.go)")
			return stations
		}

		stn := &Station{
			name:        summaryStation.StandardName,
			description: summaryStation.Description,
			timezone:    summaryStation.Timezone,
			latitude:    lat,
			longitude:   long,
		}
		stations = append(stations, stn)
	}

	return stations
}

func main() {
	loc := searchLocation("3000")

	conn := openFtpServer()
	obsURL := "/anon/gen/fwo/IDV60920.xml"
	resp, err := conn.Retr(obsURL)
	if err != nil {
		log.Fatal(err)
	}

	data, err := io.ReadAll(resp)
	if err != nil {
		log.Println(data)
	}

	stnData := UnmarshalXML(data)
	stn := getClosetStation(loc, stnData)
	log.Println("closets station: ", stn)

	closeFtpServer(conn)
}
