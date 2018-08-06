package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/ashiddo11/sqs-exporter/collector"
)

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func main() {
	var (
		listenAddress = flag.String("web.listen-address", ":9434", "Address to listen on for web interface and telemetry.")
		metricsPath   = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
	)
	flag.Parse()

	http.Handle(*metricsPath, collector.MetricHandler{})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
	        <html>
	        <head><title>SQS Exporter</title></head>
	        <body>
	        <h1>SQS Exporter</h1>
	        <p><a href='` + *metricsPath + `'>Metrics</a></p>
	        </body>
	        </html>`))
	})
	log.Println("Starting exporter on", *listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, logRequest(http.DefaultServeMux)))
}
