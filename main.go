package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"strconv"
	"io/ioutil"
	"time"
)

func main() {
	dir := "/sys/bus/w1/devices/"
	files, _ := ioutil.ReadDir(dir)

	for ; ; {
		for _, v := range files {

			if (strings.Contains(v.Name(), "w1_bus_master")) {
				continue
			}

			file, err := os.Open(dir + v.Name() + "/w1_slave")
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				text := scanner.Text()
				p := strings.Index(text, "t=")
				if (p != -1) {
					ts := text[p+2:]
					t, _ := strconv.ParseFloat(ts, 64)
					f := float64(t / 1000.0)
					fmt.Printf("Temp: %s %.3f\n", v.Name()[len(v.Name())-6:], f)
				}
			}

			if err := scanner.Err(); err != nil {
				log.Fatal(err)
			}
		}

		time.Sleep(time.Minute)
	}
}
