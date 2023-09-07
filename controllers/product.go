package controllers

import (
	"path/filepath"
	"strconv"

	"github.com/abdulmanafc2001/First-Project-Ecommerce/database"
	"github.com/abdulmanafc2001/First-Project-Ecommerce/models"
	"github.com/gin-gonic/gin"
)
// AddProduct adds a new product to the database.
// @Summary Add a new product
// @Description Adds a new product to the database.
// @Tags products
// @Accept json
// @Produce json
// @Param product body models.Product true "Product information"
// @Success 200 {string} SuccessResponse
// @Failure 400 {string} ErrorResponse
// @Router /admin/addproduct [post]
func AddProduct(c *gin.Context) {
	var product models.Product

	if err := c.Bind(&product); err != nil {
		c.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}
	var dtproduct models.Product
	database.DB.Where("product_name=?", product.Product_Name).First(&dtproduct)

	if product.Product_Name == dtproduct.Product_Name {
		c.JSON(400, gin.H{
			"error": "this product already exist",
		})
		return
	}
	var catagory models.Catagory
	database.DB.First(&catagory, product.Catagory_ID)
	var brand models.Brand
	database.DB.First(&brand, product.Catagory_ID)

	result := database.DB.Create(&models.Product{
		Product_Name: product.Product_Name,
		Descreption:  product.Descreption,
		Stock:        product.Stock,
		Price:        product.Price,
		Catagory_ID:  product.Catagory_ID,
		Brand_ID:     product.Brand_ID,
	})
	if result.Error != nil {
		c.JSON(400, gin.H{
			"error":   result.Error.Error(),
			"message": "failed to create",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "successfully created " + product.Product_Name,
	})
}

type product struct {
	Id             uint   `json:"id"`
	Newname        string `json:"newname"`
	Newdescription string `json:"newdescription"`
	Newstock       uint   `json:"newstock"`
	Newprice       uint   `json:"newprice"`
}

// EditProduct updates an existing product in the database.
// @Summary Update an existing product
// @Description Updates an existing product in the database.
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param productdetails body product true "Product details to update"
// @Success 200 {string} SuccessResponse
// @Failure 400 {string} ErrorResponse
// @Router /admin/editproduct [put]
func EditProduct(c *gin.Context) {
	
	var productdetails product
	if err := c.Bind(&productdetails); err != nil {
		c.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}
	var check models.Product
	database.DB.Where("product_name", productdetails.Newname).First(&check)
	if check.Product_Name == productdetails.Newname {
		c.JSON(400, gin.H{
			"err": "This name already in our products",
		})
		return
	}
	result := database.DB.Model(&models.Product{}).Where("id", productdetails.Id).Updates(map[string]interface{}{"product_name": productdetails.Newname, "descreption": productdetails.Newdescription, "stock": productdetails.Newstock, "price": productdetails.Newprice})
	if result.Error != nil {
		c.JSON(400, gin.H{
			"err": result.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "successfully updated " + productdetails.Newname,
	})
}

// DeleteProduct deletes a product from the database by its ID.
// @Summary Delete a product by ID
// @Description Deletes a product from the database by its ID.
// @Tags products
// @Accept json
// @Produce json
// @Param product_id path int true "Product ID to delete"
// @Success 200 {string} SuccessResponse
// @Failure 400 {string} ErrorResponse
// @Router /admin/deleteproduct/{product_id} [post]
func DeleteProduct(c *gin.Context) {
	idstr := c.Param("product_id")
	id, _ := strconv.Atoi(idstr)
	var product models.Product
	database.DB.First(&product, id)

	result := database.DB.Delete(&models.Product{}, id)
	if result.RowsAffected == 0 {
		c.JSON(400, gin.H{
			"error": "Failed to find product",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "successfully deleted " + product.Product_Name,
	})
}

// AddImage uploads one or more images for a product and associates them with the product in the database.
// @Summary Upload product images
// @Description Uploads one or more images for a product and associates them with the product in the database.
// @Tags products
// @Accept multipart/form-data
// @Produce json
// @Param product_id formData int true "Product ID to associate images with"
// @Param files formData file true "One or more image files to upload"
// @Success 200 {string} SuccessResponse
// @Failure 400 {string} ErrorResponse
// @Router /admin/addimage [post]
func AddImage(c *gin.Context) {
	// imagepath, err := c.FormFile("image")

	// if err != nil {
	// 	c.JSON(400, gin.H{
	// 		"error": "This file path not exist",
	// 	})
	// 	return
	// }

	// prodId, err := strconv.Atoi(c.PostForm("product_id"))
	// if err != nil {
	// 	fmt.Println("conversion error")
	// 	return
	// }
	// var product models.Product
	// err = database.DB.First(&product, prodId).Error
	// if err != nil {
	// 	c.JSON(400, gin.H{
	// 		"error": "Failed to find this product",
	// 	})
	// 	return
	// }
	// extention := filepath.Ext(imagepath.Filename)
	// image := uuid.New().String() + extention
	// err = c.SaveUploadedFile(imagepath, "./public/images"+image)
	// if err != nil {
	// 	c.JSON(400, gin.H{
	// 		"error": "Failed to save image",
	// 	})
	// 	return
	// }
	// err = database.DB.Create(&models.Image{
	// 	Product_id: uint(prodId),
	// 	Image:      image,
	// }).Error
	// if err != nil {
	// 	c.JSON(400, gin.H{
	// 		"error": "Database error",
	// 	})
	// 	return
	// }
	// c.JSON(200, gin.H{
	// 	"message": "successfully uploaded image",
	// })
	prodId, err := strconv.Atoi(c.PostForm("product_id"))
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Integer conversion error",
		})
		return
	}
	var product models.Product
	err = database.DB.First(&product, prodId).Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": "This product does not exist",
		})
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	files := form.File["files"]
	for _, file := range files {
		filename := filepath.Base(file.Filename)
		if err := c.SaveUploadedFile(file, "./public/"+filename); err != nil {
			c.JSON(400, gin.H{
				"error": "image save error",
			})
			return
		}

		err = database.DB.Create(&models.Image{
			Product_id: uint(prodId),
			Image:      filename,
		}).Error
		if err != nil {
			c.JSON(400, gin.H{
				"error": "database error",
			})
			return
		}
	}
	c.JSON(200, gin.H{
		"message": "successfully uploaded image",
	})
}
