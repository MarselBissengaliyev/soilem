package middleware

import "github.com/MarselBissengaliyev/soilem/internal/service"

type Middleware struct {
	services *service.Service
}

func NewMiddleware(services *service.Service) *Middleware {
	return &Middleware{services}
}
