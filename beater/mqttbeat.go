package beater

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher/bc/publisher"

	MQTT "github.com/eclipse/paho.mqtt.golang"

	"gopkg.in/vmihailenco/msgpack.v2"

	"github.com/nathan-k-/mqttbeat/config"
)

// Mqttbeat represent a mqtt beat object
type Mqttbeat struct {
	done          chan struct{}
	beatConfig    config.Config
	elasticClient publisher.Client
	mqttClient    MQTT.Client
}

// Prepare mqtt client
func setupMqttClient(bt *Mqttbeat) {
	mqttClientOpt := MQTT.NewClientOptions()
	mqttClientOpt.AddBroker(bt.beatConfig.BrokerURL)
	logp.Info("BROKER url " + bt.beatConfig.BrokerURL)

	bt.mqttClient = MQTT.NewClient(mqttClientOpt)
}

// New function creates our mqtt beater
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	config := config.DefaultConfig
	if err := cfg.Unpack(&config); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &Mqttbeat{
		done:       make(chan struct{}),
		beatConfig: config,
	}
	setupMqttClient(bt)

	if token := bt.mqttClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	logp.Info("MQTT Client connected")

	subscriptions := ParseTopics(bt.beatConfig.TopicsSubscribe)
	//bt.beatConfig.TopicsSubscribe

	// Mqtt client - Subscribe to every topic in the config file, and bind with message handler
	if token := bt.mqttClient.SubscribeMultiple(subscriptions, bt.onMessage); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	return bt, nil
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

	event["beat"] = common.MapStr{"index": "mqttbeat", "type": "message"}
	event["@timestamp"] = common.Time(time.Now())
	event["topic"] = msg.Topic()
	// Finally sending the message to elasticsearch
	bt.elasticClient.PublishEvent(event)
	logp.Info("Event sent")
}

// Run is used to start this beater, once configured and connected
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

// Stop is used to close this beater
func (bt *Mqttbeat) Stop() {
	bt.mqttClient.Disconnect(250)
	bt.elasticClient.Close()
	close(bt.done)
}

// DecodePayload will try to decode the payload. If every check fails, it will
// return the payload as a string
func DecodePayload(payload []byte) common.MapStr {
	event := make(common.MapStr)

	// A msgpack payload must be a json-like object
	err := msgpack.Unmarshal(payload, &event)
	if err == nil {
		logp.Info("Payload decoded - msgpack")
		return event
	}

	err = json.Unmarshal(payload, &event)
	if err == nil {
		logp.Info("Payload decoded - json")
		return event
	}

	// default case
	event["payload"] = string(payload)
	logp.Info("Payload decoded - text")
	return event
}

// ParseTopics will parse the config file and return a map with topic:QoS
func ParseTopics(topics []string) map[string]byte {
	subscriptions := make(map[string]byte)
	for _, value := range topics {
		// Fist, spliting the string topic?qos
		topic, qosStr := strings.Split(value, "?")[0], strings.Split(value, "?")[1]
		// Then, parsing the qos to an int
		qosInt, err := strconv.ParseInt(qosStr, 10, 0)
		if err != nil {
			panic("Error parsing topics")
		}
		// Finally, filling the subscriptions map
		subscriptions[topic] = byte(qosInt)
	}
	fmt.Println(subscriptions)
	return subscriptions
}
