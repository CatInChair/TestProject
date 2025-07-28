package routes

import (
	db_controller "project/db"

	"github.com/gin-gonic/gin"
)

// GET /api/wallet/{address}/balance
// Returns specified wallet balance
//
// req.params.address: (string, required)
//
// res.code
//
//	200:
//
//		res.body (application/json): {"error": null, "balance": :value:}
//
//	404:
//
//		res.body (application/json): {"error": "wallet not found"}
func GetBalance(router *gin.Engine) {
	router.GET("/api/wallet/:address/balance", func(c *gin.Context) {
		address, _ := c.Params.Get("address")

		balance, err := db_controller.GetBalance(address)

		if err != nil {
			c.JSON(404, gin.H{"error": "wallet nof found"})
			return
		}

		c.JSON(200, gin.H{"balance": balance, "error": nil})
	})
}
