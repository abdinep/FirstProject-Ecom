package handlers

import (
	"ecom/initializers"
	"ecom/models"
	"fmt"

	"github.com/gin-gonic/gin"
)

// @Summary List Order
// @Description List all Orders
// @Tags Admin-OrderManagement
// @Accept json
// @Produce json
// @Success 200 {json} json	"Listed all Orders"
// @Failure 401 {json} json "Order not Found or Failed to List Orders"
// @Router /admin/order [get]
func Admin_View_order(c *gin.Context) {
	var order []models.Order
	var orderData []gin.H
	count := 0
	if err := initializers.DB.Preload("Address.User").Find(&order); err.Error != nil {
		c.JSON(401, gin.H{
			"error":  "failed to fetch order",
			"status": 401,
		})
		return
	}
	for _, value := range order {
		orderData = append(orderData, gin.H{
			"id":             value.ID,
			"user":           value.Address.User.Name,
			"address":        value.Address.Address,
			"appliedCoupon":  value.CouponCode,
			"orderPrice":     value.OrderPrice,
			"paymentMethod":  value.PaymentMethod,
			"deliveryCharge": value.DeliveryCharge,
			"orderDate":      value.OrderDate,
		})
		count++
	}
	c.JSON(200, gin.H{
		"data":        orderData,
		"totalOrders": count,
		"status":      200,
	})
}

// @Summary List OrderDetails
// @Description List all OrderDtails  of a specific Order
// @Tags Admin-OrderManagement
// @Accept json
// @Produce json
// @Param ID path int true "Order ID"
// @Success 200 {json} json	"Listed all OrderDetails"
// @Failure 401 {json} json "Order not Found or Failed to List Orders"
// @Router /admin/order/details/{ID} [get]
func ViewOrderDetails(c *gin.Context) {
	var orderitem []models.OrderItem
	var orderitemsList []gin.H
	orderid := c.Param("ID")
	if err := initializers.DB.Where("order_id = ?", orderid).Joins("Product.Category").Joins("Order.Address.User").Find(&orderitem); err.Error != nil {
		c.JSON(401, gin.H{
			"error":  "Produt not found",
			"status": 401,
		})
		return
	}
	subTotal := 0
	// orderPrice := 0
	for _, view := range orderitem {
		subTotal = view.OrderQuantity * view.Product.Price
		orderitemsList = append(orderitemsList, gin.H{
			"id":              view.ID,
			"user":            view.Order.Address.User.Name,
			"orderId":         view.OrderID,
			"orderQuantity":   view.OrderQuantity,
			"productName":     view.Product.Product_Name,
			"productCategory": view.Product.Category.Name,
			"productPrice":    view.Product.Price,
			"subTotal":        view.Subtotal,
			"orderStatus":     view.Orderstatus,
		})
		subTotal += int(view.Subtotal)
	}
	c.JSON(200, gin.H{
		"data":       orderitemsList,
		"orderPrice": subTotal,
		"status":     200,
	})
}

type OrderStatus struct {
	Orderstatus string `json:"status"`
}

// @Summary Edit Order Status
// @Description Change Order status of a specific Order
// @Tags Admin-OrderManagement
// @Accept json
// @Produce json
// @Param ID path int true "Order ID"
// @Param order body models.OrderItem true "Updated order status"
// @Success 200 {json} json	"Listed all OrderDetails"
// @Failure 401 {json} json "Order not Found or Failed to List Orders"
// @Router /admin/order/{ID} [post]
func ChangeOrderStatus(c *gin.Context) {
	var order models.OrderItem
	var update models.OrderItem
	orderid := c.Param("ID")
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(401, gin.H{
			"error":  "Add status",
			"status": 401,
		})
		return
	}
	if err := initializers.DB.Where("id = ?", orderid).First(&update); err.Error != nil {
		c.JSON(401, gin.H{
			"error":  "Order not found",
			"status": 401,
		})
		fmt.Println("Order not found======>", err.Error)
		return
	}
	fmt.Println("========>", order.Orderstatus)
	update.Orderstatus = order.Orderstatus
	initializers.DB.Save(&update)
	c.JSON(200, gin.H{
		"message": "Order status changed",
		"status":  200,
	})
}
