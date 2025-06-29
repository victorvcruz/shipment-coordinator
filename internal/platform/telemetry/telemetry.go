package telemetry

import (
	"context"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/log/global"
	"os"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/victorvcruz/shipment-coordinator/internal/platform/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	otellog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.34.0"
)

type ShutdownFunc func(ctx context.Context) error

func Init(ctx context.Context, config *config.AppConfig) (ShutdownFunc, error) {
	if !config.Telemetry.Enabled {
		log.Info("Telemetry is disabled, skipping initialization")
		return nil, nil
	}

	shutdownMetrics, err := initMetricProvider(ctx, config)
	if err != nil {
		return nil, err
	}

	shutdownTracer, err := initTracerProvider(ctx, config)
	if err != nil {
		return nil, err
	}

	//shutdownLogger, err := initLoggerProvider(ctx, config)
	//if err != nil {
	//	return nil, err
	//}

	return func(ctx context.Context) error {
		log.Info("Shutting down telemetry")

		if shutdownMetrics != nil {
			if err = shutdownMetrics(ctx); err != nil {
				log.Errorf("error shutting down metrics: %v", err)
			}
		}

		if shutdownTracer != nil {
			if err = shutdownTracer(ctx); err != nil {
				log.Errorf("error shutting down tracer: %v", err)
			}
		}

		//if shutdownLogger != nil {
		//	if err = shutdownLogger(ctx); err != nil {
		//		log.Errorf("error shutting down logger: %v", err)
		//	}
		//}

		return nil
	}, nil
}

func initMetricProvider(ctx context.Context, config *config.AppConfig) (ShutdownFunc, error) {
	if !config.Telemetry.Enabled {
		log.Info("Metrics telemetry is disabled, skipping initialization")
		return nil, nil
	}

	exporter, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithEndpoint(config.Telemetry.Endpoint),
	)
	if err != nil {
		return nil, err
	}

	hostname, _ := os.Hostname()

	res, err := resource.New(ctx,
		resource.WithAttributes(
			attribute.String("host.name", hostname),
			attribute.String("env", config.Env),
			semconv.ServiceName(config.Service),
			semconv.ServiceVersion(config.Version),
		),
	)
	if err != nil {
		return nil, err
	}

	provider := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(metric.NewPeriodicReader(exporter, metric.WithInterval(15*time.Second))),
	)
	otel.SetMeterProvider(provider)

	err = InitMetrics(provider)
	if err != nil {
		return nil, err
	}

	return provider.Shutdown, nil
}

func initTracerProvider(ctx context.Context, config *config.AppConfig) (ShutdownFunc, error) {
	if !config.Telemetry.Enabled {
		log.Info("Metrics telemetry is disabled, skipping initialization")
		return nil, nil
	}

	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(config.Telemetry.Endpoint),
	)
	if err != nil {
		return nil, err
	}

	hostname, _ := os.Hostname()

	res, err := resource.New(ctx,
		resource.WithAttributes(
			attribute.String("env", config.Env),
			attribute.String("host.name", hostname),
			semconv.ServiceName(config.Service),
			semconv.ServiceVersion(config.Version),
		),
	)
	if err != nil {
		return nil, err
	}

	tracerProvider := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(res),
	)
	otel.SetTracerProvider(tracerProvider)

	return tracerProvider.Shutdown, nil
}

func initLoggerProvider(ctx context.Context, config *config.AppConfig) (ShutdownFunc, error) {
	if !config.Telemetry.Enabled {
		log.Info("Metrics telemetry is disabled, skipping initialization")
		return nil, nil
	}

	exporter, err := otlploggrpc.New(ctx,
		otlploggrpc.WithInsecure(),
		otlploggrpc.WithEndpoint(config.Telemetry.Endpoint),
	)

	if err != nil {
		return nil, err
	}

	hostname, _ := os.Hostname()

	res, err := resource.New(ctx,
		resource.WithAttributes(
			attribute.String("env", config.Env),
			attribute.String("host.name", hostname),
			semconv.ServiceName(config.Service),
			semconv.ServiceVersion(config.Version),
		),
	)
	if err != nil {
		return nil, err
	}

	loggerProvider := otellog.NewLoggerProvider(
		otellog.WithProcessor(otellog.NewBatchProcessor(exporter)),
		otellog.WithResource(res),
	)

	global.SetLoggerProvider(loggerProvider)

	return loggerProvider.Shutdown, nil
}
