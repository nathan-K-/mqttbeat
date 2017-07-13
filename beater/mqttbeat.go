package beater

import (
	"fmt"
	"time"
	"encoding/json"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"

	MQTT "github.com/eclipse/paho.mqtt.golang"

	"gopkg.in/vmihailenco/msgpack.v2"

	"github.com/nathan-k-/mqttbeat/config"
)

type Mqttbeat struct {
	done   chan struct{}
	beat_config config.Config
	elastic_client publisher.Client
	mqtt_client MQTT.Client
}

// Prepare mqtt client
func setupMqttClient(bt *Mqttbeat) {
	mqtt_client_opt := MQTT.NewClientOptions()
	mqtt_client_opt.AddBroker(bt.beat_config.Broker_url)

	bt.mqtt_client = MQTT.NewClient(mqtt_client_opt)
}

// Creates beater
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	config := config.DefaultConfig
	if err := cfg.Unpack(&config); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &Mqttbeat{
		done:   make(chan struct{}),
		beat_config: config,
	}
	setupMqttClient(bt)

	if token := bt.mqtt_client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	logp.Info("MQTT Client connected")

	// Mqtt client - Subscribe to every topic in the config file, and bind with message handler
	if token := bt.mqtt_client.SubscribeMultiple(bt.beat_config.Topics_subscribe, bt.on_message);
	token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	return bt, nil
}

func DecodePayload(payload []byte) common.MapStr {
	event := make(common.MapStr)

	// A msgpack payload must be a json-like object
	err := msgpack.Unmarshal(payload, &event)
	if  err == nil {
		logp.Info("Payload decoded - msgpack")
		return event
	}

	err = json.Unmarshal(payload, &event)
	if  err == nil {
		logp.Info("Payload decoded - json")
		return event
	}

	// default case
	event["payload"]= string(payload)
	logp.Info("Payload decoded - text")
	return event
}

// Mqtt message handler
func (bt *Mqttbeat) on_message(client MQTT.Client, msg MQTT.Message) {
	logp.Info("MQTT MESSAGE RECEIVED " + string(msg.Payload()))

	event := make(common.MapStr) // common.MapStr = map[string]interface{}

	if bt.beat_config.Decode_paylod == true {
		event = DecodePayload(msg.Payload())
	} else {
		event = make(common.MapStr)
		event["payload"] = msg.Payload()
	}

	event["beat"]= common.MapStr{"index": "mqttbeat", "type":"message"}
	event["@timestamp"] = common.Time(time.Now())
	event["topic"] = msg.Topic()
	// Finally sending the message to elasticsearch
	bt.elastic_client.PublishEvent(event)
	logp.Info("Event sent")
	}


func (bt *Mqttbeat) Run(b *beat.Beat) error {
	logp.Info("mqttbeat is running! Hit CTRL-C to stop it.")
	bt.elastic_client = b.Publisher.Connect()

	// The mqtt client is asynchronous, so here we don't have anuthing to do
	for {
		select {
		case <-bt.done:
			return nil
		}
	}
}

func (bt *Mqttbeat) Stop() {
	bt.mqtt_client.Disconnect(250)
	bt.elastic_client.Close()
	close(bt.done)
}
