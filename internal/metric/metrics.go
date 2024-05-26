package metric

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	namespace = "chat_server_space"
	appName   = "chat_server"
)

type Metrics struct {
	requestCounter prometheus.Counter
}

var metrics *Metrics

func Init(_ context.Context) error {
	metrics = &Metrics{
		requestCounter: promauto.NewCounter( // Автоматическая регистрация метрики
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: "grpc",
				Name:      appName + "_requests_total",
				Help:      "Количество запросов к серверу всего",
			},
		),
	}

	return nil
}

func IncRequestCounter() {
	metrics.requestCounter.Inc()
}
