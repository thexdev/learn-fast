package routes

import (
	"go_fiber_api/internals/handlers"
	validation "go_fiber_api/internals/utils"

	"github.com/gofiber/fiber/v2"
)

func GetAPIRoutes() *fiber.App {
	api := fiber.New()

	handler := handlers.ApiHandler{}

	api.Get("/products", handler.FetchAll)

	api.Get("/products/:id", handler.FetchProductById)

	type UpdateProductRequest struct {
		Name string `json:"name" validate:"required"`
	}
	api.Patch("/products/:id", validation.ValidateRequest[UpdateProductRequest](), handler.UpdateProductById)

	type CreateProductRequest struct {
		Name string `json:"name" validate:"required"`
	}
	api.Post("/products", validation.ValidateRequest[CreateProductRequest](), handler.CreateProduct)

	api.Delete("/products/:id", handler.DeleteProductById)

	return api
 }