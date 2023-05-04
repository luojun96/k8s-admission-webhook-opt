package metric

import (
	"log"

	"github.com/prometheus/client_golang/prometheus"
)

const module = "metric"

// Metric is the interface that wraps the basic methods of a metric.
type Metric interface {
	*prometheus.CounterVec
}

var (
	counterMap = make(map[string]*prometheus.CounterVec)
)

func NewCounter(namespace, subsystem, name, help string, labels []string) *prometheus.CounterVec {
	key := namespace + subsystem + name
	if _, ok := counterMap[key]; !ok {
		counterMap[key] = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      name,
				Help:      help,
			},
			labels,
		)
		prometheus.MustRegister(counterMap[key])
	}
	log.Println(module, "NewCounter", key)
	return counterMap[key]
}
