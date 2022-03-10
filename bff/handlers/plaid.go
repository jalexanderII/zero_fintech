package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/jalexanderII/zero_fintech/bff/client"
	"github.com/jalexanderII/zero_fintech/bff/models"
)

func CreateLinkToken(plaidClient *client.PlaidClient, ctx context.Context) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		type Input struct {
			Username string `json:"username"`
		}
		var input Input
		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
		}

		linkToken, err := plaidClient.LinkTokenCreate(ctx, input.Username, c.Cookies(input.Username))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
		}

		return c.JSON(fiber.Map{"status": "success", "message": "Successfully received link token from plaid", "link_token": linkToken})
	}
}

func UpdatePlaidToken(plaidClient *client.PlaidClient, ctx context.Context) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		_, err := plaidClient.UpdatePlaidToken(ctx, "", models.Token{})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
		}
		return nil
	}
}

func ExchangePublicToken(plaidClient *client.PlaidClient, ctx context.Context) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		_, err := plaidClient.ExchangePublicToken(ctx, "", models.Token{})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
		}
		return nil
	}
}

func GetAccountDetails(plaidClient *client.PlaidClient, ctx context.Context) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		_, err := plaidClient.GetAccountDetails(ctx, "", models.Token{})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
		}
		return nil
	}
}
