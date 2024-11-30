package main

import (
	"flag"
	"log/slog"
	"mikrotik-prometheus-exporter/prometheus_exporter"
	"os"
	"strconv"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	monitorInterfaces := flag.Bool("monitor-interfaces", false, "Enable interface monitoring")
	monitorLte := flag.Bool("monitor-lte", false, "Enable lte monitoring")
	lteInterfaceId := flag.String("lte-interface-id", "lte1", "Name of the LTE interface to monitor")

	flag.Parse()

	mikrotikApiHost := os.Getenv("MIKROTIK_API_HOST")
	mikrotikApiPort := os.Getenv("MIKROTIK_API_PORT")
	mikrotikApiUsername := os.Getenv("MIKROTIK_API_USERNAME")
	mikrotikApiPassword := os.Getenv("MIKROTIK_API_PASSWORD")

	if mikrotikApiHost == "" || mikrotikApiPort == "" || mikrotikApiUsername == "" || mikrotikApiPassword == "" {
		logger.Error("one or more of the environment variables are empty")
		return
	}

	mikrotikApiPortInteger, err := strconv.Atoi(mikrotikApiPort)
	if err != nil {
		logger.Error("couldn't convert port number to integer", "error", err)
		return
	}

	config := prometheus_exporter.Config{
		MikrotikApi: prometheus_exporter.MikrotikApi{
			Host:     mikrotikApiHost,
			Port:     mikrotikApiPortInteger,
			Username: mikrotikApiUsername,
			Password: mikrotikApiPassword,
		},
		InterfacesMonitoring: prometheus_exporter.InterfacesMonitoring{
			Enabled: *monitorInterfaces,
		},
		LteMonitoring: prometheus_exporter.LteMonitoring{
			Enabled:       *monitorLte,
			InterfaceName: *lteInterfaceId,
		},
	}

	prometheus_exporter.RunHTTPServer(logger, config)
}
