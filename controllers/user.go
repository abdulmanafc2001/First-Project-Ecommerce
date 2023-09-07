package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/abdulmanafc2001/First-Project-Ecommerce/database"
	helper "github.com/abdulmanafc2001/First-Project-Ecommerce/helpers"
	"github.com/abdulmanafc2001/First-Project-Ecommerce/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

// declaring a new validator for validating struct
var validate = validator.New()

// hashing password ------------------->
func hashPassword(password string) (string, error) {
	byte, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		fmt.Println("Failed to hash password")
		return "", errors.New("failed to hash password")
	}
	return string(byte), nil
}

// comparing password -------------------->
func compareHashPassword(hashedPassword string, password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return false
	}

	return true
}

// declaring variable for accessing otp validate function ----------------------->

// Signup registers a new user.
// @Summary Register a new user
// @Description Register a new user with the provided information.
// @Tags authentication
// @Accept json
// @Produce json
// @Param user body models.User true "User registration information"
// @Success 200 {string} SuccessfulResponse "User registration successful"
// @Failure 400 {string} ErrorResponse "Bad request"
// @Failure 409 {string} ErrorResponse "Conflict - Username or phone number already exists"
// @Failure 500 {string} ErrorResponse "Internal server error"
// @Router /user/signup [post]
func Signup(c *gin.Context) {
	var user models.User
	if err := c.Bind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to bind json",
		})
		return
	}
	//validate struct
	if err := validate.Struct(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	//check username already exist in database
	var dtuser models.User
	database.DB.Where("user_name=?", user.User_Name).First(&dtuser)
	if dtuser.User_Name == user.User_Name {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "This username already taken",
		})
		return
	}
	//check Phone number already exist in database
	database.DB.Where("phone_number=?", user.Phone_Number).First(&dtuser)
	if user.Phone_Number == dtuser.Phone_Number {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "This Phone number already taken Please enter valid phone number",
		})
		return
	}
	//checking length of phone number
	if len(user.Phone_Number) != 10 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Please give valid Phone number",
		})
		return
	}
	//password hashing
	password, err := hashPassword(user.Password)
	if err != nil {
		fmt.Println(err)
		return
	}
	otp := helper.GenerateOtp()
	fmt.Println(otp)
	otpstring := strconv.Itoa(otp)
	helper.SendOtp(otpstring, user.Email)

	referalcode := helper.RandomStringGenerator()

	if user.Referal_Code != "" {
		var userdata models.User
		err := database.DB.Where("referal_code = ?", user.Referal_Code).First(&userdata).Error
		if err != nil {
			c.JSON(400, gin.H{
				"error": "Failed to find referal code",
			})
			return
		}

		err = database.DB.Model(&models.User{}).Where("user_id=?", userdata.User_ID).Update("wallet", userdata.Wallet+200).Error
		if err != nil {
			c.JSON(400, gin.H{
				"error": "Updation error",
			})
			return
		}
		err = database.DB.Create(&models.User{
			First_Name:   user.First_Name,
			Last_Name:    user.Last_Name,
			User_Name:    user.User_Name,
			Password:     password,
			Email:        user.Email,
			Otp:          uint(otp),
			Phone_Number: user.Phone_Number,
			Created_at:   time.Now(),
			Referal_Code: referalcode,
			Wallet:       100,
		}).Error
		if err != nil {
			c.JSON(400, gin.H{
				"error": "database creation error",
			})
			return
		}
		c.JSON(200, gin.H{
			"message": "email sended your account click to validate /signup/validate",
		})

		return
	}
	database.DB.Create(&models.User{
		First_Name:   user.First_Name,
		Last_Name:    user.Last_Name,
		User_Name:    user.User_Name,
		Password:     password,
		Email:        user.Email,
		Otp:          uint(otp),
		Phone_Number: user.Phone_Number,
		Created_at:   time.Now(),
		Referal_Code: referalcode,
	})
	//message for validate otp
	c.JSON(200, gin.H{
		"message": "email sended your account click to validate /signup/validate",
	})
}

type validatdata struct {
	Email string `json:"email"`
	Otp   int    `json:"otp"`
}

// ValidateOtp validates the OTP sent to a user's email during registration.
// @Summary Validate OTP
// @Description Validate the OTP received via email during user registration.
// @Tags authentication
// @Accept json
// @Produce json
// @Param validate body validatdata true "Email and OTP to validate"
// @Success 200 {string} SuccessfulResponse "Account validation successful"
// @Failure 400 {string} ErrorResponse "Bad request"
// @Failure 404 {string} ErrorResponse "User not found or OTP doesn't match"
// @Failure 500 {string} ErrorResponse "Internal server error"
// @Router /user/signup/validate [post]
func ValidateOtp(c *gin.Context) {

	var validate validatdata
	if err := c.Bind(&validate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to bind",
		})
		return
	}

	var user models.User
	err := database.DB.Where("email=?", validate.Email).First(&user).Error

	//checking otp
	if user.Otp == uint(validate.Otp) && err == nil {
		err = database.DB.Model(&models.User{}).Where("email=?", validate.Email).Update("validate", true).Error
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "successfully created " + user.User_Name + " Account",
		})
	} else {
		database.DB.Where("validate=?", false).Delete(&models.User{})
		c.JSON(400, gin.H{
			"error": "Failed to find otp and email please restart from signup",
		})
	}
}

// Login authenticates a user by checking their username and password.
// @Summary User Login
// @Description Authenticate a user by verifying their username and password.
// @Tags authentication
// @Accept json
// @Produce json
// @Param credentials body userDetail true "User credentials (Username and Password)"
// @Success 200 {string} SuccessResponse "Login successful"
// @Failure 400 {string} ErrorResponse
// @Failure 401 {string} ErrorResponse
// @Failure 500 {string} ErrorResponse
// @Router /user/login [post]
func Login(c *gin.Context) {
	type userDetail struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var userCredentials userDetail
	if err := c.Bind(&userCredentials); err != nil {
		fmt.Println(err)
		return
	}
	//finding with username in database
	var user models.User
	database.DB.Where("user_name=?", userCredentials.Username).First(&user)
	//checking user is blocked or not
	if user.IsBlocked {
		c.JSON(401, gin.H{
			"error": "Unautharized access user is blocked",
		})
		return
	}
	//comparing password with database
	status := compareHashPassword(user.Password, userCredentials.Password)
	//checking password and username
	if !status || userCredentials.Username != user.User_Name {
		c.JSON(401, gin.H{
			"error": "Unautharized access Please check username or password",
		})
		return
	}
	//generating token with jwt
	token, err := helper.GenerateJWTToken(user.User_Name, "user", user.Email, int(user.User_ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	//set jwt in browser
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("jwt_token", token, 3600*24, "", "", true, true)
	//success message
	c.JSON(200, gin.H{
		"message": user.User_Name + " successfully logged",
	})
}

// Logout logs out the currently authenticated user.
// @Summary User Logout
// @Description Log out the currently authenticated user by clearing the JWT token cookie.
// @Tags authentication
// @Security ApiKeyAuth
// @Success 200 {string} SuccessResponse "Logout successful"
// @Router /user/logout [get]
func Logout(c *gin.Context) {
	//get data from user with middleware
	user, _ := c.Get("user")
	//set cookie to nill and expiration to -1
	c.SetCookie("jwt_token", "", -1, "", "", false, false)
	c.JSON(200, gin.H{
		"message": user.(models.User).User_Name + " successfully logout",
	})
}

// ListProducts returns a paginated list of products with additional details.
// @Summary List Products
// @Description Get a paginated list of products including product name, description, stock, price, brand name, and image.
// @Tags products
// @Produce json
// @Param page query int true "Page number for pagination (1-based)"
// @Param limit query int true "Number of products to return per page"
// @Success 200 {json} SuccessResponse
// @Failure 400 {json} JSON "Invalid page or limit values"
// @Failure 404 {json} JSON "No products found"
// @Router /products [get]
func LisProducts(c *gin.Context) {
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
	type datas struct {
		Product_name string
		Descreption  string
		Stock        int
		Price        int
		Brand_name   string
		Image        string
	}
	var products []datas
	var count int64
	var product models.Product

	result := database.DB.Table("products").Select("products.product_name,products.descreption,products.stock,products.price,brands.brand_name,images.image").
		Joins("INNER JOIN brands ON brands.brand_id=products.brand_id").Joins("INNER JOIN images ON images.product_id=products.id").
		Limit(limit).Offset(offset).
		Scan(&products)
	_ = database.DB.Find(&product).Count(&count)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Total Products Found": count,
	})

	c.JSON(http.StatusOK, gin.H{
		"Products": products,
	})
}

// ProductDetail returns details of a specific product by its name.
// @Summary Get Product Details
// @Description Get details of a product including its name, price, stock, and description.
// @Tags products
// @Produce json
// @Param product_name query string true "Name of the product to retrieve details for"
// @Success 200 {json} SuccessResponse "Product details"
// @Failure 400 {string} ErrorResponse
// @Failure 404 {string} ErrorResponse
// @Router /user/listproductsquery [get]
func ProductDetail(c *gin.Context) {

	ProductName := c.Query("product_name")
	type ProductDetails struct {
		Product_Name string
		Price        int
		Stock        int
		Descreption  string
	}
	//var data Data
	var productDetails ProductDetails
	var product []ProductDetails

	result := database.DB.Table("products").Select("product_name", "price", "stock", "descreption").Where("product_name = ?", ProductName).Find(&productDetails)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"Message": "Product not found",
		})
		return
	}
	product = append(product, productDetails)
	c.JSON(http.StatusOK, gin.H{
		"Product details": product,
	})
}
