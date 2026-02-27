package handler

import (
	"fiber-clean-transaction/internal/dto"
	"fiber-clean-transaction/internal/usecase"
	"fiber-clean-transaction/pkg/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type CategoryHandler struct {
	CategoryUsecase *usecase.CategoryUsecase
}

func NewCategoryHandler(uc *usecase.CategoryUsecase) *CategoryHandler {
	return &CategoryHandler{
		CategoryUsecase: uc,
	}
}

func (h *CategoryHandler) GetAllFilter(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	metaRequest := &dto.MetaRequest{
		Page:        page,
		Limit:       limit,
		Search:      c.Query("search", ""),
		OrderColumn: c.Query("order_column", "id"),
		OrderDir:    c.Query("order_dir", "asc"),
	}

	data, meta, err := h.CategoryUsecase.GetAllFilter(c.UserContext(), metaRequest)
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

func (h *CategoryHandler) GetCategory(c *fiber.Ctx) error {
	id, errx := strconv.Atoi(c.Params("id"))
	if errx != nil {
		return ResponseError(c, utils.BadRequest("Invalid category ID"))
	}

	data, err := h.CategoryUsecase.FindById(c.UserContext(), uint(id))
	if err != nil {
		return ResponseError(c, err)
	}

	category := &dto.CategoryResponse{
		Id:        data.Id,
		Code:      data.Code,
		Name:      data.Name,
		Status:    data.Status,
		UpdatedAt: data.UpdatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"status":  fiber.StatusOK,
		"message": "Success",
		"data":    category,
	})
}

func (h *CategoryHandler) CreateCategory(c *fiber.Ctx) error {
	categoryRequest := new(dto.CategoryRequest)
	if err := c.BodyParser(categoryRequest); err != nil {
		error := utils.BadRequest(err.Error())
		return ResponseError(c, error)
	}

	err := h.CategoryUsecase.Create(c.UserContext(), categoryRequest)
	if err != nil {
		return ResponseError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"status":  fiber.StatusCreated,
		"message": "Data created successfully",
	})
}

func (h *CategoryHandler) UpdateCategory(c *fiber.Ctx) error {
	id, errx := strconv.Atoi(c.Params("id"))
	if errx != nil {
		return ResponseError(c, utils.BadRequest("Invalid category ID"))
	}

	categoryRequest := new(dto.CategoryRequest)
	if err := c.BodyParser(categoryRequest); err != nil {
		error := utils.BadRequest(err.Error())
		return ResponseError(c, error)
	}

	err := h.CategoryUsecase.Update(c.UserContext(), uint(id), categoryRequest)
	if err != nil {
		return ResponseError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"status":  fiber.StatusOK,
		"message": "Data updated successfully",
	})
}

func (h *CategoryHandler) DeleteCategory(c *fiber.Ctx) error {
	id, errx := strconv.Atoi(c.Params("id"))
	if errx != nil {
		return ResponseError(c, utils.BadRequest("Invalid category ID"))
	}

	err := h.CategoryUsecase.Delete(c.UserContext(), uint(id))
	if err != nil {
		return ResponseError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"status":  fiber.StatusOK,
		"message": "Data deleted successfully",
	})
}
