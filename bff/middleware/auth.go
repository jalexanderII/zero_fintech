package middleware

import (
	"github.com/gofiber/fiber/v2"
	jwtMiddleware "github.com/gofiber/jwt/v2"
	"github.com/jalexanderII/zero_fintech/services/core/config"
)

// Protected func for specify routes group with JWT authentication.
// See: https://github.com/gofiber/jwt
func Protected() func(*fiber.Ctx) error {
	// Create config for JWT authentication middleware.
	jconfig := jwtMiddleware.Config{
		SigningKey:   []byte(config.GetEnv("JWT_SECRET_KEY")),
		ContextKey:   "jwt", // used in private routes
		ErrorHandler: jwtError,
	}

	return jwtMiddleware.New(jconfig)
}

func jwtError(c *fiber.Ctx, err error) error {
	// Return status 401 and failed authentication error.
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Return status 401 and failed authentication error.
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": true,
		"msg":   err.Error(),
	})
}
