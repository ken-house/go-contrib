package openTelemetry

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTracer(t *testing.T) {
	cfg := TracerConfig{
		ExporterData: struct {
			Kind string `json:"kind" mapstructure:"kind"`
			Url  string `json:"url" mapstructure:"url"`
		}(struct {
			Kind string `json:"kind"`
			Url  string `json:"url"`
		}{
			Kind: "jaeger",
			Url:  "http://10.0.98.16:14268/api/traces",
		}),
		ResourceData: struct {
			ServiceName    string `json:"service_name" mapstructure:"service_name"`
			ServiceVersion string `json:"service_version" mapstructure:"service_version"`
		}(struct {
			ServiceName    string `json:"service_name"`
			ServiceVersion string `json:"service_version"`
		}{
			ServiceName:    "go_example",
			ServiceVersion: "4.0.0",
		}),
	}
	tp, clean, err := NewTracerProvider(cfg)
	if err != nil {
		assert.Error(t, err, nil)
	}
	defer clean()
	tracer := tp.GetTracer("go_example_tracer")
	fmt.Println(tracer)
	assert.Error(t, err, nil)
}
