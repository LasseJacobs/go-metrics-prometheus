Prometheus
--------------

This is a reporter for the [go-metrics](https://github.com/rcrowley/go-metrics)
that will posts metrics to [Prometheus](https://prometheus.io).

Note: this codebase is very young and might contain bugs.

### Why use this module?
It is a very minimal module with no external dependencies other that the go-metrics module. This module does not require
the standard Prometheus go client which itself is quite a large module containing lots of other code not related to 
exposing metrics to Prometheus.

This module is quite small so feel free to internalize it into your codebase and adjust it to your needs.

### Usage

```go
import "github.com/rcrowley/go-metrics"
import "github.com/LasseJacobs/go-metrics-prometheus"

http.Handle("/metrics", prometheus.MakePrometheusHandler(nil, nil))

fmt.Printf("Starting server at port 8080\n")
if err := http.ListenAndServe(":8080", nil); err != nil {
log.Fatal(err)
}
```
See `example/simple-server` for a full working example. Alternatively you can run `docker-compose up --build`, this will start a
test web-app that exposes some metrics, and a prometheus server that will monitor these metrics.

By supplying `nil` to the `MakePrometheusHandler` the default metrics registry will be used. 
