package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/jalexanderII/zero_fintech/bff/client"
	"github.com/jalexanderII/zero_fintech/bff/models"
	"github.com/jalexanderII/zero_fintech/bff/shared"
	"github.com/jalexanderII/zero_fintech/gen/Go/core"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Link will call CreateLinkToken to get a link token, and then call ExchangePublicToken to get an access token
// will be saved to db along with account and transaction details upon success
func Link(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"Username": c.Params("username"),
		"Purpose":  c.Params("purpose"),
	})
}

func CreateLinkToken(plaidClient *client.PlaidClient, ctx context.Context) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		type Input struct {
			Username string         `json:"username"`
			Purpose  models.Purpose `json:"purpose"`
		}
		var input Input
		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(err.Error())
		}

		linkTokenResp, err := plaidClient.LinkTokenCreate(ctx, input.Username, input.Purpose)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Failure to create link token", "data": err})
		}

		shared.CreateCookie(c, fmt.Sprintf("%v_link_token", input.Username), linkTokenResp.Token)
		shared.CreateCookie(c, input.Username, linkTokenResp.UserId)
		id, err := primitive.ObjectIDFromHex(linkTokenResp.UserId)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Failure to get ObjectId from Hex", "data": err})
		}

		plaidClient.SetLinkToken(&models.Token{
			User:  &models.User{ID: id, Username: input.Username},
			Value: linkTokenResp.Token,
		})

		return c.JSON(fiber.Map{"status": "success", "message": "Successfully received link token from plaid", "link_token": linkTokenResp.Token})
	}
}

func ExchangePublicToken(plaidClient *client.PlaidClient, ctx context.Context) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		type Input struct {
			Username    string               `json:"username"`
			PublicToken string               `json:"public_token"`
			TokenId     string               `json:"tokenId,omitempty"`
			MetaData    models.PlaidMetaData `json:"meta_data,omitempty"`
		}
		var input Input
		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
		}
		fmt.Printf("METADATA: %+v", input.MetaData)

		user, err := plaidClient.GetUser(ctx, input.Username, "")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Failure to exchange for token", "data": err})
		}

		token, err := plaidClient.ExchangePublicToken(ctx, input.PublicToken)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Failure to exchange for token", "data": err})
		}

		token.User = user
		token.Institution = input.MetaData.Institution.Name
		token.InstitutionID = input.MetaData.Institution.InstitutionId
		dbToken, err := plaidClient.GetUserToken(ctx, user)
		if err == mongo.ErrNoDocuments || c.Method() == http.MethodPost {
			if err = plaidClient.SaveToken(ctx, token); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Failure to create access token", "data": err})
			}

			err = GetandSaveAccountDetails(plaidClient, ctx, token, c)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Failure to get and save account details", "data": err})
			}
		} else {
			if err = plaidClient.UpdateToken(ctx, dbToken.ID, token.Value, token.ItemId); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Failure to update access token", "data": err})
			}
			// TODO Add method that fetches latest not duplicate account and transaction details
		}
		return c.JSON(fiber.Map{"status": "success", "message": "Access token created successfully", "token": input})
	}
}

func GetandSaveAccountDetails(plaidClient *client.PlaidClient, ctx context.Context, token *models.Token, c *fiber.Ctx) error {
	accountDetails, err := plaidClient.GetAccountDetails(ctx, token)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Failure to get account details", "data": err})
	}

	accounts := accountDetails.GetAccounts()
	plaidAccToDBAccId := make(map[string]string)
	transactions := accountDetails.GetTransactions()

	for _, account := range accounts {
		req := &core.CreateAccountRequest{Account: account}
		dbAccount, err := plaidClient.CoreClient.CreateAccount(ctx, req)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Failure to save account account", "data": err})
		}
		plaidAccToDBAccId[dbAccount.PlaidAccountId] = dbAccount.AccountId
	}

	for _, transaction := range transactions {
		transaction.AccountId = plaidAccToDBAccId[transaction.PlaidAccountId]
		req := &core.CreateTransactionRequest{Transaction: transaction}
		_, err = plaidClient.CoreClient.CreateTransaction(ctx, req)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Failure to save account transaction", "data": err})
		}
	}

	return nil
}
