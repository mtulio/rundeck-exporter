package rclient

import (
	"encoding/json"
	"fmt"
)

// MetricType is the type of metric.
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

// GetMetricValueCounter return the value for the Counter type.
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
		return 0.0, fmt.Errorf("Error getting Metric Value Counter")
	}

}

// GetMetricValueGauge return the value for the Gauge type.
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
		return 0.0, fmt.Errorf("Error getting Metric Value Gauges")
	}
}

// GetMetricValueMeter return the value for the Meter type.
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
		case "Count":
			return float64(d.Count), nil
		case "M15Rate":
			return float64(d.M15Rate), nil
		case "M1Rate":
			return float64(d.M1Rate), nil
		case "M5Rate":
			return float64(d.M5Rate), nil
		case "MeanRate":
			return float64(d.MeanRate), nil
		default:
			return 0.0, nil
		}
	} else {
		return 0.0, fmt.Errorf("Error getting Metric Value Meters")
	}
}

// GetMetricValueTimer return the value for the Timer type.
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
		case "Count":
			return float64(d.Count), nil
		case "Max":
			return float64(d.Max), nil
		case "Mean":
			return float64(d.Mean), nil
		case "Min":
			return float64(d.Min), nil
		case "P50":
			return float64(d.P50), nil
		case "P75":
			return float64(d.P75), nil
		case "P95":
			return float64(d.P95), nil
		case "P98":
			return float64(d.P98), nil
		case "P99":
			return float64(d.P99), nil
		case "P999":
			return float64(d.P999), nil
		case "Stddev":
			return float64(d.Stddev), nil
		case "M15Rate":
			return float64(d.M15Rate), nil
		case "M1Rate":
			return float64(d.M1Rate), nil
		case "M5Rate":
			return float64(d.M5Rate), nil
		case "MeanRate":
			return float64(d.MeanRate), nil
		default:
			return 0.0, nil
		}
	} else {
		return 0.0, fmt.Errorf("Error getting Metric Value Timer.")
	}
}

// GetDimensions return the dimensions available on each of Types.
func (rc *RClient) GetDimensions(mtype string) []string {

	var ds []string

	switch mtype {
	case "Meter":
		ds = []string{"Count", "M15Rate", "M1Rate", "M5Rate", "MeanRate"}
	case "Timers":
		ds = []string{"Count", "Max", "Mean", "Min", "P50", "P75", "P95", "P98",
			"P99", "P999", "Stddev", "M15Rate", "M1Rate", "M5Rate", "MeanRate"}
	default:
		return ds
	}

	return ds
}
