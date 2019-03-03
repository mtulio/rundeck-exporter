package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/lusis/go-rundeck/pkg/rundeck"
)

var (
	cfg                = config{}
	defExpListenAddr   = ":9801"
	defExpMetricsPath  = "/metrics"
	defCollectInterval = 60
	defEnvAPIURL       = "RUNDECK_API_URL"
	defEnvAPIToken     = "RUNDECK_API_TOKEN"
	defEnvAPIUser      = "RUNDECK_PASS"
	defEnvAPIPass      = "RUNDECK_USER"
	defEnvAPIVersion   = "RUNDECK_API_VERSION"
	// apiClient          *rundeck.Client
	// err error
)

// usage returns the command line usage sample.
func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [options]\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(2)
}

func init() {

	cfg.expListenAddr = flag.String("web.listen-address", defExpListenAddr, "Address on which to expose metrics and web interface.")
	cfg.expMetricsPath = flag.String("web.telemetry-path", defExpMetricsPath, "Path under which to expose metrics.")

	cfg.apiURL = flag.String("rundeck.url", "", "API URL")
	cfg.apiUser = flag.String("rundeck.user", "", "API USER")
	cfg.apiPass = flag.String("rundeck.pass", "", "API_PASS")
	cfg.apiToken = flag.String("rundeck.token", "", "API token")
	cfg.apiVersion = flag.String("rundeck.version", "", "API version")

	cfg.collectorInterval = flag.Int("metrics.interval", defCollectInterval, "Interval in seconds to retrieve metrics from API")

	flag.Usage = usage
	flag.Parse()

	if *cfg.apiURL == "" {
		*cfg.apiURL = os.Getenv(defEnvAPIURL)
	}

	if *cfg.apiToken == "" {
		*cfg.apiToken = os.Getenv(defEnvAPIToken)
	}

	if *cfg.apiUser == "" {
		*cfg.apiUser = os.Getenv(defEnvAPIUser)
	}

	if *cfg.apiPass == "" {
		*cfg.apiPass = os.Getenv(defEnvAPIPass)
	}

	if *cfg.apiVersion == "" {
		v := os.Getenv(defEnvAPIVersion)
		if v == "" {
			v = "18"
		}
		*cfg.apiVersion = v
	}

	apiConfig := &rundeck.ClientConfig{
		BaseURL:      *cfg.apiURL,
		Token:        *cfg.apiToken,
		APIVersion:   *cfg.apiVersion,
		AuthMethod:   "basic",
		Username:     *cfg.apiUser,
		Password:     *cfg.apiPass,
		OverridePath: true,
	}
	fmt.Println(apiConfig)

	apiClient, _ := rundeck.NewClient(apiConfig)

	// p, e := apiClient.ListProjects()
	// fmt.Println(e)
	// fmt.Println(p)

	// i, e2 := apiClient.GetSystemInfo()
	// fmt.Println(e2)
	// fmt.Println(i.SystemInfoResponse.System.Metrics.ContentType)
	// fmt.Println(i.SystemInfoResponse.System.Metrics.HRef)
	// fmt.Println(i.SystemInfoResponse.System.Stats.Uptime.Duration)

	// type respData struct {
	// 	Version string `json:"version"`
	// }

	m, e3 := apiClient.GetMetrics()
	fmt.Println(e3)

	type dataInCount struct {
		Count int `json:"count"`
	}
	for k, v := range m.Counters {
		var d dataInCount
		metricName := strings.Replace(k, ".", "_", -1)
		if !strings.HasPrefix(metricName, "rundeck") {
			continue
		}
		// fmt.Println(k, v)

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

	type dataInGauges struct {
		Value float64 `json:"value"`
	}
	for k, v := range m.Gauges {
		var d dataInGauges
		metricName := strings.Replace(k, ".", "_", -1)
		if !strings.HasPrefix(metricName, "rundeck") {
			metricName = "rundeck_" + metricName
		}
		// fmt.Println(k, v)

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

	type dataInMeters struct {
		Count    int     `json:"count"`
		M15Rate  float64 `json:"m15_rate"`
		M1Rate   float64 `json:"m1_rate"`
		M5Rate   float64 `json:"m5_rate"`
		MeanRate float64 `json:"mean_rate"`
	}
	for k, v := range m.Meters {
		var d dataInMeters
		metricName := strings.Replace(k, ".", "_", -1)
		if !strings.HasPrefix(metricName, "rundeck") {
			metricName = "rundeck_servlet_" + metricName
		}
		// fmt.Println(k, v)

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

	type dataInTimers struct {
		Count    int     `json:"count"`
		Max      float64 `json:"max"`
		Mean     float64 `json:"mean"`
		P50      float64 `json:"p50"`
		P75      float64 `json:"p75"`
		P95      float64 `json:"p95"`
		P98      float64 `json:"p98"`
		P99      float64 `json:"p99"`
		P999     float64 `json:"p999"`
		Stddev   float64 `json:"stddev"`
		M15Rate  float64 `json:"m15_rate"`
		M1Rate   float64 `json:"m1_rate"`
		M5Rate   float64 `json:"m5_rate"`
		MeanRate float64 `json:"mean_rate"`
	}
	for k, v := range m.Timers {
		var d dataInTimers
		metricName := strings.Replace(k, ".", "_", -1)
		if !strings.HasPrefix(metricName, "rundeck") {
			metricName = "rundeck_" + metricName
		}
		// fmt.Println(k, v)

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

	// initPromCollector()
}
