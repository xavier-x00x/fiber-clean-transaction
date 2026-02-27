package handler

import (
	"fiber-clean-transaction/internal/dto"
	"fiber-clean-transaction/internal/usecase"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type UnitHandler struct {
	Usecase *usecase.UnitUsecase
}

func NewUnitHandler(uc *usecase.UnitUsecase) *UnitHandler {
	return &UnitHandler{
		Usecase: uc,
	}
}

func (h *UnitHandler) GetAllFilter(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	metaRequest := &dto.MetaRequest{
		Page:        page,
		Limit:       limit,
		Search:      c.Query("search", ""),
		OrderColumn: c.Query("order_column", "id"),
		OrderDir:    c.Query("order_dir", "asc"),
	}

	data, meta, err := h.Usecase.GetAllFilter(c.UserContext(), metaRequest)
	if err != nil {
		return ResponseError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"status":  fiber.StatusOK,
		"message": "Success",
		"data":    data,
		"meta":    meta,
	})
}
