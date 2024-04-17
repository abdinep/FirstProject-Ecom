package handlers

import (
	"crypto/rand"
	Paymentgateways "ecom/PaymentGateways"
	handlers "ecom/handlers/Admin"
	"ecom/initializers"
	"ecom/models"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

var Couponcheck models.Coupon

type checkoutdetails struct {
	AddressID     int    `json:"address_id"`
	CouponCode    string `json:"coupon_code"`
	PaymentMethod string `json:"payment_method"`
}

// @Summary Checkout and place an order
// @Description Checkout the items in the user's cart, apply discounts, handle payment, and place the order
// @Tags User-Checkout
// @Accept json
// @Produce json
// @Param requestBody body checkoutdetails true "Order details"
// @Success 200 {json} json "Order Placed Successfully"
// @Failure 401 {json} json "Failed to handle payment submission"
// @Failure 401 {json} json "Failed to place order"
// @Router /checkout [post]
func Checkout(c *gin.Context) {
	var order checkoutdetails
	var orderItems models.OrderItem
	var cart []models.Cart
	var address models.Address
	// var Payment models.Payment
	var Grandtotal int

	userid := c.GetUint("userID")
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(401, gin.H{
			"error":  "failed to fetch data",
			"status": 401,
		})
		return
	}
	//========================== Cart Details ==========================================
	if err := initializers.DB.Joins("Product").Where("user_id = ?", userid).Find(&cart); err.Error != nil {
		c.JSON(401, gin.H{
			"error":  "Failed to fetch data from Cart DB",
			"status": 401,
		})
		fmt.Println("Failed to fetch data from Cart DB=====>", err.Error)
		return
	}
	// ========================= Stock Check ============================================
	var pname string
	for _, view := range cart {

		quantity_price := int(view.Quantity) * (view.Product.Price - int(handlers.OfferCalc(view.Product_Id, c)))

		if int(view.Quantity) > view.Product.Quantity {
			pname = view.Product.Product_Name
			fmt.Println("pname------->", pname)
			c.JSON(200, gin.H{
				"message": "Product Out of Stock" + pname,
				"status":  200,
			})
			return
		}
		Grandtotal += quantity_price
	}

	//========================= Cheking Coupon =========================================
	var coup string
	if order.CouponCode != "" {

		if err := initializers.DB.Where("code = ? AND coundition < ?", order.CouponCode, Grandtotal).First(&Couponcheck); err.Error != nil {
			fmt.Println("code===>", order.CouponCode, "condition====>", Couponcheck.Coundition, "Grandtotal=====>", Grandtotal, "couponcheck=====>", Couponcheck)
			c.JSON(401, gin.H{
				"Error":  "Invalid Coupon",
				"status": 401,
			})
			coup = "No coupon added"
			fmt.Println("Invalid Coupon=====>", err.Error)
			return
		} else {
			Grandtotal -= int(Couponcheck.Discount)
			coup = order.CouponCode
			// c.JSON(200, "Coupon Added")
			fmt.Println("check==================>", Couponcheck.Code, order.CouponCode)
		}
	}
	fmt.Println("coupon=====>", Couponcheck.Code, order.CouponCode)
	//========================== Address choosing ======================================
	if err := initializers.DB.Where("user_id = ? AND id = ?", userid, order.AddressID).First(&address); err.Error != nil {
		c.JSON(401, gin.H{
			"error":  "Address not found",
			"status": 401,
		})
		fmt.Println("Address not found=======>", err.Error)
		return
	}
	//========================= Creating Random OrderID ================================
	const charset = "123456789"
	randomBytes := make([]byte, 4)
	_, err := rand.Read(randomBytes)
	if err != nil {
		fmt.Println(err)
	}
	for i, b := range randomBytes {
		randomBytes[i] = charset[b%byte(len(charset))]
	}
	orderIdstring := string(randomBytes)
	orderId, _ := strconv.Atoi(orderIdstring)
	fmt.Println("------>", orderId)
	//========================== Transaction ==============================================

	tx := initializers.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	//======================== COD checking ================================================
	if order.PaymentMethod == "COD" {
		if Grandtotal > 1000 {
			c.JSON(401, gin.H{
				"Error":  "COD not applicable for payment more than 1000",
				"status": 401,
			})
			return
		}
	}
	DeliveryCharge := 0
	if Grandtotal < 1000 && Grandtotal > 0 {
		DeliveryCharge = 40
		Grandtotal += DeliveryCharge
	}
	//======================== Payment Gateway ==============================================
	fmt.Println("orderid==>", orderId, "grandtotal==>", Grandtotal)
	if order.PaymentMethod == "UPI" {
		OrderPaymentID, err := Paymentgateways.HandlePaymentSubmission(orderId, Grandtotal)
		if err != nil {
			c.JSON(401, gin.H{
				"error":  err,
				"status": 401,
			})
			tx.Rollback()
			return
		} else {
			c.JSON(200, gin.H{
				"message":   "Continue to Payment",
				"paymentID": OrderPaymentID,
				"status":    200,
			})
			return
		}
		fmt.Println("orderpayment:==>", OrderPaymentID)
		fmt.Println("receipt====>", orderId)
		if err := tx.Create(&models.Payment{
			OrderID:       OrderPaymentID,
			Receipt:       orderId,
			PaymentStatus: "not done",
			PaymentAmount: Grandtotal,
		}); err.Error != nil {
			c.JSON(401, gin.H{
				"Error":  "Failed to upload payment details",
				"status": 401,
			})
			fmt.Println("Failed to upload payment details", err.Error)
			tx.Rollback()
		}
	}
	//========================= Order Table management ====================================
	// ID,_ := strconv.Atoi(userid)
	orderdata := models.Order{
		PaymentMethod:  order.PaymentMethod,
		AddressID:      order.AddressID,
		CouponCode:     coup,
		OrderPrice:     Grandtotal,
		DeliveryCharge: DeliveryCharge,
		OrderDate:      time.Now(),
		UserID:         int(userid),
	}
	if err := tx.Create(&orderdata); err.Error != nil {
		tx.Rollback()
		c.JSON(401, gin.H{
			"error":  "Failed to Place Order",
			"status": 401,
		})
		fmt.Println("Failed to Place Order=====>", err.Error)
		return
	}
	for _, view := range cart {
		subTotal := int(view.Quantity) * view.Product.Price
		orderItems = models.OrderItem{
			ProductID:     view.Product_Id,
			OrderID:       uint(orderId),
			OrderQuantity: int(view.Quantity),
			Subtotal:      float64(subTotal),
			Orderstatus:   "Pending",
		}
		if err := tx.Create(&orderItems); err.Error != nil {
			tx.Rollback()
			c.JSON(401, gin.H{
				"error":  "Failed to Place Order",
				"status": 401,
			})
			fmt.Println("Failed to Place Order=====>", err.Error)
			return
		}
		//========================= Stock management =======================================
		view.Product.Quantity -= int(view.Quantity)
		if err := initializers.DB.Save(&view.Product); err.Error != nil {
			c.JSON(401, gin.H{
				"error":  "Failed to update product stock",
				"status": 401,
			})
			fmt.Println("Failed to update product stock======>", err.Error)
		} else {
			fmt.Println("Stock Updated=====>")
		}
	}

	if err := initializers.DB.Where("user_id = ?", userid).Delete(&models.Cart{}); err.Error != nil {
		c.JSON(401, gin.H{
			"error":  "Failed to delete order",
			"status": 401,
		})
		return
	}
	if err := tx.Commit(); err.Error != nil {
		tx.Rollback()
		c.JSON(401, gin.H{
			"error":  "Failed to commit transaction",
			"status": 401,
		})
		fmt.Println("Failed to commit transaction=====>", err.Error)
		return
	}
	if order.PaymentMethod != "UPI" {
		c.JSON(200, gin.H{
			"Delivery Charges": DeliveryCharge,
			"message":          "Order Placed Succesfully",
			"Grand Total ":     Grandtotal,
			"status":           200,
		})
	}

}
