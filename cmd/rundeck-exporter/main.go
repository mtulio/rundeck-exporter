package main

import (
	"net/http"

	"github.com/prometheus/log"
)

func main() {
	log.Infoln("Starting exporter ")

	// http.HandleFunc(*cfg.expMetricsPath, handler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>"` + exporterDescription + `"</title></head>
			<body>
			<h1>` + exporterDescription + `</h1>
			<p><br> The metrics is available on the path:
			<a href="` + *cfg.expMetricsPath + `">Metrics</a></p>
			</body>
			</html>`))
	})

	log.Info("Beginning to serve on port " + *cfg.expListenAddr)
	log.Fatal(http.ListenAndServe(*cfg.expListenAddr, nil))

}
