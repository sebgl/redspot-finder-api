package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

func addRoutes(r *gin.Engine) {

	r.GET("/", status)
	r.GET("/api/search", search)
	r.GET("/api/last", last)
}

func status(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":  "up",
		"message": "I'm up, captain !",
	})
}

func search(c *gin.Context) {
	query := c.Query("q")
	result, err := SearchPlaylists(query)
	if err != nil {
		log.WithError(err).WithField("query", query).Error("Unable to search playlists")
		return
	}
	c.JSON(200, result)
}

func last(c *gin.Context) {
	result, err := LastPlaylists()
	if err != nil {
		log.WithError(err).Error("Unable to retrieve last playlists")
		return
	}
	c.JSON(200, result)
}
