package otel

import (
	"context"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	spanTrace "go.opentelemetry.io/otel/trace"

	"github.com/Allen-Career-Institute/go-bff-commons/v1/config"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/pkg/logger"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/pkg/utils"
)

type Handler struct {
	Cnf config.Config
	log logger.Logger
}

func NewHandler(cnf *config.Config, log logger.Logger) *Handler {
	return &Handler{Cnf: *cnf, log: log}
}

func (h *Handler) InitOtelProviders(endpoint1 string, endpoint2 string) {
	ctx := context.Background()

	metricResource := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(h.Cnf.Server.App.Name),
		semconv.ServiceVersion(os.Getenv(utils.Tag)),
		semconv.DeploymentEnvironment(os.Getenv(utils.Env)),
		attribute.KeyValue{
			Key:   utils.ServiceEnv,
			Value: attribute.StringValue(os.Getenv(utils.Env)),
		},
	)

	var metricOptions []otlpmetrichttp.Option

	temporality, err := h.Cnf.DynamicConfig.Get(utils.OtelGrpcTemporality)
	if err != nil {
		h.log.Errorf("failed to get OtelGrpcTemporality from dynamic config, Err: %v", err)
		temporality = utils.Delta
	}

	switch temporality {
	case utils.Delta:
		h.log.Infof("temporality: using %v temporality for metric exporter", utils.Delta)
		metricOptions = []otlpmetrichttp.Option{
			otlpmetrichttp.WithInsecure(),
			otlpmetrichttp.WithEndpoint(endpoint1),
			otlpmetrichttp.WithTemporalitySelector(deltaSelector),
		}
	case utils.Cumulative:
		h.log.Infof("temporality: using %v temporality for metric exporter", utils.Cumulative)
		metricOptions = []otlpmetrichttp.Option{
			otlpmetrichttp.WithInsecure(),
			otlpmetrichttp.WithEndpoint(endpoint2),
			otlpmetrichttp.WithTemporalitySelector(cumulativeSelector),
		}
	}

	metricExporter, err := otlpmetrichttp.New(ctx, metricOptions...)
	if err != nil {
		panic(err)
	}

	meterProviderOptions := []metric.Option{
		metric.WithReader(metric.NewPeriodicReader(metricExporter, metric.WithInterval(time.Second*10))),
		metric.WithResource(metricResource),
	}

	meterProvider := metric.NewMeterProvider(meterProviderOptions...)
	otel.SetMeterProvider(meterProvider)
	// The MeterProvider is configured and registered globally. You can now run
	// your code instrumented with the OpenTelemetry API that uses the global
	// MeterProvider without having to pass this MeterProvider instance. Or,
	// you can pass this instance directly to your instrumented code if it
	// accepts a MeterProvider instance.

	traceExporter1, err := otlptracehttp.New(ctx, otlptracehttp.WithInsecure(), otlptracehttp.WithEndpoint(endpoint1))
	if err != nil {
		panic(err)
	}

	traceResource := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(utils.TracerServiceName),
		semconv.ServiceVersion(os.Getenv(utils.Tag)),
		semconv.DeploymentEnvironment(os.Getenv(utils.Env)),
		attribute.KeyValue{
			Key:   utils.ServiceEnv,
			Value: attribute.StringValue(os.Getenv(utils.Env)),
		},
	)

	traceProviderOptions := []trace.TracerProviderOption{
		trace.WithIDGenerator(xray.NewIDGenerator()),
		trace.WithBatcher(traceExporter1),
		trace.WithResource(traceResource),
	}

	tracerProvider := trace.NewTracerProvider(traceProviderOptions...)
	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.TraceContext{})
}

func Trace(c echo.Context, name string) (echo.Context, spanTrace.Span) {
	parentCtx := utils.GetRequestCtx(c)
	tracer := spanTrace.SpanFromContext(parentCtx).TracerProvider().Tracer(utils.SpanTracer)
	childCtx, span := tracer.Start(parentCtx, name)
	c.SetRequest(c.Request().WithContext(childCtx))

	return c, span
}

func deltaSelector(kind metric.InstrumentKind) metricdata.Temporality {
	switch kind {
	case metric.InstrumentKindCounter,
		metric.InstrumentKindHistogram,
		metric.InstrumentKindObservableGauge,
		metric.InstrumentKindObservableCounter:
		return metricdata.DeltaTemporality
	case metric.InstrumentKindUpDownCounter,
		metric.InstrumentKindObservableUpDownCounter:
		return metricdata.CumulativeTemporality
	}

	return metricdata.CumulativeTemporality
}

func cumulativeSelector(kind metric.InstrumentKind) metricdata.Temporality {
	return metricdata.CumulativeTemporality
}
