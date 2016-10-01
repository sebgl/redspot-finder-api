package main

import (
	"flag"
	"fmt"

	"gopkg.in/olivere/elastic.v3"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

var (
	port  = flag.Int("port", 8082, "Port to bind to")
	esURL = flag.String("es", "http://localhost:9200", "Elasticsearch url in the form http://host:port")

	es *elastic.Client
)

func main() {

	es = createESClient(*esURL)

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(cORSMiddleware())
	addRoutes(r)

	log.WithField("port", *port).Info("Starting server")
	err := r.Run(fmt.Sprintf(":%d", *port))
	if err != nil {
		log.WithError(err).Error("Unable to run server")
	}
}

func createESClient(URL string) *elastic.Client {
	client, err := elastic.NewClient(elastic.SetURL(URL))
	if err != nil {
		log.WithError(err).Fatal("Unable to create elasticsearch client")
	}
	return client
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
