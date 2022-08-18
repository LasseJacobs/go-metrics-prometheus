package main

import (
	"fmt"
	"github.com/LasseJacobs/go-metrics-prometheus/prometheus"
	"github.com/rcrowley/go-metrics"
	"log"
	"net/http"
)

type countCollector struct {
	m metrics.Counter
}

func (c *countCollector) Register(r metrics.Registry) {
	r.Register("metric_load_count", c.m)
}

func (c *countCollector) Collect() {
	c.m.Inc(1)
}

func main() {
	c := metrics.NewCounter()
	metrics.Register("foo", c)
	c.Inc(47)

	g := metrics.NewGauge()
	metrics.Register("bar", g)
	g.Update(47)

	/*r := */
	//metrics.NewRegistry()
	//g := metrics.NewRegisteredFunctionalGauge("cache-evictions", r, func() int64 { return cache.getEvictionsCount() })

	s := metrics.NewExpDecaySample(1028, 0.015) // or metrics.NewUniformSample(1028)
	h := metrics.NewHistogram(s)
	metrics.Register("baz", h)
	h.Update(47)

	m := metrics.NewMeter()
	metrics.Register("quux", m)
	m.Mark(47)

	t := metrics.NewTimer()
	metrics.Register("bang", t)
	t.Time(func() {})
	t.Update(47)

	col := &countCollector{m: metrics.NewCounter()}

	http.Handle("/metrics", prometheus.MakePrometheusHandler(nil, []prometheus.Collector{col}))

	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
