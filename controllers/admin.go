package controllers

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/abdulmanafc2001/First-Project-Ecommerce/database"
	helper "github.com/abdulmanafc2001/First-Project-Ecommerce/helpers"
	"github.com/abdulmanafc2001/First-Project-Ecommerce/models"
	"github.com/gin-gonic/gin"
)

type userDetail struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AdminLogin authenticates the admin user and generates an authentication token.
// @Summary Admin login
// @Description Authenticate an admin user and generate a JWT token
// @Tags admin
// @Accept json
// @Produce json
// @Param userCredentials body userDetail true "Admin credentials"
// @Success 200 {json} SuccessResponse
// @Failure 401 {json} ErrorResponse
// @Failure 500 {json} ErrorResponse
// @Router /admin/login [post]
func AdminLogin(c *gin.Context) {
	go func(){
		<-c.Done()
	}()
	time.Sleep(time.Second * 10)
	var userCredentials userDetail
	if err := c.Bind(&userCredentials); err != nil {
		fmt.Println(err)
		return
	}
	//geting data from env files
	username := os.Getenv("ADMIN")
	password := os.Getenv("ADMIN_PASSWORD")
	//checking username and password
	if username != userCredentials.Username || password != userCredentials.Password {
		c.JSON(401, gin.H{
			"error": "Unautharized access incorrect username or password",
		})
		return
	}
	//generate token
	token, err := helper.GenerateJWTToken(username, "admin", "", 0)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Failed to generate token",
		})
		return
	}
	//set token into browser
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("jwt_admin", token, 3600*24, "", "", true, true)
	//success message
	c.JSON(200, gin.H{
		"message": "Successfully logged admin",
	})
}

// ListUsers lists all users.
// @Summary List users
// @Description Retrieve a list of all users
// @Tags admin
// @Accept json
// @Produce json
// @Success 200 {json} UsersResponse
// @Failure 400 {json} ErrorResponse
// @Router /admin/viewusers [get]
func ListUsers(c *gin.Context) {
	type user struct {
		User_ID      uint
		First_Name   string
		Last_Name    string
		User_Name    string
		Email        string
		Is_Blocked   bool
		Phone_Number string
		Wallet       uint
	}
	var users []user
	//searching for the data from database
	result := database.DB.Table("users").Select("user_id,first_name,last_name,user_name,email,is_blocked,phone_number,wallet").Scan(&users)

	if result.Error != nil {
		c.JSON(400, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"users": users,
	})
}

// AdminLogout logs out the admin user.
// @Summary Admin logout
// @Description Log out the admin user by clearing the JWT token cookie
// @Tags admin
// @Accept json
// @Produce json
// @Success 200 {json} SuccessResponse
// @Router /admin/logout [get]
func AdminLogout(c *gin.Context) {
	//set cookie to nil
	c.SetCookie("jwt_admin", "", -1, "", "", false, true)
	c.JSON(200, gin.H{
		"message": "successfully admin logout",
	})
}

// BlockUser blocks a user by their user ID.
// @Summary Block a user
// @Description Block a user by their user ID
// @Tags admin
// @Accept json
// @Produce json
// @Param user_id path int true "User ID to block"
// @Success 200 {json} SuccessResponse
// @Failure 400 {json} ErrorResponse
// @Router /admin/blockuser/{user_id} [post]
func BlockUser(c *gin.Context) {
	//get param
	id := c.Param("user_id")
	intid, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	//checking user already blocked or not
	var user models.User
	result := database.DB.First(&user, intid)

	if result.Error != nil {
		c.JSON(400, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	if user.IsBlocked {
		c.JSON(400, gin.H{
			"error": "This user already blocked",
		})
		return
	}
	//blocking the user
	result = database.DB.Model(&models.User{}).Where("user_id=?", intid).Update("is_blocked", true)

	if result.Error != nil {
		c.JSON(400, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "successfully blocked " + user.User_Name,
	})
}

// UnblockUser unblocks a user by their user ID.
// @Summary Unblock a user
// @Description Unblock a user by their user ID
// @Tags admin
// @Accept json
// @Produce json
// @Param user_id path int true "User ID to unblock"
// @Success 200 {json} SuccessResponse
// @Failure 400 {json} ErrorResponse
// @Router /admin/unblockuser/{user_id} [post]
func UnblockUser(c *gin.Context) {
	id := c.Param("user_id")
	intid, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	var user models.User
	result := database.DB.First(&user, intid)

	if result.Error != nil {
		c.JSON(400, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	if !user.IsBlocked {
		c.JSON(400, gin.H{
			"error": "This user is already unbloacked",
		})
		return
	}
	result = database.DB.Model(&models.User{}).Where("user_id", intid).Update("is_blocked", false)

	if result.Error != nil {
		c.JSON(400, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "successfully unbloacked " + user.User_Name,
	})

}
