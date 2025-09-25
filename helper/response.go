package helper

import "github.com/gofiber/fiber/v2"

type Response struct {
    Success bool        `json:"success"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
}

func SuccessResponse(c *fiber.Ctx, message string, data interface{}) error {
	response := Response{
		Success: true,
		Message: message,
		Data:    data,
	}
	
	LogResponse(c, 200, response)
	
	return c.JSON(response)
}

func ErrorResponse(c *fiber.Ctx, statusCode int, message string) error {
	response := Response{
		Success: false,
		Error:   message,
	}
	
	LogResponse(c, statusCode, response)
	
	return c.Status(statusCode).JSON(response)
}

func CreatedResponse(c *fiber.Ctx, message string, data interface{}) error {
	response := Response{
		Success: true,
		Message: message,
		Data:    data,
	}
	
	LogResponse(c, 201, response)
	
	return c.Status(201).JSON(response)
}
