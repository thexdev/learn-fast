package handlers

import (
	"context"
	"fmt"
	"go_fiber_api/internals/db"
	"go_fiber_api/internals/models"
	response "go_fiber_api/internals/utils/http"
	"log"

	"github.com/gofiber/fiber/v2"
)

type ApiHandler struct {
}

func (h *ApiHandler) FetchAll(c *fiber.Ctx) error {
	pool, err := db.Pool()
	if err != nil {
		log.Println(err)
		return response.InternalServerError(c)
	}
	defer pool.Close()

	var products []models.Product

	rows, err := pool.Query(context.Background(), "select id, name, uuid, created_at from products")
	if err != nil {
		log.Println(err)
		return response.InternalServerError(c)
	}

	for rows.Next() {
		var product models.Product

		err := rows.Scan(&product.ID, &product.Name, &product.UUID, &product.CreatedAt)
		if err != nil {
			log.Println(err)
			return response.InternalServerError(c)
		}

		products = append(products, product)
	}

	return c.JSON(fiber.Map{
		"data": products,
	})
}

func (h *ApiHandler) FetchProductById(c *fiber.Ctx) error {
	pool, err := db.Pool()
	if err != nil {
		log.Println(err)
		return response.InternalServerError(c)
	}
	defer pool.Close()

	var product models.Product

	err = pool.QueryRow(
		context.Background(),
		"select id, name, uuid, created_at from products where id = $1",
		c.Params("id"),
	).Scan(&product.ID, &product.Name, &product.UUID, &product.CreatedAt)
	if err != nil && db.RecordNotFound(err) {
		return response.NotFound(c, fmt.Sprintf(
			"product with id %s not found", c.Params("id"),
		))
	}
	if err != nil && !db.RecordNotFound(err) {
		log.Println(err)
		return response.InternalServerError(c)
	}

	return c.JSON(fiber.Map{ "data": product })
}

func (h *ApiHandler) CreateProduct(c *fiber.Ctx) error {
	pool, err := db.Pool()
	if err != nil {
		log.Println(err)
		return response.InternalServerError(c)
	}
	defer pool.Close()

	var payload struct {
		Name string
	}
	c.BodyParser(&payload)

	result, err := pool.Exec(
		context.Background(),
		"insert into products (name) values ($1)", payload.Name,
	)
	if err != nil {
		log.Println(err)
		return response.InternalServerError(c)
	}

	if result.RowsAffected() == 0 {
		return response.InternalServerError(c)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "product created",
	})
}

func (h *ApiHandler) UpdateProductById(c *fiber.Ctx) error {
	pool, err := db.Pool()
	if err != nil {
		log.Println(err)
		return response.InternalServerError(c)
	}
	defer pool.Close()

	result, err := pool.Exec(
		context.Background(),
		"update products set name = $1 where id = $2", "coba", 1,
	)
	if err != nil {
		log.Println(err)
		return response.InternalServerError(c)
	}

	// check if any row has actually updated
	if result.RowsAffected() == 0 {
		return response.NotFound(c, "product not founnd")
	}

	return c.JSON(fiber.Map{ "message": "product updated" })
}

// Delete product record using specified id
func (h *ApiHandler) DeleteProductById(c *fiber.Ctx) error {
	pool, err := db.Pool()
	if err != nil {
		log.Println(err)
		return response.InternalServerError(c)
	}
	defer pool.Close()

	conn, err := pool.Acquire(context.Background())
	if err != nil {
		log.Println(err)
		return response.InternalServerError(c)
	}
	defer conn.Release()

	result, err := conn.Exec(context.Background(), "delete from products where id = $1", c.Params("id"))
	if err != nil {
		log.Println(err)
		return response.InternalServerError(c)
	}

	if result.RowsAffected() == 0 {
		return response.NotFound(c, "product not found")
	}

	return c.JSON(fiber.Map{ "message": "product deleted" })
}
