package response

import "github.com/gofiber/fiber/v2"

func InternalServerError(c *fiber.Ctx) error {
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"message": "internal server error",
	}) 
}

func NotFound(c *fiber.Ctx, msg string) error {
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"message": msg,
	})
}
