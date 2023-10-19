package main

import (
	"context"

	"github.com/go-chi/chi/v5"

	"github.com/h3ll0kitt1/loyality-system/internal/config"
	"github.com/h3ll0kitt1/loyality-system/internal/graceful"
	"github.com/h3ll0kitt1/loyality-system/internal/handlers"
	"github.com/h3ll0kitt1/loyality-system/internal/logger"
	"github.com/h3ll0kitt1/loyality-system/internal/middleware"
	"github.com/h3ll0kitt1/loyality-system/internal/repository"
	"github.com/h3ll0kitt1/loyality-system/internal/service"
)

func main() {

	cfg := config.NewConfig()
	cfg.Parse()

	log := logger.NewLogger()

	log.Infow("info",
		"Launch server on addr", cfg.Server.HostPort,
		"Launch postgresql on addr ", cfg.DatabaseDSN,
	)

	r, err := repository.NewRepository(context.Background(), cfg.DatabaseDSN, log)
	if err != nil {
		log.Panic(err)
	}

	s := service.NewService(r, log)
	h := handlers.NewHandlers(s, log)

	router := chi.NewRouter()
	register(router, h)

	//go updater(s)

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
