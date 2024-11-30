package prometheus_exporter

type InterfacesMonitoring struct {
	Enabled bool
}

type LteMonitoring struct {
	Enabled       bool
	InterfaceName string
}

type MikrotikApi struct {
	Host     string
	Port     int
	Username string
	Password string
}

type Config struct {
	MikrotikApi          MikrotikApi
	InterfacesMonitoring InterfacesMonitoring
	LteMonitoring        LteMonitoring
}
