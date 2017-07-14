package main

import (
	"os"

	"github.com/elastic/beats/libbeat/beat"

	"github.com/nathan-k-/mqttbeat/beater"
)

func main() {
	err := beat.Run("mqttbeat", "1.0.0", beater.New)
	if err != nil {
		os.Exit(1)
	}
}
