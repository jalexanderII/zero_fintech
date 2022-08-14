package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/jalexanderII/zero_fintech/bff/client"
	"github.com/jalexanderII/zero_fintech/bff/shared"
)

func Login(authClient *client.AuthClient) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		type LoginData struct {
			Username string `json:"username"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		var input LoginData

		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Error parsing request body", "data": err})
		}
		authClient.Username = input.Username
		authClient.Email = input.Email
		authClient.Password = input.Password

		resp, err := authClient.Login()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error on login request", "data": err})
		}

		// create a cookie to authenticate user
		shared.CreateCookie(c, "AuthToken", resp.GetToken())
		shared.CreateCookie(c, input.Username, resp.GetUserId())
		fmt.Printf("Current Cookies UserId: %v\n", c.Cookies(input.Username))

		return c.JSON(fiber.Map{"status": "success", "message": "Success login", "data": resp})
	}
}

func Logout(authClient *client.AuthClient) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		type LogoutData struct {
			Username string `json:"username"`
		}
		var input LogoutData
		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Error on login request", "data": err})
		}

		// create a cookie to authenticate user
		shared.DeleteCookie(c, "AuthToken")
		shared.DeleteCookie(c, input.Username)
		shared.DeleteCookie(c, fmt.Sprintf("%v_link_token", input.Username))

		return c.JSON(fiber.Map{"status": "success", "message": "Success logout", "data": input.Username})
	}
}

// SignUp gets username, email, and password from request body, writes it to an AuthClient and then calls SignUp
func SignUp(authClient *client.AuthClient) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		type SignUpData struct {
			Username    string `json:"username"`
			Email       string `json:"email"`
			Password    string `json:"password"`
			PhoneNumber string `json:"phone_number"`
		}
		var input SignUpData
		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Error parsing request body", "data": err})
		}
		authClient.Username = input.Username
		authClient.Email = input.Email
		authClient.Password = input.Password
		authClient.PhoneNumber = input.PhoneNumber

		resp, err := authClient.SignUp()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error on sign up request", "data": err})
		}

		// create a cookie to authenticate user
		shared.CreateCookie(c, "AuthToken", resp.GetToken())
		shared.CreateCookie(c, authClient.Username, resp.GetUserId())
		fmt.Printf("Current Cookies UserId: %v\n", c.Cookies(input.Username))

		return c.JSON(fiber.Map{"status": "success", "message": "Success SignUp", "data": resp})
	}
}
