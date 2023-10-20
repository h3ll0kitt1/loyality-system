package main

import (
	"context"
	// "encoding/json"
	// "fmt"
	"log"
	// "net/http"
	// "time"

	"github.com/go-chi/chi/v5"

	"github.com/h3ll0kitt1/loyality-system/internal/config"
	//"github.com/h3ll0kitt1/loyality-system/internal/domain"
	"github.com/h3ll0kitt1/loyality-system/internal/graceful"
	"github.com/h3ll0kitt1/loyality-system/internal/handlers"
	"github.com/h3ll0kitt1/loyality-system/internal/logger"
	"github.com/h3ll0kitt1/loyality-system/internal/middleware"
	"github.com/h3ll0kitt1/loyality-system/internal/repository"
	"github.com/h3ll0kitt1/loyality-system/internal/service"
)

func main() {

	cfg := config.NewConfig()
	if err := cfg.Parse(); err != nil {
		log.Fatal("parse config:", err)
	}

	log := logger.NewLogger()

	log.Infow("info",
		"Launch server on addr", cfg.Server.HostPort,
		"Launch postgresql on addr ", cfg.DatabaseDSN,
	)

	r, err := repository.NewRepository(context.Background(), cfg.DatabaseDSN, log)
	if err != nil {
		log.Panic(err)
	}

	s := service.NewService(r, cfg, log)
	h := handlers.NewHandlers(s, log)

	router := chi.NewRouter()
	register(router, h)

	//go updater(context.Background(), s, cfg)

	graceful.StartServer(router, cfg.Server.HostPort)
}

func register(r *chi.Mux, h handlers.Handlers) {

	r.Route("/", func(r chi.Router) {
		r.Route("/api/user", func(r chi.Router) {

			r.Post("/register", h.RegisterUser)
			r.Post("/login", h.AuthUser)

			r.Group(func(r chi.Router) {

				r.Use(middleware.CheckAuth)

				r.Route("/orders", func(r chi.Router) {
					r.Post("/", h.LoadOrder)
					r.Get("/", h.GetOrdersInfo)
				})

				r.Route("/balance", func(r chi.Router) {
					r.Get("/", h.GetBonusInfo)
					r.Post("/withdraw", h.WithdrawBonus)
				})

				r.Get("/withdrawals", h.GetBonusOperationsInfo)
			})

		})
	})
	r.NotFound(h.ErrorNotFound)
}

// func updater(ctx context.Context, s handlers.Services, cfg *config.Config) {

// 	ticker := time.NewTicker(cfg.CheckInterval)

// 	for range ticker.C {
// 		ordersToUpdate, err := s.GetOrdersForUpdate(ctx, 100) // get only NEW / PROCESSING
// 		if err != nil {
// 			log.Printf("get orders failed: %s", err)
// 			return
// 		}

// 		for _, order := range ordersToUpdate {

// 			requestURL := fmt.Sprintf("http://%s/api/orders/%d", cfg.Server.HostPort, order.Number)
// 			resp, err := http.Get(requestURL)
// 			if err != nil {
// 				log.Printf("make GET request failed: %s", err)
// 				return
// 			}

// 			if resp.StatusCode == http.StatusNoContent {
// 				log.Println("order is not registered in the system")
// 				continue
// 			}

// 			if resp.StatusCode == http.StatusTooManyRequests {
// 				log.Println("too many requests")
// 				// parse resp - retry - wait for time
// 				continue
// 			}

// 			if resp.StatusCode == http.StatusInternalServerError {
// 				log.Println("internal error")
// 				continue
// 			}

// 			var orderInfo domain.OrderInfoRequest
// 			if err := json.NewDecoder(resp.Body).Decode(&orderInfo); err != nil {
// 				log.Println("json error")
// 				continue
// 			}

// 			if err := s.UpdateOrderInfo(ctx, orderInfo); err != nil {
// 				log.Println("update error")
// 				continue
// 			}
// 		}
// 	}
// }
