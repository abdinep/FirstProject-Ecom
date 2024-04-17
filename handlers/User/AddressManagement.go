package handlers

import (
	"ecom/initializers"
	"ecom/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ================================ Address management ============================================
type  addressAddHandler struct {
	Address string `json:"user_address"`
    City    string `json:"user_city"`
    State   string `json:"user_state"`
    Pincode int    `json:"user_pincode"`
    Country string `json:"user_country"`
    Type    string `json:"address_type"`
}
// @Summary Add Address
// @Description User  can add a new address to their account
// @Tags User-AddressManagement
// @Accept json
// @Produce  json
// @Param address body addressAddHandler true "Address details"
// @Success 200 {json} json	"New Address added"
// @Failure 401 {json} json "Existing Address"
// @Router /user/address [post]
func Add_Address(c *gin.Context) {
	var address addressAddHandler
	if err := c.ShouldBindJSON(&address); err != nil {
		c.JSON(401, gin.H{
			"error":"failed to fetch data",
			"status":401,
		})
	}
	UserId := int(c.GetUint("userID"))
	if err := initializers.DB.Create(&models.Address{
		Address: address.Address,
		City: address.City,
		State: address.State,
		Pincode: address.Pincode,
		Country: address.Country,
		Type: address.Type,
		UserId: UserId,
	}); err.Error != nil {
		c.JSON(401,gin.H{
			"error":"Existing Address",
			"status":401,
		} )
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message":"New Address added",
			"statis":200,
		})
	}
}
type  addressUpdateHandler struct {
	Address string `json:"user_address"`
    City    string `json:"user_city"`
    State   string `json:"user_state"`
    Pincode int    `json:"user_pincode"`
    Country string `json:"user_country"`
    Type    string `json:"address_type"`
}
// @Summary Edit Address
// @Description User  can edit their address
// @Tags User-AddressManagement
// @Accept json
// @Produce  json
// @Param ID path int true "User ID"
// @Param address body addressUpdateHandler true "Address details"
// @Success 200 {json} json	"Updated Address"
// @Failure 401 {json} json "failed to edit address"
// @Router /user/address/{ID} [patch]
func Edit_Address(c *gin.Context) {
	var address addressUpdateHandler
	id := c.Param("ID")
	if err := initializers.DB.First(&address, id); err.Error != nil {
		c.JSON(401, gin.H{
			"error":"Failed to fetch address",
			"status":401,
		})
	} else {
		if err := c.ShouldBindJSON(&address); err != nil {
			c.JSON(401, gin.H{
				"error":"failed to bind address",
				"status":401,
			})
			return
		}
		if err := initializers.DB.Save(&address); err.Error != nil {
			c.JSON(401,gin.H{
				"error": "failed to edit address",
				"status":401,
			})
		} else {
			c.JSON(200, gin.H{
				"message":"Updated Address",
				"status":200,
			})
		}
	}

}
// @Summary Delete Address
// @Description User  can delete  his own address by providing the valid address id
// @Tags User-AddressManagement
// @Accept json
// @Produce  json
// @Param user_id path int true "User ID"
// @Success 200 {json} json	"Deleted Address"
// @Failure 401 {json} json "failed to delete address"
// @Router /user/address/{ID} [delete]
func Delete_Address(c *gin.Context) {
	var address models.Address
	id := c.Param("ID")
	if err := initializers.DB.First(&address, id); err.Error != nil {
		c.JSON(401,gin.H{
			"error":"failed to fetch data",
			"status":401,
		} )
	} else {
		if err := initializers.DB.Delete(&address); err.Error != nil {
			c.JSON(401, gin.H{
				"error": "Failed to delete address",
				"status": 401,
			})
		} else {
			c.JSON(200, gin.H{
				"message":"Deleted Successfully!",
				"status":200,
			})
		}
	}

}
// @Summary List Address
// @Description User can list all  addresses 
// @Tags User-AddressManagement
// @Accept json
// @Produce  json
// @Success 200 {json} json	"Listed Address"
// @Failure 401 {json} json "failed to list address"
// @Router /user/address [get]
func View_Address(c *gin.Context) {
	var address []models.Address
	var listaddress []gin.H
	userID := c.GetUint("userID")
	if err := initializers.DB.Find(&address).Where("UserId = ?", userID); err.Error != nil {
		c.JSON(401, gin.H{
			"error":"Failed to find address",
			"status":401,
		})
		fmt.Println("--------->",err.Error, address)
	} else {
		for _, view := range address {
			if view.UserId == int(userID) {

				listaddress=append(listaddress, gin.H{
					"Address_Type": view.Type,
					"Address_ID":   view.ID,
					"User_Address": view.Address,
					"User_City":    view.City,
					"User_State":   view.State,
					"User_Pincode": view.Pincode,
					"User_Country": view.Country,
					"User_Phone":   view.Phone,
				})
			}
		}
		c.JSON(200,gin.H{
			"data": listaddress,
			"status":200,
		})
	}
}

// ========================================= END ==================================================
