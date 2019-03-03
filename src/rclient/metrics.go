package rclient

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ShowMetrics parse metrics from Http WEB app show it out.
func (rc *RClient) ShowMetrics() error {

	fmt.Println(rc.SOAP)
	if rc.SOAP == nil {
		fmt.Println("SOA is not defined")
		return nil
	}

	m, err := rc.SOAP.GetMetrics()
	if err != nil {
		panic(err)
	}

	for k, v := range m.Counters {
		var d dataInCount
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

	for k, v := range m.Gauges {
		var d dataInGauges
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

	for k, v := range m.Meters {
		var d dataInMeters
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

	for k, v := range m.Timers {
		var d dataInTimers
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
