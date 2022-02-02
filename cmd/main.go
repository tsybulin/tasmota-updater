package main

import (
	"log"
	"os"
	"strconv"
	"time"

	tau "github.com/tsybulin/tasmota-updater"
)

var version string

func init() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: tasmota-updater version")
	}

	version = os.Args[1]

	_, err := os.Open("tasmota-minimal.bin.gz")
	if err != nil {
		log.Fatal(err)
	}

	_, err = os.Open("tasmota.bin.gz")
	if err != nil {
		log.Fatal(err)
	}

	if v, ok := os.LookupEnv("MQTT_HOST"); ok {
		tau.MQTT_HOST = v
	} else {
		log.Fatal("Environment var MQTT_HOST not found")
	}

	if v, ok := os.LookupEnv("MQTT_LOGIN"); ok {
		tau.MQTT_LOGIN = v
	} else {
		log.Fatal("Environment var MQTT_LOGIN not found")
	}

	if v, ok := os.LookupEnv("MQTT_PASSWORD"); ok {
		tau.MQTT_PASSWORD = v
	} else {
		log.Fatal("Environment var MQTT_PASSWORD not found")
	}

	if v, ok := os.LookupEnv("MQTT_TIMEOUT"); ok {
		i, err := strconv.Atoi(v)
		if err != nil {
			log.Fatal(err)
		}
		tau.MQTT_TIMEOUT_SEC = time.Second * time.Duration(i)
		log.Print("MQTT_TIMEOUT = ", i, " seconds")
	}

	if v, ok := os.LookupEnv("HTTP_TIMEOUT"); ok {
		i, err := strconv.Atoi(v)
		if err != nil {
			log.Fatal(err)
		}
		tau.HTTP_TIMEOUT_SEC = time.Second * time.Duration(i)
		log.Print("HTTP_TIMEOUT = ", i, " seconds")
	}
}

func main() {
	tau.Discover()
	tau.Update(version)
}
