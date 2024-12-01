package prometheus_exporter

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
)

func RunHTTPServer(logger *slog.Logger, config Config) {
	var wg sync.WaitGroup

	quitLteMetrics := make(chan bool, 1)
	quitInterfaceMetrics := make(chan bool, 1)

	if config.LteMonitoring.Enabled {
		wg.Add(1)
		go CollectLteMetrics(logger, config, &wg, quitLteMetrics)
	}
	if config.InterfacesMonitoring.Enabled {
		wg.Add(1)
		go CollectInterfaceMetrics(logger, config, &wg, quitInterfaceMetrics)
	}

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/dashboard/interfaces.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/dashboards/interfaces.json")
	})
	http.HandleFunc("/dashboard/lte.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/dashboards/lte.json")
	})

	sigIntChannel := make(chan os.Signal, 1)
	signal.Notify(sigIntChannel, os.Interrupt)
	go func() {
		<-sigIntChannel

		quitLteMetrics <- true
		quitInterfaceMetrics <- true

		wg.Wait()
		os.Exit(0)
	}()

	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		logger.Error("failed to run http server", "error", err)
		return
	}
}
