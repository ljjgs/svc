package common

import (
	"github.com/opentracing/opentracing-go"
	jarger "github.com/uber/jaeger-client-go/config"
	"io"
	"time"
)

func NewTracer(serverName string, addr string) (opentracing.Tracer, io.Closer, error) {
	cfg := &jarger.Configuration{
		ServiceName: serverName,
		Sampler: &jarger.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &jarger.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
			LocalAgentHostPort:  addr,
		},
	}
	tracer, closer, err := cfg.NewTracer()
	return tracer, closer, err
}
