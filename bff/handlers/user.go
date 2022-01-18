package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/jalexanderII/zero_fintech/gen/Go/core"
)

// User To be used as a serializer
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email" validate:"required,email"`
}

// CreateResponseUser Takes in a model and returns a serializer
func CreateResponseUser(userModel *core.User) User {
	return User{ID: userModel.Id, Username: userModel.Username, Email: userModel.Email}
}

func ListUsers(client core.CoreClient, ctx context.Context) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		listUsers, err := client.ListUsers(ctx, &core.ListUserRequest{})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
		}

		responseUsers := make([]User, len(listUsers.GetUsers()))
		for idx, user := range listUsers.GetUsers() {
			responseUsers[idx] = CreateResponseUser(user)
		}

		return c.Status(fiber.StatusOK).JSON(responseUsers)
	}
}

func GetUser(client core.CoreClient, ctx context.Context) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		getUser, err := client.GetUser(ctx, &core.GetUserRequest{Id: c.Params("id")})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
		}

		return c.Status(fiber.StatusOK).JSON(CreateResponseUser(getUser))
	}
}

func UpdateUser(client core.CoreClient, ctx context.Context) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		type UpdateUserResponse struct {
			Username string `json:"username"`
			Email    string `json:"email"`
		}

		var updateUserResponse UpdateUserResponse
		if err := c.BodyParser(&updateUserResponse); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
		}

		updateUser, err := client.UpdateUser(ctx, &core.UpdateUserRequest{Id: c.Params("id"), User: &core.User{
			Username: updateUserResponse.Username,
			Email:    updateUserResponse.Email,
		}})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
		}

		return c.Status(fiber.StatusOK).JSON(CreateResponseUser(updateUser))
	}
}

func DeleteUser(client core.CoreClient, ctx context.Context) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		response, err := client.DeleteUser(ctx, &core.DeleteUserRequest{Id: c.Params("id")})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": response.GetStatus(), "data": CreateResponseUser(response.GetUser())})
	}
}
