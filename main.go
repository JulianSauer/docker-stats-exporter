package main

import (
	"docker-stats-exporter/metrics"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/metrics", metrics.MetricsHandler)
	fmt.Println("Listening on :9100")
	http.ListenAndServe(":9100", nil)
}
