package main

import (
	"github.com/go-chi/chi/v5"

	"github.com/h3ll0kitt1/loyality-system/internal/config"
	"github.com/h3ll0kitt1/loyality-system/internal/graceful"
	"github.com/h3ll0kitt1/loyality-system/internal/handlers"
	"github.com/h3ll0kitt1/loyality-system/internal/middleware"
	"github.com/h3ll0kitt1/loyality-system/internal/repository"
	"github.com/h3ll0kitt1/loyality-system/internal/service"
)

func main() {

	cfg := config.NewConfig()
	cfg.Parse()

	//log := logger.NewLogger(cfg)

	r, err := repository.NewRepository(cfg.DatabaseDSN)
	if err != nil {
		return
	}

	s := service.NewService(r)
	h := handlers.NewHandlers(s)

	router := chi.NewRouter()
	register(router, h)

	graceful.StartServer(router, cfg.Server.HostPort)
}

func register(router *chi.Mux, h handlers.Handlers) {

	router.Route("/", func(r chi.Router) {
		router.Route("/api/user", func(router chi.Router) {

			router.Post("/register", h.RegisterUser)
			router.Post("/login", h.AuthUser)

			r.Group(func(r chi.Router) {
				r.Use(middleware.CheckAuth)

				router.Route("/orders", func(router chi.Router) {
					router.Post("/", h.LoadOrder)
					router.Get("/", h.GetOrdersInfo)
				})

				router.Route("/balance", func(router chi.Router) {
					router.Post("/", h.GetBonusInfo)
					router.Get("/withdraw", h.WithdrawBonus)
				})

				router.Get("/withdrawals", h.GetBonusOperationsInfo)
			})

		})
	})

	router.NotFound(h.ErrorNotFound)
}
