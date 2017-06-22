package common

import (
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"time"
)

type config struct {
	MQTT configMQTT

	DevicePath      string
	MonitorInterval duration
	LogFile         string
}

type configMQTT struct {
	Broker string
	Name   string
	Path   string
	QOS    byte
}

type duration struct {
	time.Duration
}

func (d *duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}

var ConfigRoot *config

func LoadConfig() {
	buf, err := ioutil.ReadFile("Config.toml")
	if err != nil {
		panic("Could not read Config.toml")
	}

	var conf config
	if _, err := toml.Decode(string(buf), &conf); err != nil {
		panic("Could not parse Config.toml: " + err.Error())
	}

	ConfigRoot = &conf
}
