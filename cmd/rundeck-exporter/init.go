package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mtulio/rundeck-exporter/src/rclient"
)

var (
	cfg = config{}
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

	if cfg.apiPass == nil {
		*cfg.apiPass = os.Getenv(defEnvAPIPass)
	}

	if *cfg.apiVersion == "" {
		v := os.Getenv(defEnvAPIVersion)
		if v == "" {
			v = "18"
		}
		*cfg.apiVersion = v
	}

	if (*cfg.apiUser == "" || *cfg.apiPass == "") && (*cfg.apiToken == "") {
		panic("Error, auth Token, or User && Password was not provided")
	}

	rconf, ecfg := rclient.NewConfig()
	if ecfg != nil {
		emsg := fmt.Errorf("Unable to create the client")
		panic(emsg)
	}
	rconf.Base.BaseURL = *cfg.apiURL
	rconf.Base.APIVersion = *cfg.apiVersion

	if *cfg.apiUser != "" && *cfg.apiPass != "" {
		rconf.EnableHTTP = true
		rconf.Base.AuthMethod = "basic"
		rconf.Base.Username = *cfg.apiUser
		rconf.Base.Password = *cfg.apiPass
	}

	if *cfg.apiToken != "" {
		rconf.EnableAPI = true
		rconf.Base.Token = *cfg.apiToken
	}

	// Init clint
	rcli, err := rclient.NewClient(rconf)
	if err != nil {
		panic(err)
	}

	if err := rcli.ListProjects(); err != nil {
		fmt.Println("Unable to list projects")
	}

	if err := rcli.ShowMetrics(); err != nil {
		fmt.Println("Unable to show Metrics")
	}

	// initPromCollector()
}
