package handlers

import (
	handlers "ecom/handlers/Admin"
	"ecom/initializers"
	"ecom/models"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary Add	 new product to the database
// @Description add a new product with given name, price and quantity in the cart
// @Tags User-CartManagement
// @Accept json
// @Produce  json
// @Param ID path int true "Product ID"
// @Success 200 {json} JSON	"New product added to cart"
// @Failure 401 {json} JSON "failed to add to cart"
// @Router /cart/{ID} [post]
func Add_Cart(c *gin.Context) {
	var cart models.Cart
	var dbcart models.Cart

	id := c.Param("ID")
	fmt.Println("id-------->", id)
	Product_Id, _ := strconv.Atoi(id)
	cart.Product_Id = Product_Id
	cart.User_id = int(c.GetUint("userID"))
	fmt.Println("productidd------->", cart.Product_Id, cart.User_id)
	if err := initializers.DB.Where("product_id = ? AND user_id = ?", cart.Product_Id, cart.User_id).First(&dbcart); err.Error != nil {
		cart.Quantity = 1
		if err := initializers.DB.Create(&cart); err.Error != nil {
			c.JSON(401, gin.H{
				"error":  "failed to add to cart",
				"status": 401,
			})
			fmt.Println("failed to add to cart=====>", err)
		} else {
			c.JSON(200, gin.H{
				"message": "New product added to your cart",
				"status":  200,
			})

		}
		fmt.Println("error---------->", err.Error)
		fmt.Println("-------->", dbcart.Product_Id)
	} else {
		c.JSON(401, gin.H{
			"error":  "Product already added to Cart",
			"status": 401,
		})
	}
}

// @Summary Add quantity to a product in the cart
// @Description Add one more quantity to a product in the user's cart
// @Accept json
// @Produce json
// @Tags User-CartManagement
// @Param ID path string true "Product ID"
// @Success 200 {json} json "Added one more quantity"
// @Failure 500 {json} json "Cant add more quantity or No stock"
// @Router /cart/addquantity/{ID} [patch]
func Add_Quantity_Cart(c *gin.Context) {
	var add models.Cart
	// var userid models.Cart
	var Product models.Product
	id := c.Param("ID")
	userid := c.GetUint("userID")
	if err := initializers.DB.First(&Product, id); err.Error != nil {
		c.JSON(401, gin.H{
			"error":  "Failed to fetch data from product DB",
			"status": 401,
		})
		fmt.Println("Failed to fetch data from product DB", err.Error)
		return
	}
	if err := initializers.DB.Where("product_id = ? AND user_id = ? ", id, userid).First(&add); err.Error != nil {
		c.JSON(401, gin.H{
			"error":  "Failed to fetch data",
			"status": 401,
		})
		fmt.Println("Failed to fetch data=====>", err.Error)
	} else {
		add.Quantity += 1
		if add.Quantity <= 5 && add.Quantity <= uint(Product.Quantity) {
			if err := initializers.DB.Model(&add).Updates(models.Cart{Quantity: add.Quantity}); err.Error != nil {
				c.JSON(401, gin.H{
					"error":  "Failed to update quantity ",
					"status": 401,
				})
				fmt.Println("Failed to update quantity", err.Error)
			} else {
				c.JSON(200, gin.H{
					"Quantity": add.Quantity,
					"Message":  "Added one more quantity",
					"status":   200,
				})
			}
		} else {
			c.JSON(401, gin.H{
				"error":  "Cant add more quantity or No stock",
				"status": 401,
			})
		}
	}
}

// @Summary Remove quantity from a product in the cart
// @Description Remove one quantity from a product in the user's cart
// @Accept json
// @Produce json
// @Tags User-CartManagement
// @Param ID path string true "Product ID"
// @Success 200 {json} json "Removed one more quantity"
// @Failure 500 {json} json "Failed to fetch data" or "Failed to update quantity"
// @Router /cart/removequantity/{ID} [patch]
func Remove_Quantity_cart(c *gin.Context) {
	// var remove models.Cart
	var dbremove models.Cart
	id := c.Param("ID")
	userid := c.GetUint("userID")
	if err := initializers.DB.Where("product_id = ? AND user_id = ? ", id, userid).First(&dbremove); err.Error != nil {
		c.JSON(401, gin.H{
			"error":  "Failed to fetch data",
			"status": 401,
		})
		fmt.Println("Failed to fetch data=====>", err.Error, userid, id)
	} else {
		dbremove.Quantity -= 1
		if err := initializers.DB.Model(&dbremove).Updates(models.Cart{Quantity: dbremove.Quantity}); err.Error != nil {
			c.JSON(401, gin.H{
				"error":  "Failed to update quantity",
				"status": 401,
			})
			fmt.Println("Failed to update quantity", err.Error)
		} else {
			c.JSON(200, gin.H{
				"Quantity": dbremove.Quantity,
				"Message":  "removed one more quantity",
				"status":   200,
			})
		}
	}
}

// @Summary View user's cart
// @Description View the products in the user's cart along with discounts and total amount
// @Accept json
// @Produce json
// @Tags User-CartManagement
// @Success 200 {json} json "User's cart"
// @Failure 500 {json} json "Product not found"
// @Router /cart [get]
func View_Cart(c *gin.Context) {
	var cart []models.Cart
	var listcart []gin.H
	var quantity_price int
	var Grandtotal = 0
	id := c.GetUint("userID")
	if err := initializers.DB.Joins("Product.Offer").Where("user_id = ?", id).Find(&cart); err.Error != nil {
		c.JSON(401, gin.H{
			"error":  "Product not found",
			"status": 401,
		})
		fmt.Println("product not found=====>", err.Error)
		return
	}

	count := 0
	var offer float64
	// id, _ := strconv.Atoi(id)
	discount := 0
	for _, view := range cart {
		if view.User_id == int(id) {
			quantity_price = int(view.Quantity) * (view.Product.Price - int(handlers.OfferCalc(view.Product_Id, c)))
			fmt.Println("quantityprice---------->", quantity_price)
			Grandtotal += quantity_price
			fmt.Println("Grandtotal------------->", Grandtotal)
			count += 1
			discount += int(view.Quantity) * int(handlers.OfferCalc(view.Product_Id, c))
			fmt.Println("discount--------------->", discount)
			offer += offer
			listcart = append(listcart, gin.H{
				"product_id":     view.Product_Id,
				"product_name":   view.Product.Product_Name,
				"quantity":       view.Quantity,
				"price":          view.Product.Price,
				"totalCartItems": count,
			})

		}
	}
	c.JSON(200, gin.H{
		"data":       listcart,
		"discount":   discount,
		"grandTotal": Grandtotal - int(offer),
	})
}

// @Summary Remove a product from the user's cart
// @Description Remove a product from the user's cart based on the product ID
// @Accept json
// @Produce json
// @Tags User-CartManagement
// @Param ID path string true "Product ID"
// @Success 200 {json} json "Product removed from cart"
// @Failure 500 {json} json "Failed to fetch data" or "Can't delete the product"
// @Router /cart/{ID} [delete]
func Remove_Cart_Product(c *gin.Context) {
	var remove models.Cart
	id := c.Param("ID")
	userid := c.GetUint("userID")
	if err := initializers.DB.Where("product_id = ? AND user_id = ? ", id, userid).First(&remove); err.Error != nil {
		c.JSON(500, "Failed to fetch data")
		fmt.Println("Failed to fetch data=====>", err.Error)
	} else {
		if err := initializers.DB.Delete(&remove); err.Error != nil {
			c.JSON(500, "Cant delete the product")
			fmt.Println("Cant delete the product=====>", err.Error)
		} else {
			c.JSON(200, gin.H{
				"message": "Product removed from cart",
				"status":  200,
			})
		}
	}
}
