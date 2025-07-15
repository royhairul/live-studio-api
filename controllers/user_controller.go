package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/royhairul/live-studio-api/database"
	"github.com/royhairul/live-studio-api/dto"
	"github.com/royhairul/live-studio-api/models"
	"github.com/royhairul/live-studio-api/services/host"
	"golang.org/x/crypto/bcrypt"
)

func UserIndex(c *gin.Context) {
	// 1. Ambil ID dari URL dan konversi ke uint

	id, exists := c.Get("superadmin_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	// 2. Ambil user_id dari user_relation di mana parent_id = parentID
	var relatedUserIDs []uint
	if err := database.DB.
		Table("user_relations").
		Where("parent_id = ?", id).
		Pluck("child_id", &relatedUserIDs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil relasi user", "details": err.Error()})
		return
	}

	// 3. Ambil data user berdasarkan relatedUserIDs
	var users []models.User
	if err := database.DB.
		Preload("Role").
		Where("id IN (?)", relatedUserIDs).
		Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": users})
}

func UserCreate(c *gin.Context) {
	var req dto.CreateUserDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Ambil ID user yang membuat
	superadminID, exists := c.Get("superadmin_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}
	currentUserID, ok := superadminID.(uint)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format ID tidak sesuai"})
		return
	}

	// Buat user baru
	user := models.User{
		Name:     req.Name,
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		RoleID:   req.RoleID,
	}

	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Refresh user ID (opsional, untuk yakin)
	if err := tx.Preload("Role").First(&user, user.ID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Println(user.ID)

	if user.Role.Name == "host" {
		dto := dto.CreateHostDTO{
			Name:     req.Name,
			Phone:    req.Phone,
			StudioID: uint16(req.StudioID),
			UserID:   uint16(user.ID),
		}

		if err := host.CreateHost(&dto); err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create host record"})
			return
		}
	}

	relation := models.UserRelation{
		ParentID: currentUserID,
		ChildID:  user.ID,
		Relation: "created",
	}

	if err := tx.Create(&relation).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create relation"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"message": "User created", "data": user})
}

func UserShow(c *gin.Context) {
	// Ambil ID superadmin dari context
	superadminID, exists := c.Get("superadmin_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	// Ambil ID user dari parameter URL
	idParam := c.Param("id")
	userID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID user tidak valid"})
		return
	}

	// Pastikan user tersebut memang terkait dengan superadmin yang sedang login
	var relation models.UserRelation
	if err := database.DB.Where("parent_id = ? AND child_id = ?", superadminID, userID).First(&relation).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "User tidak ditemukan atau bukan milik Anda"})
		return
	}

	// Ambil data user beserta role-nya
	var user models.User
	if err := database.DB.Preload("Role").First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
		return
	}

	// Siapkan respons
	response := gin.H{
		"user": user,
	}

	// Jika user adalah host, cari data host terkait
	if user.Role.Name == "host" {
		var host models.Host
		if err := database.DB.Where("user_id = ?", user.ID).First(&host).Error; err == nil {
			response["host"] = host
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

func UserUpdate(c *gin.Context) {
	var req dto.UpdateUserDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ambil ID user dari URL
	userID, exist := c.Get("superadmin_id")
	if !exist {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	// Ambil data user berdasarkan ID
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
		return
	}

	// Update data user
	user.Username = req.Username
	user.Email = req.Email
	user.RoleID = req.RoleID

	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		user.Password = string(hashedPassword)
	}

	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui user", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated", "data": user})
}

func UserDelete(c *gin.Context) {
	// Ambil ID superadmin dari context
	superadminID, exists := c.Get("superadmin_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	// Ambil ID user dari parameter URL
	idParam := c.Param("id")
	userID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID user tidak valid"})
		return
	}

	// Pastikan user tersebut memang terkait dengan superadmin yang sedang login
	var relation models.UserRelation
	if err := database.DB.Where("parent_id = ? AND child_id = ?", superadminID, userID).First(&relation).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "User tidak ditemukan atau bukan milik Anda"})
		return
	}

	// Hapus user berdasarkan ID
	if err := database.DB.Delete(&models.User{}, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus user", "details": err.Error()})
		return
	}

	// Hapus juga relasi user_relations terkait (opsional tapi direkomendasikan)
	_ = database.DB.Where("child_id = ?", userID).Delete(&models.UserRelation{})

	c.JSON(http.StatusOK, gin.H{"message": "User berhasil dihapus"})
}
