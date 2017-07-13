// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import "time"

type Config struct {
	Broker_url string `config:"broker_url"`
	Topics_subscribe map[string]byte `config:"topics_subscribe"`
	Decode_paylod bool `config:"decode_payload"`
	Period          time.Duration `config:"period"`
}

var DefaultConfig = Config{
	Broker_url: "tcp://localhost:1883",
	Topics_subscribe: map[string]byte{"/test/mqttbeat/#":1},
	Decode_paylod: true,
}
