package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/royhairul/live-studio-api/services/account"
	"github.com/royhairul/live-studio-api/services/shopee"
)

func FinanceGetLiveReport(c *gin.Context) {
	// Ambil query params dari request
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")
	timeDim := c.DefaultQuery("timeDim", "1d")
	endDate := c.DefaultQuery("endDate", "2025-03-01")
	// endDate := c.Query("endDate")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
		return
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pageSize parameter"})
		return
	}
	if endDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "endDate parameter is required"})
		return
	}

	// Ambil semua akun dari database/service
	accounts, err := account.GetAccountAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get accounts"})
		return
	}

	client := shopee.NewClient("https://creator.shopee.co.id", "")

	// Data gabungan dari semua akun
	var combinedResults []map[string]interface{}

	for _, acc := range *accounts {
		if acc.Cookies == "" {
			continue // skip akun tanpa cookie
		}

		resp, err := client.GetLiveList(acc.Cookies, page, pageSize, timeDim, endDate)
		if err != nil {
			continue // skip akun yang error → kalau mau bisa dikumpulkan juga error-nya
		}

		var parsed map[string]interface{}
		if err := json.Unmarshal(resp, &parsed); err != nil {
			continue // skip jika response tidak bisa di-parse
		}

		// Optional → tambahkan identitas akun agar tahu dari akun mana
		resultWithAccount := map[string]interface{}{
			"account_id": acc.ID,
			"shop_name":  acc.Name,
			"data":       parsed,
		}

		combinedResults = append(combinedResults, resultWithAccount)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    combinedResults,
	})
}
