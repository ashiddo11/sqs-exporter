package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/ashiddo11/sqs-exporter/collector"
)

var addr = flag.String("", ":9434", "The address to listen on for HTTP requests.")

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func main() {
	flag.Parse()
	http.Handle("/metrics", collector.MetricHandler{})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
	        <html>
	        <head><title>SQS Exporter</title></head>
	        <body>
	        <h1>SQS Exporter</h1><p><a href="/metrics">Metrics</a></p>
	        </body>
	        </html>`))
	})
	log.Println("Starting exporter on", *addr)
	log.Fatal(http.ListenAndServe(*addr, logRequest(http.DefaultServeMux)))
}
