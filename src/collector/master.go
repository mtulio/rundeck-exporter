package collector

import (
	"sync"
	"time"

	"github.com/apex/log"
	"github.com/mtulio/rundeck-exporter/src/rclient"
	"github.com/prometheus/client_golang/prometheus"
)

// Master implements the prometheus.Collector interface.
type Master struct {
	Collectors        map[string]Collector
	RClient           *rclient.RClient
	CollectorInterval int
}

// Collector is the interface a collector has to implement.
type Collector interface {
	// Get new metrics and expose them via prometheus registry.
	Update(ch chan<- prometheus.Metric) error
}

const (
	// Namespace defines the common namespace to be used by all metrics.
	namespace       = "rundeck"
	defaultEnabled  = true
	defaultDisabled = false
)

var (
	scrapeDurationDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "scrape", "collector_duration_seconds"),
		"rundeck_exporter: Duration of a collector scrape.",
		[]string{"collector"},
		nil,
	)
	scrapeSuccessDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "scrape", "collector_success"),
		"rundeck_exporter: Whether a collector succeeded.",
		[]string{"collector"},
		nil,
	)
)

// NewCollectorMaster creates a new NodeCollector.
func NewCollectorMaster(rcli *rclient.RClient, collectInterval int) (*Master, error) {
	var err error
	err = nil
	collectors := make(map[string]Collector)
	collectors["metrics"], err = NewCollectorMetrics(rcli, collectInterval)
	if err != nil {
		panic(err)
	}
	// collectors["sysinfo"], err = NewCollectorSysInfo(rcli, metrics...)
	// if err != nil {
	// 	panic(err)
	// }

	return &Master{
		Collectors: collectors,
		RClient:    rcli,
	}, nil
}

//Describe is a Prometheus implementation to be called by collector.
//It essentially writes all descriptors to the prometheus desc channel.
func (cm *Master) Describe(ch chan<- *prometheus.Desc) {

	//Update this section with the each metric you create for a given collector
	ch <- scrapeDurationDesc
	ch <- scrapeSuccessDesc
}

//Collect implements required collect function for all promehteus collectors
func (cm *Master) Collect(ch chan<- prometheus.Metric) {
	wg := sync.WaitGroup{}
	wg.Add(len(cm.Collectors))
	for name, c := range cm.Collectors {
		go func(name string, c Collector) {
			execute(name, c, ch)
			wg.Done()
		}(name, c)
	}
	wg.Wait()
}

// execute calls Update() function on subsystem to gather metrics
func execute(name string, c Collector, ch chan<- prometheus.Metric) {
	begin := time.Now()
	err := c.Update(ch)
	duration := time.Since(begin)
	var success float64

	if err != nil {
		log.Errorf("ERROR: %s collector failed after %fs: %s", name, duration.Seconds(), err)
		success = 0
	} else {
		log.Debugf("OK: %s collector succeeded after %fs.", name, duration.Seconds())
		success = 1
	}
	ch <- prometheus.MustNewConstMetric(scrapeDurationDesc, prometheus.GaugeValue, duration.Seconds(), name)
	ch <- prometheus.MustNewConstMetric(scrapeSuccessDesc, prometheus.GaugeValue, success, name)
}
