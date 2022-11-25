package openTelemetry

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMeterProvider_GetMeter(t *testing.T) {
	cfg := MeterConfig{
		ExporterData: struct {
			Kind string `json:"kind" mapstructure:"kind"`
		}{
			Kind: "prometheus",
		},
		ResourceData: struct {
			ServiceName    string `json:"service_name" mapstructure:"service_name"`
			ServiceVersion string `json:"service_version" mapstructure:"service_version"`
		}{
			ServiceName:    "go_example",
			ServiceVersion: "v4.0.0",
		},
	}

	mp, clean, err := NewMeterProvider(cfg)
	if err != nil {
		assert.Fail(t, "错误："+err.Error())
	}
	defer clean()
	meter := mp.GetMeter("go_example_meter")
	fmt.Println(meter)
	assert.True(t, true, "验证成功")
}
