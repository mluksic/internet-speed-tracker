package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type InternetSpeed struct {
	Timestamp string  `json:"timestamp"`
	Download  float64 `json:"download"`
	Upload    float64 `json:"upload"`
	Ping      float64 `json:"ping"`
}

func main() {
	// Set up InfluxDB client
	client := influxdb2.NewClient("http://localhost:8086", "21341")
	writeAPI := client.WriteAPIBlocking("8680519cd55eb937", "internet_speed")

	// Run speed test every 5 minutes and store the results in InfluxDB
	for {
		cmd := exec.Command("echo", `{"download": 6214.59018058,"timestamp": "2023-02-17T00:04:13.135999Z","ping": 50.464,"upload": 528.04054276}`)
		resultsJson, err := cmd.Output()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Println(string(resultsJson))
		var internetSpeed InternetSpeed

		if err := json.Unmarshal([]byte(resultsJson), &internetSpeed); err != nil {
			log.Fatal(err)
			return
		}
		fmt.Println(internetSpeed.Timestamp, internetSpeed.Download, internetSpeed.Upload, internetSpeed.Ping)

		// Write the data to InfluxDB
		p := influxdb2.NewPoint("measurement",
			map[string]string{},
			map[string]interface{}{
				"ping":     internetSpeed.Ping,
				"download": internetSpeed.Download,
				"upload":   internetSpeed.Upload,
			},
			time.Now())
		writeAPI.WritePoint(context.Background(), p)

		time.Sleep(5 * time.Second)
	}
}
