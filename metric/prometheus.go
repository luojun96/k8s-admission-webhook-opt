package metric

import "github.com/prometheus/client_golang/prometheus"

const module = "metric"

// Metric is the interface that wraps the basic methods of a metric.
type Metric interface {
	*prometheus.CounterVec
}

var (
	counterMap = make(map[string]*prometheus.CounterVec)
)
