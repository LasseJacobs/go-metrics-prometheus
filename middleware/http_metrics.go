package middleware

import (
	"github.com/rcrowley/go-metrics"
	"net/http"
	"time"
)

// ResponseLatency returns a metric handler.
func ResponseLatency(next http.Handler) (http.Handler, error) {
	fn, err := CustomResponseLatency(metrics.DefaultRegistry)
	return fn(next), err
}

func CustomResponseLatency(r metrics.Registry) (func(next http.Handler) http.Handler, error) {
	s := metrics.NewExpDecaySample(1028, 0.015) // or metrics.NewUniformSample(1028)
	h := metrics.NewHistogram(s)
	err := r.Register("response_latency", h)
	if err != nil {
		return nil, err
	}

	okCount := metrics.NewCounter()
	err = r.Register("2XX_status", okCount)
	if err != nil {
		return nil, err
	}
	userErrCount := metrics.NewCounter()
	err = r.Register("4XX_status", userErrCount)
	if err != nil {
		return nil, err
	}
	sysErrCount := metrics.NewCounter()
	err = r.Register("5XX_status", sysErrCount)
	if err != nil {
		return nil, err
	}

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			t1 := time.Now()
			defer func() {
				h.Update(time.Since(t1).Nanoseconds())
			}()
			lrw := NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(lrw, r)

			if lrw.Status() >= 500 {
				sysErrCount.Inc(1)
			} else if lrw.Status() >= 400 {
				userErrCount.Inc(1)
			} else {
				okCount.Inc(1)
			}
		}
		return http.HandlerFunc(fn)
	}, nil
}
