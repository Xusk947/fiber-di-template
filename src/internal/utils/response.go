package utils

import (
	"github.com/gofiber/fiber/v2"
)

// Response represents a standardized API response structure
// Example:
// {
//   "success": true,
//   "data": { "user": { "id": 1, "name": "John" } },
//   "error": null,
//   "meta": { "page": 1, "page_size": 10, "total_pages": 5, "total_count": 42 }
// }
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
	Meta    *MetaInfo   `json:"meta,omitempty"`
}

// ErrorInfo represents detailed error information
// Example:
// {
//   "code": "validation_error",
//   "message": "Invalid input parameters",
//   "details": { "field": "email", "issue": "format" }
// }
type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

// MetaInfo represents metadata for paginated responses
// Example:
// {
//   "page": 1,
//   "page_size": 10,
//   "total_pages": 5,
//   "total_count": 42
// }
type MetaInfo struct {
	Page       int `json:"page,omitempty"`
	PageSize   int `json:"page_size,omitempty"`
	TotalPages int `json:"total_pages,omitempty"`
	TotalCount int `json:"total_count,omitempty"`
}

// NewResponse creates a new response with default values
// Example:
//
//	resp := utils.NewResponse()
func NewResponse() *Response {
	return &Response{
		Success: true,
	}
}

// WithData adds data to the response
// Example:
//
//	resp := utils.NewResponse().WithData(user)
func (r *Response) WithData(data interface{}) *Response {
	r.Data = data
	return r
}

// WithError adds error information to the response
// Example:
//
//	resp := utils.NewResponse().WithError("not_found", "User not found", nil)
func (r *Response) WithError(code string, message string, details any) *Response {
	r.Success = false
	r.Error = &ErrorInfo{
		Code:    code,
		Message: message,
		Details: details,
	}
	return r
}

// WithMeta adds pagination metadata to the response
// Example:
//
//	resp := utils.NewResponse().WithMeta(1, 10, 5, 42)
func (r *Response) WithMeta(page, pageSize, totalPages, totalCount int) *Response {
	r.Meta = &MetaInfo{
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
		TotalCount: totalCount,
	}
	return r
}

// Send writes the response to the Fiber context
// Example:
//
//	return utils.NewResponse().WithData(user).Send(c, fiber.StatusOK)
func (r *Response) Send(c *fiber.Ctx, status int) error {
	return c.Status(status).JSON(r)
}

// Success creates a success response with optional data
// Example:
//
//	return utils.Success(c, user)
func Success(c *fiber.Ctx, data interface{}) error {
	return NewResponse().WithData(data).Send(c, fiber.StatusOK)
}

// SuccessWithMeta creates a success response with data and pagination metadata
// Example:
//
//	return utils.SuccessWithMeta(c, users, 1, 10, 5, 42)
func SuccessWithMeta(c *fiber.Ctx, data interface{}, page, pageSize, totalPages, totalCount int) error {
	return NewResponse().
		WithData(data).
		WithMeta(page, pageSize, totalPages, totalCount).
		Send(c, fiber.StatusOK)
}

// Error creates an error response
// Example:
//
//	return utils.Error(c, fiber.StatusBadRequest, "validation_error", "Invalid input", validationErrors)
func Error(c *fiber.Ctx, status int, code, message string, details any) error {
	return NewResponse().
		WithError(code, message, details).
		Send(c, status)
}

// BadRequest is a helper for 400 Bad Request responses
// Example:
//
//	return utils.BadRequest(c, "validation_error", "Invalid email format", nil)
func BadRequest(c *fiber.Ctx, code, message string, details any) error {
	return Error(c, fiber.StatusBadRequest, code, message, details)
}

// NotFound is a helper for 404 Not Found responses
// Example:
//
//	return utils.NotFound(c, "user_not_found", "User not found", nil)
func NotFound(c *fiber.Ctx, code, message string, details any) error {
	return Error(c, fiber.StatusNotFound, code, message, details)
}

// InternalServerError is a helper for 500 Internal Server Error responses
// Example:
//
//	return utils.InternalServerError(c, err)
func InternalServerError(c *fiber.Ctx, err error) error {
	return Error(c, fiber.StatusInternalServerError, "internal_error", "An internal server error occurred", nil)
}

// Unauthorized is a helper for 401 Unauthorized responses
// Example:
//
//	return utils.Unauthorized(c, "Invalid or expired token")
func Unauthorized(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusUnauthorized, "unauthorized", message, nil)
}

// Forbidden is a helper for 403 Forbidden responses
// Example:
//
//	return utils.Forbidden(c, "Insufficient permissions to access this resource")
func Forbidden(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusForbidden, "forbidden", message, nil)
}

// SuccessOnly returns a success response with no data
// Example:
//
//	return utils.SuccessOnly(c)
func SuccessOnly(c *fiber.Ctx) error {
	return NewResponse().Send(c, fiber.StatusOK)
}

// Created is a helper for 201 Created responses
// Example:
//
//	return utils.Created(c, newUser)
func Created(c *fiber.Ctx, data interface{}) error {
	return NewResponse().WithData(data).Send(c, fiber.StatusCreated)
}

// NoContent is a helper for 204 No Content responses
// Example:
//
//	return utils.NoContent(c)
func NoContent(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}
