package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jalexanderII/zero_fintech/services/core/gen/core"
)

// Login get user and password
func Login(c *fiber.Ctx) error {
	ctx, cancel := NewClientContext(timeout)
	defer cancel()

	type LoginData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var input LoginData

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Error on login request", "data": err})
	}

	resp, err := coreClient.Login(ctx, &core.LoginRequest{Username: input.Username, Password: input.Password})
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Success login", "data": resp.GetToken()})
}

func SignUp(c *fiber.Ctx) error {
	ctx, cancel := NewClientContext(timeout)
	defer cancel()

	type UserData struct {
		Username string `json:"username"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password"`
	}
	var input UserData
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	resp, err := coreClient.SignUp(ctx, &core.SignupRequest{Username: input.Username, Email: input.Email, Password: input.Password})
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Success SignUp", "data": resp.GetToken()})
}
