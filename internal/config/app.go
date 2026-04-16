package config

import (
	"fiber-clean-transaction/internal/delivery/http/middleware"
	"fiber-clean-transaction/internal/delivery/routes"
	"fiber-clean-transaction/internal/domain/infrastructure"
	"fiber-clean-transaction/internal/transaction"
	"fiber-clean-transaction/pkg/validation"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gorm.io/gorm"
)

// khusus untuk setup gorm,
// bila menggunakan ORM lain, cukup ubah bagian ini saja tanpa perlu mengubah bagian lain di aplikasi

type BootstrapConfig struct {
	DB  *gorm.DB
	App *fiber.App
}

func Bootstrap(conf *BootstrapConfig) {

	// Global repositories
	seqRepo := infrastructure.NewNumberSequenceRepository(conf.DB)

	// Shared dependencies
	validator := validation.NewValidatorHelper(conf.DB)

	// harus di ubah bila ORM diganti,
	// pastikan interface UnitOfWork tetap sama agar tidak perlu mengubah bagian lain di aplikasi
	uow := transaction.NewGormUnitOfWork(conf.DB)
	// End of shared dependencies

	// Create container
	handlerContainer := &routes.HandlerContainer{
		DB:         conf.DB,
		SeqRepo:    seqRepo,
		Validator:  validator,
		UnitOfWork: uow,
	}

	routeContainer := &routes.RouteContainer{
		App:            conf.App,
		AuthMiddleware: middleware.NewAuthMiddleware(),
	}

	// Setup global middleware
	setupGlobalMiddleware(routeContainer.App)

	// Auto-register all route modules
	routes.RegisterAllRoutes(*handlerContainer, *routeContainer)
}

func setupGlobalMiddleware(app *fiber.App) {
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Content-Type, Authorization",
		AllowCredentials: true,
	}))
}
