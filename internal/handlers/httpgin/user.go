package httpgin

import (
	"csTrade/internal/domain/user"
	"csTrade/internal/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (uh *UserHandler) CreateUser(c *gin.Context) {
	var req user.UserCreateReq
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error1": err.Error()})
		return
	}

	err := uh.service.CreateUser(c.Request.Context(), &req)
	if err != nil {
		c.JSON(400, gin.H{"error2": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "ok"})
}
