package routes

import (
	"fiber-clean-transaction/internal/delivery/http/handler"
	"fiber-clean-transaction/internal/domain/infrastructure"
	"fiber-clean-transaction/internal/usecase"
)

type AuthRoutes struct {
	handler *handler.AuthHandler
}

func (r *AuthRoutes) GetModuleName() string {
	return "auth"
}

func (r *AuthRoutes) RegisterHandler(container HandlerContainer) {
	// Initialize dependencies khusus untuk route ini
	repo := infrastructure.NewUserRepository(container.DB)
	uc := usecase.NewUserUsecase(repo)
	r.handler = handler.NewAuthHandler(uc)
}

func (r *AuthRoutes) RegisterRoutes(rc RouteContainer) {

	// Setup API group
	api := rc.App.Group("/api")

	route := api.Group("/auth")
	route.Post("/register", r.handler.Register)
	route.Post("/login", r.handler.Login)

	route.Post("/google-register", r.handler.GoogleRegister)
	route.Post("/google-auth", r.handler.GoogleAuth)
	route.Post("/refresh", r.handler.Refresh)
	route.Post("/logout", r.handler.Logout)

	route.Use(rc.AuthMiddleware)
	route.Get("/profile", r.handler.Profile)
}
