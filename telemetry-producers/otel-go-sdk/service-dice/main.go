package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/resource"
)

const (
	periodicReaderIntervalSeconds = 5
	grpc                          = true    // False will use the stdout logger as exporter. True will use the otlp gRPC exporter.
	temporality                   = "delta" // Values: "delta" or "cumulative"
)

func main() {
	// // [TRACING] 1. Set up propagator.
	// prop := newPropagator()
	// otel.SetTextMapPropagator(prop)

	// // [TRACING] 2. Set up propagator.
	// tracerProvider, err := newTraceProvider()
	// if err != nil {
	// 	panic(err)
	// }
	// otel.SetTracerProvider(tracerProvider)

	// 1. Create resource
	res, err := newResource()
	if err != nil {
		panic(err)
	}

	// 2. Create a meter provider.
	meterProvider, err := newMeterProvider(res, grpc)
	if err != nil {
		panic(err)
	}

	// 3. Handle shutdown properly so nothing leaks.
	defer func() {
		if err := meterProvider.Shutdown(context.Background()); err != nil {
			log.Println(err)
		}
		log.Println("Shutting down meter provider")
	}()

	// 3. Register as global meter provider so that it can be used via otel.Meter
	// and accessed using otel.GetMeterProvider.
	otel.SetMeterProvider(meterProvider)

	log.Println("Starting server. Go to localhost:8080/rolldice")
	http.HandleFunc("/rolldice", rolldice)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

// func newPropagator() propagation.TextMapPropagator {
// 	return propagation.NewCompositeTextMapPropagator(
// 		propagation.TraceContext{},
// 		propagation.Baggage{},
// 	)
// }

// func newTraceProvider() (*trace.TracerProvider, error) {
// 	ctx := context.Background() // Create a context for the gRPC exporter client
// 	var traceExporter *otlptrace.Exporter
// 	var err error

// 	log.Println("Using gRPC traceExporter")
// 	traceExporter, err = otlptracegrpc.New(
// 		ctx,
// 		otlptracegrpc.WithEndpoint("localhost:4317"),
// 		otlptracegrpc.WithInsecure(),
// 	)
// 	if err != nil {
// 		return nil, err
// 	}

// 	traceProvider := trace.NewTracerProvider(
// 		trace.WithBatcher(traceExporter,
// 			trace.WithBatchTimeout(5*time.Second)),
// 	)
// 	return traceProvider, nil
// }

func newResource() (*resource.Resource, error) {
	return resource.New(
		context.Background(),
		resource.WithFromEnv(),      // Discover and provide attributes from OTEL_RESOURCE_ATTRIBUTES and OTEL_SERVICE_NAME environment variables.
		resource.WithTelemetrySDK(), // Discover and provide information about the OpenTelemetry SDK used.
		resource.WithContainer(),    // Discover and provide container information.
		// resource.WithProcess(),      // Discover and provide process information.
		// resource.WithOS(),           // Discover and provide OS information.
		// resource.WithHost(),         // Discover and provide host information.
		resource.WithAttributes(
			attribute.Key("service.name").String("service-dice"), // Add custom resource attributes.
		),
	)
}

func TemporalitySelector(temporality string) metric.TemporalitySelector {
	return func(kind metric.InstrumentKind) metricdata.Temporality {
		if temporality == "delta" {
			return metricdata.DeltaTemporality
		} else if temporality == "cumulative" {
			return metricdata.CumulativeTemporality
		}
		return metricdata.CumulativeTemporality // Default to cumulative
	}
}

func newMeterProvider(res *resource.Resource, grpc bool) (*metric.MeterProvider, error) {
	ctx := context.Background() // Create a context for the gRPC exporter client
	var metricExporter metric.Exporter
	var err error

	if grpc == true {
		log.Println("Using gRPC exporter")
		metricExporter, err = otlpmetricgrpc.New(
			ctx,
			otlpmetricgrpc.WithEndpoint("localhost:4317"),
			otlpmetricgrpc.WithInsecure(),
			otlpmetricgrpc.WithTemporalitySelector(TemporalitySelector(temporality)),
		)
	} else {
		log.Println("Using stdout exporter")
		metricExporter, err = stdoutmetric.New( // stdout for debug purposes. Otherwise we use otlpmetricgrpc to export metrics to the OTEL Collector
			stdoutmetric.WithPrettyPrint(),
			stdoutmetric.WithTemporalitySelector(TemporalitySelector(temporality)),
		)
	}
	if err != nil {
		return nil, err
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(metric.NewPeriodicReader(
			metricExporter,
			metric.WithInterval(periodicReaderIntervalSeconds*time.Second),
		)), // PeriodicReader is a Reader that continuously collects and exports metric data at a set interval.
	)
	return meterProvider, nil
}
