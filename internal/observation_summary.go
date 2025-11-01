package internal

import (
	"encoding/xml"
	"io"
	"log"
	"strconv"
)

// Observations - State/Territory summaries (http://www.bom.gov.au/catalogue/data-feeds.shtml)
const NSW_ACT = "/anon/gen/fwo/IDN60920.xml"
const NT = "/anon/gen/fwo/IDD60920.xml"
const QLD = "/anon/gen/fwo/IDQ60920.xml"
const SA = "/anon/gen/fwo/IDS60920.xml"
const TAS = "/anon/gen/fwo/IDT60920.xml"
const VIC = "/anon/gen/fwo/IDV60920.xml"
const WA = "/anon/gen/fwo/IDW60920.xml"

var AusObservationState = map[string]string{
	"New South Wales":              NSW_ACT,
	"Australian Capital Territory": NSW_ACT,
	"Northern Territory":           NT,
	"Queensland":                   QLD,
	"South Australia":              SA,
	"Tasmania":                     TAS,
	"Victoria":                     VIC,
	"Western Australia":            WA,
}

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
			Period             struct {
				Level struct {
					BasicElement []struct {
						Text           string `xml:",chardata"`
						Units          string `xml:"units,attr"`
						Type           string `xml:"type,attr"`
						StartTimeLocal string `xml:"start-time-local,attr"`
						EndTimeLocal   string `xml:"end-time-local,attr"`
						Duration       string `xml:"duration,attr"`
						StartTimeUtc   string `xml:"start-time-utc,attr"`
						EndTimeUtc     string `xml:"end-time-utc,attr"`
						Instance       string `xml:"instance,attr"`
						TimeUtc        string `xml:"time-utc,attr"`
						TimeLocal      string `xml:"time-local,attr"`
					} `xml:"element"`
				} `xml:"level"`
			} `xml:"period"`
		} `xml:"station"`
	} `xml:"observations"`
}

type BasicElement struct {
	Units string
	Type  string
	Value string
}

type Element struct {
	StartTimeLocal string
	EndTimeLocal   string
	Duration       string
	StartTimeUTC   string
	EndTimeUTC     string

	Units string
	Type  string
	Value string
}

type Station struct {
	// station tag
	name        string
	description string
	timezone    string
	latitude    float64
	longitude   float64
	// element tag (elements)
	ApparentTemp BasicElement
	DeltaTemp    BasicElement
	GustKmh      BasicElement
	WindGustSpd  BasicElement
	AirTemp      BasicElement
	DewPoint     BasicElement
	Pres         BasicElement
	MslPres      BasicElement
	QnhPres      BasicElement
	RelHumidity  BasicElement
	VisKm        BasicElement
	WindDir      string
	WindDirDeg   BasicElement
	WindSpdKmh   BasicElement
	WindSpd      BasicElement
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

		// elements := map[string]string{}
		elements := map[string]BasicElement{}
		for _, e := range summaryStation.Period.Level.BasicElement {
			if e.Type != "" {
				elements[e.Type] = BasicElement{
					Units: e.Units,
					Type:  e.Type,
					Value: e.Text,
				}
			}
		}

		stn := &Station{
			name:        summaryStation.StandardName,
			description: summaryStation.Description,
			timezone:    summaryStation.Timezone,
			latitude:    lat,
			longitude:   long,
			// values form elements can be store in Station
			ApparentTemp: elements["apparent_temp"],
			DeltaTemp:    elements["delta_t"],
			GustKmh:      elements["gust_kmh"],
			WindGustSpd:  elements["wind_gust_spd"],
			AirTemp:      elements["air_temperature"],
			DewPoint:     elements["dew_point"],
			Pres:         elements["pres"],
			MslPres:      elements["msl_pres"],
			QnhPres:      elements["qnh_pres"],
			RelHumidity:  elements["rel-humidity"],
			VisKm:        elements["vis_km"],
			WindDir:      elements["wind_dir"].Value,
			WindDirDeg:   elements["wind_dir_deg"],
			WindSpdKmh:   elements["wind_spd_kmh"],
			WindSpd:      elements["wind_spd"],
		}
		stations = append(stations, stn)
	}

	return stations
}

func GetObservationSummary(postcode string) *Station {
	loc := SearchLocation(postcode)
	if *loc == (Location{}) {
		log.Fatalln("Not a valid postcode (main.go)")
		return &Station{}
	}

	conn := OpenFtpServer()
	obsURL := AusObservationState[loc.State]
	resp, err := conn.Retr(obsURL)
	if err != nil {
		log.Fatal(err)
	}

	data, err := io.ReadAll(resp)
	if err != nil {
		log.Println(data)
	}

	stnData := UnmarshalXML(data)
	stn := GetClosetStation(loc, stnData)
	log.Println("closets station:", stn)

	CloseFtpServer(conn)

	return stn
}
