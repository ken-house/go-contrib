package openTelemetry

import (
	"context"
	"log"

	"go.opentelemetry.io/otel/trace"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/ken-house/go-contrib/utils/env"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

type TracerProvider interface {
	GetTracer(tracerName string) trace.Tracer
}

type tracerProvider struct {
	*sdktrace.TracerProvider
}

func NewTracerProvider(cfg TracerConfig) (*tracerProvider, func(), error) {
	tp, err := initTracerProvider(cfg)
	if err != nil {
		return nil, nil, err
	}
	return &tracerProvider{
			tp,
		}, func() {
			if err := tp.Shutdown(context.Background()); err != nil {
				log.Printf("Error shutting down tracer provider: %v", err)
			}
		}, nil
}

type TracerConfig struct {
	ExporterData struct {
		Kind string `json:"kind"`
		Url  string `json:"url"`
	} `json:"exporter_data"`
	ResourceData struct {
		ServiceName    string `json:"service_name"`
		ServiceVersion string `json:"service_version"`
	} `json:"resource_data"`
}

// 定义Jaeger导出器
func newJaegerExporter(url string) (*jaeger.Exporter, error) {
	return jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
}

// 定义Tracer资源
func newResource(serviceName string, serviceVersion string) *resource.Resource {
	rs, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
			semconv.ServiceVersionKey.String(serviceVersion),
			attribute.String(env.RunMode, env.Mode()),
		),
	)
	return rs
}

// 初始化TracerProvider
func initTracerProvider(cfg TracerConfig) (*sdktrace.TracerProvider, error) {
	exporter, err := newJaegerExporter(cfg.ExporterData.Url)
	if err != nil {
		log.Println("导出器初始化失败")
		return nil, err
	}

	// 创建追踪提供对象
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(newResource(cfg.ResourceData.ServiceName, cfg.ResourceData.ServiceVersion)),
	)

	// 声明为全局对象
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp, err
}

// GetTracer 创建一个分布式追踪对象
func (tp *tracerProvider) GetTracer(tracerName string) trace.Tracer {
	return tp.Tracer(tracerName)
}
