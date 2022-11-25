package openTelemetry

import (
	"context"
	"log"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/gin-gonic/gin"

	"go.opentelemetry.io/otel/metric"

	"go.opentelemetry.io/otel/exporters/prometheus"
	sdkMetric "go.opentelemetry.io/otel/sdk/metric"
)

type MeterProvider interface {
	GetMeter(meterName string) metric.Meter
}

type meterProvider struct {
	*sdkMetric.MeterProvider
}

type MeterConfig struct {
	ExporterData struct {
		Kind string `json:"kind" mapstructure:"kind"`
	} `json:"exporter_data" mapstructure:"exporter_data"`
	ResourceData struct {
		ServiceName    string `json:"service_name" mapstructure:"service_name"`
		ServiceVersion string `json:"service_version" mapstructure:"service_version"`
	} `json:"resource_data" mapstructure:"resource_data"`
}

func NewMeterProvider(cfg MeterConfig) (*meterProvider, func(), error) {
	mp, err := initMeterProvider(cfg)
	if err != nil {
		return nil, nil, err
	}
	return &meterProvider{
			mp,
		}, func() {
			if err := mp.Shutdown(context.Background()); err != nil {
				log.Printf("Error shutting down meter provider: %v", err)
			}
		}, nil
}

// newPrometheusExporter 定义一个prometheus导出器
func newPrometheusExporter() (*prometheus.Exporter, error) {
	return prometheus.New()
}

// 初始化MeterProvider
func initMeterProvider(cfg MeterConfig) (*sdkMetric.MeterProvider, error) {
	exporter, err := newPrometheusExporter()
	if err != nil {
		log.Println("导出器初始化失败")
		return nil, err
	}

	mp := sdkMetric.NewMeterProvider(
		sdkMetric.WithReader(exporter),
		sdkMetric.WithResource(newResource(cfg.ResourceData.ServiceName, cfg.ResourceData.ServiceVersion)),
	)
	return mp, err
}

// GetMeter 创建一个指标监控对象
func (mp *meterProvider) GetMeter(meterName string) metric.Meter {
	return mp.Meter(meterName)
}

// MeterPrometheusForGin 为Gin接入Prometheus
func (mp *meterProvider) MeterPrometheusForGin(router *gin.Engine) {
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
}
