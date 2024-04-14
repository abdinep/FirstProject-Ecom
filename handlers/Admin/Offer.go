package handlers

import (
	"ecom/initializers"
	"ecom/models"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)
type addsoffer struct{
	OfferName string  `json:"OfferName"`
    Amount    float64 `json:"Amount"`
	Expire    time.Time `json:"Expire"`
	ProductId int
	Created   time.Time
}
// @Summary Add an offer for a product
// @Description Add an offer for a specific product by providing its ID and offer details
// @Tags Admin-OfferManagement
// @Accept json
// @Produce json
// @Param ID path string true "Product ID"
// @Param offer body addsoffer true "Offer details"
// @Success 200 {json} json	"Added new Offer"
// @Failure 401 {json} json "Offer not Found or Failed to Add New Offer"
// @Router /admin/offer/{ID} [post]
func AddOffer(c *gin.Context) {
	var offer addsoffer
	var product models.Product
	productid := c.Param("ID")
	c.ShouldBindJSON(&offer)
	if err := initializers.DB.Where("id = ?", productid).First(&product); err.Error != nil {
		c.JSON(500, gin.H{"Error": "Product not available"})
		fmt.Println("Product not available======>", err.Error)
	} else {
		offer.Created = time.Now()
		offer.ProductId, _ = strconv.Atoi(productid)
		if err := initializers.DB.Create(&offer); err.Error != nil {
			c.JSON(500, gin.H{"Error": "Failed to add offer"})
			fmt.Println("Failed to add offer=====>", err.Error)
		} else {
			c.JSON(200, gin.H{"Message": "Offer Added for the Product"})
		}
	}
}
// @Summary List offer
// @Description List all offer for a specific product
// @Tags Admin-OfferManagement
// @Accept json
// @Produce json
// @Success 200 {json} json	"Listed all Offer"
// @Failure 401 {json} json "Offer not Found or Failed to List Offer"
// @Router /admin/offer [get]
func ViewOffer(c *gin.Context) {
	var offer []models.Offer
	if err := initializers.DB.Find(&offer); err.Error != nil {
		c.JSON(401, gin.H{"error": "Offer not found"})
		return
	} else {
		c.JSON(200,gin.H{
			"data":offer,
		})
	}

}

func OfferCalc(productid int, c *gin.Context) float64 {
	var offercheck models.Product
	var Discount float64 = 0
	if err := initializers.DB.Joins("Offer").First(&offercheck,"products.id = ?", productid); err.Error != nil {
		return Discount
	} else {
		percentage := offercheck.Offer.Amount
		ProductAmount := offercheck.Price
		Discount = (percentage * float64(ProductAmount)) / 100
		fmt.Println("discount:", Discount)
	}
	return Discount
}