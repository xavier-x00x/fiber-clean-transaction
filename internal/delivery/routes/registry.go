package routes

import (
	"fiber-clean-transaction/internal/domain/repository"
	"log"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Container interface for dependency injection
type HandlerContainer struct {
	DB      *gorm.DB
	SeqRepo repository.NumberSequenceRepository
}

type RouteContainer struct {
	App            *fiber.App
	AuthMiddleware fiber.Handler
}

// RouteModule interface that all route modules must implement
type RouteModule interface {
	GetModuleName() string
	RegisterHandler(c HandlerContainer)
	RegisterRoutes(c RouteContainer)
}

var registeredModules []RouteModule

// RegisterModule registers a route module
func RegisterModule(module RouteModule) {
	registeredModules = append(registeredModules, module)
}

// RegisterAllRoutes automatically registers all registered modules
func RegisterAllRoutes(hc HandlerContainer, rc RouteContainer) {

	for _, module := range registeredModules {
		log.Printf("Registering route module: %s", module.GetModuleName())
		module.RegisterHandler(hc)
		module.RegisterRoutes(rc)
		log.Printf("Route module registered: %s", module.GetModuleName())
	}
}

// Auto-register modules in init functions
func init() {
	// Modules will register themselves
	RegisterModule(&AuthRoutes{})
	RegisterModule(&StoreRoutes{})
	RegisterModule(&CategoryRoutes{})
	RegisterModule(&PermissionRoutes{})
	RegisterModule(&RoleRoutes{})
}
