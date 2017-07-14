// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import "time"

type Config struct {
	BrokerUrl string `config:"broker_url"`
	TopicsSubscribe []string `config:"topics_subscribe"`
	DecodePaylod bool `config:"decode_payload"`
	Period          time.Duration `config:"period"`
}

var DefaultConfig = Config{
	BrokerUrl: "tcp://localhost:1883",
	TopicsSubscribe: []string{"/test/mqttbeat/#?1"},
	DecodePaylod: true,
}
