package net

import (
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"gogarden/common"
	"gocmn"
	"time"
	"errors"
)

type Message struct {
	Path    string
	Message []byte
}

var client MQTT.Client

var comms chan Message
var exit chan byte

var connecting = false

//define a function for the default message handler
var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
}

func Setup() {
	comms = make(chan Message, 32)
	exit = make(chan byte)
}

func Connect() error {

	if connecting {
		for !connecting {
			gocmn.Log.Debug("Already connecting, waiting..")
			time.Sleep(1 * time.Second)
		}

		if client.IsConnected() {
			return nil
		} else {
			return errors.New("MQTT client is not connected")
		}
	}

	connecting = true

	cfg := common.ConfigRoot.MQTT

	opts := MQTT.NewClientOptions().AddBroker(cfg.Broker)
	opts.SetClientID(cfg.Name)
	opts.SetDefaultPublishHandler(f)

	//create and start a client using the above ClientOptions
	client = MQTT.NewClient(opts)
	token := client.Connect()
	token.Wait()

	connecting = false

	if token.Error() != nil {
		return token.Error()
	}

	return nil
}

func Disconnect() {
	exit <- 0
	if client.IsConnected() {
		client.Disconnect(250)
	}
}

func ListenForMessages() {
Loop:
	for {
		select {
		case m := <-comms:
			if client.IsConnected() == false {
				gocmn.Log.Warning("MQTT disconnected, trying to reconnect..")
				if err := Connect(); err != nil {
					gocmn.Log.Error("Client not connected, cannot publish message")
					return
				}
			}
			go publishMessage(m)
		case <-exit:
			gocmn.Log.Info("Stop listening for messages")
			break Loop
		}
	}
}

func GetCommsChan() *chan Message {
	return &comms
}

func publishMessage(m Message) {
	gocmn.Log.Debug("Publising message..")
	token := client.Publish(fmt.Sprintf(common.ConfigRoot.MQTT.Path, m.Path), common.ConfigRoot.MQTT.QOS, false, m.Message)
	token.Wait()
	gocmn.Log.Debug("Message published")
	if token.Error() != nil {
		gocmn.Log.Error("Error publishing message: " + token.Error().Error())
	}
}
