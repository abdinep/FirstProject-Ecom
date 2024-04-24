package handlers

import (
	handlers "ecom/handlers/Admin"
	"ecom/initializers"
	"ecom/models"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// @Summary View orders placed by user
// @Description View orders placed by the authenticated user
// @Tags User-OrderManagement
// @Accept json
// @Produce json
// @Success 200 {json} json "Orders listed successfully"
// @Failure 401 {json} json "Failed to fetch order data"
// @Router /user/profile/order [get]
func View_Orders(c *gin.Context) {
	var order []models.Order
	var listOrder []gin.H
	userID := c.GetUint("userID")
	if err := initializers.DB.Preload("User").Preload("Address").Where("user_id = ?", userID).Find(&order); err.Error != nil {
		c.JSON(401, gin.H{
			"error":  "Failed to fetch order data",
			"status": 401,
		})
		fmt.Println("Failed to fetch order data=====>", err.Error)
		return
	}
	for _, view := range order {
		listOrder = append(listOrder, gin.H{
			"order_id":         view.ID,
			"payment_method":   view.PaymentMethod,
			"order_date":       view.OrderDate,
			"delivery_charger": view.DeliveryCharge,
			"coupon_code":      view.CouponCode,
			"address_id":       view.AddressID,
			"address":          view.Address.Address,
		})
	}
	c.JSON(200, gin.H{
		"data":   listOrder,
		"status": 200,
	})

}

// @Summary View order details placed by user
// @Description View order details placed by the authenticated user
// @Tags User-OrderManagement
// @Accept json
// @Produce json
// @Param ID path string true "Order ID"
// @Success 200 {json} json "Order details listed successfully"
// @Failure 401 {json} json "Failed to fetch order details"
// @Router /user/profile/orderdetails/{ID} [get]
func View_Order_Details(c *gin.Context) {
	var orderitems []models.OrderItem
	var GrandTotal int
	OrderID := c.Param("ID")

	if err := initializers.DB.Where("order_id = ?", OrderID).Preload("Product").Preload("Order").Find(&orderitems); err.Error != nil {
		c.JSON(401, gin.H{
			"error":  "Failed to fetch data",
			"status": 401,
		})
		fmt.Println("Failed to fetch data=====>", err.Error)
	} else {
		subTotal := 0
		var result float64
		var TotalOffer float64
		var orderitemlist []gin.H
		for _, view := range orderitems {
			subTotal = view.OrderQuantity * view.Product.Price
			// count += 1
			GrandTotal += subTotal
			result = (handlers.OfferCalc(view.ProductID, c)) * float64(view.OrderQuantity)
			orderitemlist = append(orderitemlist, gin.H{
				"id":             view.ID,
				"order_id":       view.OrderID,
				"product_id":     view.ProductID,
				"product_name":   view.Product.Product_Name,
				"product_price":  view.Product.Price,
				"order_quantity": view.OrderQuantity,
				"order_status":   view.Orderstatus,
				"order_date":     view.Order.OrderDate,
				"offer":          result,
				"Subtotal":       view.Subtotal,
			})
			TotalOffer += result
		}
		c.JSON(200, gin.H{
			"data":          orderitemlist,
			"discountPrice": TotalOffer,
			"grand_total":   GrandTotal - int(TotalOffer),
			"status":        200,
		})
	}
}

// @Summary Cancel order placed by user
// @Description cancel order placed by the authenticated user
// @Tags User-OrderManagement
// @Accept json
// @Produce json
// @Param ID path string true "Order ID"
// @Success 200 {json} json "Order canceled successfully"
// @Failure 401 {json} json "Failed to fetch order details"
// @Router /user/profile/order/{ID} [patch]
func Cancel_Orders(c *gin.Context) {
	var order models.Order
	var orderitem models.OrderItem
	// var coup models.Coupon
	var wallet models.Wallet
	orderID := c.Param("ID")
	fmt.Println("orderID=======>", orderID)
	if err := initializers.DB.Where("id = ?", orderID).First(&orderitem); err.Error != nil {
		c.JSON(401, gin.H{
			"error":  "Order not exist",
			"status": 401,
		})
		fmt.Println("Order not exist=======>", err.Error)
	} else {
		if orderitem.Orderstatus == "Order cancelled" {
			c.JSON(200, gin.H{
				"message": "order already cancelled",
				"status":  200,
			})
			return
		}
		canceledAmount := orderitem.Subtotal
		var paymentid models.Payment
		initializers.DB.Where("order_id = ?", orderID).First(&paymentid)
		// result := Paymentgateways.RefundCancelledAmount(paymentid.PaymentID,int(canceledAmount))
		// fmt.Println("result=======>",result)
		if err := initializers.DB.Model(&orderitem).Updates(&models.OrderItem{
			Orderstatus: "Order cancelled",
		}); err.Error != nil {
			c.JSON(401, gin.H{
				"error":  "Order not cancelled",
				"status": 401,
			})
			fmt.Println("Order not cancelled========>", err.Error)
		} else {
			c.JSON(200, gin.H{
				"message": "Order cancelled successfully",
				"status":  200,
			})
			initializers.DB.First(&order, orderitem.OrderID)
			if err := initializers.DB.Where("user_id = ?", order.UserID).First(&wallet); err.Error != nil {
				c.JSON(401, gin.H{
					"error":  "Failed to add Money to wallet",
					"status": 401,
				})
				fmt.Println("Failed to add Money to wallet======>", err.Error)
				return
			}
			//=========================== Adding cancelled amount to wallet ==========================================
			wallet.Balance += int(canceledAmount)
			wallet.Updated_at = time.Now()
			initializers.DB.Save(&wallet)

			order.OrderPrice += int(Couponcheck.Discount)
			if err := initializers.DB.Model(&order).Updates(&models.Order{
				OrderPrice: order.OrderPrice - int(canceledAmount),
			}); err.Error != nil {
				fmt.Println("Failed to Update Order Amount========>", err.Error)
			}
			// orderprice := order.OrderPrice - int(canceledAmount)
			// GrandTotal -= orderprice
			// initializers.DB.First(&coup, "code = ?", order.CouponCode)
			// if order.OrderPrice > coup.Coundition {
			// 	order.OrderPrice -= int(Couponcheck.Discount)
			// 	// GrandTotal -= int(Couponcheck.Discount)
			// 	c.JSON(200, gin.H{
			// 		"message": "Coupon applied",
			// 	})

			initializers.DB.Save(&order)
		}
	}
}
