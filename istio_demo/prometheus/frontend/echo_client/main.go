package main

import (
	"client/conf"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var config conf.Conf
var parserErr error
var protocol conf.Protocol

func init() {
	// var confFile string
	// flag.StringVar(&confFile, "c", os.Args[1], "config file")
	// flag.Parse()

	// config, parserErr = conf.ConfParser(confFile, &protocol)
	// if parserErr != nil {
	// 	log.Fatalf("parser config failed:", parserErr.Error())
	// }

	prometheus.MustRegister(uptimeGauge)
}

func recordMetrics() {
	go func() {
		for {
			uptimeGauge.Inc()
			time.Sleep(2 * time.Second)
		}
	}()
}

var (
	uptimeGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "echo_client",
		Name:      "uptime_seconds",
		Help:      "How long the system has been up",
	})
)

func main() {
	recordMetrics()
	start := time.Now().Second()
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":5000", nil)
	started := time.Now().Second()
	uptimeGauge.Set(float64(started - start))
}
