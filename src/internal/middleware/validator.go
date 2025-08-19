package middleware

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gitlab.stat4market.com/reelsmarket/fiber-di-server-template/src/internal/model"
	"gitlab.stat4market.com/reelsmarket/fiber-di-server-template/src/internal/utils"
)

// Global validator instance
var validate = validator.New()

// ValidateStruct validates a struct against its validation tags
// and returns a slice of validation errors if any
//
// Example:
//
//	errors := ValidateStruct(user)
//	if len(errors) > 0 {
//	    return utils.BadRequest(c, "validation_error", "Invalid input", errors)
//	}
func ValidateStruct(payload interface{}) []*model.ValidationError {
	var errors []*model.ValidationError
	err := validate.Struct(payload)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element model.ValidationError
			element.Field = err.Field()
			element.Tag = err.Tag()
			element.Value = err.Param()
			element.Message = getErrorMsg(err)
			errors = append(errors, &element)
		}
	}
	return errors
}

// getErrorMsg returns a user-friendly error message based on the validation tag
func getErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Please enter a valid email address"
	case "min":
		return "Does not meet minimum length requirement"
	case "max":
		return "Exceeds maximum length limit"
	case "eqfield":
		return "Fields do not match"
	case "len":
		return "Must be exactly " + fe.Param() + " characters long"
	case "alphanum":
		return "Must contain only alphanumeric characters"
	case "numeric":
		return "Must contain only numeric characters"
	case "uuid":
		return "Must be a valid UUID"
	default:
		return "Validation error on field: " + fe.Field()
	}
}

// RequestValidator is a generic middleware that validates request bodies
// against the provided struct type
//
// Example:
//
//	app.Post("/users", RequestValidator(model.CreateUserRequest{}), handlers.CreateUser)
func RequestValidator[T any](payload T) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Create a new instance of the payload type
		var requestData T

		// Parse the request body
		if err := c.BodyParser(&requestData); err != nil {
			return utils.BadRequest(c, "invalid_request", "Invalid request body format", nil)
		}

		// Validate the struct
		if errors := ValidateStruct(requestData); len(errors) > 0 {
			return utils.BadRequest(c, "validation_error", "Validation failed", errors)
		}

		// Store the validated payload in context locals
		c.Locals("payload", requestData)

		return c.Next()
	}
}

// QueryValidator is a generic middleware that validates query parameters
// against the provided struct type
//
// Example:
//
//	app.Get("/users", QueryValidator(model.ListUsersQuery{}), handlers.ListUsers)
func QueryValidator[T any](payload T) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Create a new instance of the payload type
		var queryData T

		// Parse the query parameters
		if err := c.QueryParser(&queryData); err != nil {
			return utils.BadRequest(c, "invalid_query", "Invalid query parameters", nil)
		}

		// Validate the struct
		if errors := ValidateStruct(queryData); len(errors) > 0 {
			return utils.BadRequest(c, "validation_error", "Validation failed", errors)
		}

		// Store the validated query in context locals
		c.Locals("query", queryData)

		return c.Next()
	}
}

// ParamsValidator is a generic middleware that validates URL parameters
// against the provided struct type
//
// Example:
//
//	app.Get("/users/:id", ParamsValidator(model.UserParams{}), handlers.GetUser)
func ParamsValidator[T any](payload T) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Create a new instance of the payload type
		var paramsData T

		// Parse the URL parameters
		if err := c.ParamsParser(&paramsData); err != nil {
			return utils.BadRequest(c, "invalid_params", "Invalid URL parameters", nil)
		}

		// Validate the struct
		if errors := ValidateStruct(paramsData); len(errors) > 0 {
			return utils.BadRequest(c, "validation_error", "Validation failed", errors)
		}

		// Store the validated parameters in context locals
		c.Locals("params", paramsData)

		return c.Next()
	}
}
