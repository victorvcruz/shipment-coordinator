package telemetry

import (
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
)

var (
	OrderCreatedCounter    metric.Int64Counter = noop.Int64Counter{}
	OrderUpdatedCounter    metric.Int64Counter = noop.Int64Counter{}
	ContractCreatedCounter metric.Int64Counter = noop.Int64Counter{}
)

func InitMetrics(meterProvider metric.MeterProvider) error {
	meter := meterProvider.Meter("shipment-coordinator")

	var err error

	OrderCreatedCounter, err = meter.Int64Counter(
		"order_created_total",
		metric.WithDescription("Total number of orders created"),
	)
	if err != nil {
		return err
	}

	OrderUpdatedCounter, err = meter.Int64Counter(
		"order_updated_total",
		metric.WithDescription("Total number of orders updated"),
	)
	if err != nil {
		return err
	}

	ContractCreatedCounter, err = meter.Int64Counter(
		"contract_created_total",
		metric.WithDescription("Total number of contracts created"),
	)
	if err != nil {
		return err
	}

	return nil
}
