package validation

import (
	"errors"
	"fmt"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	"github.com/gofiber/fiber/v2"
)

func ValidateRequest[T any]() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var data T
		if err := c.BodyParser(&data); err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		validate := validator.New()

		english := en.New()
		uni := ut.New(english, english)

		translator, found := uni.GetTranslator("en")
		if !found {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "internal server error",
			})
		}

		if err := enTranslations.RegisterDefaultTranslations(validate, translator); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "internal server error",
			})
		}

		if err := validate.Struct(data); err != nil {
			var errs validator.ValidationErrors
			errors.As(err, &errs)

			// Translate all errors
			translatedErrors := errs.Translate(translator)

			fmt.Println(translatedErrors)

			for field, message := range translatedErrors {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": map[string]string {
						"field": field,
						"message": message,
					},
				})
			}
		}

		if err := validator.New().Struct(data); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Next()
	}
}