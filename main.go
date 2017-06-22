package main

import (
	"gogarden/net"
	"gogarden/common"
	"gogarden/sensor"
	"time"
	"os"
)

func main() {
	common.LoadConfig()
	common.InitLogger()

	if err := net.Connect(); err != nil {
		common.Log.Fatal("Could not connect to MQTT broker")
		os.Exit(1)
	}
	common.Log.Debug("Connected to MQTT broker")
	defer net.Disconnect()

	go net.ListenForMessages()
	go sensor.MonitorTemperatures()

	common.Log.Debug("Monitoring..")
	for {
		time.Sleep(time.Second)
	}
}
