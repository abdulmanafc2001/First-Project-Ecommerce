package controllers

import (
	"net/http"
	"strconv"

	"github.com/abdulmanafc2001/First-Project-Ecommerce/database"
	"github.com/abdulmanafc2001/First-Project-Ecommerce/models"
	"github.com/gin-gonic/gin"
)
type datas struct {
	Productid uint `json:"productid"`
}

// AddToWishlist allows the authenticated user to add a product to their wishlist.
// @Summary Add Product to Wishlist
// @Description Add a product to the user's wishlist.
// @Tags wishlist
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body datas true "Product ID"
// @Success 200 {string} SuccessResponse
// @Failure 400 {string} ErrorResponse
// @Failure 401 {string} ErrorResponse
// @Router /user/addtowishlist [post]
func AddToWishtList(c *gin.Context) {
	user, _ := c.Get("user")
	userId := user.(models.User).User_ID
	
	var data datas
	if err := c.Bind(&data); err != nil {
		c.JSON(500, gin.H{
			"error": "Binding error",
		})
		return
	}
	var wishlist models.Wishlist
	row := database.DB.Where("product_id=?", data.Productid).First(&wishlist).RowsAffected
	if row > 0 {
		c.JSON(400, gin.H{
			"error": "This product already in wishlist",
		})
		return
	}
	var product models.Product
	err := database.DB.First(&product, data.Productid).Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": "This produc doesnot exist in database",
		})
		return
	}
	err = database.DB.Create(&models.Wishlist{
		Product_ID: data.Productid,
		User_ID:    userId,
	}).Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Database error",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "successfully " + product.Product_Name + " added to wishlist",
	})
}

// ListWishlist retrieves the authenticated user's wishlist items.
// @Summary List Wishlist Items
// @Description Get the user's wishlist items paginated.
// @Tags wishlist
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page path int true "Page number for pagination"
// @Success 200 {json} SuccessResponse
// @Failure 400 {json} ErrorResponse
// @Failure 401 {json} ErrorResponse
// @Router /user/listwishlist/{page} [get]
func ListWishlist(c *gin.Context) {
	user, _ := c.Get("user")
	userId := user.(models.User).User_ID
	page, _ := strconv.Atoi(c.Param("page"))
	limit := 3
	offset := (page - 1) * limit

	type wishlist struct {
		Product_Name string
		Price        uint
	}
	var wishlists []wishlist
	err := database.DB.Table("wishlists").Select("products.product_name,products.price").
		Joins("INNER JOIN products ON products.id=wishlists.product_id").
		Where("user_id=?", userId).
		Limit(limit).Offset(offset).
		Scan(&wishlists).Error

	if err != nil {
		c.JSON(400, gin.H{
			"error": "database error",
		})
		return
	}
	c.JSON(200, gin.H{
		"wishlist": wishlists,
	})
}

type carts1 struct {
	Wishlist_ID uint `json:"wishlistid"`
	Quantity    uint `json:"quantity"`
}
// AddToCartFromWishlist adds a product from the user's wishlist to the cart.
// @Summary Add Product from Wishlist to Cart
// @Description Add a product from the wishlist to the cart.
// @Tags cart
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param cart body carts1 true "cart details"
// @Success 200 {json} SuccessResponse
// @Failure 400 {json} ErrorResponse
// @Failure 401 {json} ErrorResponse
// @Router /wishlist/addtocart [post]
func AddTocartFromWishlist(c *gin.Context) {
	user, _ := c.Get("user")
	userId := user.(models.User).User_ID

	
	var cart carts1
	if err := c.Bind(&cart); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to bind",
		})
		return
	}
	var wishlist models.Wishlist
	err := database.DB.Where("wishtlist_id=? AND user_id=?", cart.Wishlist_ID, userId).First(&wishlist).Error

	if err != nil {
		c.JSON(400, gin.H{
			"error": "Product didn't find in wishlis",
		})
		return
	}
	var product models.Product
	err = database.DB.First(&product, wishlist.Product_ID).Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Failed to find this product",
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
		err = database.DB.Where("product_id=? AND user_id=?", wishlist.Product_ID, userId).First(&dbcart).Error

		var result models.Catagory_Offer
		database.DB.Where("catagory_id=?", product.Catagory_ID).First(&result)

		if err != nil {
			err = database.DB.Create(&models.Cart{
				Product_ID:  wishlist.Product_ID,
				Quantity:    cart.Quantity,
				Price:       (product.Price) - ((product.Price) * result.Percentage / 100),
				Catagory_ID: product.Catagory_ID,
				Total_Price: (product.Price * cart.Quantity) - ((product.Price * cart.Quantity) * result.Percentage / 100),
				User_ID:     userId,
				Catagory_Offer:    (product.Price * cart.Quantity) * result.Percentage / 100,
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
		database.DB.Table("carts").Select("discount").Where("product_id=? AND user_id=?", wishlist.Product_ID, userId).Scan(&discount)
		total := ((product.Price) - ((product.Price) * result.Percentage / 100)) * cart.Quantity
		err = database.DB.Model(&models.Cart{}).Where("product_id=? AND user_id=?", wishlist.Product_ID, userId).Updates(map[string]interface{}{"quantity": dbcart.Quantity + cart.Quantity, "total_price": dbcart.Total_Price + total, "catagory_offer": discount + ((product.Price) * result.Percentage / 100)}).Error
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
		err = database.DB.Where("product_id=? AND user_id=?", wishlist.Product_ID, userId).First(&dbcart).Error

		if err != nil {
			err = database.DB.Create(&models.Cart{
				Product_ID:  wishlist.Product_ID,
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
		err = database.DB.Model(&models.Cart{}).Where("product_id=? AND user_id=?", wishlist.Product_ID, userId).Updates(map[string]interface{}{"quantity": dbcart.Quantity + cart.Quantity, "total_price": dbcart.Total_Price + total}).Error
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

// RemoveFromWishlist removes a product from the user's wishlist.
// @Summary Remove Product from Wishlist
// @Description Remove a product from the wishlist.
// @Tags wishlist
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param wishlist_id path int true "Wishlist item ID to remove"
// @Success 200 {json} SuccessResponse
// @Failure 400 {json} ErrorResponse
// @Failure 401 {json} ErrorResponse
// @Router /wishlist/delete/{wishlist_id} [post]
func RemoveFromWishlist(c *gin.Context) {
	user, _ := c.Get("user")
	userId := user.(models.User).User_ID
	wishlistId, err := strconv.Atoi(c.Param("wishlist_id"))
	if err != nil {
		c.JSON(400, gin.H{
			"errro": "integer converting error",
		})
		return
	}
	err = database.DB.Where("wishtlist_id=? AND user_id=?", wishlistId, userId).Delete(&models.Wishlist{}).Error
	if err != nil {
		c.JSON(400,gin.H{
			"error":"failed to find in wishlist",
		})
		return
	}
	c.JSON(200,gin.H{
		"message":"successfully deleted from wishlist",
	})
}
