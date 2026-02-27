package routes

import (
	"fiber-clean-transaction/internal/delivery/http/handler"
	"fiber-clean-transaction/internal/domain/infrastructure"
	"fiber-clean-transaction/internal/transaction"
	"fiber-clean-transaction/internal/usecase"
	"fiber-clean-transaction/pkg/validation"
)

type PermissionRoutes struct {
	handler *handler.PermissionHandler
}

func (r *PermissionRoutes) GetModuleName() string {
	return "permission"
}

func (r *PermissionRoutes) RegisterHandler(c HandlerContainer) {
	// Initialize dependencies khusus untuk route ini
	tr := transaction.NewGormUnitOfWork(c.DB)
	vh := validation.NewValidatorHelper(c.DB)
	rp := infrastructure.NewPermissionRepository(c.DB)
	uc := usecase.NewPermissionUsecase(rp, tr, vh)
	r.handler = handler.NewPermissionHandler(uc)
}

func (r *PermissionRoutes) RegisterRoutes(c RouteContainer) {
	// Setup API group
	api := c.App.Group("/api")

	route := api.Group("/permissions", c.AuthMiddleware)
	route.Get("/", r.handler.GetAllFilter)
	route.Get("/:id", r.handler.GetPermission)
	route.Post("/", r.handler.CreatePermission)
	route.Put("/:id", r.handler.UpdatePermission)
	route.Delete("/:id", r.handler.DeletePermission)
}
