# Prometheus Push Proxy

Prometheus Push Proxy, PPP, enables application to push Prometheus to this Proxy to be exposed for scraping. The use case is to get Prometheus exp data from applications behind a firewall.

## Endpoints

### Prometheus push endpoint
- Push from a specific instance, the instance id has to be unique and conform http URL standard
HTTP Methed - `POST`
```
/v1/proxy/{instance}
```
### Prometheus scrape endpoint
- Scrape a sepcific instance
HTTP Methed - `GET`
```
/proxy-metrics/{instance}
```

or 
- Scrape all the instances

```
/proxy-metrics
```

## Docker

```
docker run -d -it -v $HOME/go/src/github.com/kafkaesque-io/prometheus-pushproxy/config/default_config.yml:/root/config/default_config.yml -p 8981:8981 --name=prometheus-pushproxy prometheus-pushproxy
```

## Helm chart

```
export PPP_CHART_DIR="$HOME/go/src/github.com/kafkaesque-io/prometheus-pushproxy/prometheus-pushproxy-chart"

helm3 install --debug --dry-run prod-prometheus-pushproxy --namespace monitoring --values values.yaml $PPP_CHART_DIR
```