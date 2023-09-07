package routes

import (
	"github.com/abdulmanafc2001/First-Project-Ecommerce/controllers"
	"github.com/abdulmanafc2001/First-Project-Ecommerce/middleware"
	"github.com/gin-gonic/gin"
)

func AdminRoutes(r *gin.Engine) {
	router := r.Group("/admin")
	{
		//login and logout admin
		router.POST("/login", controllers.AdminLogin)//
		router.GET("/logout", middleware.AdminAuth, controllers.AdminLogout)//

		//user controllers
		router.GET("/viewusers", middleware.AdminAuth, controllers.ListUsers)//
		router.POST("/blockuser/:user_id", middleware.AdminAuth, controllers.BlockUser)//
		router.POST("/unblockuser/:user_id", middleware.AdminAuth, controllers.UnblockUser)//

		//catagory controllers
		router.POST("/addcatagory", middleware.AdminAuth, controllers.AddCatagory)//
		router.GET("/viewcatagories", middleware.AdminAuth, controllers.ListCatagories)//
		router.POST("/blockcatagory/:catagory_id", middleware.AdminAuth, controllers.BlockCatagory)//
		router.POST("/unblockcatagory/:catagory_id", middleware.AdminAuth, controllers.UnBlockCatagory)//

		//product controllers
		router.POST("/addproduct", middleware.AdminAuth, controllers.AddProduct)
		router.PUT("/editproduct", middleware.AdminAuth, controllers.EditProduct)
		router.POST("/deleteproduct/:product_id", middleware.AdminAuth, controllers.DeleteProduct)
		router.POST("/addimage", middleware.AdminAuth, controllers.AddImage)

		//add brands
		router.POST("/addbrand", middleware.AdminAuth, controllers.AddBrand)

		//order management
		router.GET("/viewallorders", middleware.AdminAuth, controllers.ViewOrders)//
		router.POST("/cancelorder/:order_id", middleware.AdminAuth, controllers.CancelOrder)//
		router.PATCH("/changestatus/:order_id", middleware.AdminAuth, controllers.ChangeStatus) //

		//download sales report
		router.GET("/salesreport", middleware.AdminAuth, controllers.SalesReport)
		router.GET("/salesreport/xlsx", middleware.AdminAuth, controllers.DownloadExel)
		router.GET("/salesreport/pdf", middleware.AdminAuth, controllers.Downloadpdf)

		//add coupoun
		router.POST("/addcoupon", middleware.AdminAuth, controllers.AddCoupon)
		router.GET("/listcoupons", middleware.AdminAuth, controllers.ListCoupons)
		router.PUT("/cancelcoupon/:coupon_id", middleware.AdminAuth, controllers.CancelCoupon)

		router.GET("/dashboard", middleware.AdminAuth, controllers.AdminDashboard)

		//catagoryoffers
		router.POST("/addcatagoryoffer", middleware.AdminAuth, controllers.AddCatagoryOffer)//
		router.GET("/listcatagoryoffer", middleware.AdminAuth, controllers.ListCatagoryOffer)//
		router.POST("/cancellcatagoryoffer/:offer_id", middleware.AdminAuth, controllers.DeleteCatagoryOffer)//

	}
}
