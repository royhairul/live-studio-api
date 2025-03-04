package controllers

import (
	"live-studio-api/database"
	"live-studio-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GET Products
func ProductIndex(c *gin.Context) {
	var products []models.Product
	database.DB.Find(&products)
	
	// Check if product list is empty
	if len(products) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "List Product is empty.",
			"data": products,
		})
		return // Prevent sending multiple responses
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "List All Product",
		"data": products,
	})
}

// POST Create a Product
func ProductCreate(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to create product",
			"error": err.Error(),
		})
		return
	}
	database.DB.Create(&product)
	c.JSON(http.StatusCreated, gin.H{
		"message": "Product successfully created",
		"data": product,
	})
}

func ProductUpdate() {}
func ProductDelete() {}