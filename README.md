# mikrotik-prometheus-exporter

Simple project to export Mikrotik interface & lte metrics to Prometheus.

Docker image: `sokolimedia/mikrotik-prometheus-exporter:latest`

Project exports http api on `:9000` with metrics at `/metrics` url.

Cli options:
* `--monitor-interfaces`: enable exporting interface metrics (default: disabled)
* `--monitor-lte`: enable exporting LTE metrics (default: disabled)
* `--lte-interface-id <id>`: name of the LTE interface (default: `lte1`)

Required environmental variables:
 * `MIKROTIK_API_HOST`: host of your RouterOS device
 * `MIKROTIK_API_PORT`: port of the api interface on your RouterOS
 * `MIKROTIK_API_USERNAME`: username for your RouterOS account
 * `MIKROTIK_API_PASSWORD`: password for your RouterOS account
