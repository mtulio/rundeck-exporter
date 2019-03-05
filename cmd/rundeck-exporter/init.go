package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mtulio/rundeck-exporter/src/collector"
	"github.com/mtulio/rundeck-exporter/src/rclient"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/log"
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

	apiURL := flag.String("rundeck.url", "", "API URL")
	apiUser := flag.String("rundeck.user", "", "API USER")
	apiPass := flag.String("rundeck.pass", "", "API_PASS")
	apiToken := flag.String("rundeck.token", "", "API token")
	apiVersion := flag.String("rundeck.version", "", "API version")

	cfg.collectorInterval = flag.Int("metrics.interval", defCollectorInterval, "Interval in seconds to retrieve metrics from API")

	flag.Usage = usage
	flag.Parse()

	if *apiURL == "" {
		*apiURL = os.Getenv(defEnvAPIURL)
	}

	if *apiToken == "" {
		*apiToken = os.Getenv(defEnvAPIToken)
	}

	if *apiUser == "" {
		*apiUser = os.Getenv(defEnvAPIUser)
	}

	if apiPass == nil {
		*apiPass = os.Getenv(defEnvAPIPass)
	}

	if *apiVersion == "" {
		v := os.Getenv(defEnvAPIVersion)
		if v == "" {
			v = "18"
		}
		*apiVersion = v
	}

	if cfg.collectorInterval == nil {
		*cfg.collectorInterval = defCollectorInterval
	}

	if (*apiUser == "" || *apiPass == "") || (*apiToken == "") {
		emsg := fmt.Errorf("#ERR> unable to find credentials, User and Passord, or Toekn")
		fmt.Println(emsg)
		os.Exit(1)
	}

	rconf, ecfg := rclient.NewConfig()
	if ecfg != nil {
		emsg := fmt.Errorf("Unable to create the client")
		fmt.Println(emsg)
		os.Exit(1)
	}
	rconf.Base.BaseURL = *apiURL
	rconf.Base.APIVersion = *apiVersion

	if *apiUser == "" || *apiPass == "" {
		emsg := fmt.Errorf("unable to create the client. Missing User and Passwords")
		fmt.Println(emsg)
		os.Exit(1)
	} else {
		rconf.EnableHTTP = true
		rconf.Base.AuthMethod = "basic"
		rconf.Base.Username = *apiUser
		rconf.Base.Password = *apiPass
	}

	if *apiToken != "" {
		rconf.EnableAPI = true
		rconf.Base.Token = *apiToken
	}

	// Init clint
	rcli, err := rclient.NewClient(rconf)
	if err != nil {
		panic(err)
	}
	cfg.rcli = rcli

	// Sample to show the metrics on the initialization:
	// if err := rcli.UpdateMetrics(); err != nil {
	// 	fmt.Println("Unable to update Metrics: ", err)
	// }
	// if err := rcli.ShowMetrics(); err != nil {
	// 	fmt.Println("Unable to show Metrics: ", err)
	// }

	initPromCollector()
}

func initPromCollector() error {
	var err error
	err = nil
	if cfg.prom == nil {
		cfg.prom = new(configProm)
	}

	cfg.prom.Collector, err = collector.NewCollectorMaster(cfg.rcli, *cfg.collectorInterval)
	if err != nil {
		log.Warnln("Init Prom: Couldn't create collector: ", err)
		return err
	}

	cfg.prom.Registry = prometheus.NewRegistry()
	err = cfg.prom.Registry.Register(cfg.prom.Collector)
	if err != nil {
		log.Errorln("Init Prom: Couldn't register collector:", err)
		return err
	}

	cfg.prom.Gatherers = &prometheus.Gatherers{
		prometheus.DefaultGatherer,
		cfg.prom.Registry,
	}
	return nil
}
