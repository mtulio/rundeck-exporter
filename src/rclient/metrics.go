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
		return 0.0, fmt.Errorf("Error getting Metric Value Timer: ", ok)
	}
	return 0.0, nil
}

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
