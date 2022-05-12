# POC-GOLDEN-SIGNALS

## Latency
```
sum(greeting_seconds_sum)/sum(greeting_seconds_count)  //Average
histogram_quantile(0.95, sum(rate(greeting_seconds_bucket[5m])) by (le)) //Percentile p95
```

## Request rate
```
sum(rate(greeting_seconds_count{}[2m]))  //Including errors
rate(greeting_seconds_count{code="200"}[2m])  //Only 200 OK requests
```

## Errors per second
```
sum(rate(greeting_seconds_count{}[2m]))  //Including errors
rate(greeting_seconds_count{code="200"}[2m])  //Only 200 OK requests
```

## Saturation
```
100 - (avg by (instance) (irate(node_cpu_seconds_total{}[5m])) * 100)
```
