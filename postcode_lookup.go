package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
)

const OPENSTREETMAP_URL = "https://nominatim.openstreetmap.org/search"
const COUNTRY = "australia"
const FORMAT = "json"

type GeoSearchResult struct {
	Latitude    string `json:"lat"`
	Longitude   string `json:"lon"`
	DisplayName string `json:"display_name"`
}

func searchLocation(postcode string) {
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
	}
	if resp.StatusCode != http.StatusOK {
		log.Println("Error: Status code not 200 OK, response status code", resp.Status)
	}
	defer resp.Body.Close()

	var results []*GeoSearchResult
	err = json.NewDecoder(resp.Body).Decode(&results)
	if err != nil {
		log.Println("Decode error:", err)
	}

	for _, r := range results {
		log.Println(r)
	}
}
