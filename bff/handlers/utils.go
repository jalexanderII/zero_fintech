package handlers

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// CreateCookie makes a valid httponly cookie
func CreateCookie(c *fiber.Ctx, name, value string) error {
	cookie := &fiber.Cookie{
		Name:     name,
		Value:    value,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
	}

	// Set cookie
	c.Cookie(cookie)
	return nil
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
