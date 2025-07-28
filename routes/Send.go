package routes

import (
	db_controller "project/db"

	"github.com/gin-gonic/gin"
)

type Movement struct {
	From   string  `form:"from" json:"from" binding:"required"`
	To     string  `form:"to" json:"to" binding:"required"`
	Amount float64 `form:"amount" json:"amount" binding:"required"`
}

// POST /api/send
// Creates new transaction
// req.body (application/json)
//
//	{
//			from: (string, required)
//			to: (string, required)
//			amount: (float64, required)
//	}
//
// res.code
//
//	200:
//
//		res.body (application/json): {"error": null}
//
//	400:
//
//		res.body (application/json): {"error": (error)}
func Send(router *gin.Engine) {
	router.POST("/api/send", func(c *gin.Context) {
		var data Movement

		if err := c.BindJSON(&data); err != nil {
			c.JSON(400, gin.H{"error": "unable to convert request body"})
			return
		}

		if err := db_controller.Send(data.From, data.To, data.Amount); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"error": nil})
	})
}
