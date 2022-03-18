package shared

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stripe/stripe-go/v72"
)

// CreateCookie makes a valid httponly cookie
func CreateCookie(c *fiber.Ctx, name, value string) {
	cookie := &fiber.Cookie{
		Name:     name,
		Value:    value,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
	}

	// Set cookie
	c.Cookie(cookie)
}

// DeleteCookie removes existing cookie
func DeleteCookie(c *fiber.Ctx, name string) {
	c.Cookie(&fiber.Cookie{
		Name: name,
		// Set expiry date to the past
		Expires:  time.Now().Add(-(time.Hour * 2)),
		HTTPOnly: true,
	})
}

// GetPlaidErrorCode will get the error code from the error message and return it as a string
func GetPlaidErrorCode(err error) string {
	errorMessage := err.Error()

	// first get the index of the substring code
	start := strings.Index(errorMessage, ", code: ") + 8

	// get the end by creating a substring and getting the index of the first comma
	end := strings.Index(errorMessage[start:], ", ") + start

	// return the substring with the window of indexes
	return errorMessage[start:end]
}

// GetStripeErrorCode will get the error code from the error message and return it as a string
func GetStripeErrorCode(err error) string {
	// Try to safely cast a generic error to a stripe.Error so that we can get at
	// some additional Stripe-specific information about what went wrong.
	if stripeErr, ok := err.(*stripe.Error); ok {
		// The Code field will contain a basic identifier for the failure.
		switch stripeErr.Code {
		case stripe.ErrorCodeCardDeclined:
		case stripe.ErrorCodeExpiredCard:
		case stripe.ErrorCodeIncorrectCVC:
		case stripe.ErrorCodeIncorrectZip:
		}

		// The Err field can be coerced to a more specific error type with a type
		// assertion. This technique can be used to get more specialized
		// information for certain errors.
		if cardErr, ok := stripeErr.Err.(*stripe.CardError); ok {
			return fmt.Sprintf("Card was declined with code: %v\n", cardErr.DeclineCode)
		} else {
			return fmt.Sprintf("Other Stripe error occurred: %v\n", stripeErr.Error())
		}
	}

	return fmt.Sprintf("Other error occurred: %v\n", err.Error())
}
