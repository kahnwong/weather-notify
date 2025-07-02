package main

import (
	"github.com/rs/zerolog/log"
)

func main() {
	// check weather
	outputMessage := WeatherForecast()
	log.Info().Msg(outputMessage)

	// send notification
	if outputMessage != "" {
		err := notify(outputMessage)
		if err != nil {
			log.Error().Msg("Error sending notification")
		}
	} else {
		log.Info().Msg("Your location is not affected with no running water.")
	}
}
