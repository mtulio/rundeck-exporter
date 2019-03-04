package rclient

import (
	"encoding/json"
	"fmt"
	"strings"
)

type MetricType int

const (
	// MTypeCounter is an metric type Counter
	MTypeCounter MetricType = iota
	// MTypeGauge is an metric type Gauge
	MTypeGauge
	// MTypeHistogram is an metric type Histogram
	MTypeHistogram
	// MTypeMeter is an metric type Meter
	MTypeMeter
	// MTypeTimer is an metric type Timer
	MTypeTimer
)

// Metric is an representation of one metric
type Metric struct {
	Name  string
	Value float64
	// Type  metricType
}

// Metrics is an slice of Metric
type Metrics map[string]*Metric

// UpdateMetrics retrieve metrics from Rundeck and make it available
func (rc *RClient) UpdateMetrics() error {
	if rc.SOAP == nil {
		err := fmt.Errorf("Client SOA is not initializated")
		return err
	}

	m, err := rc.SOAP.GetMetrics()
	if err != nil {
		panic(err)
	}
	rc.Metrics = m
	return nil
}

// ShowMetrics parse metrics from Http WEB app show it out.
func (rc *RClient) ShowMetrics() error {

	if rc.SOAP == nil {
		err := fmt.Errorf("Client SOA is not initializated")
		return err
	}

	if rc.Metrics == nil {
		err := fmt.Errorf("Metrics was not initializaded. Use UpdateMetrics()")
		return err
	}

	for k, v := range rc.Metrics.Counters {
		var d DataInMetricCount
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

		fmt.Println(metricName+"_total : ", d.Count)
	}

	for k, v := range rc.Metrics.Gauges {
		var d DataInMetricGauges
		metricName := strings.Replace(k, ".", "_", -1)
		if !strings.HasPrefix(metricName, "rundeck") {
			metricName = "rundeck_" + metricName
		}

		b, e := json.Marshal(v)
		if e != nil {
			fmt.Println(e)
		}

		e = json.Unmarshal(b, &d)
		if e != nil {
			fmt.Println(e)
		}

		fmt.Println(metricName+"_total : ", d.Value)
	}

	for k, v := range rc.Metrics.Meters {
		var d DataInMetricMeters
		metricName := strings.Replace(k, ".", "_", -1)
		if !strings.HasPrefix(metricName, "rundeck") {
			metricName = "rundeck_servlet_" + metricName
		}

		b, e := json.Marshal(v)
		if e != nil {
			fmt.Println(e)
		}

		e = json.Unmarshal(b, &d)
		if e != nil {
			fmt.Println(e)
		}

		fmt.Println(metricName+"_total : ", d.Count)
		fmt.Println(metricName+"_rate : ", d.M1Rate, d.M5Rate, d.M15Rate, d.MeanRate)
	}

	for k, v := range rc.Metrics.Timers {
		var d DataInMetricTimers
		metricName := strings.Replace(k, ".", "_", -1)
		if !strings.HasPrefix(metricName, "rundeck") {
			metricName = "rundeck_" + metricName
		}

		b, e := json.Marshal(v)
		if e != nil {
			fmt.Println(e)
		}

		e = json.Unmarshal(b, &d)
		if e != nil {
			fmt.Println(e)
		}

		fmt.Println(metricName+"_total : ", d.Count)
		fmt.Println(metricName+"_rate : ", d.M1Rate, d.M5Rate, d.M15Rate, d.MeanRate)
	}

	return nil
}

func (rc *RClient) GetMetricValueCounter(metricName string) (float64, error) {
	var d DataInMetricCount

	if v, ok := rc.Metrics.Counters[metricName]; ok {
		b, e := json.Marshal(v)
		if e != nil {
			fmt.Println(e)
		}

		e = json.Unmarshal(b, &d)
		if e != nil {
			fmt.Println(e)
		}
		return float64(d.Count), nil
	} else {
		e := fmt.Errorf("Error getting Metric Value Counter: %v", ok)
		return 0.0, e
	}

}

func (rc *RClient) GetMetricValueGauge(metricName string) (float64, error) {
	var d DataInMetricGauges

	if v, ok := rc.Metrics.Gauges[metricName]; ok {
		b, e := json.Marshal(v)
		if e != nil {
			fmt.Println(e)
		}

		e = json.Unmarshal(b, &d)
		if e != nil {
			fmt.Println(e)
		}

		return float64(d.Value), nil
	} else {
		return 0.0, fmt.Errorf("Error getting Metric Value Gauges: ", ok)
	}

}

func (rc *RClient) GetMetricValueMeter(metricName, dimension string) (float64, error) {
	var d DataInMetricMeters

	if v, ok := rc.Metrics.Meters[metricName]; ok {
		b, e := json.Marshal(v)
		if e != nil {
			fmt.Println(e)
		}

		e = json.Unmarshal(b, &d)
		if e != nil {
			fmt.Println(e)
		}
		switch dimension {
		case "count":
			return float64(d.Count), nil
		default:
			return 0.0, nil
		}
	} else {
		return 0.0, fmt.Errorf("Error getting Metric Value Meters: ", ok)
	}
	return 0.0, nil
}

func (rc *RClient) GetMetricValueTimer(metricName, dimension string) (float64, error) {
	var d DataInMetricTimers

	if v, ok := rc.Metrics.Timers[metricName]; ok {
		b, e := json.Marshal(v)
		if e != nil {
			fmt.Println(e)
		}

		e = json.Unmarshal(b, &d)
		if e != nil {
			fmt.Println(e)
		}
		switch dimension {
		case "count":
			return float64(d.Count), nil
		default:
			return 0.0, nil
		}
	} else {
		return 0.0, fmt.Errorf("Error getting Metric Value Timer: ", ok)
	}
	return 0.0, nil
}
