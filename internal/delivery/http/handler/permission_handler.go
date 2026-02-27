package handler

import (
	"fiber-clean-transaction/internal/dto"
	"fiber-clean-transaction/internal/usecase"
	"fiber-clean-transaction/pkg/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type PermissionHandler struct {
	PermissionUsecase *usecase.PermissionUsecase
}

func NewPermissionHandler(store *usecase.PermissionUsecase) *PermissionHandler {
	return &PermissionHandler{
		PermissionUsecase: store,
	}
}

func (h *PermissionHandler) GetAllFilter(c *fiber.Ctx) error {

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	metaRequest := &dto.MetaRequest{
		Page:        page,
		Limit:       limit,
		Search:      c.Query("search", ""),
		OrderColumn: c.Query("order_column", "id"),
		OrderDir:    c.Query("order_dir", "asc"),
	}

	data, meta, err := h.PermissionUsecase.GetAllFilter(c.UserContext(), metaRequest)
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

func (h *PermissionHandler) GetPermission(c *fiber.Ctx) error {

	id, errx := strconv.Atoi(c.Params("id"))
	if errx != nil {
		return ResponseError(c, utils.BadRequest("Invalid permission ID"))
	}

	data, err := h.PermissionUsecase.FindById(c.UserContext(), uint(id))
	if err != nil {
		return ResponseError(c, err)
	}

	store := &dto.PermissionResponse{
		ID:        data.ID,
		Path:      data.Path,
		Name:      data.Name,
		UpdatedAt: data.UpdatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"status":  fiber.StatusOK,
		"message": "Success",
		"data":    store,
	})
}

func (h *PermissionHandler) CreatePermission(c *fiber.Ctx) error {

	permissionRequest := new(dto.PermissionRequest)

	if err := c.BodyParser(&permissionRequest); err != nil {
		error := utils.BadRequest(err.Error())
		return ResponseError(c, error)
	}

	err := h.PermissionUsecase.Create(c.UserContext(), permissionRequest)
	if err != nil {
		return ResponseError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"status":  fiber.StatusCreated,
		"message": "Data created successfully",
	})
}

func (h *PermissionHandler) UpdatePermission(c *fiber.Ctx) error {

	id, _ := c.ParamsInt("id")

	permissionRequest := new(dto.PermissionRequest)

	if err := c.BodyParser(&permissionRequest); err != nil {
		errx := utils.BadRequest(err.Error())
		return ResponseError(c, errx)
	}

	err := h.PermissionUsecase.Update(c.UserContext(), uint(id), permissionRequest)

	if err != nil {
		return ResponseError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"status":  fiber.StatusOK,
		"message": "Data updated successfully",
	})
}

func (h *PermissionHandler) DeletePermission(c *fiber.Ctx) error {

	id, _ := c.ParamsInt("id")

	err := h.PermissionUsecase.Delete(c.UserContext(), uint(id))

	if err != nil {
		return ResponseError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"status":  fiber.StatusOK,
		"message": "Data deleted successfully",
	})
}
