package server

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humafiber"
	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/victorvcruz/shipment-coordinator/internal/platform/config"
	log "go.uber.org/zap"
)

func RouterSetup(cfg *config.AppConfig, handler *Handler) huma.API {
	humaConfig := huma.DefaultConfig("Shipment Coordinator", "1.0.0")
	humaConfig.Servers = []*huma.Server{
		{Description: "Development", URL: "http://localhost:" + cfg.Server.Port},
	}

	app := fiber.New()
	app.Use(fiberzap.New(fiberzap.Config{
		Logger: log.L(),
	}))
	app.Use(otelfiber.Middleware())

	api := humafiber.New(app, humaConfig)

	RegisterRoutes(api, handler)

	return api
}

func RegisterRoutes(api huma.API, handler *Handler) {
	huma.Register(api, huma.Operation{
		Method:        http.MethodPost,
		Path:          "/api/v1/orders",
		Summary:       "Create a new order",
		Description:   "Creates a new package order with the provided details",
		Tags:          []string{"Orders"},
		DefaultStatus: http.StatusCreated,
		Errors:        []int{400, 500},
	}, handler.CreateOrder)

	huma.Register(api, huma.Operation{
		Method:        http.MethodGet,
		Path:          "/api/v1/orders/{id}",
		Summary:       "Get an order by ID",
		Description:   "Retrieves the details of an order by its ID",
		Tags:          []string{"Orders"},
		DefaultStatus: http.StatusOK,
		Errors:        []int{404, 500},
	}, handler.GetOrder)

	huma.Register(api, huma.Operation{
		Method:        http.MethodPatch,
		Path:          "/api/v1/orders/{id}",
		Summary:       "Update order status",
		Description:   "Updates the status of an existing order",
		Tags:          []string{"Orders"},
		DefaultStatus: http.StatusOK,
		Errors:        []int{400, 500},
	}, handler.UpdateOrderStatus)

	huma.Register(api, huma.Operation{
		Method:        http.MethodPost,
		Path:          "/api/v1/carriers",
		Summary:       "Create a new carrier",
		Description:   "Creates a new carrier with the provided details",
		Tags:          []string{"Carriers"},
		DefaultStatus: http.StatusCreated,
		Errors:        []int{400, 500},
	}, handler.CreateCarrier)

	huma.Register(api, huma.Operation{
		Method:        http.MethodGet,
		Path:          "/api/v1/shipping/quotes/{order_id}",
		Summary:       "Get shipping quotes",
		Description:   "Retrieves shipping quotes for all carriers based on the order ID",
		Tags:          []string{"Shipping"},
		DefaultStatus: http.StatusOK,
		Errors:        []int{404, 500},
	}, handler.GetQuotes)

	huma.Register(api, huma.Operation{
		Method:        http.MethodPost,
		Path:          "/api/v1/shipping/contracts",
		Summary:       "Create a shipping contract",
		Description:   "Creates a shipping contract with the selected carrier for the specified order",
		Tags:          []string{"Shipping"},
		DefaultStatus: http.StatusOK,
		Errors:        []int{400, 500},
	}, handler.ContractCarrier)
}
