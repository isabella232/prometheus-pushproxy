# Prometheus Push Proxy

Prometheus Push Proxy, PPP, enables application to push Prometheus to this Proxy to be exposed for scraping. The use case is to get Prometheus exp data from applications behind a firewall.

## Docker

## Helm chart

```
export PPP_CHART_DIR="$HOME/go/src/github.com/kafkaesque-io/prometheus-pushproxy/prometheus-pushproxy-chart"

helm3 install --debug --dry-run prod-prometheus-pushproxy --namespace monitoring --values values.yaml $PPP_CHART_DIR
```