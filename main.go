package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type InternetSpeed struct {
	Timestamp string  `json:"timestamp"`
	Download  float64 `json:"download"`
	Upload    float64 `json:"upload"`
	Ping      float64 `json:"ping"`
}

func recordInternetSpeed() {
	cmd := exec.Command("speedtest-cli", `--json`)
	resultsJson, err := cmd.Output()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(string(resultsJson))

	if err := json.Unmarshal([]byte(resultsJson), &internetSpeed); err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(internetSpeed.Timestamp, internetSpeed.Download, internetSpeed.Upload, internetSpeed.Ping)
	downloadGauge.Set(internetSpeed.Download)
	uploadGauge.Set(internetSpeed.Upload)
	pingGauge.Set(internetSpeed.Ping)
}

var (
	downloadGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "speedtest",
			Name:      "download",
			Help:      "Download speed gauge, measured in MB/s.",
		})
	uploadGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "speedtest",
			Name:      "upload",
			Help:      "Upload speed gauge, measured in MB/s.",
		})
	pingGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "speedtest",
			Name:      "latency",
			Help:      "Ping gauge, measured in ms.",
		})
	internetSpeed InternetSpeed
)

func main() {
	// // Set up InfluxDB client
	// client := influxdb2.NewClient("http://localhost:8086", "21341")
	// writeAPI := client.WriteAPIBlocking("8680519cd55eb937", "internet_speed")

	// // Run speed test every 5 minutes and store the results in InfluxDB
	// for {
	// 	cmd := exec.Command("echo", `{"download": 6214.59018058,"timestamp": "2023-02-17T00:04:13.135999Z","ping": 50.464,"upload": 528.04054276}`)
	// 	resultsJson, err := cmd.Output()
	// 	if err != nil {
	// 		fmt.Println("Error:", err)
	// 		return
	// 	}

	// 	fmt.Println(string(resultsJson))
	// 	var internetSpeed InternetSpeed

	// 	if err := json.Unmarshal([]byte(resultsJson), &internetSpeed); err != nil {
	// 		log.Fatal(err)
	// 		return
	// 	}
	// 	fmt.Println(internetSpeed.Timestamp, internetSpeed.Download, internetSpeed.Upload, internetSpeed.Ping)

	// 	// Write the data to InfluxDB
	// 	p := influxdb2.NewPoint("measurement",
	// 		map[string]string{},
	// 		map[string]interface{}{
	// 			"ping":     internetSpeed.Ping,
	// 			"download": internetSpeed.Download,
	// 			"upload":   internetSpeed.Upload,
	// 		},
	// 		time.Now())
	// 	writeAPI.WritePoint(context.Background(), p)

	// 	time.Sleep(5 * time.Second)
	// }

	prometheus.MustRegister(downloadGauge)
	prometheus.MustRegister(uploadGauge)
	prometheus.MustRegister(pingGauge)

	go func() {
		for {
			recordInternetSpeed()
			time.Sleep(60 * time.Second)
		}
	}()
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":9000", nil)
}
