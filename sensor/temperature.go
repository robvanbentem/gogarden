package sensor

import (
	"io/ioutil"
	"strings"
	"os"
	"log"
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
			log.Fatal(err)
		}

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			text := scanner.Text()
			p := strings.Index(text, "t=")
			if p != -1 {
				ts := text[p+2:]
				t, _ := strconv.ParseFloat(ts, 64)

				name := v.Name()[len(v.Name())-6:]
				temp := float64(t / 1000.0)

				readouts = append(readouts, TempReadout{name, temp, time.Now().Format(time.RFC3339)})
			}
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

		file.Close()
	}

	for _, temp := range readouts {
		msg, _ := json.Marshal(temp)
		*net.GetCommsChan() <- net.Message{"temp", msg}
	}
}
