package helper

import (
	"encoding/json"
	"log"
	"os"
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

func ensureLogDir() {
	logDir := `D:\hp\Documents\Kuliah Panji\smt5\Pemrograman Backend\crudprojectgo\logs`
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		os.MkdirAll(logDir, 0755)
	}
}

func LogResponse(c *fiber.Ctx, statusCode int, responseBody interface{}) {
	log.Printf("[DEBUG] LogResponse called - Method: %s, URL: %s, Status: %d", c.Method(), c.OriginalURL(), statusCode)
	
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
	
	logFileName := `D:\hp\Documents\Kuliah Panji\smt5\Pemrograman Backend\crudprojectgo\logs\response.log`
	
	// Convert to JSON
	logJSON, err := json.Marshal(logEntry)
	if err != nil {
		log.Printf("Error marshaling log entry: %v", err)
		return
	}
	
	log.Printf("[DEBUG] Attempting to write to: %s", logFileName)
	
	// Write to file
	file, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Error opening log file: %v", err)
		return
	}
	defer file.Close()
	
	_, writeErr := file.WriteString(string(logJSON) + "\n")
	if writeErr != nil {
		log.Printf("Error writing to log file: %v", writeErr)
		return
	}
	
	log.Printf("[DEBUG] Successfully logged response to file")
}
