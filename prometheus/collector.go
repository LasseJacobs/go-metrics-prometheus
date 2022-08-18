package prometheus

import "github.com/rcrowley/go-metrics"

type Collector interface {
	Register(r metrics.Registry)
	Collect()
}
