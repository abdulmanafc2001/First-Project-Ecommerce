package controllers

import (
	"strconv"

	"github.com/abdulmanafc2001/First-Project-Ecommerce/database"
	"github.com/abdulmanafc2001/First-Project-Ecommerce/models"
	"github.com/gin-gonic/gin"
)

// ListOrders retrieves a list of orders for the authenticated user, paginated and with details about each order.
// @Summary List user orders
// @Description Retrieves a list of orders for the authenticated user, paginated and with details about each order.
// @Tags orders
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param page query integer false "Page number for pagination" default(1)
// @Param limit query integer false "Number of orders per page" default(10)
// @Security ApiKeyAuth
// @Success 200 {json} SuccessResponse
// @Failure 500 {string} ErrorResponse
// @Failure 400 {string} ErrorResponse
// @Router /user/listorders [get]
func ListOrders(c *gin.Context) {
	user, _ := c.Get("user")
	userId := user.(models.User).User_ID
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		c.JSON(500, gin.H{
			"error": "query value error",
		})
		return
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		c.JSON(500, gin.H{
			"error": "query value error",
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
		Image        string
	}
	//Fetching data from database and inner joins product table knowing product details
	var orders []orderDetails
	err = database.DB.Table("order_items").
		Select("products.product_name,order_items.quantity,order_items.price,order_items.total_price,order_items.status,order_items.order_id,order_items.brand,order_items.catagory,order_items.address_id,images.image").
		Joins("INNER JOIN products ON products.id=order_items.product_id").Joins("INNER JOIN images ON images.product_id=order_items.product_id").
		Where("user_id=?", userId).Limit(limit).Offset(offset).
		Scan(&orders).Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Scanning error",
		})
		return
	}
	//showing data in to user
	c.JSON(200, gin.H{
		"orders": orders,
	})
}

// ListOrdersWithBrand retrieves a list of orders for the authenticated user and a specific brand, paginated and with details about each order.
// @Summary List user orders with a specific brand
// @Description Retrieves a list of orders for the authenticated user and a specific brand, paginated and with details about each order.
// @Tags orders
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param brandid query integer true "Brand ID for filtering orders"
// @Param page query integer false "Page number for pagination" default(1)
// @Param limit query integer false "Number of orders per page" default(10)
// @Security ApiKeyAuth
// @Success 200 {json} SuccessResponse
// @Failure 500 {string} ErrorResponse
// @Failure 400 {string} ErrorResponse
// @Router /user/listorderswithbrand [get]
func ListOrdersWithBrand(c *gin.Context) {
	brandId, err := strconv.Atoi(c.Query("brandid"))
	if err != nil {
		c.JSON(500, gin.H{
			"error": "query value error",
		})
		return
	}
	user, _ := c.Get("user")
	userId := user.(models.User).User_ID
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		c.JSON(500, gin.H{
			"error": "query value error",
		})
		return
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		c.JSON(500, gin.H{
			"error": "query value error",
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
		Image        string
	}
	//Fetching data from database and inner joins product table knowing product details
	var orders []orderDetails
	err = database.DB.Table("order_items").
		Select("products.product_name,order_items.quantity,order_items.price,order_items.total_price,order_items.status,order_items.order_id,order_items.brand,order_items.catagory,order_items.address_id,images.image").
		Joins("INNER JOIN products ON products.id=order_items.product_id").Joins("INNER JOIN images ON images.product_id=order_items.product_id").
		Where("order_items.user_id=? AND brand_id=?", userId, brandId).Limit(limit).Offset(offset).
		Scan(&orders).Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Scanning error",
		})
		return
	}
	//showing data in to user
	c.JSON(200, gin.H{
		"orders": orders,
	})
}

// ListOrdersWithCatagory retrieves a list of orders for the authenticated user and a specific category, paginated and with details about each order.
// @Summary List user orders with a specific category
// @Description Retrieves a list of orders for the authenticated user and a specific category, paginated and with details about each order.
// @Tags orders
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param catagoryid query integer true "Category ID for filtering orders"
// @Param page query integer false "Page number for pagination" default(1)
// @Param limit query integer false "Number of orders per page" default(10)
// @Security ApiKeyAuth
// @Success 200 {json} SuccessResponse
// @Failure 500 {string} ErrorResponse
// @Failure 400 {string} ErrorResponse
// @Router /user/listorderswithcatagory [get]
func ListOrdersWithCatagory(c *gin.Context) {
	catagoryId, err := strconv.Atoi(c.Query("catagoryid"))
	if err != nil {
		c.JSON(500, gin.H{
			"error": "query value error1",
		})
		return
	}
	user, _ := c.Get("user")
	userId := user.(models.User).User_ID
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		c.JSON(500, gin.H{
			"error": "query value error",
		})
		return
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		c.JSON(500, gin.H{
			"error": "query value error",
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
		Image        string
	}
	//Fetching data from database and inner joins product table knowing product details
	var orders []orderDetails
	err = database.DB.Table("order_items").
		Select("products.product_name,order_items.quantity,order_items.price,order_items.total_price,order_items.status,order_items.order_id,order_items.brand,order_items.catagory,order_items.address_id,images.image").
		Joins("INNER JOIN products ON products.id=order_items.product_id").Joins("INNER JOIN images ON images.product_id=order_items.product_id").
		Where("order_items.user_id=? AND catagory_id=?", userId, catagoryId).Limit(limit).Offset(offset).
		Scan(&orders).Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Scanning error",
		})
		return
	}
	//showing data in to user
	c.JSON(200, gin.H{
		"orders": orders,
	})
}

// ListOrderDesc retrieves a list of orders for the authenticated user, sorted in descending order of price, and paginated with details about each order.
// @Summary List user orders in descending order of price
// @Description Retrieves a list of orders for the authenticated user, sorted in descending order of price, paginated and with details about each order.
// @Tags orders
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param page query integer false "Page number for pagination" default(1)
// @Param limit query integer false "Number of orders per page" default(10)
// @Security ApiKeyAuth
// @Success 200 {json} SuccessResponse
// @Failure 500 {string} ErrorResponse
// @Failure 400 {string} ErrorResponse
// @Router /user/listorderdesc [get]
func ListOrderDesc(c *gin.Context) {
	user, _ := c.Get("user")
	userId := user.(models.User).User_ID

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		c.JSON(500, gin.H{
			"error": "query value error",
		})
		return
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		c.JSON(500, gin.H{
			"error": "query value error",
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
		Image        string
	}
	var orders []orderDetails
	err = database.DB.Table("order_items").
		Select("products.product_name,order_items.quantity,order_items.price,order_items.total_price,order_items.status,order_items.order_id,order_items.brand,order_items.catagory,order_items.address_id,images.image").
		Joins("INNER JOIN products ON products.id=order_items.product_id").Joins("INNER JOIN images ON images.product_id=order_items.product_id").
		Where("user_id=?", userId).Order("order_items.price desc").Limit(limit).Offset(offset).
		Scan(&orders).Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Scanning error",
		})
		return
	}
	//showing data in to user
	c.JSON(200, gin.H{
		"orders": orders,
	})
}

// ListOrderAsc retrieves a list of orders for the authenticated user, sorted in ascending order of price, and paginated with details about each order.
// @Summary List user orders in ascending order of price
// @Description Retrieves a list of orders for the authenticated user, sorted in ascending order of price, paginated and with details about each order.
// @Tags orders
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param page query integer false "Page number for pagination" default(1)
// @Param limit query integer false "Number of orders per page" default(10)
// @Security ApiKeyAuth
// @Success 200 {json} SuccessResponse
// @Failure 500 {string} ErrorResponse
// @Failure 400 {string} ErrorResponse
// @Router /user/listorderasc [get]
func ListOrderAsc(c *gin.Context) {
	user, _ := c.Get("user")
	userId := user.(models.User).User_ID

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		c.JSON(500, gin.H{
			"error": "query value error",
		})
		return
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		c.JSON(500, gin.H{
			"error": "query value error",
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
		Image        string
	}
	var orders []orderDetails
	err = database.DB.Table("order_items").
		Select("products.product_name,order_items.quantity,order_items.price,order_items.total_price,order_items.status,order_items.order_id,order_items.brand,order_items.catagory,order_items.address_id,images.image").
		Joins("INNER JOIN products ON products.id=order_items.product_id").Joins("INNER JOIN images ON images.product_id=order_items.product_id").
		Where("user_id=?", userId).Order("order_items.price asc").Limit(limit).Offset(offset).
		Scan(&orders).Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Scanning error",
		})
		return
	}
	//showing data in to user
	c.JSON(200, gin.H{
		"orders": orders,
	})
}

// CancelOrderWithId cancels an order for the authenticated user based on the provided order ID.
// @Summary Cancel user order by order ID
// @Description Cancels an order for the authenticated user based on the provided order ID and updates the order status and related data accordingly.
// @Tags orders
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param order_id path integer true "Order ID to be cancelled"
// @Security ApiKeyAuth
// @Success 200 {json} SuccessResponse
// @Failure 400 {string} ErrorResponse 
// @Router /user/cancelorder/{order_id} [post]
func CancelOrderWithId(c *gin.Context) {
	//getting user details from middlewares
	user, _ := c.Get("user")
	userId := user.(models.User).User_ID
	orderId, err := strconv.Atoi(c.Param("order_id"))
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	//fetching the data from database
	var order models.Order
	err = database.DB.Where("user_id=? AND order_id=?", userId, orderId).First(&order).Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Your order didn't find",
		})
		return
	}

	var payments models.Payment
	database.DB.First(&payments, order.Payment_ID)
	//checking if order is already cancelled or not
	if order.Status == "cancelled" {
		c.JSON(400, gin.H{
			"error": "This order alredy cancelled",
		})
		return
	}
	//updating the status in to cancelled
	err = database.DB.Model(&models.Order{}).Where("user_id=? AND order_id=?", userId, orderId).Update("status", "cancelled").Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": "updation error",
		})
		return
	}

	err = database.DB.Model(&models.OrderItem{}).Where("order_id=?", orderId).Update("status", "cancelled").Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": "updation error",
		})
		return
	}

	var orderItems []models.OrderItem
	database.DB.Where("order_id=?", orderId).Find(&orderItems)

	totalprices := 0
	var product models.Product
	for _, v := range orderItems {
		totalprices += int(v.Total_Price)
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
	if payments.Payment_Type == "RAZOR PAY" {
		err = database.DB.Model(&models.User{}).Where("user_id=?", userId).Update("wallet", user.(models.User).Wallet+uint(totalprices)).Error
		if err != nil {
			c.JSON(400, gin.H{
				"error": "Database error",
			})
			return
		}
	}

	c.JSON(200, gin.H{
		"message": "successfully cancelled order",
	})
}
