// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import "time"

// Config represents every needed configuration fields
type Config struct {
	BrokerURL       string        `config:"broker_url"`
	BrokerUsername  string        `config:"broker_username"`
	BrokerPassword  string        `config:"broker_password"`
	TopicsSubscribe []string      `config:"topics_subscribe"`
	DecodePaylod    bool          `config:"decode_payload"`
	Period          time.Duration `config:"period"`
}

// DefaultConfig will be used if no config file is founded
var DefaultConfig = Config{
	BrokerURL:       "tcp://localhost:1883",
	BrokerUsername:	 "",
	BrokerPassword:	 "",
	TopicsSubscribe: []string{"/test/mqttbeat/#?1"},
	DecodePaylod:    true,
}
