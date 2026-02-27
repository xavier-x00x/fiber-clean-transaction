package routes

import (
	"fiber-clean-transaction/internal/delivery/http/handler"
	"fiber-clean-transaction/internal/domain/infrastructure"
	"fiber-clean-transaction/internal/transaction"
	"fiber-clean-transaction/internal/usecase"
	"fiber-clean-transaction/pkg/validation"
)

type StoreRoutes struct {
	handler *handler.StoreHandler
}

func (r *StoreRoutes) GetModuleName() string {
	return "store"
}

func (r *StoreRoutes) RegisterHandler(c HandlerContainer) {
	// Initialize dependencies khusus untuk route ini
	tr := transaction.NewGormRunner(c.DB)
	vh := validation.NewValidatorHelper(c.DB)
	rp := infrastructure.NewStoreRepository(c.DB)
	uc := usecase.NewStoreUsecase(rp, tr, vh)
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
