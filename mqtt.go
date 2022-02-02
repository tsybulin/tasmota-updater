package tasmotaupdater

import (
	"encoding/json"
	"log"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const MQTT_CLIENT = "tasmota-updater"

var (
	MQTT_HOST        = "1.1.1.1:1883"
	MQTT_LOGIN       = "xxx"
	MQTT_PASSWORD    = "123456"
	MQTT_TIMEOUT_SEC = 3 * time.Second
)

func discoveryHandler(client mqtt.Client, message mqtt.Message) {
	tasmota := tasmota{}
	if err := json.Unmarshal(message.Payload(), &tasmota); err != nil {
		log.Print("Mqtt Unmarshall error: ", err)
	} else {
		tasmotas[tasmota.Ip] = tasmota
	}
}

func Discover() {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "123"
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(MQTT_HOST)
	opts.SetClientID(MQTT_CLIENT + "-" + hostname)
	opts.SetUsername(MQTT_LOGIN)
	opts.SetPassword(MQTT_PASSWORD)
	opts.SetKeepAlive(2 * time.Second)
	opts.SetPingTimeout(1 * time.Second)

	opts.SetDefaultPublishHandler(func(client mqtt.Client, message mqtt.Message) {
		log.Print("Mqtt.defaultHandler ", message.Topic(), " ", message.Payload())
	})

	opts.OnConnect = func(client mqtt.Client) {
		token := client.Subscribe("tasmota/discovery/+/config", 0, discoveryHandler)
		token.Wait()

	}

	mqttClient := mqtt.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		log.Print("GoBest.Mqtt Connect error: ", token.Error())
	}

	time.Sleep(MQTT_TIMEOUT_SEC)

	mqttClient.Disconnect(1000)
}
