package handlers

import (
	"github.com/MarselBissengaliyev/soilem/internal/middleware"
	"github.com/MarselBissengaliyev/soilem/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	m := middleware.NewMiddleware(h.services)

	v1 := router.Group("/api/v1")

	initUserRoutes(v1, h, m)

	return router
}

func initUserRoutes(rg *gin.RouterGroup, h *Handler, m *middleware.Middleware) {
	users := rg.Group("/users")
	{
		users.GET("/:user_name", h.getUserByUserName)

		auth := users.Group("/auth")
		{
			auth.POST("/registration", h.registration)
			auth.POST("/login", h.login)
		}

		privateRoutes := auth.Group("/", m.Authenticate)
		{
			privateRoutes.POST("/send-smscode", h.sendSMSCode)
			privateRoutes.POST("/send-emailcode", h.sendEmailCode)
			privateRoutes.PUT("/confirm-smscode", h.confirmSMSCode)
			privateRoutes.PUT("/confirm-emailcode", h.confirmEmailCode)
			privateRoutes.GET("/logout", h.logout)
		}
	}
}
