package handler

import (
	"context"
	"fiber-clean-transaction/internal/contextkeys"
	"fiber-clean-transaction/internal/dto"
	"fiber-clean-transaction/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

type MenuHandler struct {
	MenuUsecase *usecase.MenuUsecase
}

func NewMenuHandler(uc *usecase.MenuUsecase) *MenuHandler {
	return &MenuHandler{
		MenuUsecase: uc,
	}
}

func (h *MenuHandler) GetMenusByRole(c *fiber.Ctx) error {
	userClaims := contextkeys.GetUserC(c.UserContext())
	if userClaims == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	roleName := userClaims.Role
	if roleName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing role name",
		})
	}

	menus, err := h.MenuUsecase.GetMenusByRole(context.Background(), roleName)
	if err != nil {
		return ResponseError(c, err)
	}

	// Convert to response DTO
	menuResponses := make([]dto.NuxtMenuResponse, 0, len(menus))
	for _, menu := range menus {
		menuResponses = append(menuResponses, dto.NuxtMenuResponse{
			Type:   menu.Type,
			Code:   menu.Code,
			Title:  menu.Title,
			Icon:   menu.Icon,
			Parent: menu.Parent,
			To:     menu.To,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"status":  fiber.StatusOK,
		"message": "Success",
		"data":    menuResponses,
	})
}
