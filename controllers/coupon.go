package controllers

import (
	"strconv"
	"time"

	"github.com/abdulmanafc2001/First-Project-Ecommerce/database"
	"github.com/abdulmanafc2001/First-Project-Ecommerce/models"
	"github.com/gin-gonic/gin"
)
type coupons struct {
	Coupon_Code  string `json:"couponcode"`
	Days         uint   `json:"days"`
	Type         string `json:"type"`
	Value        uint   `json:"value"`
	Max_Discount uint   `json:"maxdiscount"`
	Min_Discount uint   `json:"mindiscount"`
}

// AddCoupon creates a new coupon with the provided details.
// @Summary Create a new coupon
// @Description Creates a new coupon with the specified attributes.
// @Tags coupons
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer Token"
// @Param coupon body coupons true "Coupon details"
// @Success 200 {json} SuccessResponse
// @Failure 400 {json} ErrorResponse
// @Router /admin/addcoupon [post]
func AddCoupon(c *gin.Context) {
	
	var coupon coupons
	if err := c.Bind(&coupon); err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	var coupon1 models.Coupon
	database.DB.Where("coupon_code=?", coupon.Coupon_Code).First(&coupon1)

	if coupon.Coupon_Code == coupon1.Coupon_Code {
		c.JSON(400, gin.H{
			"error": "This coupon code already exist in database",
		})
		return
	}

	if len(coupon.Coupon_Code) < 5 || len(coupon.Coupon_Code) > 10 {
		c.JSON(400, gin.H{
			"error": "Coupon code must be lenght between 5 to 10",
		})
		return
	}
	if coupon.Type == "fixed" || coupon.Type == "percentage" {
		database.DB.Create(&models.Coupon{
			Coupon_Code:   coupon.Coupon_Code,
			Starting_Time: time.Now(),
			Ending_Time:   time.Now().Add(time.Hour * 24 * time.Duration(coupon.Days)),
			Value:         coupon.Value,
			Type:          coupon.Type,
			Max_Discount:  coupon.Max_Discount,
			Min_Discount:  coupon.Min_Discount,
		})
		c.JSON(200, gin.H{
			"success": "successfully created coupon",
		})
	} else {
		c.JSON(400, gin.H{
			"error": "This type not applicable",
		})
	}
}

// ListCoupons retrieves a list of all available coupons.
// @Summary Retrieve a list of coupons
// @Description Retrieves a list of all available coupons.
// @Tags coupons
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer Token"
// @Success 200 {json} SuccessResponse
// @Failure 500 {json} ErrorResponse
// @Router /admin/listcoupons [get]
func ListCoupons(c *gin.Context) {
	var coupons []models.Coupon
	err := database.DB.Find(&coupons).Error
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Couldn't find coupons",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": coupons,
	})
}

// CancelCoupon cancels a coupon by marking it as canceled and updating the ending time.
// @Summary Cancel a coupon by ID
// @Description Cancels a coupon by marking it as canceled and updating the ending time.
// @Tags coupons
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer Token"
// @Param coupon_id path int true "Coupon ID to cancel"
// @Success 200 {json} SuccessResponse
// @Failure 400 {json} ErrorResponse
// @Router /admin/cancelcoupon/{coupon_id} [put]
func CancelCoupon(c *gin.Context) {
	cid, err := strconv.Atoi(c.Param("coupon_id"))
	if err != nil {
		c.JSON(400, gin.H{
			"error": "String convertion error",
		})
		return
	}
	var cou models.Coupon
	err = database.DB.First(&cou, cid).Error

	if err != nil {
		c.JSON(400, gin.H{
			"error": "Failed to find coupon please try different id",
		})
		return
	}

	err = database.DB.Model(&models.Coupon{}).Where("coupon_id=?", cid).Updates(map[string]interface{}{"cancel": true, "ending_time": time.Now()}).Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "successfully cancelled coupon",
	})
}

type coupontype struct {
	Coupon_Code string `json:"couponcode"`
}

// ApplyCoupon applies a coupon to the user's cart if it's valid and not expired.
// @Summary Apply a coupon to the user's cart
// @Description Applies a coupon to the user's cart if it's valid and not expired.
// @Tags coupons
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer Token"
// @Param couponcode body coupontype true "Coupon code to apply"
// @Success 200 {json} SuccessResponse
// @Failure 400 {json} ErrorResponse
// @Router /user/applycoupon [post]
func ApplyCoupon(c *gin.Context) {
	user, _ := c.Get("user")
	userId := user.(models.User).User_ID

	
	var coupon coupontype
	if err := c.Bind(&coupon); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	var coupon1 models.Coupon
	row := database.DB.Where("coupon_code=?", coupon.Coupon_Code).First(&coupon1).RowsAffected
	if row == 0 {
		c.JSON(400, gin.H{
			"error": "Failed to find coupon",
		})
		return
	}
	if time.Now().Unix() > (coupon1.Ending_Time).Unix() {
		c.JSON(400, gin.H{
			"error": "coupon expired",
		})
		return
	}
	var cart []models.Cart
	row = database.DB.Where("user_id=? AND coupon_applied = true", userId).Find(&cart).RowsAffected
	if row >= 1 {
		c.JSON(400, gin.H{
			"error": "Coupon already aplied",
		})
		return
	}
	var totalprice uint
	err := database.DB.Table("carts").Select("SUM(total_price)").Where("user_id=?", userId).Scan(&totalprice).Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Cart is empty",
		})
		return
	}
	var cart1 []models.Cart
	err = database.DB.Where("user_id=?", userId).Find(&cart1).Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": "cart is empty",
		})
		return
	}
	if coupon1.Cancel {
		c.JSON(400, gin.H{
			"error": "This coupon is not valid",
		})
		return
	}

	if coupon1.Type == "percentage" {

		for _, v := range cart1 {
			discount := (v.Total_Price * coupon1.Value / 100)
			if discount > coupon1.Max_Discount {
				discount = coupon1.Max_Discount
			} else if discount < coupon1.Min_Discount {
				discount = coupon1.Min_Discount
			}
			err := database.DB.Model(&models.Cart{}).Where("user_id=? AND id=?", userId, v.ID).Updates(map[string]interface{}{"total_price": v.Total_Price - discount, "coupon_discount": discount, "coupon_applied": true}).Error
			if err != nil {
				c.JSON(400, gin.H{
					"error": err.Error(),
				})
				return
			}
		}
		c.JSON(200, gin.H{
			"message": "Coupon applied successfully",
		})
	} else {
		for _, v := range cart1 {
			discount := coupon1.Value
			err := database.DB.Model(&models.Cart{}).Where("user_id=? AND id=?", userId, v.ID).Updates(map[string]interface{}{"total_price": v.Total_Price - discount, "coupon_discount": discount, "coupon_applied": true}).Error
			if err != nil {
				c.JSON(400, gin.H{
					"error": err.Error(),
				})
				return
			}
		}
	
	c.JSON(200, gin.H{
		"message": "Coupon applied successfully",
	})
	}
}
