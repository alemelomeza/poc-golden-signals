# POC-GOLDEN-SIGNALS

## Services

| Service    | Port  |
| ---------- | ----- |
| Grafana    | :3000 |
| Prometheus | :9090 |
| App        | :1313 |

## Endpoints application

| HTTP Method | URL        | Description        |
| ----------- | ---------- | ------------------ |
| GET         | /          | Root               |
| GET         | /metrics   | Prometheus metrics |
| GET         | /counter   | Increase counter   |
| GET         | /histogram | Increase histogram |

## Latency
```
sum(myapp_greeting_seconds_sum)/sum(myapp_greeting_seconds_count)  //Average
histogram_quantile(0.95, sum(rate(myapp_greeting_seconds_bucket[5m])) by (le)) //Percentile p95
```

## Request rate
```
sum(rate(myapp_greeting_seconds_count{}[2m]))  //Including errors
rate(myapp_greeting_seconds_count{code="200"}[2m])  //Only 200 OK requests
```

## Errors per second
```
sum(rate(myapp_greeting_seconds_count{}[2m]))  //Including errors
rate(myapp_greeting_seconds_count{code="200"}[2m])  //Only 200 OK requests
```

## Saturation
```
100 - (avg by (instance) (irate(node_cpu_seconds_total{}[5m])) * 100)
```
