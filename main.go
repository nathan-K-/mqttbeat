package main

import (
	"os"

	"github.com/elastic/beats/libbeat/beat"

	"github.com/nathan-k-/mqttbeat/beater"
)

func main() {
	err := beat.Run("mqttbeat", "", beater.New)
	if err != nil {
		os.Exit(1)
	}
}
