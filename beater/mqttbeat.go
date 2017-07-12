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

	if token := bt.mqtt_client.SubscribeMultiple(bt.beat_config.Topics_subscribe, bt.on_message);
	token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	return bt, nil
}

func decodePayload(payload []byte) common.MapStr {
	event := make(common.MapStr)
	sent := false

	err := msgpack.Unmarshal(payload, &event)
	if  err == nil {
		err = json.Unmarshal(payload, &event)
		if err == nil {
			sent = true
			logp.Info("Payload decoded - msgpack + json")
			return event
		} else {
			sent = true
			logp.Info("Payload decoded - msgpack")
			return event
		}
	}

	err = json.Unmarshal(payload, &event)
	if  sent == false && err == nil {
		sent = true
		logp.Info("Payload decoded - json)")
		return event
	}

	if sent == false {
		sent = true
		event["payload"]= string(payload)
		logp.Info("Payload decoded - text")
		//TODO handle case bytes ... ? (string fail)
	}
	return event
}

//define a function for the default message handler
// see https://discuss.elastic.co/t/how-to-append-a-json-string-to-libbeat-event/34020/2
func (bt *Mqttbeat) on_message(client MQTT.Client, msg MQTT.Message) {
	logp.Info("MQTT MESSAGE RECEIVED " + string(msg.Payload()))
	event := make(common.MapStr)
	if bt.beat_config.Decode_paylod == true {
		event = decodePayload(msg.Payload())
	} else {
		event = make(common.MapStr)
		event["payload"] = msg.Payload()
	}

	event["beat"]= common.MapStr{"index": "mymqtt", "type":"message"}
	event["@timestamp"] = common.Time(time.Now())
	event["topic"] = msg.Topic()
	bt.elastic_client.PublishEvent(event)
	logp.Info("Event sent")
}


func (bt *Mqttbeat) Run(b *beat.Beat) error {
	logp.Info("mqttbeat is running! Hit CTRL-C to stop it.")
	bt.elastic_client = b.Publisher.Connect()

	for {
		select {
		case <-bt.done:
			return nil
		}
	}
}

func (bt *Mqttbeat) Stop() {
	bt.elastic_client.Close()
	close(bt.done)
}
