package beater

import (
	"github.com/elastic/beats/libbeat/common"
	"gopkg.in/vmihailenco/msgpack.v2"
	"reflect"
	"testing"
)

func TestDecodeMsgpackJson(t *testing.T) {

	reference := make(common.MapStr)
	reference["hello"] = "world"
	reference["answer"] = 42.0 // floats, not int, because of json unmarshal

	input, _ := msgpack.Marshal(&reference)
	output := DecodePayload(input)

	if !reflect.DeepEqual(reference, output) {
		t.Error("Not equals")
	}
}

func TestDecodeJson(t *testing.T) {

	reference := make(common.MapStr)
	reference["hello"] = "world"
	reference["answer"] = 42.0

	output := DecodePayload([]byte(`{"hello":"world", "answer":42}`))

	if !reflect.DeepEqual(reference, output) {
		t.Error("Not equals")
	}
}

func TestDecodeText(t *testing.T) {
	payload := "Bonjour, monde!"

	reference := make(common.MapStr)
	reference["payload"] = payload

	output := DecodePayload([]byte(payload))

	if !reflect.DeepEqual(reference, output) {
		t.Error("Not equals")
	}
}

func TestParseTopic(t *testing.T) {
	input := []string{"some/topic?0", "some/ohter/topic?2", "final/topic?1"}

	reference := map[string]byte{"some/topic": 0, "some/ohter/topic": 2, "final/topic": 1}

	output := ParseTopics(input)

	if !reflect.DeepEqual(reference, output) {
		t.Error("Not equals")
	}
}
