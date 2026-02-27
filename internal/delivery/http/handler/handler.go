package handler

import (
	"errors"
	"fiber-clean-transaction/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

func ResponseError(c *fiber.Ctx, err error) error {
	var domainErr *utils.DomainError
	if errors.As(err, &domainErr) {

		if domainErr.StatusCode == fiber.StatusUnprocessableEntity {
			return c.Status(domainErr.StatusCode).JSON(fiber.Map{
				"success": false,
				"status":  domainErr.StatusCode,
				"message": domainErr.Message,
				"error":   domainErr.Code,
				"errors":  domainErr.Errors,
			})
		}

		if domainErr.StatusCode != fiber.StatusInternalServerError {
			return c.Status(domainErr.StatusCode).JSON(fiber.Map{
				"success": false,
				"status":  domainErr.StatusCode,
				"message": domainErr.Message,
				"error":   domainErr.Code,
			})
		}
	}

	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"success": false,
		"status":  fiber.StatusInternalServerError,
		"message": "Something went wrong",
		"error":   err.Error(),
	})
}
