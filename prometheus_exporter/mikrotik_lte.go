package prometheus_exporter

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"gopkg.in/routeros.v2"
	"log/slog"
	"strings"
	"sync"
	"time"
)

var lteLabels = []string{"host", "technology", "model"}
var cqiMetric = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "lte_modem_cqi"}, lteLabels)
var sinrMetric = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "lte_modem_sinr"}, lteLabels)
var rsrqMetric = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "lte_modem_rsrq"}, lteLabels)
var rsrpMetric = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "lte_modem_rsrp"}, lteLabels)
var lteModemLastUpdate = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "lte_modem_last_update"}, lteLabels)

func collectLte(logger *slog.Logger, config Config) {
	logger.Info("collecting LTE metrics")

	address := fmt.Sprintf("%s:%d", config.MikrotikApi.Host, config.MikrotikApi.Port)
	client, err := routeros.Dial(address, config.MikrotikApi.Username, config.MikrotikApi.Password)
	if err != nil {
		logger.Error("couldn't connect to mikrotik", "error", err)
		return
	}
	defer client.Close()

	idArgument := fmt.Sprintf("=.id=%s", config.LteMonitoring.InterfaceName)
	reply, err := client.Run("/interface/lte/monitor", idArgument, "=duration=1s")
	if err != nil {
		logger.Error("couldn't execute api command", "error", err)
	}

	sentence := *reply.Re[0]

	technology, err1 := getKey(sentence, "access-technology")
	model, err2 := getKey(sentence, "model")
	model = strings.Trim(model, "\"")
	if err1 != nil {
		logger.Error("couldn't get technology", "error", err)
		return
	}
	if err2 != nil {
		logger.Error("couldn't get model", "error", err)
		return
	}

	cqi, err := getKeyAsFloat(sentence, "cqi")
	if err != nil {
		logger.Error("couldn't get cqi", "error", err)
		return
	}
	sinr, err := getKeyAsFloat(sentence, "sinr")
	if err != nil {
		logger.Error("couldn't get sinr", "error", err)
		return
	}
	rsrq, err := getKeyAsFloat(sentence, "rsrq")
	if err != nil {
		logger.Info("couldn't get rsrq", "error", err)
		return
	}
	rsrp, err := getKeyAsFloat(sentence, "rsrp")
	if err != nil {
		logger.Info("couldn't get rsrp", "error", err)
		return
	}

	labels := prometheus.Labels{
		"host":       config.MikrotikApi.Host,
		"technology": technology,
		"model":      model,
	}
	cqiMetric.With(labels).Set(cqi)
	sinrMetric.With(labels).Set(sinr)
	rsrqMetric.With(labels).Set(rsrq)
	rsrpMetric.With(labels).Set(rsrp)
	lteModemLastUpdate.With(labels).SetToCurrentTime()
}

func CollectLteMetrics(logger *slog.Logger, config Config, wg *sync.WaitGroup, quitChannel chan bool) {
	defer wg.Done()
	logger = logger.With("exporter", "mikrotik-interface")
	ticker := time.NewTicker(15 * time.Second)
	for {
		select {
		case <-quitChannel:
			logger.Info("closing lte metrics gracefully")
			return
		case <-ticker.C:
			collectLte(logger, config)
		}
	}
}
