package middlewares

import (
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

/*
LogIPMiddleware allow to find out users IP
that need to be detected if someone will try
to send something on server and etc
*/
func LogIPMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		file, err := os.OpenFile("../logs/ip_logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Printf("Error opening or creating file: %v", err)
			c.Next()
			return
		}
		defer file.Close()

		if _, err := file.WriteString(clientIP + "\n"); err != nil {
			log.Printf("Error writing to file: %v", err)
		}

		c.Next()
	}
}
