package middleware

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
)

type RequestLog struct {
	Timestamp   time.Time   `json:"timestamp"`
	Method      string      `json:"method"`
	URL         string      `json:"url"`
	Headers     interface{} `json:"headers"`
	Body        string      `json:"body"`
	UserAgent   string      `json:"user_agent"`
	IP          string      `json:"ip"`
}

func RequestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Log request
		logRequest(c)
		
		// Continue to next handler
		return c.Next()
	}
}

func logRequest(c *fiber.Ctx) {
	// Ensure logs directory exists
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		os.MkdirAll("logs", 0755)
	}
	
	requestLog := RequestLog{
		Timestamp: time.Now(),
		Method:    c.Method(),
		URL:       c.OriginalURL(),
		Headers:   c.GetReqHeaders(),
		Body:      string(c.Body()),
		UserAgent: c.Get("User-Agent"),
		IP:        c.IP(),
	}
	
	// Create log file name with current date
	logFileName := fmt.Sprintf("logs/requests_%s.log", time.Now().Format("2006-01-02"))
	
	// Convert to JSON
	logJSON, err := json.Marshal(requestLog)
	if err != nil {
		return
	}
	
	// Write to file
	file, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer file.Close()
	
	file.WriteString(string(logJSON) + "\n")
}
