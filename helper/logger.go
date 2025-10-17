package helper

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
)

type LogEntry struct {
	Timestamp    time.Time   `json:"timestamp"`
	Method       string      `json:"method"`
	URL          string      `json:"url"`
	StatusCode   int         `json:"status_code"`
	ResponseBody interface{} `json:"response_body"`
	UserAgent    string      `json:"user_agent"`
	IP           string      `json:"ip"`
}

// path global log
var logDir = filepath.Join(".", "logs")
var logFile = filepath.Join(logDir, "response.log")

func ensureLogDir() {
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		if err := os.MkdirAll(logDir, 0755); err != nil {
			log.Printf("Error creating log directory: %v", err)
		}
	}
}

func LogResponse(c *fiber.Ctx, statusCode int, responseBody interface{}) {
	log.Printf("[DEBUG] LogResponse called - Method: %s, URL: %s, Status: %d",
		c.Method(), c.OriginalURL(), statusCode)

	ensureLogDir()

	logEntry := LogEntry{
		Timestamp:    time.Now(),
		Method:       c.Method(),
		URL:          c.OriginalURL(),
		StatusCode:   statusCode,
		ResponseBody: responseBody,
		UserAgent:    c.Get("User-Agent"),
		IP:           c.IP(),
	}

	logJSON, err := json.Marshal(logEntry)
	if err != nil {
		log.Printf("Error marshaling log entry: %v", err)
		return
	}

	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Error opening log file: %v", err)
		return
	}
	defer file.Close()

	if _, err := file.WriteString(string(logJSON) + "\n"); err != nil {
		log.Printf("Error writing to log file: %v", err)
		return
	}

	log.Printf("[DEBUG] Successfully logged response to %s", logFile)
}
