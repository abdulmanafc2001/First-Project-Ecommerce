package controllers

import (
	"time"

	"github.com/abdulmanafc2001/First-Project-Ecommerce/database"
	"github.com/abdulmanafc2001/First-Project-Ecommerce/models"
	"github.com/gin-gonic/gin"
)

// CheckOutCOD processes a cash-on-delivery (COD) order.
// @Summary Process a cash-on-delivery order
// @Description Processes a cash-on-delivery order for the authenticated user.
// @Tags orders
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer Token"
// @Param request body models.Order true "Order details"
// @Success 200 {json} SuccessResponse
// @Failure 400 {json} ErrorResponse
// @Router /user/checkoutcod [post]
func CheckOutCOD(c *gin.Context) {
	//getting user details from middleware
	user, _ := c.Get("user")
	userId := user.(models.User).User_ID
	//searching for database all cart data
	var cartdata []models.Cart
	err := database.DB.Where("user_id=?", userId).Find(&cartdata).Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Please check your cart",
		})
		return
	}
	//getting total price of cart
	var totalprice uint
	err = database.DB.Table("carts").Select("SUM(total_price)").Where("user_id=?", userId).Scan(&totalprice).Error
	if err != nil {
		c.JSON(400, gin.H{
			"error":   "Failed to find total price",
			"message": "cart is empty",
		})
		return
	}
	//checking stock level
	var product models.Product
	for _, v := range cartdata {
		database.DB.First(&product, v.Product_ID)
		if product.Stock-v.Quantity < 0 {
			c.JSON(400, gin.H{
				"error": "Please check quantity",
			})
			return
		}
	}

	var order models.Order
	if err := c.Bind(&order); err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	database.DB.Create(&models.Payment{
		Payment_Type:   "COD",
		Total_Amount:   totalprice,
		Payment_Status: "Pending",
		User_ID:        userId,
		Date:           time.Now(),
	})
	var payment models.Payment
	database.DB.Last(&payment)
	var address models.Address
	err = database.DB.Where("user_id=? AND address_id=?", userId, order.Address_ID).First(&address).Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Failed to find address choose different id",
		})
		return
	}
	err = database.DB.Create(&models.Order{
		User_ID:     userId,
		Address_ID:  order.Address_ID,
		Total_Price: totalprice,
		Payment_ID:  payment.Payment_ID,
		Status:      "processing",
	}).Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	var cartbrand struct {
		Brand_Name    string
		Catagory_Name string
	}

	var order1 models.Order
	database.DB.Last(&order1)
	for _, cartdata := range cartdata {
		database.DB.Table("products").Select("brands.brand_name,catagories.catagory_name").
			Joins("INNER JOIN brands ON brands.brand_id=products.brand_id").
			Joins("INNER JOIN catagories ON catagories.catagory_id=products.catagory_id").Where("id=?", cartdata.Product_ID).Scan(&cartbrand)
		err = database.DB.Create(&models.OrderItem{
			Order_ID:    order1.Order_ID,
			User_ID:     userId,
			Product_ID:  cartdata.Product_ID,
			Address_ID:  order.Address_ID,
			Brand:       cartbrand.Brand_Name,
			Catagory:    cartbrand.Catagory_Name,
			Quantity:    cartdata.Quantity,
			Price:       cartdata.Price,
			Total_Price: cartdata.Total_Price,
			Discount:    cartdata.Catagory_Offer + cartdata.Coupon_Discount,
			Cart_ID:     cartdata.ID,
			Status:      "processing",
			Created_at:  time.Now(),
		}).Error
		if err != nil {
			break
		}
	}
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	//reducing the stock count in the database
	var products models.Product
	for _, v := range cartdata {
		database.DB.First(&products, v.Product_ID)
		database.DB.Model(&models.Product{}).Where("id=?", v.Product_ID).Update("stock", product.Stock-v.Quantity)
	}

	err = database.DB.Delete(&models.Cart{}, "user_id=?", userId).Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	//giving success message
	c.JSON(200, gin.H{
		"message": "successfully ordered your cart",
	})
}
