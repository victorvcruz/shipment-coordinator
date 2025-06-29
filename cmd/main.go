package main

import (
	"context"
	"errors"
	"github.com/gofiber/fiber/v2/log"
	"github.com/victorvcruz/shipment-coordinator/cmd/server"
	"github.com/victorvcruz/shipment-coordinator/internal/carrier"
	"github.com/victorvcruz/shipment-coordinator/internal/order"
	"github.com/victorvcruz/shipment-coordinator/internal/platform/config"
	"github.com/victorvcruz/shipment-coordinator/internal/platform/logger"
	"github.com/victorvcruz/shipment-coordinator/internal/platform/postgres"
	"github.com/victorvcruz/shipment-coordinator/internal/platform/postgres/migrations"
	"github.com/victorvcruz/shipment-coordinator/internal/platform/telemetry"
	"github.com/victorvcruz/shipment-coordinator/internal/shipping"
	"net/http"
)

func main() {
	ctx := context.Background()

	cfg, err := config.Load("development")
	if err != nil {
		log.Fatal("failed to load configuration ", err)
	}

	l, err := logger.NewLogger(cfg)
	if err != nil {
		log.Fatal("failed to load configuration ", err)
	}
	defer l.Sync() //nolint:errcheck

	shutdown, err := telemetry.Init(ctx, cfg)
	if err != nil {
		log.Fatal("failed to initialize telemetry ", err)
	}
	defer shutdown(ctx) //nolint:errcheck

	db, err := postgres.Connect(ctx, cfg)
	if err != nil {
		log.Fatal("failed to connect to database ", err)
	}

	if err := migrations.Setup(*cfg); err != nil {
		log.Fatal("failed to run migrations ", err)
	}

	orderRepository := order.NewRepository(db)

	orderService := order.NewService(orderRepository)

	carrierRepository := carrier.NewRepository(db)

	carrierService := carrier.NewService(carrierRepository)

	shippingRepository := shipping.NewRepository(db)

	shippingService := shipping.NewService(orderRepository, carrierRepository, shippingRepository)

	handler := server.NewHandler(orderService, carrierService, shippingService)

	api := server.RouterSetup(cfg, handler)

	svr := http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: api.Adapter(),
	}

	log.Debug("starting server ", "port ", cfg.Server.Port)
	if serverErr := svr.ListenAndServe(); serverErr != nil &&
		!errors.Is(serverErr, http.ErrServerClosed) {
		log.Fatal("failed to start server ", serverErr)
	}
}
