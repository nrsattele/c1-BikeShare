package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
)

var data = make(map[string]Station)

/**
 * Gets the CSV data from github then analyzes it
 */
func main() {
	resp, err := http.Get("https://nrsattele.github.io/c1-BikeShare/data/metro-bike-share-trip-data.csv")
	if err != nil {
		log.Fatal(err)
	}

	reader := bufio.NewReader(resp.Body)
	var header = true
	for {
		line, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		if !header {
			var fields = strings.Split(string(line), ",")
			analyze(fields, true)
			analyze(fields, false)

		} else {
			header = false
		}
	}

	// Exports to CSV
	file, err := os.Create("data/preprocessed.csv")
	if err != nil {
		log.Fatal("Cannot create to file", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	headers := []string{"stationID", "incoming", "leaving", "avgTimeSpent", "percentRoundTrip", "percentPassholder", "latitude", "longitude"}

	err = writer.Write(headers)
	if err != nil {
		log.Fatal("Cannot wrtie to file", err)
	}

	//Sorts Data by Station ID
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys) //sort by key

	for _, key := range keys {
		station := data[key]
		s := []string{station.stationID, fmt.Sprintf("%v", station.incoming), fmt.Sprintf("%v", station.leaving), fmt.Sprintf("%f", station.avgTimeSpent), fmt.Sprintf("%f", station.percentRoundTrip), fmt.Sprintf("%f", station.percentPassholder), station.stationLatitude, station.stationLongitude}
		err := writer.Write(s)
		if err != nil {
			log.Fatal("Cannot wrtie to file", err)
		}
	}

	// Exports to GeoJSON
	geo := GeoJSON{
		Type:     "FeatureCollection",
		Features: []Feature{},
	}

	for _, station := range data {
		var feat Feature
		feat.Constructor(station)
		geo.Features = append(geo.Features, feat)
	}

	gJSON, err := json.Marshal(geo)
	if err != nil {
		log.Fatal("Error converting geoJSON to JSON object")
	}

	err = ioutil.WriteFile("data/output.geojson", gJSON, 0644)
	if err != nil {
		log.Fatal("Error writing to file")
	}
	fmt.Printf("%+v", gJSON)

}

//Trip ID, Duration, Start Time, End Time, Starting Station ID,	Starting Station Latitude,
//	6 - Starting Station Longitude,	Ending Station ID,	Ending Station Latitude, Ending Station Longitude
//	10 - Bike ID, Plan Duration, Trip Route Category, Passholder Type, Starting Lat-Long, Ending Lat-Long

/**
 * Analyzes each line of data.
 * Requires that line is comma separated data with no more entires than size of headers
 * analyzeStarting is a boolean that indicates if this should update the starting station entry or the ending station entry
 * Updates
 */
func analyze(fields []string, analyzeStarting bool) {
	var station Station
	var ok bool
	if analyzeStarting {
		// Don't want to analyze if we don't have the station ID
		if len(fields[4]) < 1 {
			return
		}
		station, ok = data[fields[4]]
	} else {
		// Don't want to analyze if we don't have the station ID
		if len(fields[7]) < 1 {
			return
		}
		station, ok = data[fields[7]]
	}

	//Checks if data has already recorded this station
	if ok {
		//increment counter of bikes leaving or entering this station
		if analyzeStarting {
			station.leaving = station.leaving + 1
		} else {
			station.incoming = station.incoming + 1
		}

		//Updates average time spent on bike when coming into or out of this station
		timeSpent, err := strconv.ParseFloat(fields[11], 64)
		if err != nil {
			fmt.Println("Incorrect Time Spent")
		} else {
			numTrips := float64(station.incoming+station.leaving) + 1
			totalTime := (station.avgTimeSpent * numTrips) + timeSpent
			station.avgTimeSpent = totalTime / numTrips
		}

		// Breaks down the percent round trip into a fraction and updates the numerator (numRoundTrip) and denominator
		numRoundTrip := station.percentRoundTrip * float64(station.incoming+station.leaving)
		if fields[12] == "Round Trip" {
			numRoundTrip++
		}
		station.percentRoundTrip = numRoundTrip / float64(station.incoming+station.leaving+1)

		// Breaks down percent pass into a fraction and updates the numerator (numRoundTrip) and denominator
		numPass := station.percentPassholder * float64(station.incoming+station.leaving)
		if fields[13] != "Walk-up" {
			numPass++
		}
		station.percentPassholder = numPass / float64(station.incoming+station.leaving+1)

		//Updates the latitude and longitude of the station if not correct
		if station.stationLatitude == "" || station.stationLatitude == "0" {
			if analyzeStarting {
				station.stationLatitude = fields[5]
			} else {
				station.stationLatitude = fields[8]
			}
		}
		if station.stationLongitude == "" || station.stationLongitude == "0" {
			if analyzeStarting {
				station.stationLongitude = fields[6]
			} else {
				station.stationLongitude = fields[9]
			}
		}

		data[station.stationID] = station
	} else {
		//Creates new station
		timeSpent, err := strconv.ParseFloat(fields[11], 64)
		if err != nil {
			fmt.Println("Incorrect Time Spent")
		}
		pRoundTrip := 0.0
		if fields[12] == "Round Trip" {
			pRoundTrip = 1.0
		}
		pPass := 0.0
		if fields[13] != "Walk-up" {
			pPass = 1.0
		}

		in := 0
		out := 0
		id := ""
		lat := ""
		lon := ""
		if analyzeStarting {
			out = 1
			in = 0
			id = fields[4]
			lat = fields[5]
			lon = fields[6]
		} else {
			out = 0
			in = 1
			id = fields[7]
			lat = fields[8]
			lon = fields[9]
		}

		station := Station{
			stationID:         id,
			incoming:          in,
			leaving:           out,
			avgTimeSpent:      timeSpent,
			percentRoundTrip:  pRoundTrip,
			percentPassholder: pPass,
			stationLatitude:   lat,
			stationLongitude:  lon,
		}

		data[station.stationID] = station
	}
}

// Station Comment
type Station struct {
	stationID         string
	incoming          int
	leaving           int
	avgTimeSpent      float64
	percentRoundTrip  float64
	percentPassholder float64
	stationLatitude   string
	stationLongitude  string
	//Average time bikes are taken from station
	//avgTimeOfDay time.Time
}

// GeoJSON Comment
type GeoJSON struct {
	Type     string    `json:"type"`
	Features []Feature `json:"features"`
}

// Feature Comment
type Feature struct {
	Type       string `json:"type"`
	Properties struct {
		Name              string  `json:"name"`
		Incoming          int     `json:"incoming"`
		Leaving           int     `json:"leaving"`
		AvgTimeSpent      float64 `json:"avgTimeSpent"`
		PercentRoundTrip  float64 `json:"percentRoundTrip"`
		PercentPassholder float64 `json:"percentPassholder"`
	} `json:"properties"`
	Geometry struct {
		Type        string    `json:"type"`
		Coordinates []float64 `json:"coordinates"`
	} `json:"geometry"`
}

// Constructor Comment
func (f *Feature) Constructor(s Station) {
	f.Type = "Feature"
	f.Properties.Name = s.stationID
	f.Geometry.Type = "Point"
	f.Properties.Incoming = s.incoming
	f.Properties.Leaving = s.leaving
	f.Properties.AvgTimeSpent = s.avgTimeSpent
	f.Properties.PercentRoundTrip = s.percentRoundTrip
	f.Properties.PercentPassholder = s.percentPassholder

	// Converts coordinates from string to float
	latFloat, err := strconv.ParseFloat(s.stationLatitude, 64)
	if err != nil {
		fmt.Println("Cannot convert Latitude, " + s.stationLatitude + " to a float")
		latFloat = 0.0
	}
	lonFloat, err := strconv.ParseFloat(s.stationLongitude, 64)
	if err != nil {
		fmt.Println("Cannot convert Longitude, " + s.stationLongitude + " to a float")
		latFloat = 0.0
	}
	f.Geometry.Coordinates = []float64{latFloat, lonFloat}
}
