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

	c.JSON(http.StatusOK, products)
}

// POST Create a Product
func ProductCreate(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	database.DB.Create(&product)
	c.JSON(http.StatusCreated, product)
}

func ProductUpdate() {}
func ProductDelete() {}