package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
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
