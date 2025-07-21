package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/royhairul/live-studio-api/internal/domains/permission/service"
)

type PermissionControllerImpl struct {
	// TODO: add dependencies
	service service.PermissionService
}

func NewPermissionController(service service.PermissionService) PermissionController {
	return &PermissionControllerImpl{service}
}

// Create implements PermissionController.
func (p *PermissionControllerImpl) Create(ctx *gin.Context) {
	panic("unimplemented")
}

// Delete implements PermissionController.
func (p *PermissionControllerImpl) Delete(ctx *gin.Context) {
	panic("unimplemented")
}

// FindAll implements PermissionController.
func (p *PermissionControllerImpl) FindAll(ctx *gin.Context) {
	panic("unimplemented")
}

// FindByID implements PermissionController.
func (p *PermissionControllerImpl) FindByID(ctx *gin.Context) {
	panic("unimplemented")
}

// Update implements PermissionController.
func (p *PermissionControllerImpl) Update(ctx *gin.Context) {
	panic("unimplemented")
}
