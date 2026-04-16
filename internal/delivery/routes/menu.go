package routes

import (
	"fiber-clean-transaction/internal/delivery/http/handler"
	"fiber-clean-transaction/internal/domain/infrastructure"
	"fiber-clean-transaction/internal/usecase"
)

type MenuRoutes struct {
	handler *handler.MenuHandler
}

func (r *MenuRoutes) GetModuleName() string {
	return "menu"
}

func (r *MenuRoutes) RegisterHandler(c HandlerContainer) {
	menuRepo := infrastructure.NewMenuRepository(c.DB)
	uc := usecase.NewMenuUsecase(menuRepo)
	r.handler = handler.NewMenuHandler(uc)
}

func (r *MenuRoutes) RegisterRoutes(c RouteContainer) {
	api := c.App.Group("/api")

	route := api.Group("/menus", c.AuthMiddleware)
	route.Get("/", r.handler.GetMenusByRole)
}

func init() {
	RegisterModule(&MenuRoutes{})
}
