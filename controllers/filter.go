package controllers

import (
	"github.com/abdulmanafc2001/First-Project-Ecommerce/database"
	"github.com/abdulmanafc2001/First-Project-Ecommerce/models"
	"github.com/gin-gonic/gin"
)

// SortWithAscending sorts and retrieves a list of products in ascending order of their prices.
// @Summary Sort products by ascending price
// @Description Retrieves a list of products sorted in ascending order of their prices.
// @Tags products
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Security ApiKeyAuth
// @Success 200 {json} SuccessResponse
// @Failure 400 {string} ErrorResponse
// @Router /user/ascendingfilter [get]
func SortWithAscending(c *gin.Context) {
	type product struct {
		ID            uint
		Product_Name  string
		Price         uint
		Brand_Name    string
		Catagory_Name string
		Stock         uint
	}
	var products []product
	err := database.DB.Table("products").Select("products.id,products.product_name,products.price,brands.brand_name,catagories.catagory_name,products.stock").
		Joins("INNER JOIN brands ON brands.brand_id=products.brand_id").Joins("INNER JOIN catagories ON catagories.catagory_id=products.catagory_id").
		Order("products.price asc").Scan(&products).Error

	if err != nil {
		c.JSON(400, gin.H{
			"error": "database error",
		})
		return
	}
	c.JSON(200, gin.H{
		"product": products,
	})
}

// SortWithDescending sorts and retrieves a list of products in descending order of their prices.
// @Summary Sort products by descending price
// @Description Retrieves a list of products sorted in descending order of their prices.
// @Tags products
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Security ApiKeyAuth
// @Success 200 {json} SuccessResponse
// @Failure 400 {string} ErrorResponse
// @Router /user/descendingfilter [get]
func SortWithDescending(c *gin.Context) {
	type product struct {
		ID            uint
		Product_Name  string
		Price         uint
		Brand_Name    string
		Catagory_Name string
		Stock         uint
	}
	var products []product
	err := database.DB.Table("products").Select("products.id,products.product_name,products.price,brands.brand_name,catagories.catagory_name,products.stock").
		Joins("INNER JOIN brands ON brands.brand_id=products.brand_id").Joins("INNER JOIN catagories ON catagories.catagory_id=products.catagory_id").
		Order("products.price desc").Scan(&products).Error

	if err != nil {
		c.JSON(400, gin.H{
			"error": "database error",
		})
		return
	}
	c.JSON(200, gin.H{
		"product": products,
	})
}

// FilterWithBrands filters products by brand name.
// @Summary Filter products by brand name
// @Description Retrieves a list of products associated with the specified brand name.
// @Tags products
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param brand_name query string true "Brand name to filter by"
// @Security ApiKeyAuth
// @Success 200 {json} SuccessResponse
// @Failure 400 {string} ErrorResponse
// @Router /user/brandfilter [get]
func FilterWithBrands(c *gin.Context) {
	brand := c.Query("brand_name")

	var branddb models.Brand
	err := database.DB.Where("brand_name=?", brand).First(&branddb).Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": "This brand not exist",
		})
		return
	}
	type product struct {
		ID            uint
		Product_Name  string
		Price         uint
		Brand_Name    string
		Catagory_Name string
		Stock         uint
	}
	var products []product
	err = database.DB.Table("products").Select("products.id,products.product_name,products.price,brands.brand_name,catagories.catagory_name,products.stock").
		Joins("INNER JOIN brands ON brands.brand_id=products.brand_id").Joins("INNER JOIN catagories ON catagories.catagory_id=products.catagory_id").
		Where("products.brand_id=?", branddb.Brand_ID).Scan(&products).Error

	if err != nil {
		c.JSON(400, gin.H{
			"error": "database error",
		})
		return
	}
	c.JSON(200, gin.H{
		"products": products,
	})

}

// FilterWithCatagories filters products by category name.
// @Summary Filter products by category name
// @Description Retrieves a list of products associated with the specified category name.
// @Tags products
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param catagory_name query string true "Category name to filter by"
// @Security ApiKeyAuth
// @Success 200 {json} SuccessResponse
// @Failure 400 {string} ErrorResponse
// @Router /user/catagoryfilter [get]
func FilterWithCatagories(c *gin.Context) {
	catagory := c.Query("catagory_name")

	var catagorydb models.Catagory
	err := database.DB.Where("catagory_name=?", catagory).First(&catagorydb).Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": "This catagory not exist",
		})
		return
	}
	type product struct {
		ID            uint
		Product_Name  string
		Price         uint
		Brand_Name    string
		Catagory_Name string
		Stock         uint
	}
	var products []product
	err = database.DB.Table("products").Select("products.id,products.product_name,products.price,brands.brand_name,catagories.catagory_name,products.stock").
		Joins("INNER JOIN brands ON brands.brand_id=products.brand_id").Joins("INNER JOIN catagories ON catagories.catagory_id=products.catagory_id").
		Where("products.catagory_id=?", catagorydb.Catagory_ID).Scan(&products).Error

	if err != nil {
		c.JSON(400, gin.H{
			"error": "database error",
		})
		return
	}
	c.JSON(200, gin.H{
		"products": products,
	})
}
