package net

import (
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"gogarden/common"
)

type Message struct {
	Path    string
	Message []byte
}

var client MQTT.Client

var comms chan Message
var exit chan byte

//define a function for the default message handler
var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func Connect() error {
	cfg := common.ConfigRoot.MQTT

	opts := MQTT.NewClientOptions().AddBroker(cfg.Broker)
	opts.SetClientID(cfg.Name)
	opts.SetDefaultPublishHandler(f)

	//create and start a client using the above ClientOptions
	client = MQTT.NewClient(opts)
	token := client.Connect()
	token.Wait()

	if token.Error() != nil {
		return token.Error()
	}

	comms = make(chan Message, 10)
	exit = make(chan byte)

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
			token := client.Publish(fmt.Sprintf(common.ConfigRoot.MQTT.Path, m.Path), common.ConfigRoot.MQTT.QOS, false, m.Message)
			token.Wait()
			if token.Error() != nil {
				common.Log.Error("Error publishing message: " + token.Error().Error())
			}
		case <-exit:
			break Loop
		}
	}
}

func GetCommsChan() *chan Message {
	return &comms
}
