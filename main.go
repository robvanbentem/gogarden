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
	defer net.Disconnect()
	common.Log.Info("Connected to MQTT broker")

	go net.ListenForMessages()
	go sensor.MonitorTemperatures()

	common.Log.Info("Monitoring..")
	for {
		time.Sleep(time.Second)
	}
}
