package controllers

import (
	"strconv"

	"github.com/abdulmanafc2001/First-Project-Ecommerce/database"
	"github.com/abdulmanafc2001/First-Project-Ecommerce/models"
	"github.com/gin-gonic/gin"
)

// ViewOrders retrieves a list of orders with details, optionally paginated.
// @Summary View orders
// @Description Retrieve a list of orders with product details
// @Tags adminorder
// @Accept json
// @Produce json
// @Param page query int false "Page number for pagination (default is 1)"
// @Param limit query int false "Number of items per page (default is 10)"
// @Success 200 {json} SuccessResponse
// @Failure 400 {json} ErrorResponse
// @Failure 500 {json} ErrorResponse
// @Router /admin/viewallorders [get]
func ViewOrders(c *gin.Context) {
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Query error",
		})
		return
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		c.JSON(500, gin.H{
			"error": "query error",
		})
		return
	}
	offset := (page - 1) * limit

	type orderDetails struct {
		Product_Name string
		Quantity     string
		Price        uint
		Total_Price  uint
		Status       string
		Order_ID     uint
		Brand        string
		Catagory     string
		Address_ID   uint
	}
	//Fetching data from database and inner joins product table knowing product details
	var orders []orderDetails
	err = database.DB.Table("order_items").
		Select("products.product_name,order_items.quantity,order_items.price,order_items.total_price,order_items.status,order_items.order_id,order_items.brand,order_items.catagory,order_items.address_id").
		Joins("INNER JOIN products ON products.id=order_items.product_id").Limit(limit).Offset(offset).
		Scan(&orders).Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Scanning error",
		})
		return
	}

	c.JSON(200, gin.H{
		"orders": orders,
	})

}

// CancelOrder cancels an order by its order ID.
// @Summary Cancel an order
// @Description Cancel an order by its order ID and update associated records
// @Tags adminorder
// @Accept json
// @Produce json
// @Param order_id path int true "Order ID to cancel"
// @Success 200 {json} SuccessResponse
// @Failure 400 {json} ErrorResponse
// @Router /admin/cancelorder/{order_id} [post]
func CancelOrder(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("order_id"))
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Integer convertion error",
		})
		return
	}
	var order models.Order
	err = database.DB.Where("order_id=?", id).First(&order).Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Failed to find this order",
		})
		return
	}

	if order.Status == "cancelled" {
		c.JSON(400, gin.H{
			"error": "This order already cancelled",
		})
		return
	}

	err = database.DB.Model(&models.Order{}).Where("order_id=?", id).Update("status", "cancelled").Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Updation error",
		})
		return
	}
	err = database.DB.Model(&models.OrderItem{}).Where("order_id=?", id).Update("status", "cancelled").Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Updation error",
		})
		return
	}
	var orderItems []models.OrderItem
	if err = database.DB.Where("order_id=?", id).Find(&orderItems).Error; err != nil {
		c.JSON(400, gin.H{
			"error": "Failed to find order items",
		})
		return
	}
	var product models.Product
	for _, v := range orderItems {
		database.DB.First(&product, v.Product_ID)
		database.DB.Model(&models.Product{}).Where("id=?", v.Product_ID).Update("stock", product.Stock+v.Quantity)
	}
	err = database.DB.Model(&models.Payment{}).Where("payment_id=?", order.Payment_ID).Update("payment_status", "cancelled").Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": "updation error",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "successfully cancelled order",
	})

}

type data struct {
	Status string
}

// ChangeStatus updates the status of an order by its order ID.
// @Summary Change order status
// @Description Update the status of an order by its order ID
// @Tags adminorder
// @Accept json
// @Produce json
// @Param order_id path int true "Order ID to update status"
// @Param status body data true "New status value: 'shipped', 'pending', or 'cancelled'"
// @Success 200 {json} SuccessResponse
// @Failure 400 {json} ErrorResponse
// @Router /admin/changestatus/{order_id} [patch]
func ChangeStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("order_id"))
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Integer convertion error",
		})
		return
	}
	var order models.Order
	err = database.DB.Where("order_id=?", id).First(&order).Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Failed to find this order",
		})
		return
	}
	var input data
	if err = c.Bind(&input); err != nil {
		c.JSON(400, gin.H{
			"error": "Binding error",
		})
		return
	}
	if input.Status == "shipped" || input.Status == "pending" || input.Status == "cancelled" {
		if input.Status == "cancelled" {
			if order.Status == "cancelled" {
				c.JSON(400, gin.H{
					"error": "This order already cancelled",
				})
				return
			}

			err = database.DB.Model(&models.Order{}).Where("order_id=?", id).Update("status", "cancelled").Error
			if err != nil {
				c.JSON(400, gin.H{
					"error": "Updation error",
				})
				return
			}
			err = database.DB.Model(&models.OrderItem{}).Where("order_id=?", id).Update("status", "cancelled").Error
			if err != nil {
				c.JSON(400, gin.H{
					"error": "Updation error",
				})
				return
			}
			var orderItems []models.OrderItem
			if err = database.DB.Where("order_id=?", id).Find(&orderItems).Error; err != nil {
				c.JSON(400, gin.H{
					"error": "Failed to find order items",
				})
				return
			}
			var product models.Product
			for _, v := range orderItems {
				database.DB.First(&product, v.Product_ID)
				database.DB.Model(&models.Product{}).Where("id=?", v.Product_ID).Update("stock", product.Stock+v.Quantity)
			}
			err = database.DB.Model(&models.Payment{}).Where("payment_id=?", order.Payment_ID).Update("payment_status", "cancelled").Error
			if err != nil {
				c.JSON(400, gin.H{
					"error": "updation error",
				})
				return
			}

			c.JSON(200, gin.H{
				"message": "successfully cancelled order",
			})
		} else {
			err = database.DB.Model(&models.Order{}).Where("order_id=?", id).Update("status", input.Status).Error
			if err != nil {
				c.JSON(400, gin.H{
					"error": "Updation error",
				})
				return
			}
			err = database.DB.Model(&models.OrderItem{}).Where("order_id=?", id).Update("status", input.Status).Error
			if err != nil {
				c.JSON(400, gin.H{
					"error": "Updation error",
				})
				return
			}
			c.JSON(200, gin.H{
				"message": "successfully updated status",
			})
		}
	} else {
		c.JSON(400, gin.H{
			"error": "this status not applicable",
		})
	}
}
