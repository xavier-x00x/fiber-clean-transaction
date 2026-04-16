package handler

import (
	"fiber-clean-transaction/internal/dto"
	"fiber-clean-transaction/internal/usecase"
	"fiber-clean-transaction/pkg/utils"
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

func (h *UnitHandler) GetUnit(c *fiber.Ctx) error {
	ID, errx := strconv.Atoi(c.Params("id"))
	if errx != nil {
		return ResponseError(c, utils.BadRequest("Invalid unit id"))
	}

	data, err := h.Usecase.FindByID(c.UserContext(), uint(ID))
	if err != nil {
		return ResponseError(c, err)
	}

	unit := &dto.UnitResponse{
		ID:        data.ID,
		Code:      data.Code,
		Name:      data.Name,
		Status:    data.Status,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"status":  fiber.StatusOK,
		"message": "Success",
		"data":    unit,
	})
}

func (h *UnitHandler) CreateUnit(c *fiber.Ctx) error {
	unitRequest := new(dto.UnitRequest)

	if err := c.BodyParser(&unitRequest); err != nil {
		return ResponseError(c, utils.BadRequest(err.Error()))
	}

	err := h.Usecase.Create(c.UserContext(), unitRequest)
	if err != nil {
		return ResponseError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"status":  fiber.StatusCreated,
		"message": "Data created successfully",
	})
}

func (h *UnitHandler) UpdateUnit(c *fiber.Ctx) error {
	ID, errx := c.ParamsInt("id")
	if errx != nil {
		return ResponseError(c, utils.BadRequest("Invalid unit id"))
	}

	unitRequest := new(dto.UnitRequest)

	if err := c.BodyParser(&unitRequest); err != nil {
		return ResponseError(c, utils.BadRequest(err.Error()))
	}

	err := h.Usecase.Update(c.UserContext(), uint(ID), unitRequest)
	if err != nil {
		return ResponseError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"status":  fiber.StatusOK,
		"message": "Data updated successfully",
	})
}

func (h *UnitHandler) DeleteUnit(c *fiber.Ctx) error {
	ID, errx := c.ParamsInt("id")
	if errx != nil {
		return ResponseError(c, utils.BadRequest("Invalid unit id"))
	}

	err := h.Usecase.Delete(c.UserContext(), uint(ID))
	if err != nil {
		return ResponseError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"status":  fiber.StatusOK,
		"message": "Data deleted successfully",
	})
}
