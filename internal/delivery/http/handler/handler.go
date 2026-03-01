package handler

import (
	"errors"
	"fiber-clean-transaction/pkg/utils"
	"fiber-clean-transaction/pkg/validation"

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

	var validationErr *validation.ResponseError
	if errors.As(err, &validationErr) {
		return c.Status(validationErr.StatusCode).JSON(fiber.Map{
			"success": false,
			"status":  validationErr.StatusCode,
			"message": validationErr.Message,
			"error":   validationErr.Code,
			"errors":  validationErr.Errors,
		})
	}

	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"success": false,
		"status":  fiber.StatusInternalServerError,
		"message": "Something went wrong",
		"error":   err.Error(),
	})
}
