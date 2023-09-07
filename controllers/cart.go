package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/abdulmanafc2001/First-Project-Ecommerce/database"
	"github.com/abdulmanafc2001/First-Project-Ecommerce/models"
	"github.com/gin-gonic/gin"
)

type carts struct {
	ProductId uint `json:"productid"`
	Quantity  uint `json:"quantity"`
}

// AddToCart adds a product to the user's shopping cart.
// @Summary Add product to cart
// @Description Add a product to the user's shopping cart with quantity
// @Tags carts
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param cart body carts true "cart data"
// @Success 200 {json} SuccessResponse
// @Failure 400 {json} ErrorResponse
// @Router /user/addtocart [post]
func AddToCart(c *gin.Context) {
	user, _ := c.Get("user")
	userId := user.(models.User).User_ID

	var cart carts
	if err := c.Bind(&cart); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to bind",
		})
		c.Abort()
		return
	}
	var product models.Product
	err := database.DB.First(&product, cart.ProductId).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "This product not found",
		})
		return
	}
	var id []uint
	database.DB.Table("catagory_offers").Select("id").Where("offer=true").Scan(&id)
	count := 0
	for _, v := range id {
		if v == product.Catagory_ID {
			count = 1
			break
		}
	}
	if count == 1 {
		var dbcart models.Cart
		err = database.DB.Where("product_id=? AND user_id=?", cart.ProductId, userId).First(&dbcart).Error

		var result models.Catagory_Offer
		database.DB.Where("catagory_id=?", product.Catagory_ID).First(&result)

		if err != nil {
			err = database.DB.Create(&models.Cart{
				Product_ID:     cart.ProductId,
				Quantity:       cart.Quantity,
				Price:          product.Price,
				Catagory_ID:    product.Catagory_ID,
				Total_Price:    (product.Price * cart.Quantity) - ((product.Price * cart.Quantity) * result.Percentage / 100),
				User_ID:        userId,
				Catagory_Offer: (product.Price * cart.Quantity) * result.Percentage / 100,
			}).Error
			if err != nil {
				c.JSON(400, gin.H{
					"error": err.Error(),
				})
				return
			}
			c.JSON(200, gin.H{
				"message": "success fully added to cart",
			})
			return
		}
		var discount uint
		database.DB.Table("carts").Select("catagory_offer").Where("product_id=? AND user_id=?", cart.ProductId, userId).Scan(&discount)
		total := ((product.Price) - ((product.Price) * result.Percentage / 100)) * cart.Quantity
		err = database.DB.Model(&models.Cart{}).Where("product_id=? AND user_id=?", cart.ProductId, userId).Updates(map[string]interface{}{"quantity": dbcart.Quantity + cart.Quantity, "total_price": dbcart.Total_Price + total, "catagory_offer": discount + ((product.Price) * result.Percentage / 100)}).Error
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"message": "successfully updated cart",
		})
	} else {

		var dbcart models.Cart
		err = database.DB.Where("product_id=? AND user_id=?", cart.ProductId, userId).First(&dbcart).Error

		if err != nil {
			err = database.DB.Create(&models.Cart{
				Product_ID:  cart.ProductId,
				Quantity:    cart.Quantity,
				Price:       product.Price,
				Catagory_ID: product.Catagory_ID,
				Total_Price: product.Price * cart.Quantity,
				User_ID:     userId,
			}).Error
			if err != nil {
				c.JSON(400, gin.H{
					"error": err.Error(),
				})
				return
			}
			c.JSON(200, gin.H{
				"message": "success fully added to cart",
			})
			return
		}
		total := product.Price * cart.Quantity
		err = database.DB.Model(&models.Cart{}).Where("product_id=? AND user_id=?", cart.ProductId, userId).Updates(map[string]interface{}{"quantity": dbcart.Quantity + cart.Quantity, "total_price": dbcart.Total_Price + total}).Error
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"message": "successfully updated cart",
		})
	}
}

// ListCart lists the products in the user's shopping cart.
// @Summary List products in cart
// @Description Retrieve a list of products in the user's shopping cart with pagination.
// @Tags carts
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer Token"
// @Param page path int true "Page number for pagination"
// @Success 200 {json} SuccessResponse
// @Failure 400 {json} ErrorResponse
// @Router /user/viewcart/{page} [get]
func ListCart(c *gin.Context) {
	user, _ := c.Get("user")
	userId := user.(models.User).User_ID

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		c.JSON(500, gin.H{
			"error": "query error",
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

	type cart struct {
		ID              uint `json:"productid"`
		Catagory_ID     uint
		Product_name    string
		Quantity        uint
		Price           string
		Total_price     uint
		Coupon_Discount uint
		Catagory_Offer  uint
		Image           string
	}
	var totalprice int
	var carts []cart
	err = database.DB.Table("carts").
		Select("products.id,products.catagory_id,products.product_name,carts.quantity,carts.price,carts.total_price,carts.coupon_discount,carts.catagory_offer,images.image").
		Joins("INNER JOIN products ON products.id=carts.product_id").Joins("INNER JOIN images ON images.product_id=carts.product_id").Where("carts.user_id=?", userId).
		Limit(limit).Offset(offset).
		Scan(&carts).Error

	if err != nil {
		c.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}
	database.DB.Table("carts").Select("SUM(total_price)").Where("user_id=?", userId).Scan(&totalprice)
	c.JSON(200, gin.H{
		"total price": totalprice,
		"cartitems":   carts,
	})
}

// DeleteFromCart deletes a product from the user's shopping cart.
// @Summary Delete a product from cart
// @Description Delete a product from the user's shopping cart by providing the cart item's ID.
// @Tags carts
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer Token"
// @Param cart_id path int true "ID of the cart item to delete"
// @Success 200 {json} SuccessResponse
// @Failure 400 {json} ErrorResponse
// @Router /user/deletefromcart/{cart_id} [post]
func DeleteFromCart(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("cart_id"))
	user, _ := c.Get("user")
	userId := user.(models.User).User_ID

	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	fmt.Println(userId, id)
	result := database.DB.Where("id=? AND user_id=?", id, userId).Delete(&models.Cart{})
	fmt.Println(result.RowsAffected)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "successfully deleted from cart",
	})
}

type cart struct {
	CartId      uint `json:"cartid"`
	NewQuantity uint `json:"newquantity"`
}

// UpdateCartQuantity updates the quantity of a product in the user's shopping cart.
// @Summary Update cart item quantity
// @Description Update the quantity of a product in the user's shopping cart by providing the cart item's ID and the new quantity.
// @Tags carts
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer Token"
// @Param updateData body carts true "cart data"
// @Success 200 {json} SuccessResponse
// @Failure 400 {json} ErrorResponse
// @Router /user/updatecartquantity [put]
func UpdateCartQuantity(c *gin.Context) {
	user, _ := c.Get("user")
	id := user.(models.User).User_ID
	var updateData cart
	if err := c.Bind(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	if updateData.NewQuantity <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "quantity must postive value",
		})
		return
	}
	var dtcart models.Cart
	err := database.DB.First(&dtcart, updateData.CartId).Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
	}
	err = database.DB.Model(&models.Cart{}).Where("user_id=? AND id=?", id, updateData.CartId).Updates(map[string]interface{}{"quantity": updateData.NewQuantity, "total_price": updateData.NewQuantity * dtcart.Price}).Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "successfully updated quantity",
	})
}
