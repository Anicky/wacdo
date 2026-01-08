package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.Use(cors.Default()) // All origins allowed by default

	router.GET("/status", func(context *gin.Context) {
		context.JSON(200, gin.H{
			"message": "OK",
		})
	})

	// Listens on 0.0.0.0:8080 by default
	err := router.Run()

	if err != nil {
		log.Fatal("Failed to run router.")
	}
}
