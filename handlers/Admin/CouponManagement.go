package handlers

import (
	"ecom/initializers"
	"ecom/models"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

type addcoupon struct {
	Code       string    `json:"code"`
	Discount   float64   `json:"discount"`
	Coundition int       `json:"condition"`
	ValidFrom  time.Time `json:"validfrom"`
	ValidTo    time.Time `json:"validto"`
}

// @Summary Add Coupon
// @Description Admin  can Add new coupons
// @Tags Admin-CouponManagement
// @Accept json
// @Produce  json
// @Param coupon body addcoupon true "Coupon details"
// @Success 200 {json} json	"Added new Coupon"
// @Failure 401 {json} json "Coupo not Found or Failed to Add New Coupon"
// @Router /admin/coupon [post]
func Coupon(c *gin.Context) {
	var coupon addcoupon
	if err := c.ShouldBindJSON(&coupon); err != nil {
		c.JSON(500, "Failed to fetch data")
	} else {
		if err := initializers.DB.Create(&coupon); err.Error != nil {
			c.JSON(500, "Coupon already exist")
			fmt.Println("Coupon already exist", err.Error)
		} else {
			c.JSON(200, "New Coupon added")
		}
	}
}
// @Summary List Coupon
// @Description Admin  can View all the coupons
// @Tags Admin-CouponManagement
// @Accept json
// @Produce  json
// @Success 200 {json} json	"Listed all Coupons"
// @Failure 401 {json} json "Coupon not Found or Failed to List Coupon"
// @Router /admin/coupon [get]
func ListCoupon(c *gin.Context) {
	var coupon []models.Coupon
	if err := initializers.DB.Find(&coupon); err.Error != nil {
		c.JSON(500, gin.H{
			"error": err.Error,
		})
		return
	}
	c.JSON(200, gin.H{
		"data": coupon,
	})
}
// @Summary Delete Coupon
// @Description Admin  can Delete coupons with selected Coupon ID
// @Tags Admin-CouponManagement
// @Accept json
// @Produce  json
// @Param id path int true "Coupon ID"
// @Success 200 {json} json	"Deleted Coupons"
// @Failure 401 {json} json "Coupon not Found or Failed to Delete Coupon"
// @Router /admin/coupon/{ID} [delete]
func DeleteCoupon(c *gin.Context) {
	var coup models.Coupon
	couponId := c.Param("ID")
	if err := initializers.DB.Where("id = ?", couponId).Delete(&coup); err.Error != nil {
		c.JSON(400, gin.H{
			"error": "Coupon not found",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "Coupon Deleted",
	})

}
