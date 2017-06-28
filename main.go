package main

import (
	"gogarden/net"
	"gogarden/common"
	"gogarden/sensor"
	"time"
	"os"
	"gocmn"
)

func main() {
	common.LoadConfig()
	gocmn.InitLogger(common.ConfigRoot.LogFile)

	if err := net.Connect(); err != nil {
		gocmn.Log.Fatal("Could not connect to MQTT broker")
		os.Exit(1)
	}
	defer net.Disconnect()
	gocmn.Log.Info("Connected to MQTT broker")

	go net.ListenForMessages()
	go sensor.MonitorTemperatures()

	gocmn.Log.Info("Monitoring..")
	for {
		time.Sleep(time.Second)
	}
}
