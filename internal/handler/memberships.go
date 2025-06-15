package handler

import (
	"music_catalog/internal/helper"
	"music_catalog/internal/models"

	"github.com/gofiber/fiber/v2"
)

type MembershipsUc interface {
	SignIn(req models.SignReq) error
}

type memberships struct {
	usecase MembershipsUc
	rg      fiber.Router
}

func NewMembershipsHandler(usecase MembershipsUc, rg fiber.Router) *memberships {
	return &memberships{
		usecase: usecase,
		rg:      rg,
	}
}

func (h *memberships) SetupRoutes() {
	h.rg.Post(helper.SignIn, h.SignIn)
}

func (h *memberships) SignIn(c *fiber.Ctx) error {

	var payload models.SignReq

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Invalid request body",
		})
	}

	err := h.usecase.SignIn(payload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  fiber.StatusCreated,
		"message": "Success",
	})
}
