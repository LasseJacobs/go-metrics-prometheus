package prometheus

import (
	"github.com/rcrowley/go-metrics"
	"io"
	"net/http"
)

var (
	percentiles = []float64{0.50, 0.75, 0.95, 0.99}
)

func MakePrometheusHandler(r metrics.Registry, collectors []Collector) http.Handler {
	var registry = r
	if registry == nil {
		registry = metrics.DefaultRegistry
	}

	if collectors != nil {
		for _, c := range collectors {
			c.Register(registry)
		}
	}

	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if collectors != nil {
			for _, c := range collectors {
				c.Collect()
			}
		}

		err := snapshot(registry, w)
		if err != nil {
			http.Error(w, "Internal Server Issue", 500)
			return
		}
	})
}

func snapshot(r metrics.Registry, writer io.Writer) error {
	r.Each(func(metricName string, i interface{}) {
		//todo name, tags := r.splitNameAndTags(metricName)
		var name = metricName

		switch metric := i.(type) {
		case metrics.Counter:
			v := metric.Count()
			writeType(name, "counter", writer)
			writeCounter(name, v, writer)

		case metrics.Gauge:
			writeType(name, "gauge", writer)
			writeGauge(name, float64(metric.Value()), writer)

		case metrics.GaugeFloat64:
			writeType(name, "gauge", writer)
			writeGauge(name, metric.Value(), writer)

		case metrics.Histogram:
			ms := metric.Snapshot()
			if ms.Count() == 0 {
				return
			}
			writeType(name, "summary", writer)
			writeGauge(name+"_count", float64(ms.Count()), writer)
			writeGauge(name+"_sum", float64(ms.Sum()), writer)

			if len(percentiles) > 0 {
				values := ms.Percentiles(percentiles)
				for i, p := range percentiles {
					//{quantile="0.01"}
					writeGaugeQuantile(name, values[i], p, writer)
				}
			}

		case metrics.Meter:
			ms := metric.Snapshot()
			writeType(name, "gauge", writer)
			writeGauge(name, ms.Rate1(), writer)

		case metrics.Timer:
			ms := metric.Snapshot()
			if ms.Count() == 0 {
				return
			}
			writeType(name, "summary", writer)
			writeGauge(name+"_count", float64(ms.Count()), writer)
			writeGauge(name+"_sum", float64(ms.Sum()), writer)

			if len(percentiles) > 0 {
				values := ms.Percentiles(percentiles)
				for i, p := range percentiles {
					//{quantile="0.01"}
					writeGaugeQuantile(name, values[i], p, writer)
				}
			}
		}
	})

	return nil
}
