package main

import (
	"github.com/gin-gonic/gin"

	db_controller "project/db"
	"project/routes"
)

// There's nothing interesting here.
func main() {
	db_controller.Init()
	router := gin.Default()

	routes.GetBalance(router)
	routes.GetLast(router)
	routes.Send(router)

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello from Dockerized Gin!",
		})
	})

	router.Run()
}
