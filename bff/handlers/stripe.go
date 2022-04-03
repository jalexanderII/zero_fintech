package handlers

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/customer"
	"github.com/stripe/stripe-go/v72/paymentsource"
)

// VerifyParams Alias for deserializing json
type VerifyParams struct {
	Amounts     [2]int64 `json:"amounts"`
	Customer    string   `json:"customer"`
	BankAccount string   `json:"bank_account"`
}

func verifyHandler(c *fiber.Ctx) (err error) {
	params := &VerifyParams{}
	if err := c.BodyParser(&params); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}
	verifyParams := &stripe.SourceVerifyParams{
		Amounts:  params.Amounts,
		Customer: stripe.String(params.Customer),
	}
	ba, err := paymentsource.Verify(params.BankAccount, verifyParams)
	if err != nil {
		if stripeErr, ok := err.(*stripe.Error); ok {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Failure to verify payment source", "data": stripeErr})
		}
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": ba})
}

// CreateCustomerParams Alias for deserializing json
type CreateCustomerParams struct {
	BankAccount string `json:"bank_account"`
}

func createCustomerHandler(c *fiber.Ctx) (err error) {
	params := &CreateCustomerParams{}
	if err := c.BodyParser(&params); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}
	cusParams := &stripe.CustomerParams{}
	err = cusParams.SetSource(params.BankAccount)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Failure setting customer stripe bank account", "data": err})
	}
	cus, err := customer.New(cusParams)
	if err != nil {
		if stripeErr, ok := err.(*stripe.Error); ok {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Failure creating new stripe customer", "data": stripeErr})
		}
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": cus})
}

// ExchangeTokenParams Alias for deserializing json
type ExchangeTokenParams struct {
	PublicToken string `json:"public_token"`
	AccountID   string `json:"account_id"`
}

// PublicKeys Alias for deserializing json
type PublicKeys struct {
	StripeKey string `json:"stripe_key"`
	PlaidKey  string `json:"plaid_key"`
}

func publicKeyHandler(c *fiber.Ctx) (err error) {
	data := PublicKeys{
		StripeKey: os.Getenv("STRIPE_SECRET_KEY"),
		PlaidKey:  os.Getenv("PLAID_PUBLIC_KEY"),
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": data})
}
