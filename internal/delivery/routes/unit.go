package routes

import (
	"fiber-clean-transaction/internal/delivery/http/handler"
	"fiber-clean-transaction/internal/domain/infrastructure"
	"fiber-clean-transaction/internal/usecase"
)

type UnitRoutes struct {
	handler *handler.UnitHandler
}

func (r *UnitRoutes) GetModuleName() string {
	return "unit"
}

func (r *UnitRoutes) RegisterHandler(c HandlerContainer) {
	rp := infrastructure.NewUnitRepository(c.DB)
	uc := usecase.NewUnitUsecase(rp, c.SeqRepo, c.UnitOfWork, c.Validator)
	r.handler = handler.NewUnitHandler(uc)
}

func (r *UnitRoutes) RegisterRoutes(c RouteContainer) {
	// Setup API group
	api := c.App.Group("/api")

	unit := api.Group("/units", c.AuthMiddleware)
	unit.Get("/", r.handler.GetAllFilter)
	unit.Get("/:id", r.handler.GetUnit)
	unit.Post("/", r.handler.CreateUnit)
	unit.Put("/:id", r.handler.UpdateUnit)
	unit.Delete("/:id", r.handler.DeleteUnit)
}
