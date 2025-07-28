package routes

import (
	db_controller "project/db"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GET /api/transactions
// Returns last time transactions
//
// req.query.count: (int64, required)
//
// res.code
//
//	200:
//
//		res.body (application/json): {"error": null}
//
//	400:
//
//		res.body (application/json): {"error": :error:}
func GetLast(router *gin.Engine) {
	router.GET("/api/transactions", func(c *gin.Context) {
		countQuery, isCountProvided := c.GetQuery("count")

		if !isCountProvided {
			c.JSON(400, gin.H{"error": "\"count\" query param is not provided"})
			return
		}

		count, err := strconv.ParseInt(countQuery, 10, 64)

		if err != nil {
			c.JSON(400, gin.H{"error": "provided count is not a number"})
		}

		data, err := db_controller.GetLast(count)

		if err != nil {
			c.JSON(500, gin.H{"error": "database error"})
			return
		}

		c.JSON(200, gin.H{"error": nil, "transactions": data})
	})
}
