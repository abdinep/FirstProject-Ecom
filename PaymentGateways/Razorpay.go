package Paymentgateways

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"ecom/initializers"
	"ecom/models"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/razorpay/razorpay-go"
)

func HandlePaymentSubmission(orderid int, amount int) (string, error) {
	// Get payment details from the request body (replace with your actual logic)
	fmt.Println("paymentorderid==>", orderid, "paymentamount==>", amount)
	keyid := os.Getenv("KEYID")
	secretKey := os.Getenv("SECRETKEY")
	razorpayClient := razorpay.NewClient(keyid, secretKey)

	paymentDetails := map[string]interface{}{
		"amount":   amount * 100,
		"currency": "INR",
		"receipt":  strconv.Itoa(orderid),
	}
	// Create the order using Razorpay API
	order, err := razorpayClient.Order.Create(paymentDetails, nil)
	if err != nil {
		return "", err
	}
	razororderid, _ := order["id"].(string)
	fmt.Println("razororderid==>", razororderid)
	return razororderid, nil
}
func PaymentTemplate(c *gin.Context) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyaWQiOjEsImVtYWlsIjoiYWJkQGdtYWlsLmNvbSIsInVzZXJuYW1lIjoiIiwicm9sZSI6ImFkbWluIiwiZXhwIjoxNzExNjAyNjU4fQ.0HWfjGnoqXkhiHUKhQ1rYBtdQI4dJEciWuiU-lIhACg"
	fmt.Println("token----->", token)
	c.SetCookie("jwtToken", token, int((time.Hour * 1).Seconds()), "/", "localhost", false, true)
	c.HTML(http.StatusOK, "Razorpay.html", gin.H{})
}
func PaymentDetailsFromFrontend(c *gin.Context) {
	var Paymentdetails = make(map[string]string)
	var Payment models.Payment
	if err := c.ShouldBindJSON(&Paymentdetails); err != nil {
		c.JSON(500, gin.H{"Error": "Invalid Request"})
		fmt.Println("Invalid Request", err)
	}
	fmt.Println("====>", Paymentdetails)
	err := RazorPaymentVerification(Paymentdetails["signatureID"], Paymentdetails["order_Id"], Paymentdetails["paymentID"])
	if err != nil {
		fmt.Println("====>", err)
		return
	}

	fmt.Println("======", Paymentdetails["order_Id"])
	if err := initializers.DB.Where("order_id = ?", Paymentdetails["order_Id"]).First(&Payment); err.Error != nil {
		c.JSON(500, gin.H{"Error": "OrderID not found"})
		fmt.Println("OrderID not found====>", err.Error)
		return
	}
	fmt.Println("-------", Payment)
	Payment.PaymentID = Paymentdetails["paymentID"]
	Payment.PaymentStatus = "Done"
	if err := initializers.DB.Model(&Payment).Updates(&models.Payment{
		PaymentID:     Payment.PaymentID,
		PaymentStatus: Payment.PaymentStatus,
	}); err.Error != nil {
		c.JSON(500, gin.H{"Error": "Failed to update paymentID"})
		fmt.Println("Failed to update paymentID=======>", err.Error)
	} else {
		c.JSON(200, gin.H{"Message": "Payment Done"})
		fmt.Println("Payment Done")
	}
}
func RazorPaymentVerification(sign, orderId, paymentId string) error {
	secretKey := os.Getenv("SECRETKEY")
	signature := sign
	secret := secretKey
	data := orderId + "|" + paymentId
	h := hmac.New(sha256.New, []byte(secret))
	_, err := h.Write([]byte(data))
	if err != nil {
		panic(err)
	}
	sha := hex.EncodeToString(h.Sum(nil))
	if subtle.ConstantTimeCompare([]byte(sha), []byte(signature)) != 1 {
		return errors.New("PAYMENT FAILED")
	} else {
		return nil
	}
}

// func RefundCancelledAmount(paymentID string, amount int) error {
// 	// Prepare refund request body
// 	refundData := map[string]interface{}{
// 		"amount": amount,
// 	}
// 	refundBody, err := json.Marshal(refundData)
// 	if err != nil {
// 		return err
// 	}

// 	// Create HTTP client
// 	client := &http.Client{}

// 	// Prepare refund request
// 	refundURL := fmt.Sprintf("https://api.razorpay.com/v1/payments/%s/refund", paymentID)
// 	req, err := http.NewRequest("POST", refundURL, bytes.NewBuffer(refundBody))
// 	if err != nil {
// 		return err
// 	}
// 	req.Header.Set("Content-Type", "application/json")

// 	// Set Razorpay API key for authentication
// 	req.SetBasicAuth("rzp_test_CRHoZP9WQjbjhm", "0M4F9wxvzeoSpMsLuTSycQan")

// 	// Send refund request
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return err
// 	}
// 	defer resp.Body.Close()

// 	// Read response body
// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return err
// 	}

// 	// Handle response
// 	if resp.StatusCode != http.StatusOK {
// 		return fmt.Errorf("refund request failed: %s", string(body))
// 	}

// 	fmt.Println("Refund successful")
// 	return nil
// }
