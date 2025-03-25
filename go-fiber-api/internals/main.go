package main

import (
	"go_fiber_api/internals/routes"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app := fiber.New(fiber.Config{
		AppName: "Go Fiber API",
	})

	app.Mount("/apix", routes.GetAPIRoutes())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Get("/errors", func (c *fiber.Ctx) error {
		return fiber.NewError(789, "xixixi")
	})

	app.Get("/:value", func (c *fiber.Ctx) error {
		return c.SendString("vlaue: " + c.Params("value"))
	})

	app.Get("/:name?", func (c *fiber.Ctx) error {
		if c.Params("name") != "" {
			return c.SendString("Hello " + c.Params("name"))
		}
		return c.SendString("Where is john?")
	})

	app.Get("/api/*", func (c *fiber.Ctx) error {
		return c.SendString("API path: " + c.Params("*"))
	})

	app.Listen(":3000")
}
