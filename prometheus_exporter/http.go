package prometheus_exporter

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log/slog"
	"net/http"
)

func RunHTTPServer(logger *slog.Logger, config Config) {
	if config.LteMonitoring.Enabled {
		go CollectLteMetrics(logger, config)
	}
	if config.InterfacesMonitoring.Enabled {
		go CollectInterfaceMetrics(logger, config)
	}

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/dashboard/interfaces.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/dashboards/interfaces.json")
	})
	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		logger.Error("failed to run http server", "error", err)
		return
	}
}
