package handlers

import (
	"context"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/jalexanderII/zero_fintech/bff/client"
	"github.com/jalexanderII/zero_fintech/bff/models"
	"github.com/jalexanderII/zero_fintech/gen/Go/core"
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

		linkToken, err := plaidClient.LinkTokenCreate(ctx, input.Username)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
		}

		CreateCookie(c, "link_token", linkToken)

		return c.JSON(fiber.Map{"status": "success", "message": "Successfully received link token from plaid", "link_token": linkToken})
	}
}

// Link should be accessed after createLinkToken so that a link token can be set in cookies
func Link(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"LinkToken": c.Cookies("link_token", ""),
	})
}

func ExchangePublicToken(plaidClient *client.PlaidClient, ctx context.Context) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		type Input struct {
			Username    string `json:"username"`
			PublicToken string `json:"public_token"`
			TokenId     string `json:"tokenId,omitempty"`
		}
		var input Input
		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
		}

		token, err := plaidClient.ExchangePublicToken(ctx, input.PublicToken)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
		}

		if c.Method() == http.MethodPost {
			if err = plaidClient.SaveToken(ctx, token); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Failure to create access token", "data": err})
			}
		} else {
			if err = plaidClient.UpdateToken(ctx, input.TokenId, token.Value, token.ItemId); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Failure to update access token", "data": err})
			}
		}

		err = GetandSaveAccountDetails(plaidClient, ctx, token)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
		}

		return c.JSON(fiber.Map{"status": "success", "message": "Access token created successfully", "token": token})
	}
}

func GetandSaveAccountDetails(plaidClient *client.PlaidClient, ctx context.Context, token *models.Token) error {
	accountDetails, err := plaidClient.GetAccountDetails(ctx, token)
	if err != nil {
		return err
	}

	accounts := accountDetails.GetAccounts()
	transactions := accountDetails.GetTransactions()

	for _, account := range accounts {
		req := &core.CreateAccountRequest{Account: account}
		_, err = plaidClient.CoreClient.CreateAccount(ctx, req)
		if err != nil {
			return err
		}
	}

	for _, transaction := range transactions {
		req := &core.CreateTransactionRequest{Transaction: transaction}
		_, err = plaidClient.CoreClient.CreateTransaction(ctx, req)
		if err != nil {
			return err
		}
	}

	return nil
}
