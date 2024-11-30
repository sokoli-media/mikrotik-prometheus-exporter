package prometheus_exporter

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"gopkg.in/routeros.v2"
	"log/slog"
	"time"
)

var interfaceLabels = []string{"host", "name"}
var rxBytesGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "mikrotik_interface_rx_bytes"}, interfaceLabels)
var txBytesGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "mikrotik_interface_tx_bytes"}, interfaceLabels)
var rxPacketsGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "mikrotik_interface_rx_packets"}, interfaceLabels)
var txPacketsGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "mikrotik_interface_tx_packet"}, interfaceLabels)
var fpRxBytesGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "mikrotik_interface_fp_rx_bytes"}, interfaceLabels)
var fpTxBytesGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "mikrotik_interface_fp_tx_bytes"}, interfaceLabels)
var fpRxPacketsGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "mikrotik_interface_fp_rx_packets"}, interfaceLabels)
var fpTxPacketsGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "mikrotik_interface_fp_tx_packet"}, interfaceLabels)
var lteInterfaceLastUpdate = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "mikrotik_interface_last_update"}, interfaceLabels)

func collectInterface(logger *slog.Logger, config Config) {
	logger.Info("collecting interface metrics")

	address := fmt.Sprintf("%s:%d", config.MikrotikApi.Host, config.MikrotikApi.Port)
	client, err := routeros.Dial(address, config.MikrotikApi.Username, config.MikrotikApi.Password)
	if err != nil {
		logger.Info("couldn't connect to mikrotik", "error", err)
		return
	}
	defer client.Close()

	reply, err := client.Run("/interface/print")
	if err != nil {
		logger.Info("couldn't execute api command", "error", err)
		return
	}

	for _, netInterface := range reply.Re {
		name, err := getKey(*netInterface, "name")
		if err != nil {
			logger.Error("couldn't get interface name", "error", err)
			continue
		}

		labels := prometheus.Labels{
			"host": config.MikrotikApi.Host,
			"name": name,
		}

		rxBytes, err := getKeyAsFloat(*netInterface, "rx-byte")
		if err != nil {
			logger.Error("couldn't get rx-bytes as float", "error", err)
			continue
		}
		txBytes, err := getKeyAsFloat(*netInterface, "tx-byte")
		if err != nil {
			logger.Error("couldn't get tx-bytes as float", "error", err)
			continue
		}
		rxPackets, err := getKeyAsFloat(*netInterface, "rx-packet")
		if err != nil {
			logger.Error("couldn't get rx-packet as float", "error", err)
			continue
		}
		txPackets, err := getKeyAsFloat(*netInterface, "tx-packet")
		if err != nil {
			logger.Error("couldn't get tx-packet as float", "error", err)
			continue
		}
		fpRxBytes, err := getKeyAsFloat(*netInterface, "fp-rx-byte")
		if err != nil {
			logger.Error("couldn't get fp-rx-bytes as float", "error", err)
			continue
		}
		fpTxBytes, err := getKeyAsFloat(*netInterface, "fp-tx-byte")
		if err != nil {
			logger.Error("couldn't get fp-tx-bytes as float", "error", err)
			continue
		}
		fpRxPackets, err := getKeyAsFloat(*netInterface, "fp-rx-packet")
		if err != nil {
			logger.Error("couldn't get fp-rx-packet as float", "error", err)
			continue
		}
		fpTxPackets, err := getKeyAsFloat(*netInterface, "fp-tx-packet")
		if err != nil {
			logger.Error("couldn't get fp-tx-packet as float", "error", err)
			continue
		}

		rxBytesGauge.With(labels).Set(rxBytes)
		txBytesGauge.With(labels).Set(txBytes)
		rxPacketsGauge.With(labels).Set(rxPackets)
		txPacketsGauge.With(labels).Set(txPackets)
		fpRxBytesGauge.With(labels).Set(fpRxBytes)
		fpTxBytesGauge.With(labels).Set(fpTxBytes)
		fpRxPacketsGauge.With(labels).Set(fpRxPackets)
		fpTxPacketsGauge.With(labels).Set(fpTxPackets)
		lteInterfaceLastUpdate.With(labels).SetToCurrentTime()
	}
}

func CollectInterfaceMetrics(logger *slog.Logger, config Config) {
	logger = logger.With("exporter", "mikrotik-interface")
	ticker := time.NewTicker(15 * time.Second)
	for {
		select {
		case <-ticker.C:
			logger.Info("collecting interface metrics")
			collectInterface(logger, config)
		}
	}
}
