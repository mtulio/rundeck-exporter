package main

import (
	"github.com/mtulio/azion-exporter/src/collector"
	"github.com/prometheus/client_golang/prometheus"
)

type config struct {
	Collector         *collector.CollectorMaster
	Registry          *prometheus.Registry
	Gatherers         *prometheus.Gatherers
	rundeckToken      *string
	expListenAddr     *string
	expMetricsPath    *string
	collectorInterval *int
	apiToken          *string
	apiURL            *string
	apiUser           *string
	apiPass           *string
	apiVersion        *string
}

const (
	exporterName        = "rundeck_exporter"
	exporterDescription = "Rundeck Exporter"
)
