package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/royhairul/live-studio-api/internal/domains/user/params"
	"github.com/royhairul/live-studio-api/internal/domains/user/service"
	"github.com/royhairul/live-studio-api/internal/pkg/errorhandler"
	"github.com/royhairul/live-studio-api/internal/pkg/response"
)

type userControllerImpl struct {
	service service.UserService
}

func NewUserController(s service.UserService) UserController {
	return &userControllerImpl{s}
}

func (ctrl *userControllerImpl) GetAll(c *gin.Context) {
	users, err := ctrl.service.GetAll()
	if err != nil {
		errorhandler.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.NewBaseResponse("Berhasil mengambil semua data user", users))
}

func (ctrl *userControllerImpl) GetByID(c *gin.Context) {
	id := c.Param("id")
	user, err := ctrl.service.GetByID(id)
	if err != nil {
		errorhandler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, response.NewBaseResponse("Berhasil mengambil data user", user))
}

func (ctrl *userControllerImpl) Create(c *gin.Context) {
	var req params.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorhandler.HandleError(c, err)
		return
	}
	createdUser, err := ctrl.service.Create(req)
	if err != nil {
		errorhandler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, response.NewBaseResponse("User berhasil dibuat", createdUser))
}

func (ctrl *userControllerImpl) Update(c *gin.Context) {
	id := c.Param("id")
	var req params.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewBaseResponse("Gagal memperbarui user", err))
		return
	}
	updatedUser, err := ctrl.service.Update(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewBaseResponse("Gagal memperbarui user", err))
		return
	}
	c.JSON(http.StatusOK, response.NewBaseResponse("User berhasil diperbarui", updatedUser))
}

func (ctrl *userControllerImpl) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := ctrl.service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewBaseResponse("Gagal menghapus user", err))
		return
	}
	c.JSON(http.StatusOK, response.NewBaseResponse("User berhasil dihapus", nil))
}
