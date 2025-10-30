package middleware

import (
	"crudprojectgo/helper"
	"net/url"
	"strings"
	"regexp"

	"github.com/gofiber/fiber/v2"
)

var pathTraversalPattern = regexp.MustCompile(`(\.\./|\.\.\\|//|\\)`)

func SanitizeMiddleware(c *fiber.Ctx) error {
	rawPath := c.OriginalURL()
	cleaned := strings.TrimSpace(rawPath)
	cleaned, _ = url.PathUnescape(cleaned)

	// Deteksi tanda path traversal
	if pathTraversalPattern.MatchString(cleaned) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid path detected (possible traversal attempt)",
		})
	}

	cleaned = strings.ReplaceAll(cleaned, "//", "/")
	c.Path(cleaned)

	// Bersihkan semua params & query
	for key, val := range c.AllParams() {
		c.Locals("param_"+key, helper.CleanText(val))
	}
	for key, val := range c.Queries() {
		c.Locals("query_"+key, helper.CleanText(val))
	}

	return c.Next()
}
