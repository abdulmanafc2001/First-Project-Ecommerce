package controllers

import (
	"github.com/abdulmanafc2001/First-Project-Ecommerce/database"
	"github.com/abdulmanafc2001/First-Project-Ecommerce/models"
	"github.com/gin-gonic/gin"
)

type brands struct {
	BrandName string
}

// AddBrand creates a new brand.
// @Summary Create a brand
// @Description Create a new brand with a unique name
// @Tags brands
// @Accept json
// @Produce json
// @Param brand body brands true "Brand name to create"
// @Success 200 {json} SuccessResponse
// @Failure 400 {json} ErrorResponse
// @Router /admin/addbrand [post]
func AddBrand(c *gin.Context) {
	var brand brands
	if err := c.Bind(&brand); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	var dtbrand models.Brand
	database.DB.Where("brand_name=?", brand.BrandName).First(&dtbrand)

	if dtbrand.Brand_Name == brand.BrandName {
		c.JSON(400, gin.H{
			"error": "This brand already exist",
		})
		return
	}

	database.DB.Create(&models.Brand{Brand_Name: brand.BrandName})
	c.JSON(200, gin.H{
		"message": "successfully created " + brand.BrandName,
	})
}
