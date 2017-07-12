// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import "time"

type Config struct {
	Period time.Duration `config:"period"`
	Broker_url string `config:"broker_url"`
	Topics_subscribe map[string]byte `config:"topics_subscribe"`
	Decode_paylod bool `config:"decode_payload"`
}

var DefaultConfig = Config{
	Period: 1 * time.Second,
	Broker_url: "tcp://localhost:1883",
	Topics_subscribe: map[string]byte{"/test/mqttbeat/#":1},
	Decode_paylod: true,
}
