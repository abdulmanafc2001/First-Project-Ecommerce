package controllers

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/abdulmanafc2001/First-Project-Ecommerce/database"
	"github.com/abdulmanafc2001/First-Project-Ecommerce/models"
	"github.com/gin-gonic/gin"
	razorpay "github.com/razorpay/razorpay-go"
)

// RazorPay generates a RazorPay payment order for the authenticated user's cart.
// @Summary Generate RazorPay payment order
// @Description Generates a RazorPay payment order for the authenticated user's cart based on the total cart price.
// @Tags payments
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Security ApiKeyAuth
// @Success 200 {json} SuccessResponse
// @Failure 400 {string} ErrorResponse
// @Router /user/razorpay [get]
func RazorPay(c *gin.Context) {
	id := 1
	var user models.User
	err := database.DB.First(&user, id).Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": "This user didn't find",
		})
		return
	}
	var totalprice uint
	err = database.DB.Table("carts").Select("SUM(total_price)").Where("user_id=?", id).Scan(&totalprice).Error
	if err != nil {
		c.JSON(400, gin.H{
			"error":   "Failed to find total",
			"message": "Please check your cart",
		})
		return
	}
	client := razorpay.NewClient(os.Getenv("RAZOR_kEY"), os.Getenv("RAZOR_SECRET"))
	data := map[string]interface{}{
		"amount":   totalprice * 100,
		"currency": "INR",
		"receipt":  "some_receipt_id",
	}
	body, err := client.Order.Create(data, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": err,
		})
		return
	}
	value := body["id"]
	c.HTML(http.StatusOK, "app.html", gin.H{
		"userid":     id,
		"totalprice": totalprice,
		"paymentid":  value,
	})
}

// RazorPaySuccess handles the successful completion of a RazorPay payment and processes the order.
// @Summary Process RazorPay payment and create an order
// @Description Handles the successful completion of a RazorPay payment, creates an order, deducts stock quantities, and clears the user's cart.
// @Tags payments
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Security ApiKeyAuth
// @Param order_id query string true "RazorPay order ID"
// @Param payment_id query string true "RazorPay payment ID"
// @Param signature query string true "RazorPay signature"
// @Param total query string true "Total amount paid"
// @Success 200 {string} SuccessResponse
// @Failure 400 {string} ErrorResponse
// @Router /user/payment/success [post]
func RazorPaySuccess(c *gin.Context) {
	//getting user details from middleware
	user, _ := c.Get("user")
	userId := user.(models.User).User_ID

	orderid := c.Query("order_id")
	paymentid := c.Query("payment_id")
	signature := c.Query("signature")
	totalamount := c.Query("total")

	err := database.DB.Create(&models.RazorPay{
		User_id:          uint(userId),
		RazorPayment_id:  paymentid,
		Signature:        signature,
		RazorPayOrder_id: orderid,
		AmountPaid:       totalamount,
	}).Error
	if err != nil {
		c.JSON(400, gin.H{
			"Error": err.Error(),
		})
		return
	}

	//searching for database all cart data
	var cartdata []models.Cart
	err = database.DB.Where("user_id=?", userId).Find(&cartdata).Error
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
		Payment_Type:   "RAZOR PAY",
		Total_Amount:   totalprice,
		Payment_Status: "Completed",
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
		database.DB.Table("products").Select("brand_name,catagory_name").Where("id=?", cartdata.Product_ID).Scan(&cartbrand)
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

// Success handles the successful payment response and renders a success HTML page.
// @Summary Handle successful payment response
// @Description Handles a successful payment response, rendering an HTML success page.
// @Tags payments
// @Produce html
// @Param id query int true "Payment ID"
// @Success 200 {html} HTML "Success page"
// @Failure 400 {string} ErrorResponse
// @Router /user/success [get]
func Success(c *gin.Context) {

	pid, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
	}

	c.HTML(200, "success.html", gin.H{
		"paymentid": pid,
	})
}
