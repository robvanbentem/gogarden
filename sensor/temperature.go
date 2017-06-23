package sensor

import (
	"io/ioutil"
	"strings"
	"os"
	"bufio"
	"strconv"
	"gogarden/net"
	"time"
	"encoding/json"
	"gogarden/common"
)

type TempReadout struct {
	DeviceID    string
	Temperature float64
	Date        string
}

func MonitorTemperatures() {
	go reportTemperatures()
	ticker := time.Tick(common.ConfigRoot.MonitorInterval.Duration)

	for {
		select {
		case <-ticker:
			go reportTemperatures()
		}
	}
}

func reportTemperatures() {
	dir := common.ConfigRoot.DevicePath + "/"
	files, _ := ioutil.ReadDir(dir)

	readouts := make([]TempReadout, 0)

	for _, v := range files {
		if strings.Contains(v.Name(), "w1_bus_master") {
			continue
		}

		file, err := os.Open(dir + v.Name() + "/w1_slave")
		if err != nil {
			common.Log.Error("Error reading device: " + err.Error())
			continue
		}

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			text := scanner.Text()
			p := strings.Index(text, "t=")
			if p != -1 {
				ts := text[p+2:]
				t, err := strconv.ParseFloat(ts, 64)
				if err != nil {
					common.Log.Error("Could not parse temperature: " + err.Error())
				}

				name := v.Name()[len(v.Name())-6:]
				temp := float64(t / 1000.0)

				readouts = append(readouts, TempReadout{name, temp, time.Now().Format(time.RFC3339)})
			}
		}

		if err := scanner.Err(); err != nil {
			common.Log.Error("Error while reading device: " + err.Error())
		}

		file.Close()
	}

	for _, temp := range readouts {
		common.Log.Infof("Reporting temperature for %s: %.2f", temp.DeviceID, temp.Temperature)
		msg, _ := json.Marshal(temp)
		*net.GetCommsChan() <- net.Message{"temp", msg}
	}
}
