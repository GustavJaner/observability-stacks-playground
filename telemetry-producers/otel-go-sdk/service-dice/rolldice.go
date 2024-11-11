package main

import (
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var (
	// tracer = otel.Tracer("dice-tracer")
	// 4. Create a Meter from the Meter Provider
	meter = otel.Meter("dice-meter")
	// 5. Create a counter from the Meter
	rollCount metric.Int64Counter
	rollSum   metric.Int64Counter
)

func init() {
	var err error
	rollCount, err = meter.Int64Counter("dice_rolls_total",
		metric.WithDescription("Total number of dice rolls"),
		metric.WithUnit("{roll}"))
	if err != nil {
		panic(err)
	}

	rollSum, err = meter.Int64Counter("dice_rolls_value_total",
		metric.WithDescription("Total sum of dice roll values"),
		metric.WithUnit("{roll}"))
	if err != nil {
		panic(err)
	}
}

func rolldice(w http.ResponseWriter, r *http.Request) {
	// ctx, span := tracer.Start(r.Context(), "roll")
	// defer span.End()

	roll := 1 + rand.Intn(6)

	log.Print("Rolled: ", roll)

	// rollValueAttr := attribute.Int("roll.value", roll)
	// span.SetAttributes(rollValueAttr)

	ctx := r.Context()
	rollCount.Add(ctx, 1, metric.WithAttributes(attribute.String("endpoint", "rolldice")))
	rollSum.Add(ctx, int64(roll), metric.WithAttributes(attribute.String("endpoint", "rolldice")))
	// rollSum.Add(ctx, 1, metric.WithAttributes(attribute.String("endpoint", "rolldice"), attribute.String("test", "foo")))

	resp := strconv.Itoa(roll) + "\n"
	if _, err := io.WriteString(w, resp); err != nil {
		log.Printf("Write failed: %v\n", err)
	}
}
