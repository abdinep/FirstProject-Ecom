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
		return
	}
	if err := initializers.DB.Create(&models.Coupon{
		Code:       coupon.Code,
		Discount:   coupon.Discount,
		Coundition: coupon.Coundition,
		ValidFrom:  coupon.ValidFrom,
		ValidTo:    coupon.ValidTo,
	}); err.Error != nil {
		c.JSON(401, gin.H{
			"error":  "Coupon already exist",
			"status": 401,
		})
		fmt.Println("Coupon already exist", err.Error)
	} else {
		c.JSON(200, gin.H{
			"message": "New Coupon added",
			"status":  200,
		})
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
	var couponList []gin.H
	if err := initializers.DB.Find(&coupon); err.Error != nil {
		c.JSON(401, gin.H{
			"error":  err.Error,
			"status": 401,
		})
		return
	}
	for _, coupon := range coupon {
		couponList = append(couponList, gin.H{
			"Id":        coupon.ID,
			"Code":      coupon.Code,
			"Discount":  coupon.Discount,
			"Condition": coupon.Coundition,
			"CreatedAt": coupon.CreatedAt,
			"UpdatedAt": coupon.UpdatedAt,
		})
	}
	c.JSON(200, gin.H{
		"data":   couponList,
		"status": 200,
	})
}

// @Summary Delete Coupon
// @Description Admin  can Delete coupons with selected Coupon ID
// @Tags Admin-CouponManagement
// @Accept json
// @Produce  json
// @Param ID path int true "Coupon ID"
// @Success 200 {json} json	"Deleted Coupons"
// @Failure 401 {json} json "Coupon not Found or Failed to Delete Coupon"
// @Router /admin/coupon/{ID} [delete]
func DeleteCoupon(c *gin.Context) {
	var coup models.Coupon
	couponId := c.Param("ID")
	if err := initializers.DB.Where("id = ?", couponId).Delete(&coup); err.Error != nil {
		c.JSON(400, gin.H{
			"error":  "Coupon not found",
			"status": 400,
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "Coupon Deleted",
		"status":  200,
	})

}
