package router

import (
	"github.com/go-chi/chi/v5"
	"go-contracts/service"
)

type Routes struct {
	router *chi.Mux
	svc    service.Service // 业务服务实例
}

// NewRoutes ... Construct a new route handler instance
func NewRoutes(r *chi.Mux, svc service.Service) Routes {
	return Routes{
		router: r,
		svc:    svc,
	}
}
