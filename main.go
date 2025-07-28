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

	router.Run()
}
