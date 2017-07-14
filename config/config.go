// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

type Config struct {
	BrokerUrl string `config:"broker_url"`
	TopicsSubscribe []string `config:"topics_subscribe"`
	DecodePaylod bool `config:"decode_payload"`
}

var DefaultConfig = Config{
	BrokerUrl: "tcp://localhost:1883",
	TopicsSubscribe: []string{"test/mqttbeat/#?1"},
	DecodePaylod: true,
}
