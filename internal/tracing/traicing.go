package tracing

import (
	"log"

	"github.com/uber/jaeger-client-go/config"
)

func Init(serviceName string) {
	cfg := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
	}

	_, err := cfg.InitGlobalTracer(serviceName)
	if err != nil {
		log.Fatalf("failed to init traicing: %v", err)
	}
}
