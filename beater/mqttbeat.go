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
	beatConfig config.Config
	elasticClient publisher.Client
	mqttClient MQTT.Client
}

// Prepare mqtt client
func setupMqttClient(bt *Mqttbeat) {
	mqttClientOpt := MQTT.NewClientOptions()
	mqttClientOpt.AddBroker(bt.beatConfig.BrokerUrl)

	bt.mqttClient = MQTT.NewClient(mqttClientOpt)
}

// Creates beater
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	config := config.DefaultConfig
	if err := cfg.Unpack(&config); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &Mqttbeat{
		done:   make(chan struct{}),
		beatConfig: config,
	}
	setupMqttClient(bt)

	if token := bt.mqttClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	logp.Info("MQTT Client connected")

	// Mqtt client - Subscribe to every topic in the config file, and bind with message handler
	if token := bt.mqttClient.SubscribeMultiple(bt.beatConfig.TopicsSubscribe, bt.onMessage);
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
func (bt *Mqttbeat) onMessage(client MQTT.Client, msg MQTT.Message) {
	logp.Info("MQTT MESSAGE RECEIVED " + string(msg.Payload()))

	event := make(common.MapStr) // common.MapStr = map[string]interface{}

	if bt.beatConfig.DecodePaylod == true {
		event = DecodePayload(msg.Payload())
	} else {
		event = make(common.MapStr)
		event["payload"] = msg.Payload()
	}

	event["beat"]= common.MapStr{"index": "mqttbeat", "type":"message"}
	event["@timestamp"] = common.Time(time.Now())
	event["topic"] = msg.Topic()
	// Finally sending the message to elasticsearch
	bt.elasticClient.PublishEvent(event)
	logp.Info("Event sent")
	}


func (bt *Mqttbeat) Run(b *beat.Beat) error {
	logp.Info("mqttbeat is running! Hit CTRL-C to stop it.")
	bt.elasticClient = b.Publisher.Connect()

	// The mqtt client is asynchronous, so here we don't have anuthing to do
	for {
		select {
		case <-bt.done:
			return nil
		}
	}
}

func (bt *Mqttbeat) Stop() {
	bt.mqttClient.Disconnect(250)
	bt.elasticClient.Close()
	close(bt.done)
}
