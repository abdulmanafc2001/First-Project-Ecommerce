package controllers

import (
	"net/http"

	"github.com/abdulmanafc2001/First-Project-Ecommerce/database"
	"github.com/abdulmanafc2001/First-Project-Ecommerce/models"
	"github.com/gin-gonic/gin"
)

// UserDetail returns details of the authenticated user.
// @Summary Get User Details
// @Description Get details of the authenticated user including first name, last name, username, email, phone number, and wallet balance.
// @Tags users
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {json} SuccessResponse
// @Failure 401 {string} ErrorResponse
// @Router /user/userdetail [get]
func UserDetail(c *gin.Context) {
	user, _ := c.Get("user")
	c.JSON(200, gin.H{
		"firstname": user.(models.User).First_Name,
		"lastname":  user.(models.User).Last_Name,
		"username":  user.(models.User).User_Name,
		"email":     user.(models.User).Email,
		"Phone":     user.(models.User).Phone_Number,
		"wallet":    user.(models.User).Wallet,
	})
}

type userInputs struct {
	OldPassword string `json:"oldpassword"`
	NewPassword string `json:"newpassword"`
}
// ChangePassword allows the authenticated user to change their password.
// @Summary Change Password
// @Description Change the user's password with a new one.
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body userInputs true "Old and New Passwords"
// @Success 200 {string} SuccessResponse
// @Failure 400 {string} ErrorResponse
// @Failure 401 {string} ErrorResponse
// @Router /user/changepassword [post]
func ChangePassword(c *gin.Context) {
	user, _ := c.Get("user")

	
	var userInput userInputs
	if err := c.Bind(&userInput); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	hashedPassword := user.(models.User).Password
	if status := compareHashPassword(hashedPassword, userInput.OldPassword); !status {
		c.JSON(400, gin.H{
			"error": "Invalid old Password",
		})
		return
	}
	if len(userInput.NewPassword) < 5 || len(userInput.NewPassword) > 20 {
		c.JSON(400, gin.H{
			"error": "please input validpassword with in 5 to 20 charecter",
		})
		return
	}
	password, err := hashPassword(userInput.NewPassword)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	err = database.DB.Model(&models.User{}).Where("user_id=?", user.(models.User).User_ID).Update("password", password).Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "successfully updated password",
	})
}

type inputDatas struct {
	NewFirstName  string `json:"newfirstname"`
	NewSecondName string `json:"newsecondname"`
	NewUserName   string `json:"newusername"`
	NewEmail      string `json:"newemail"`
	NewPhone      string `json:"newphone"`
}

// UpdateProfile allows the authenticated user to update their profile information.
// @Summary Update Profile
// @Description Update the user's profile with new information.
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body inputDatas true "New Profile Information"
// @Success 200 {string} SuccessResponse
// @Failure 400 {string} ErrorResponse
// @Failure 401 {string} ErrorResponse
// @Router /user/updateprofile [put]
func UpdateProfile(c *gin.Context) {
	user, _ := c.Get("user")
	userId := user.(models.User).User_ID
	
	var userInput inputDatas
	if err := c.Bind(&userInput); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Binding error",
		})
		return
	}

	if len(userInput.NewPhone) != 10 {
		c.JSON(400, gin.H{
			"error": "Invalid phone number",
		})
		return
	}

	var userdt models.User
	database.DB.Where("user_name=?", userInput.NewUserName).First(&userdt)
	if userdt.Email == userInput.NewEmail {
		c.JSON(400, gin.H{
			"error": "this username is already taken",
		})
		return
	}

	database.DB.Where("phone_number=?", userInput.NewPhone).First(&userdt)
	if userdt.Phone_Number == userInput.NewPhone {
		c.JSON(400, gin.H{
			"error": "This phone number already taken please choose another one",
		})
		return
	}

	database.DB.Where("email=?", userInput.NewEmail).First(&userdt)
	if userdt.Email == userInput.NewEmail {
		c.JSON(400, gin.H{
			"error": "This Email already taken please choose different",
		})
		return
	}

	err := database.DB.Model(&models.User{}).Where("user_id=?", userId).Updates(map[string]interface{}{
		"first_name":   userInput.NewFirstName,
		"last_name":    userInput.NewSecondName,
		"user_name":    userInput.NewUserName,
		"email":        userInput.NewEmail,
		"phone_number": userInput.NewPhone,
	}).Error

	if err != nil {
		c.JSON(400,gin.H{
			"error":err.Error(),
		})
		return
	}
	c.JSON(200,gin.H{
		"message":"successfully updated "+ userInput.NewUserName + " your data",
	})
}
