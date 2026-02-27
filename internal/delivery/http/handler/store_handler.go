package handler

import (
	"fiber-clean-transaction/internal/dto"
	"fiber-clean-transaction/internal/usecase"
	"fiber-clean-transaction/pkg/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type StoreHandler struct {
	StoreUsecase *usecase.StoreUsecase
}

func NewStoreHandler(store *usecase.StoreUsecase) *StoreHandler {
	return &StoreHandler{
		StoreUsecase: store,
	}
}

func (h *StoreHandler) GetAllFilter(c *fiber.Ctx) error {

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	metaRequest := &dto.MetaRequest{
		Page:        page,
		Limit:       limit,
		Search:      c.Query("search", ""),
		OrderColumn: c.Query("order_column", "id"),
		OrderDir:    c.Query("order_dir", "asc"),
	}

	data, meta, err := h.StoreUsecase.GetAllFilter(metaRequest)
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

func (h *StoreHandler) GetStore(c *fiber.Ctx) error {

	id, errx := strconv.Atoi(c.Params("id"))
	if errx != nil {
		return ResponseError(c, utils.BadRequest("Invalid store ID"))
	}

	data, err := h.StoreUsecase.FindById(uint(id))
	if err != nil {
		return ResponseError(c, err)
	}

	store := &dto.StoreResponse{
		ID:        data.Id,
		Code:      data.Code,
		Name:      data.Name,
		Address:   data.Address,
		Phone:     data.Phone,
		Email:     data.Email,
		Phone2:    data.Phone2,
		Email2:    data.Email2,
		Status:    data.Status,
		UpdatedAt: data.UpdatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"status":  fiber.StatusOK,
		"message": "Success",
		"data":    store,
	})
}

func (h *StoreHandler) CreateStore(c *fiber.Ctx) error {

	storeRequest := new(dto.StoreRequest)

	if err := c.BodyParser(&storeRequest); err != nil {
		error := utils.BadRequest(err.Error())
		return ResponseError(c, error)
	}

	err := h.StoreUsecase.Create(c.UserContext(), storeRequest)
	if err != nil {
		return ResponseError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"status":  fiber.StatusCreated,
		"message": "Data created successfully",
	})
}

func (h *StoreHandler) UpdateStore(c *fiber.Ctx) error {

	id, _ := c.ParamsInt("id")

	storeRequest := new(dto.StoreRequest)

	if err := c.BodyParser(&storeRequest); err != nil {
		errx := utils.BadRequest(err.Error())
		return ResponseError(c, errx)
	}

	err := h.StoreUsecase.Update(c.UserContext(), uint(id), storeRequest)

	if err != nil {
		return ResponseError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"status":  fiber.StatusOK,
		"message": "Data updated successfully",
	})
}

func (h *StoreHandler) DeleteStore(c *fiber.Ctx) error {

	id, _ := c.ParamsInt("id")

	err := h.StoreUsecase.Delete(c.UserContext(), uint(id))

	if err != nil {
		return ResponseError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"status":  fiber.StatusOK,
		"message": "Data deleted successfully",
	})
}
