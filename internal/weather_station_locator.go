package internal

import "math"

const EARTH_RADIUS = 6371

func haversine(loc *Location, stn *Station) float64 {
	lat1 := loc.Latitude * math.Pi / 180
	lon1 := loc.Longitude * math.Pi / 180
	lat2 := stn.latitude * math.Pi / 180
	lon2 := stn.longitude * math.Pi / 180

	deltaLat := lat2 - lat1
	deltaLon := lon2 - lon1
	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) + math.Cos(lat1)*math.Cos(lat2)*math.Sin(deltaLon/2)*math.Sin(deltaLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return EARTH_RADIUS * c
}

func GetClosetStation(loc *Location, stns []*Station) *Station {
	var closetStn *Station

	minDistance := math.MaxFloat64
	for _, stn := range stns {
		dist := haversine(loc, stn)

		if dist < minDistance {
			minDistance = dist
			closetStn = stn
		}
	}
	return closetStn
}
