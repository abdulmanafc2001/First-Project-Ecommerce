package routes

import (
	"github.com/abdulmanafc2001/First-Project-Ecommerce/controllers"
	"github.com/abdulmanafc2001/First-Project-Ecommerce/middleware"
	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine) {
	r.LoadHTMLGlob("templates/*.html")

	router := r.Group("/user")
	{
		router.POST("/signup", controllers.Signup)
		router.POST("/signup/validate", controllers.ValidateOtp)
		router.POST("/login", controllers.Login)
		router.GET("/logout", middleware.UserAuth, controllers.Logout)
		//product
		router.GET("/listproducts", controllers.LisProducts)
		router.GET("/listproductsquery", controllers.ProductDetail)

		//filter products
		router.GET("/brandfilter", controllers.FilterWithBrands)
		router.GET("/catagoryfilter", controllers.FilterWithCatagories)
		router.GET("/ascendingfilter", controllers.SortWithAscending)
		router.GET("/descendingfilter", controllers.SortWithDescending)

		//cart
		router.POST("/addtocart", middleware.UserAuth, controllers.AddToCart)//
		router.GET("/viewcart", middleware.UserAuth, controllers.ListCart)
		router.DELETE("/deletefromcart/:cart_id", middleware.UserAuth, controllers.DeleteFromCart)
		router.PUT("/updatecartquantity", middleware.UserAuth, controllers.UpdateCartQuantity)

		//address
		router.POST("/addaddress", middleware.UserAuth, controllers.AddAddress)
		router.PUT("/editaddress/:address_id", middleware.UserAuth, controllers.EditAddress)
		router.GET("/listaddresses", middleware.UserAuth, controllers.ListAddresses)

		//userdetails
		router.GET("/userdetail", middleware.UserAuth, controllers.UserDetail)
		router.PUT("/changepassword", middleware.UserAuth, controllers.ChangePassword)
		router.PUT("/updateprofile", middleware.UserAuth, controllers.UpdateProfile)

		//checkout with cod
		router.POST("/checkoutcod", middleware.UserAuth, controllers.CheckOutCOD)

		//list orders
		router.GET("/listorders", middleware.UserAuth, controllers.ListOrders)
		router.POST("/cancelorder/:order_id", middleware.UserAuth, controllers.CancelOrderWithId)

		//list orders with filters and sort
		router.GET("/listorderswithbrand", middleware.UserAuth, controllers.ListOrdersWithBrand)
		router.GET("/listorderswithcatagory", middleware.UserAuth, controllers.ListOrdersWithCatagory)
		router.GET("/listorderdesc", middleware.UserAuth, controllers.ListOrderDesc)
		router.GET("/listorderasc", middleware.UserAuth, controllers.ListOrderAsc)

		//razor pay rendering
		router.GET("/razorpay", middleware.UserAuth, controllers.RazorPay)
		router.POST("/payment/success", middleware.UserAuth, controllers.RazorPaySuccess)
		router.GET("/success", middleware.UserAuth, controllers.Success)

		//apply coupon
		router.POST("/applycoupon", middleware.UserAuth, controllers.ApplyCoupon)

		//wishlist management
		router.POST("/addtowishlist", middleware.UserAuth, controllers.AddToWishtList)
		router.GET("/listwishlist/:page", middleware.UserAuth, controllers.ListWishlist)
		router.POST("/wishlist/addtocart", middleware.UserAuth, controllers.AddTocartFromWishlist)
		router.POST("/wishlist/delete/:wishlist_id", middleware.UserAuth, controllers.RemoveFromWishlist)

		router.GET("/createinvoice", middleware.UserAuth, controllers.GenetatePdf)
		router.GET("/downloadinvoice", middleware.UserAuth, controllers.DownloadInvoice)
	}

}
