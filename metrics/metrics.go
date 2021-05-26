package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Metrics interface {
	IncreaseRequestCount()
	IncreaseFailedRequestCount()
	IncreaseComponentCount()
	IncreaseFailedComponentCount()
}

type metrics struct {
	requestCount         prometheus.Counter
	failedCount          prometheus.Counter
	componentCount       prometheus.Counter
	failedComponentCount prometheus.Counter
}

func (m metrics) IncreaseRequestCount() {
	m.requestCount.Inc()
}

func (m metrics) IncreaseFailedRequestCount() {
	m.failedCount.Inc()
}

func (m metrics) IncreaseComponentCount() {
	m.componentCount.Inc()
}

func (m metrics) IncreaseFailedComponentCount() {
	m.failedComponentCount.Inc()
}

func NewInstance() metrics {
	m := metrics{
		requestCount: promauto.NewCounter(prometheus.CounterOpts{
			Name: "alarm_notification_count",
			Help: "total alarm count",
		}),
		failedCount: promauto.NewCounter(prometheus.CounterOpts{
			Name: "failed_alarm_notification_count",
			Help: "failed request count",
		}),
		componentCount: promauto.NewCounter(prometheus.CounterOpts{
			Name: "component_count",
			Help: "total component count",
		}),
		failedComponentCount: promauto.NewCounter(prometheus.CounterOpts{
			Name: "failed_component_notification_count",
			Help: "failed component count",
		}),
	}
	prometheus.Register(m.requestCount)
	prometheus.Register(m.failedCount)
	prometheus.Register(m.componentCount)
	prometheus.Register(m.failedComponentCount)
	return m
}
