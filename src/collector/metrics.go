package collector

import (
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
	Client            *rclient.RClient
	Metrics           []*Metric
	CollectorInterval int
}

// Metric describe the metric attributes
type Metric struct {
	Prom        *prometheus.Desc
	Name        string
	NameRaw     string
	NameDim     string
	Value       float64
	Type        rclient.MetricType
	Labels      []string
	LabelsValue []string
}

var (
	metricCollectorName     = "collectorRundeck"
	metricCollectorDuration time.Duration
	metricCollectorSuccess  float64
)

// NewCollectorMetrics return the CollectorAnalytics object
func NewCollectorMetrics(rcli *rclient.RClient, colInterval int, msEnabled ...string) (*RMetrics, error) {

	ca := &RMetrics{
		Client:            rcli,
		CollectorInterval: colInterval,
	}
	err := ca.InitMetrics(msEnabled...)
	if err != nil {
		log.Info("collector.Metrics: error initializing metrics")
	}
	go ca.InitCollectorsUpdater()
	return ca, nil
}

// InitMetrics initialize a list of metrics names and return error if fails.
func (ca *RMetrics) InitMetrics(msEnabled ...string) error {

	if ca.Client.Metrics == nil {
		err := ca.Client.UpdateMetrics()
		if err != nil {
			em := fmt.Errorf("Error initializing metrics from the server: %s", err)
			panic(em)
		}
	}

	// Create Counter metrics
	for k := range ca.Client.Metrics.Counters {
		m := Metric{}
		mName := strings.Replace(k, ".", "_", -1)
		mName = strings.Replace(mName, "-", "_", -1)
		m.NameRaw = k
		m.Type = rclient.MTypeCounter

		if !strings.HasPrefix(mName, "rundeck") {
			m.Name = "rundeck_" + mName
		} else {
			m.Name = mName
		}

		m.Prom = prometheus.NewDesc(
			m.Name,
			"Rundeck metrics Counter",
			nil, nil,
		)
		ca.Metrics = append(ca.Metrics, &m)
	}
	for k := range ca.Client.Metrics.Gauges {
		m := Metric{}
		mName := strings.Replace(k, ".", "_", -1)
		mName = strings.Replace(mName, "-", "_", -1)
		m.NameRaw = k
		m.Type = rclient.MTypeGauge

		if !strings.HasPrefix(mName, "rundeck") {
			m.Name = "rundeck_" + mName
		} else {
			m.Name = mName
		}

		m.Prom = prometheus.NewDesc(
			m.Name,
			"Rundeck metrics Gauge",
			nil, nil,
		)
		ca.Metrics = append(ca.Metrics, &m)
	}
	for k := range ca.Client.Metrics.Meters {
		for _, d := range ca.Client.GetDimensions("Meter") {
			m := Metric{}

			mName := strings.Replace(k, ".", "_", -1)
			mName = strings.Replace(mName, "-", "_", -1)
			m.NameRaw = k
			m.NameDim = d
			m.Type = rclient.MTypeMeter
			m.Labels = []string{"type"}

			if !strings.HasPrefix(mName, "rundeck") {
				m.Name = "rundeck_" + mName
			} else {
				m.Name = mName
			}

			m.Prom = prometheus.NewDesc(
				m.Name,
				"Rundeck metrics Meter",
				m.Labels, nil,
			)
			ca.Metrics = append(ca.Metrics, &m)
		}
	}
	for k := range ca.Client.Metrics.Timers {
		for _, d := range ca.Client.GetDimensions("Timers") {
			m := Metric{}

			mName := strings.Replace(k, ".", "_", -1)
			mName = strings.Replace(mName, "-", "_", -1)
			m.NameRaw = k
			m.NameDim = d
			m.Type = rclient.MTypeTimer
			m.Labels = []string{"type"}

			if !strings.HasPrefix(mName, "rundeck") {
				m.Name = "rundeck_" + mName
			} else {
				m.Name = mName
			}

			m.Prom = prometheus.NewDesc(
				m.Name,
				"Rundeck metrics Timer",
				m.Labels, nil,
			)
			ca.Metrics = append(ca.Metrics, &m)
		}
	}
	return nil
}

// InitCollectorsUpdater start the paralel auto update for each collector
func (ca *RMetrics) InitCollectorsUpdater() {
	for {
		mBegin := time.Now()
		if err := ca.Client.UpdateMetrics(); err != nil {
			fmt.Println("Unable to update Metrics: ", err)
			metricCollectorSuccess = 0.0
		} else {
			metricCollectorSuccess = 1.0
		}
		metricCollectorDuration = time.Since(mBegin)

		ca.collectorUpdate()
		time.Sleep(time.Second * time.Duration(ca.CollectorInterval))
	}
}

func (ca *RMetrics) collectorUpdate() error {

	for _, m := range ca.Metrics {
		var err error
		if m.Type == rclient.MTypeCounter {
			m.Value, err = ca.Client.GetMetricValueCounter(m.NameRaw)
			if err != nil {
				fmt.Println("Error getting Counter metric value: ", err)
				continue
			}
		} else if m.Type == rclient.MTypeGauge {
			m.Value, err = ca.Client.GetMetricValueGauge(m.NameRaw)
			if err != nil {
				fmt.Println("Error getting Gauge metric value: ", err)
				continue
			}
		} else if m.Type == rclient.MTypeMeter {
			m.Value, err = ca.Client.GetMetricValueMeter(m.NameRaw, m.NameDim)
			if err != nil {
				fmt.Println("Error getting Meter metric value: ", err)
				continue
			}
			m.LabelsValue = []string{m.NameDim}
		} else if m.Type == rclient.MTypeTimer {
			m.Value, err = ca.Client.GetMetricValueTimer(m.NameRaw, m.NameDim)
			if err != nil {
				fmt.Println("Error getting Timer metric value: ", err)
				continue
			}
			m.LabelsValue = []string{m.NameDim}
		} else {
			fmt.Println("#>> Type not found")
		}
	}
	return nil
}

// Update implements Collector and exposes related metrics
func (ca *RMetrics) Update(ch chan<- prometheus.Metric) error {

	wg := sync.WaitGroup{}
	wg.Add(len(ca.Metrics))

	for mID := range ca.Metrics {
		go func(m *Metric, ch chan<- prometheus.Metric) {
			if m.Labels == nil {
				ch <- prometheus.MustNewConstMetric(
					m.Prom,
					prometheus.GaugeValue,
					m.Value,
				)
			} else {
				ch <- prometheus.MustNewConstMetric(
					m.Prom,
					prometheus.GaugeValue,
					m.Value,
					m.LabelsValue...,
				)
			}
			wg.Done()
		}(ca.Metrics[mID], ch)
		// ch <- prometheus.MustNewConstMetric(scrapeDurationDesc, prometheus.GaugeValue, metricCollectorDuration.Seconds(), metricCollectorName)
		// ch <- prometheus.MustNewConstMetric(scrapeSuccessDesc, prometheus.GaugeValue, 1, "rundeckMetrics")
	}

	wg.Wait()
	return nil
}
