package main

import (
	"flag"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

var (
	port = flag.Int("port", 8082, "Port to bind to")
	es   = flag.String("es", "localhost:9002", "Elasticsearch url in the form host:port")
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	r.Use(cORSMiddleware())

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "up",
			"message": "I'm up, captain !",
		})
	})
	err := r.Run(fmt.Sprintf(":%d", *port))
	if err != nil {
		log.WithError(err).Error("Unable to run server")
	}
}

func cORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		domain := "*"
		c.Writer.Header().Set("Access-Control-Allow-Origin", domain)
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}
