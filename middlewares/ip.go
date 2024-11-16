package middlewares

import (
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

/*
LogIPMiddleware logs the user's IP address along with
the request method and path to track any incoming requests
*/
func LogIPMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		requestMethod := c.Request.Method
		requestPath := c.Request.URL.Path

		file, err := os.OpenFile("../logs/ip_logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Printf("Error opening or creating file: %v", err)
			c.Next()
			return
		}
		defer file.Close()

		logEntry := clientIP + " - " + requestMethod + " " + requestPath + "\n"

		if _, err := file.WriteString(logEntry); err != nil {
			log.Printf("Error writing to file: %v", err)
		}

		c.Next()
	}
}
