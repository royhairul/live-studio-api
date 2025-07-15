package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/royhairul/live-studio-api/dto"
	"github.com/royhairul/live-studio-api/models"
	"github.com/royhairul/live-studio-api/services/role"
)

var validate = validator.New()

func RoleIndex(c *gin.Context) {
	roles, err := role.GetAllRole()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to retrieve roles",
			"error":   err.Error(),
		})
		return
	}

	// Struct untuk response yang mencakup Role + permission_names tambahan
	type RoleWithPermissions struct {
		ID              uint                `json:"id"`
		Name            string              `json:"name"`
		Permissions     []models.Permission `json:"permissions"`
		PermissionNames []string            `json:"permission_names"`
	}

	var result []RoleWithPermissions

	for _, r := range *roles {
		var permissionNames []string
		for _, p := range r.Permissions {
			permissionNames = append(permissionNames, p.Name)
		}

		result = append(result, RoleWithPermissions{
			ID:              r.ID,
			Name:            r.Name,
			Permissions:     r.Permissions,
			PermissionNames: permissionNames,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Roles retrieved successfully",
		"data":    result,
	})
}

func RoleCreate(c *gin.Context) {
	var req dto.CreateRoleDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Request",
			"error":   err.Error(),
		})
		return
	}

	if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid validation",
			"error":   err.Error(),
		})
		return
	}

	if err := role.CreateRole(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create role",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Role successfully created",
	})
}

func RoleShow(c *gin.Context) {
	id := c.Param("id")

	role, err := role.GetRoleByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get role",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Role found",
		"data":    role,
	})
}

func RoleUpdate(c *gin.Context) {
	id := c.Param("id")

	var req dto.UpdateRoleDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Request",
			"error":   err.Error(),
		})
		return
	}

	if err := role.UpdateRole(id, &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to update role",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Role successfully updated",
	})
}

func RoleDelete(c *gin.Context) {
	id := c.Param("id")

	if err := role.DeleteRole(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get role",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Role successfully deleted",
	})
}
