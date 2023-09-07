package controllers

import (
	"fmt"
	"strconv"

	"github.com/abdulmanafc2001/First-Project-Ecommerce/database"
	"github.com/abdulmanafc2001/First-Project-Ecommerce/models"
	"github.com/gin-gonic/gin"
)

// AddCategory adds a new category to the database.
// @Summary Add a new category
// @Description Add a new category to the database by providing the category's name.
// @Tags categories
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer Token"
// @Param catagory_name body models.Catagory true "Name of the category to add"
// @Success 200 {json} SuccessResponse
// @Failure 400 {json} ErrorResponse
// @Router /admin/addcatagory [post]
func AddCatagory(c *gin.Context) {
	var catagory models.Catagory

	if err := c.Bind(&catagory); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		fmt.Println("binding error")
		return
	}
	var dbcat models.Catagory
	database.DB.Where("catagory_name=?", catagory.Catagory_Name).First(&dbcat)

	//checking this catagory already exist in database
	if dbcat.Catagory_Name == catagory.Catagory_Name {
		c.JSON(400, gin.H{
			"error": "This catagory already exist",
		})
		fmt.Println("catagory already exist")
		return
	}
	//adding data into database
	result := database.DB.Create(&models.Catagory{
		Catagory_Name: catagory.Catagory_Name,
	})
	if result.Error != nil {
		c.JSON(400, gin.H{
			"error": result.Error.Error(),
		})
		fmt.Println("database error")
		return
	}
	c.JSON(200, gin.H{
		"message": "successfully created " + catagory.Catagory_Name + " catagory",
	})
}

// ListCategories lists all categories from the database.
// @Summary List all categories
// @Description Retrieves a list of all categories from the database.
// @Tags categories
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer Token"
// @Success 200 {json} CategoriesResponse
// @Failure 400 {json} ErrorResponse
// @Router /admin/viewcatagories [get]
func ListCatagories(c *gin.Context) {
	var catagories []models.Catagory

	result := database.DB.Raw("SELECT * FROM catagories").Scan(&catagories)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"error": result.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"catagories": catagories,
	})
}

// BlockCategory blocks a category by its ID.
// @Summary Block a category
// @Description Blocks a category based on its ID.
// @Tags categories
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer Token"
// @Param catagory_id path integer true "Category ID" Format(int64)
// @Success 200 {json} SuccessResponse
// @Failure 400 {json} ErrorResponse
// @Failure 401 {json} ErrorResponse
// @Router /admin/blockcatagory/{catagory_id} [post]
func BlockCatagory(c *gin.Context) {
	id := c.Param("catagory_id")
	intId, err := strconv.Atoi(id)

	if err != nil {
		c.JSON(401, gin.H{
			"err": err.Error(),
		})
		return
	}

	var catagory models.Catagory
	result := database.DB.First(&catagory, intId)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"error": result.Error.Error(),
		})
		return
	}
	if catagory.Unlist {
		c.JSON(401, gin.H{
			"error": "this catagory already blocked",
		})
		return
	}
	result = database.DB.Model(&models.Catagory{}).Where("catagory_id=?", intId).Update("unlist", true)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"error": result.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "successfully blocked " + catagory.Catagory_Name + " catagory",
	})
}

// UnBlockCategory unblocks a category by its ID.
// @Summary Unblock a category
// @Description Unblocks a category based on its ID.
// @Tags categories
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer Token"
// @Param catagory_id path integer true "Category ID" Format(int64)
// @Success 200 {json} SuccessResponse
// @Failure 400 {json} ErrorResponse
// @Failure 401 {json} ErrorResponse
// @Router /admin/unblockcatagory/{catagory_id} [post]
func UnBlockCatagory(c *gin.Context) {
	id := c.Param("catagory_id")
	intId, err := strconv.Atoi(id)

	if err != nil {
		c.JSON(401, gin.H{
			"err": err.Error(),
		})
		return
	}

	var catagory models.Catagory
	result := database.DB.First(&catagory, intId)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"error": result.Error.Error(),
		})
		return
	}
	if !catagory.Unlist {
		c.JSON(401, gin.H{
			"error": "this catagory already unblocked",
		})
		return
	}
	result = database.DB.Model(&models.Catagory{}).Where("catagory_id=?", intId).Update("unlist", false)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"error": result.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "successfully unblocked " + catagory.Catagory_Name + " catagory",
	})
}
