package controllers

import (
	"time"

	"github.com/abdulmanafc2001/First-Project-Ecommerce/database"
	"github.com/abdulmanafc2001/First-Project-Ecommerce/models"
	"github.com/gin-gonic/gin"
)

// AdminDashboard provides an overview of the admin's dashboard data, including recent sales, total sales, and product counts.
// @Summary Get admin dashboard data
// @Description Provides an overview of the admin's dashboard data, including recent sales, total sales, and product counts for the last 30 days.
// @Tags admin
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer Token"
// @Success 200 {json} successResponse
// @Failure 400 {json} ErrorResponse
// @Router /admin/dashboard [get]
func AdminDashboard(c *gin.Context) {

	type data struct {
		User_ID      uint
		Order_ID     uint
		Product_Name string
		Price        uint
		Total_Price  uint
		Quantity     uint
		Status       string
	}

	before30Days := time.Now().AddDate(0, 0, -30)

	var result []data

	err := database.DB.Model(&models.OrderItem{}).
		Select("products.product_name,order_items.user_id,order_items.order_id,order_items.price,order_items.total_price,order_items.quantity,order_items.status").
		Joins("INNER JOIN products ON products.id = order_items.product_id").Where("order_items.created_at > ?", before30Days).
		Scan(&result).Error

	if err != nil {
		c.JSON(400, gin.H{
			"error": "Failed to find data",
		})
		return
	}
	var totalsale uint
	err = database.DB.Table("order_items").Select("SUM(total_price)").Where("created_at > ?", before30Days).Scan(&totalsale).Error

	if err != nil {
		c.JSON(400, gin.H{
			"error": "Failed to find total sale",
		})
		return
	}

	type orderdata struct {
		Product_Id uint
		Count      uint
	}
	var orderdatas []orderdata
	database.DB.Table("order_items").Select("product_id,COUNT(quantity)").Group("product_id").Scan(&orderdatas)

	var users []models.User
	err = database.DB.Find(&users).Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Failed to find users",
		})
		return
	}
	num := database.DB.Find(&users).RowsAffected

	c.JSON(200, gin.H{
		"last 30 datys sale":    result,
		"total sale of 30 days": totalsale,
		"count":                 orderdatas,
		"no.of user":            num,
		"users":                 users,
	})

}
