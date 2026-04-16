package handler

import (
	"fiber-clean-transaction/internal/contextkeys"
	"fiber-clean-transaction/internal/domain/entity"
	"fiber-clean-transaction/internal/dto"
	"fiber-clean-transaction/internal/usecase"
	"fiber-clean-transaction/pkg/utils"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type RoleHandler struct {
	RoleUsecase *usecase.RoleUsecase
}

func NewRoleHandler(u *usecase.RoleUsecase) *RoleHandler {
	return &RoleHandler{
		RoleUsecase: u,
	}
}

func (h *RoleHandler) GetAllFilter(c *fiber.Ctx) error {

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	metaRequest := &dto.MetaRequest{
		Page:        page,
		Limit:       limit,
		Search:      c.Query("search", ""),
		OrderColumn: c.Query("order_column", "id"),
		OrderDir:    c.Query("order_dir", "asc"),
	}

	data, meta, err := h.RoleUsecase.GetAllFilter(c.UserContext(), metaRequest)
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

func (h *RoleHandler) GetRole(c *fiber.Ctx) error {

	ID, errx := strconv.Atoi(c.Params("id"))
	if errx != nil {
		return ResponseError(c, utils.BadRequest("Invalid role ID"))
	}

	data, err := h.RoleUsecase.FindByID(c.UserContext(), uint(ID))
	if err != nil {
		return ResponseError(c, err)
	}

	role := &dto.RoleResponse{
		ID:   int(data.ID),
		Name: data.Name,
		Permissions: func(perms []entity.Permission) []string {
			codes := make([]string, len(perms))
			for i, perm := range perms {
				codes[i] = perm.Name
			}
			return codes
		}(data.Permissions),
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"status":  fiber.StatusOK,
		"message": "Success",
		"data":    role,
	})
}

func (h *RoleHandler) CreateRole(c *fiber.Ctx) error {

	roleRequest := new(dto.RoleRequest)

	if err := c.BodyParser(&roleRequest); err != nil {
		error := utils.BadRequest(err.Error())
		return ResponseError(c, error)
	}

	err := h.RoleUsecase.Create(c.UserContext(), roleRequest)
	if err != nil {
		return ResponseError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"status":  fiber.StatusCreated,
		"message": "Data created successfully",
	})
}

func (h *RoleHandler) UpdateRole(c *fiber.Ctx) error {

	ID, errx := c.ParamsInt("id")
	if errx != nil {
		return ResponseError(c, utils.BadRequest("Invalid role ID"))
	}

	roleRequest := new(dto.RoleRequest)

	if err := c.BodyParser(&roleRequest); err != nil {
		errx := utils.BadRequest(err.Error())
		return ResponseError(c, errx)
	}

	err := h.RoleUsecase.Update(c.UserContext(), uint(ID), roleRequest)

	if err != nil {
		return ResponseError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"status":  fiber.StatusOK,
		"message": "Data updated successfully",
	})
}

func (h *RoleHandler) DeleteRole(c *fiber.Ctx) error {

	ID, errx := c.ParamsInt("id")
	if errx != nil {
		return ResponseError(c, utils.BadRequest("Invalid role ID"))
	}

	err := h.RoleUsecase.Delete(c.UserContext(), uint(ID))

	if err != nil {
		return ResponseError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"status":  fiber.StatusOK,
		"message": "Data deleted successfully",
	})
}

func (h *RoleHandler) Authorization(c *fiber.Ctx) error {
	type Request struct {
		Path string `json:"path"`
	}

	userClaims := contextkeys.GetUserC(c.UserContext())
	log.Printf("xxx : %v", userClaims.Role)
	if userClaims == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Path == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing path",
		})
	}

	println(req.Path)
	access, err := h.RoleUsecase.Authorization(userClaims.Role, req.Path)

	if err != nil {
		return ResponseError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"status":  fiber.StatusOK,
		"message": "Success",
		"data":    access,
	})

}

func (h *RoleHandler) GetPermissionsByRole(c *fiber.Ctx) error {

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
	println(roleName)
	permissions, err := h.RoleUsecase.GetPermissionsByRole(c.UserContext(), roleName)

	if err != nil {
		return ResponseError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"status":  fiber.StatusOK,
		"message": "Success",
		"data":    permissions,
	})
}
