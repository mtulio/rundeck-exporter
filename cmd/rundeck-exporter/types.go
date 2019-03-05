package main

import (
	"github.com/mtulio/rundeck-exporter/src/collector"
	"github.com/mtulio/rundeck-exporter/src/rclient"
	"github.com/prometheus/client_golang/prometheus"
)

type configProm struct {
	Collector *collector.Master
	Registry  *prometheus.Registry
	Gatherers *prometheus.Gatherers
}

type config struct {
	expListenAddr     *string
	expMetricsPath    *string
	prom              *configProm
	rcli              *rclient.RClient
	collectorInterval *int
}

const (
	exporterName         = "rundeck_exporter"
	exporterDescription  = "Rundeck Exporter"
	defExpListenAddr     = ":9801"
	defExpMetricsPath    = "/metrics"
	defCollectorInterval = 60
	defEnvAPIURL         = "RUNDECK_API_URL"
	defEnvAPIToken       = "RUNDECK_API_TOKEN"
	defEnvAPIUser        = "RUNDECK_PASS"
	defEnvAPIPass        = "RUNDECK_USER"
	defEnvAPIVersion     = "RUNDECK_API_VERSION"
)
