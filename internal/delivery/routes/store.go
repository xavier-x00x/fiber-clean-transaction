package routes

import (
	"fiber-clean-transaction/internal/delivery/http/handler"
	"fiber-clean-transaction/internal/domain/infrastructure"
	"fiber-clean-transaction/internal/usecase"
)

type StoreRoutes struct {
	handler *handler.StoreHandler
}

func (r *StoreRoutes) GetModuleName() string {
	return "store"
}

func (r *StoreRoutes) RegisterHandler(c HandlerContainer) {
	// Initialize dependencies khusus untuk route ini
	rp := infrastructure.NewStoreRepository(c.DB)
	uc := usecase.NewStoreUsecase(rp, c.UnitOfWork, c.Validator)
	r.handler = handler.NewStoreHandler(uc)
}

func (r *StoreRoutes) RegisterRoutes(c RouteContainer) {
	// Setup API group
	api := c.App.Group("/api")

	store := api.Group("/stores", c.AuthMiddleware)
	store.Get("/", r.handler.GetAllFilter)
	store.Get("/:id", r.handler.GetStore)
	store.Post("/", r.handler.CreateStore)
	store.Put("/:id", r.handler.UpdateStore)
	store.Delete("/:id", r.handler.DeleteStore)
}
