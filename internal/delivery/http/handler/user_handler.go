package handler

import (
	"fiber-clean-transaction/internal/dto"
	"fiber-clean-transaction/internal/usecase"
	"fiber-clean-transaction/pkg/utils"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	UserUsecase *usecase.UserUsecase
}

func NewUserHandler(uc *usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		UserUsecase: uc,
	}
}

func (h *UserHandler) GetAllFilter(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	metaRequest := &dto.MetaRequest{
		Page:        page,
		Limit:       limit,
		Search:      c.Query("search", ""),
		OrderColumn: c.Query("order_column", "id"),
		OrderDir:    c.Query("order_dir", "asc"),
	}

	data, meta, err := h.UserUsecase.GetAllFilter(metaRequest)
	if err != nil {
		return ResponseError(c, err)
	}

	var responses []dto.UserResponse
	for _, user := range data {
		updatedAt := time.Time{}
		if user.UpdatedAt != nil {
			updatedAt = *user.UpdatedAt
		}
		responses = append(responses, dto.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Username:  user.Username,
			Email:     user.Email,
			Role:      user.Role,
			Avatar:    user.Avatar,
			UpdatedAt: updatedAt,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"status":  fiber.StatusOK,
		"message": "Success",
		"data":    responses,
		"meta":    meta,
	})
}

func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	ID, errx := strconv.Atoi(c.Params("id"))
	if errx != nil {
		return ResponseError(c, utils.BadRequest("Invalid user ID"))
	}

	data, err := h.UserUsecase.Profile(uint(ID))
	if err != nil {
		return ResponseError(c, err)
	}

	updatedAt := time.Time{}
	if data.UpdatedAt != nil {
		updatedAt = *data.UpdatedAt
	}
	response := &dto.UserResponse{
		ID:        data.ID,
		Name:      data.Name,
		Username:  data.Username,
		Email:     data.Email,
		Role:      data.Role,
		Avatar:    data.Avatar,
		UpdatedAt: updatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"status":  fiber.StatusOK,
		"message": "Success",
		"data":    response,
	})
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	ID, errx := c.ParamsInt("id")
	if errx != nil {
		return ResponseError(c, utils.BadRequest("Invalid user ID"))
	}

	userRequest := new(dto.UserUpdateRequest)

	if err := c.BodyParser(&userRequest); err != nil {
		errx := utils.BadRequest(err.Error())
		return ResponseError(c, errx)
	}

	err := h.UserUsecase.Update(c.UserContext(), uint(ID), userRequest)
	if err != nil {
		return ResponseError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"status":  fiber.StatusOK,
		"message": "Data updated successfully",
	})
}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	ID, errx := c.ParamsInt("id")
	if errx != nil {
		return ResponseError(c, utils.BadRequest("Invalid user ID"))
	}

	err := h.UserUsecase.Delete(c.UserContext(), uint(ID))
	if err != nil {
		return ResponseError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"status":  fiber.StatusOK,
		"message": "Data deleted successfully",
	})
}
