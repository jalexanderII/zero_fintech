package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jalexanderII/zero_fintech/bff/client"
)

func Login(authClient *client.AuthClient) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		type LoginData struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		var input LoginData

		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Error on login request", "data": err})
		}
		authClient.Username = input.Username
		authClient.Password = input.Password

		token, err := authClient.Login()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error on login request", "data": err})
		}

		// create a cookie to authenticate user
		_ = CreateCookie(c, authClient.Username, authClient.Interceptor.AccessToken)

		return c.JSON(fiber.Map{"status": "success", "message": "Success login", "data": token})
	}
}

func Logout(authClient *client.AuthClient) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// create a cookie to authenticate user
		DeleteCookie(c, authClient.Interceptor.AccessToken)

		return c.JSON(fiber.Map{"status": "success", "message": "Success logout", "data": authClient.Username})
	}
}

// SignUp gets username, email, and password from request body, writes it to an AuthClient and then calls SignUp
func SignUp(authClient *client.AuthClient) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		type SignUpData struct {
			Username string `json:"username"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		var input SignUpData
		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err})
		}
		authClient.Username = input.Username
		authClient.Email = input.Email
		authClient.Password = input.Password

		token, err := authClient.SignUp()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error on sign up request", "data": err})
		}

		return c.JSON(fiber.Map{"status": "success", "message": "Success SignUp", "data": token})
	}
}
