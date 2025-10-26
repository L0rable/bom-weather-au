package internal

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const OPENSTREETMAP_URL = "https://nominatim.openstreetmap.org/search"
const COUNTRY = "australia"
const FORMAT = "json"

type GeoSearchResult struct {
	Latitude    string `json:"lat"`
	Longitude   string `json:"lon"`
	DisplayName string `json:"display_name"`
}

type Location struct {
	Postcode  string
	Suburb    string
	State     string
	Latitude  float64
	Longitude float64
}

func CheckAusPostcode(postcode string) string {
	_, err := strconv.Atoi(postcode)
	if err != nil {
		log.Println("Err:", err, "checkAusPostcode (location.go)")
		return ""
	}
	if len(postcode) < 3 || len(postcode) > 4 {
		log.Println("Err: Invalid postcode length", len(postcode), "checkAusPostcode (location.go)")
		return ""
	}

	if len(postcode) == 3 {
		postcode = "0" + postcode
	}
	return postcode
}

func SearchLocation(postcode string) *Location {
	postcode = CheckAusPostcode(postcode)
	if postcode == "" {
		log.Println("Err: Invalid postcode", postcode)
		return &Location{}
	}

	query := url.Values{}
	query.Set("postalcode", postcode)
	query.Set("country", COUNTRY)
	query.Set("format", FORMAT)

	reqURL := OPENSTREETMAP_URL + "?" + query.Encode()
	req, _ := http.NewRequest("GET", reqURL, nil)
	req.Header.Set("User-Agent", "Go-Geocoder")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Error:", err)
		return &Location{}
	}
	if resp.StatusCode != http.StatusOK {
		log.Println("Error: Status code not 200 OK, response status code", resp.Status)
		return &Location{}
	}
	defer resp.Body.Close()

	var results []*GeoSearchResult
	err = json.NewDecoder(resp.Body).Decode(&results)
	if err != nil {
		log.Println("Decode error:", err)
		return &Location{}
	}

	if len(results) == 0 {
		log.Println("Postcode is not found on OpenStreetMapApi")
		return &Location{}
	}

	for _, r := range results {
		log.Println(r)
	}

	locationData := strings.Split(results[0].DisplayName, ", ")
	// TODO: need to properly transfer the data (at least 5 pieces of data in display_name)
	lat, err1 := strconv.ParseFloat(results[0].Latitude, 64)
	long, err2 := strconv.ParseFloat(results[0].Longitude, 64)
	if err1 != nil || err2 != nil {
		log.Println("Err1 lat:", err1, " searchLocation() location.go")
		log.Println("Err2 long:", err2, " searchLocation() location.go")
		return &Location{}
	}

	location := &Location{
		Postcode:  locationData[0],
		Suburb:    locationData[1],
		State:     locationData[len(locationData)-2],
		Latitude:  lat,
		Longitude: long,
	}

	return location
}
