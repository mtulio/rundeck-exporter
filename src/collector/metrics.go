package collector

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/apex/log"
	"github.com/mtulio/rundeck-exporter/src/rclient"
	"github.com/prometheus/client_golang/prometheus"
)

// RMetrics keeps the collector info
type RMetrics struct {
	Client  *rclient.RClient
	Metrics []Metric
}

// Metric describe the metric attributes
type Metric struct {
	Prom       *prometheus.Desc
	Name       string
	fCollector func(m *Metric) error
	Value      float64
}

// NewCollectorMetrics return the CollectorAnalytics object
func NewCollectorMetrics(rcli *rclient.RClient, msEnabled ...string) (*RMetrics, error) {

	ca := &RMetrics{
		Client: rcli,
	}
	err := ca.InitMetrics(msEnabled...)
	if err != nil {
		log.Info("collector.Metrics: error initializing metrics")
	}
	go ca.InitCollectorsUpdater()
	return ca, nil
}

// Update implements Collector and exposes related metrics
func (ca *RMetrics) Update(ch chan<- prometheus.Metric) error {
	// done := make(chan bool)
	wg := sync.WaitGroup{}
	wg.Add(len(ca.Metrics))

	for mID := range ca.Metrics {
		go func(m *Metric, ch chan<- prometheus.Metric) {
			ch <- prometheus.MustNewConstMetric(
				m.Prom,
				prometheus.GaugeValue,
				m.Value,
			)
			// done <- true
			wg.Done()
		}(&ca.Metrics[mID], ch)
	}

	// wait to finish all go routines
	wg.Wait()
	// <-done
	return nil
}

// InitMetrics initialize a list of metrics names and return error if fails.
func (ca *RMetrics) InitMetrics(msEnabled ...string) error {

	for k, v := range ca.Client.Metrics.Counters {
		var d rclient.DataInMetricCount
		m := Metric{}

		metricName := strings.Replace(k, ".", "_", -1)
		if !strings.HasPrefix(metricName, "rundeck") {
			continue
		}

		b, e := json.Marshal(v)
		if e != nil {
			fmt.Println(e)
		}

		e = json.Unmarshal(b, &d)
		if e != nil {
			fmt.Println(e)
		}

		m.Name = prometheus.BuildFQName(namespace, "counter", metricName+"_total")
		fmt.Println(metricName+"_total : ", d.Count)

		m.Prom = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "counter", metricName+"_total"),
			"Teste",
			nil, nil,
		)
		ca.Metrics = append(ca.Metrics, m)
	}
	return nil
}

// InitCollectorsUpdater start the paralel auto update for each collector
func (ca *RMetrics) InitCollectorsUpdater() {
	for {
		for mID := range ca.Metrics {
			go func(m *Metric) {
				m.fCollector(m)
			}(&ca.Metrics[mID])
		}
		time.Sleep(time.Second * time.Duration(60))
	}
}

func (ca *RMetrics) collectorCounters(metric, dimension string) func(m *Metric) error {
	return func(m *Metric) error {

		for k, v := range ca.Client.Metrics.Counters {
			var d rclient.DataInMetricCount

			metricName := strings.Replace(k, ".", "_", -1)
			if !strings.HasPrefix(metricName, "rundeck") {
				continue
			}

			b, e := json.Marshal(v)
			if e != nil {
				fmt.Println(e)
			}

			e = json.Unmarshal(b, &d)
			if e != nil {
				fmt.Println(e)
			}

			m.Name = prometheus.BuildFQName(namespace, "counter", metricName+"_total")
			fmt.Println(metricName+"_total : ", d.Count)
			m.Value = float64(d.Count)

		}
		return nil
	}
}
