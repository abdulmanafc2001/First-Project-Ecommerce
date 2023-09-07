package controllers

import (
	"strconv"

	"github.com/abdulmanafc2001/First-Project-Ecommerce/database"
	"github.com/abdulmanafc2001/First-Project-Ecommerce/models"
	"github.com/gin-gonic/gin"
)
type catagoryoffers struct {
	Catagory_ID uint `json:"catagoryid"`
	Percentage  uint `json:"percentage"`
}

// AddCatagoryOffer adds a new category offer.
// @Summary Add a category offer
// @Description Adds a new category offer with a given percentage.
// @Tags category-offers
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer Token"
// @Param catagory_id body catagoryoffers true "Category ID and Offer Percentage"
// @Success 200 {json} SuccessResponse
// @Failure 400 {json} ErrorResponse
// @Failure 401 {json} ErrorResponse
// @Router /admin/addcatagoryoffer [post]
func AddCatagoryOffer(c *gin.Context) {
	
	var catagoryoffer catagoryoffers
	if err := c.Bind(&catagoryoffer); err != nil {
		c.JSON(400, gin.H{
			"error": "Binding error",
		})
		return
	}

	var offer models.Catagory_Offer
	row := database.DB.Where("catagory_id = ? AND offer = true", catagoryoffer.Catagory_ID).First(&offer).RowsAffected
	if row > 0 {
		c.JSON(400, gin.H{
			"error": "This offer already exist",
		})
		return
	}

	if catagoryoffer.Percentage < 1 && catagoryoffer.Percentage > 99 {
		c.JSON(400, gin.H{
			"error": "Invalid offer percentage",
		})
		return
	}
	var catagory models.Catagory
	err := database.DB.First(&catagory, catagoryoffer.Catagory_ID).Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": "This catagory doesn't exist in your database",
		})
		return
	}
	err = database.DB.Create(&models.Catagory_Offer{
		Catagory_Id: catagoryoffer.Catagory_ID,
		Offer:       true,
		Percentage:  catagoryoffer.Percentage,
	}).Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": "database error",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "Successfully added catagory offer",
	})
}

// ListCatagoryOffer lists category offers.
// @Summary List category offers
// @Description Retrieves a list of category offers including category names, offer status, and offer percentages.
// @Tags category-offers
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer Token"
// @Success 200 {json} SuccessResponse
// @Failure 400 {json} ErrorResponse
// @Router /admin/listcatagoryoffer [get]
func ListCatagoryOffer(c *gin.Context) {
	type catagoryoffer struct {
		Catagory_Name string
		Offer         bool
		Percentage    uint
	}
	var offers []catagoryoffer
	database.DB.Table("catagory_offers").Select("catagories.catagory_name,catagory_offers.offer,catagory_offers.percentage").
		Joins("INNER JOIN catagories ON catagories.catagory_id=catagory_offers.catagory_id").
		Scan(&offers)

	c.JSON(200, gin.H{
		"message": offers,
	})
}

// DeleteCatagoryOffer cancels a category offer.
// @Summary Cancel a category offer
// @Description Cancels a category offer by its offer ID.
// @Tags category-offers
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer Token"
// @Param offer_id path int true "Offer ID" Format(int64)
// @Success 200 {json} SuccessResponse
// @Failure 400 {json} ErrorResponse
// @Router /admin/cancellcatagoryoffer/{offer_id} [post]
func DeleteCatagoryOffer(c *gin.Context) {
	offid, err := strconv.Atoi(c.Param("offer_id"))
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Failed to find offerid",
		})
		return
	}

	var offer models.Catagory_Offer
	row := database.DB.First(&offer, offid).RowsAffected
	if row == 0 {
		c.JSON(400, gin.H{
			"error": "Failed to find offer in database",
		})
		return
	}
	if !offer.Offer {
		c.JSON(400, gin.H{
			"error": "This offer already cancelled",
		})
		return
	}
	err = database.DB.Model(&models.Catagory_Offer{}).Where("id=?", offid).Update("offer", false).Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Updation error",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "Successfully cancelled offer",
	})
}
