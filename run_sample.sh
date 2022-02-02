#!/bin/sh

export MQTT_HOST=192.168.10.10:1883
export MQTT_LOGIN=tastas
export MQTT_PASSWORD=paspas
export MQTT_TIMEOUT=3
export HTTP_TIMEOUT=30

./tasmota-updater 10.1.0
