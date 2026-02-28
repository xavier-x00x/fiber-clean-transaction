package routes

import (
	"fiber-clean-transaction/internal/delivery/http/handler"
	"fiber-clean-transaction/internal/domain/infrastructure"
	"fiber-clean-transaction/internal/transaction"
	"fiber-clean-transaction/internal/usecase"
	"fiber-clean-transaction/pkg/validation"
)

type CategoryRoutes struct {
	handler *handler.CategoryHandler
}

func (r *CategoryRoutes) GetModuleName() string {
	return "category"
}

func (r *CategoryRoutes) RegisterHandler(c HandlerContainer) {
	// Initialize dependencies khusus untuk route ini
	tr := transaction.NewGormUnitOfWork(c.DB)
	vh := validation.NewValidatorHelper(c.DB)
	rp := infrastructure.NewCategoryRepository(c.DB)
	uc := usecase.NewCategoryUsecase(rp, tr, vh)
	r.handler = handler.NewCategoryHandler(uc)
}

func (r *CategoryRoutes) RegisterRoutes(c RouteContainer) {
	// Setup API group
	api := c.App.Group("/api")

	category := api.Group("/categories", c.AuthMiddleware)
	category.Get("/", r.handler.GetAllFilter)
	category.Get("/:id", r.handler.GetCategory)
	category.Post("/", r.handler.CreateCategory)
	category.Put("/:id", r.handler.UpdateCategory)
	category.Delete("/:id", r.handler.DeleteCategory)
}
