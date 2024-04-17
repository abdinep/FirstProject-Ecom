package routers

import (
	controllers "ecom/controllers/Admin"
	handlers "ecom/handlers/Admin"
	"ecom/middleware"

	"github.com/gin-gonic/gin"
)

var roleAdmin = "Admin"

func AdminGroup(r *gin.RouterGroup) {

	r.POST("/login", controllers.Login)
	r.GET("/logout", controllers.Admin_Logout)
	//=========================== Admin user management ======================================

	r.GET("/usermanagement", middleware.JwtMiddleware(roleAdmin), handlers.List_user)
	r.PATCH("/usermanagement/:ID", middleware.JwtMiddleware(roleAdmin), handlers.Edit_User)
	r.PATCH("/usermanagement/block/:ID", middleware.JwtMiddleware(roleAdmin), handlers.Status)

	//=========================== Admin Coupon management ======================================

	r.POST("/coupon",middleware.JwtMiddleware(roleAdmin), handlers.Coupon)
	r.GET("/coupon",middleware.JwtMiddleware(roleAdmin),handlers.ListCoupon)
	r.DELETE("/coupon/:ID",middleware.JwtMiddleware(roleAdmin),handlers.DeleteCoupon)

	//=========================== Admin Product management ===================================

	r.POST("/products", middleware.JwtMiddleware(roleAdmin), handlers.Add_Product)
	r.GET("/productlist", middleware.JwtMiddleware(roleAdmin), handlers.View_Product)
	r.PATCH("/products/:ID", middleware.JwtMiddleware(roleAdmin), handlers.Edit_Product)
	r.DELETE("/products/:ID", middleware.JwtMiddleware(roleAdmin), handlers.Delete_Product)

	//=========================== Admin Category Management ==================================

	r.POST("/category/addcategory", middleware.JwtMiddleware(roleAdmin), handlers.Category)
	r.GET("/category", middleware.JwtMiddleware(roleAdmin), handlers.View_Category)
	r.PATCH("/category/edit/:ID", middleware.JwtMiddleware(roleAdmin), handlers.Edit_Category)
	r.DELETE("/category/delete/:ID", middleware.JwtMiddleware(roleAdmin), handlers.Delete_Category)
	// r.PATCH("/admin_panel/products/Recover/:ID",handlers.Undelete_Product)

	//=========================== Admin Order Management =======================================

	r.GET("/order",middleware.JwtMiddleware(roleAdmin), handlers.Admin_View_order)
	r.GET("/order/details/:ID",middleware.JwtMiddleware(roleAdmin),handlers.ViewOrderDetails)
	r.POST("/order/:ID",middleware.JwtMiddleware(roleAdmin),handlers.ChangeOrderStatus)

	//============================ Admin Offer Management ======================================

	r.POST("/offer/:ID",middleware.JwtMiddleware(roleAdmin),handlers.AddOffer)
	r.GET("/offer",middleware.JwtMiddleware(roleAdmin),handlers.ViewOffer)

	//============================= Sales Report ================================================

	r.GET("/SalesReport/downloadexcel",middleware.JwtMiddleware(roleAdmin),handlers.GenerateSalesReport)
	r.GET("/SalesReport/downloadpdf",middleware.JwtMiddleware(roleAdmin),handlers.SalesReportPDF)
	// r.GET("/dailyreport",middleware.JwtMiddleware(roleAdmin),handlers.SearchReport)

	//============================= Best Seller =================================================
	r.GET("/bestsell",middleware.JwtMiddleware(roleAdmin),handlers.BestSelling)
}
