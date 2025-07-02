package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
)

type OpenMeteoParams struct {
	Latitude  float64 `url:"latitude"`
	Longitude float64 `url:"longitude"`
	Hourly    string  `url:"hourly"`
	Timezone  string  `url:"timezone"`
}

type Weather struct {
	Latitude             float64 `json:"latitude"`
	Longitude            float64 `json:"longitude"`
	GenerationtimeMs     float64 `json:"generationtime_ms"`
	UtcOffsetSeconds     int     `json:"utc_offset_seconds"`
	Timezone             string  `json:"timezone"`
	TimezoneAbbreviation string  `json:"timezone_abbreviation"`
	Elevation            float64 `json:"elevation"`
	HourlyUnits          struct {
		Time                     string `json:"time"`
		Temperature2M            string `json:"temperature_2m"`
		PrecipitationProbability string `json:"precipitation_probability"`
	} `json:"hourly_units"`
	Hourly struct {
		Time                     []string  `json:"time"`
		Temperature2M            []float64 `json:"temperature_2m"`
		PrecipitationProbability []int     `json:"precipitation_probability"`
	} `json:"hourly"`
}

func getWeather() Weather {
	// set params
	queryParams, _ := query.Values(OpenMeteoParams{
		Latitude:  13.754,
		Longitude: 100.5014,
		Hourly:    "temperature_2m,precipitation_probability",
		Timezone:  "Asia/Bangkok",
	})

	// fetch data
	url := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?%s", queryParams.Encode())
	resp, err := http.Get(url)
	if err != nil {
		log.Println("No response from request")
	}
	defer resp.Body.Close()

	// parse response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body")
	}

	var result Weather
	if err := json.Unmarshal(body, &result); err != nil {
		log.Println("Can not unmarshal JSON")
	}

	return result
}

func getCurrentHourInt() int64 {
	hours, _, _ := time.Now().Clock()
	currentHourString := fmt.Sprintf("%d", hours)

	currentHourInt, err := strconv.ParseInt(strings.TrimSpace(currentHourString), 10, 32)
	if err != nil {
		fmt.Println("Error converting string to int:", err)
	}

	return currentHourInt
}

func WeatherForecast() string {
	weather := getWeather()

	// get current hour
	currentHourInt := getCurrentHourInt()

	//// forecasting starts at "2024-06-10T00:00" of current date
	//currentTemperature := weather.Hourly.Temperature2M[currentHourInt]

	rainProbabilityCurrent := weather.Hourly.PrecipitationProbability[currentHourInt]
	rainProbabilityDeltaPlus1 := weather.Hourly.PrecipitationProbability[currentHourInt+1]
	rainProbabilityDeltaPlus2 := weather.Hourly.PrecipitationProbability[currentHourInt+2]
	rainProbabilityDeltaPlus3 := weather.Hourly.PrecipitationProbability[currentHourInt+3]

	// set output
	msg := ""
	for index, v := range []int{rainProbabilityCurrent, rainProbabilityDeltaPlus1, rainProbabilityDeltaPlus2, rainProbabilityDeltaPlus3} {
		var icon string
		if v <= 40 {
			icon = "🌥"
		} else {
			icon = "🌧️"
		}

		msg += fmt.Sprintf("%v:00 - %s %v%%\n", int(currentHourInt)+index, icon, v)
	}

	return msg
}
