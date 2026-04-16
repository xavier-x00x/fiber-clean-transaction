package routes

import (
	"fiber-clean-transaction/internal/delivery/http/handler"
	"fiber-clean-transaction/internal/domain/infrastructure"
	"fiber-clean-transaction/internal/usecase"
)

type UserRoutes struct {
	handler *handler.UserHandler
}

// GetModuleName returns the name of the module
func (r *UserRoutes) GetModuleName() string {
	return "User"
}

// RegisterHandler initializes the handler and its dependencies
func (r *UserRoutes) RegisterHandler(c HandlerContainer) {
	repo := infrastructure.NewUserRepository(c.DB)
	userUsecase := usecase.NewUserUsecase(repo, c.Validator, c.UnitOfWork)
	r.handler = handler.NewUserHandler(userUsecase)
}

// RegisterRoutes sets up the routes for this module
func (r *UserRoutes) RegisterRoutes(c RouteContainer) {
	api := c.App.Group("/api")

	user := api.Group("/users", c.AuthMiddleware)
	user.Get("/", r.handler.GetAllFilter)
	user.Get("/:id", r.handler.GetUser)
	user.Put("/:id", r.handler.UpdateUser)
	user.Delete("/:id", r.handler.DeleteUser)
}
