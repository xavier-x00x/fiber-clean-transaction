package routes

import (
	"fiber-clean-transaction/internal/delivery/http/handler"
	"fiber-clean-transaction/internal/domain/infrastructure"
	"fiber-clean-transaction/internal/usecase"
)

type RoleRoutes struct {
	handler *handler.RoleHandler
}

// GetModuleName returns the name of the module
func (r *RoleRoutes) GetModuleName() string {
	return "Role"
}

// RegisterHandler initializes the handler and its dependencies
func (r *RoleRoutes) RegisterHandler(c HandlerContainer) {
	roleUsecase := usecase.NewRoleUsecase(usecase.RoleUsecaseDeps{
		Validator: c.Validator,
		UOW:       c.UnitOfWork,
		RoleRepo:  infrastructure.NewRoleRepository(c.DB),
		PermRepo:  infrastructure.NewPermissionRepository(c.DB),
	})
	r.handler = handler.NewRoleHandler(roleUsecase)
}

// RegisterRoutes sets up the routes for this module
func (r *RoleRoutes) RegisterRoutes(c RouteContainer) {
	api := c.App.Group("/api")

	role := api.Group("/roles", c.AuthMiddleware)
	role.Get("/", r.handler.GetAllFilter)
	role.Get("/permissions", r.handler.GetPermissionsByRole)
	role.Get("/:id", r.handler.GetRole)
	role.Post("/", r.handler.CreateRole)
	role.Post("/authorization", r.handler.Authorization)
	role.Put("/:id", r.handler.UpdateRole)
	role.Delete("/:id", r.handler.DeleteRole)
}
